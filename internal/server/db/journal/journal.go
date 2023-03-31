package journal

import (
	"bufio"
	"errors"
	"log"
	"memtracker/internal/config"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v7"
)

func NewJournal() Journal {
	cfg := config.JournalConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error while parsing env %v", err)
	}
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

type Journal struct {
	File         string
	WithRestore  bool
	ReadInterval int
	Channel      chan []byte
}

func (j Journal) Start() {
	file, err := j.openWriteFile()
	if err != nil {
		return
	}
	if j.ReadInterval == 0 {
		j.synchRead(file)
	} else {
		j.readByTimer(file)
	}
}

func (j Journal) readByTimer(file *os.File) {
	writer := bufio.NewWriter(file)
	read := time.NewTicker(time.Second * time.Duration(j.ReadInterval))
	for {
		<-read.C
		for {
			if bytes, ok := <-j.Channel; ok {
				writer.Write(append(bytes, '\n'))
				writer.Flush()
			} else {
				writer.Flush()
				break
			}
		}
	}
}

func (j Journal) synchRead(file *os.File) {
	writer := bufio.NewWriter(file)
	for {
		if bytes, ok := <-j.Channel; ok {
			writer.Write(append(bytes, '\n'))
		} else {
			writer.Flush()
			break
		}
	}

}

func (j Journal) Restore() ([][]byte, error) {
	if !j.WithRestore {
		return [][]byte{}, errors.New("restore is disabled")
	}
	file, err := j.openReadFile()
	if err != nil {
		return [][]byte{}, nil
	}
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

func (j Journal) openWriteFile() (*os.File, error) {
	file, err := os.OpenFile(j.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("can't open file %s %v", j.File, err)
		return nil, err
	}
	return file, nil
}

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
