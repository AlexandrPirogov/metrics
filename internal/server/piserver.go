package pi

import (
	"memtracker/internal/server/db/nosql/pidb"
	"net/http"
)

type Server struct {
	Server *http.Server
	DB     *pidb.DB
}
