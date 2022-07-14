package model

import (
	"database/sql"
)

type DbObj interface {
	Save(*sql.DB) error
	WebsocketEvent
}
