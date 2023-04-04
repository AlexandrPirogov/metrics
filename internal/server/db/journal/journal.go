package journal

import (
	"bufio"
	"errors"
	"log"
	"memtracker/internal/config/server"
	"os"
	"strconv"
	"time"
)

// Newjournal returns new instance of Journal
func NewJournal() Journal {
	cfg := server.JournalCfg
	readInterval := cfg.ReadInterval[:len(cfg.ReadInterval)-1]
	read, err := strconv.Atoi(string(readInterval))
	if err != nil {
		log.Fatalf("%v", err)
	}
	return Journal{
		File:         cfg.StoreFile,
		WithRestore:  cfg.Restore,
		ReadInterval: read,
		Channel:      make(chan []byte),
	}
}

// Journal writes data from channel to give file
// Works like replication/ db.log journal
type Journal struct {
	File         string
	WithRestore  bool
	ReadInterval int
	Channel      chan []byte
}

// Start make journal stats writing data to the given file in json format
//
// Pre-cond: j have file that can be modified
//
// Post-cond: data written to the file depending on the chosen mode
// There are two modes: synch mode writes permanently data to file
// delayed mode: writes data to the file once at the given period
// Returns nil if success started otherwise returns error
func (j Journal) Start() error {
	file, err := j.openWriteFile()
	if err != nil {
		return err
	}
	if j.ReadInterval == 0 {
		go func() {
			j.writeSynch(file)
		}()
	} else {
		go func() {
			j.writeDelayed(file)
		}()
	}
	return nil
}

// readByTimer writes data once in a given period from channel
//
// Pre-cond: given file to write data
//
// Post-cond: data written to the file
func (j Journal) writeDelayed(file *os.File) {
	defer file.Close()
	writer := bufio.NewWriter(file)
	read := time.NewTicker(time.Second * time.Duration(j.ReadInterval))
	for {
		<-read.C
		for {
			if bytes, ok := <-j.Channel; ok {
				writer.Write(append(bytes, '\n'))
				writer.Flush()
			} else {
				break
			}
		}
	}
}

// NewWriter writes data every time when channel got new data
//
// Pre-cond: given file to write data
//
// Post-cond: data written to the file
func (j Journal) writeSynch(file *os.File) {
	defer file.Close()
	writer := bufio.NewWriter(file)
	for {
		if bytes, ok := <-j.Channel; ok {
			writer.Write(append(bytes, '\n'))
			writer.Flush()
		} else {
			break
		}
	}

}

// Restore write data written in file to the DB
//
// Pre-cond:
//
// Post-cond: if restore is enabled fill db with data written in given file
func (j Journal) Restore() ([][]byte, error) {
	if !j.WithRestore {
		return [][]byte{}, errors.New("restore is disabled")
	}
	file, err := j.openReadFile()
	if err != nil {
		return [][]byte{}, err
	}
	defer file.Close()
	bytes := make([][]byte, 0)
	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)
	for {
		if reader.Scan() {
			bytes = append(bytes, reader.Bytes())
		} else {
			break
		}
	}
	return bytes, nil
}

// openWriteFile opens given file to write data
//
// Pre-cond: given file that can be modified
//
// Post-cond: returns pointer to opened file and error
func (j Journal) openWriteFile() (*os.File, error) {
	file, err := os.OpenFile(j.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("can't open file %s %v", j.File, err)
		return nil, err
	}
	return file, nil
}

// openWriteFile read data from given file
//
// Pre-cond: given file that can be read
//
// Post-cond: returns pointer to opened file and error
func (j Journal) openReadFile() (*os.File, error) {
	file, err := os.OpenFile(j.File, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("can't open file %s %v", j.File, err)
		return nil, err
	}
	return file, nil
}

func (j Journal) Write(record []byte) {
	j.Channel <- record
}
