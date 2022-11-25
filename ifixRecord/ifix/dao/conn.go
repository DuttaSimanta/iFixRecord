package dao

import "database/sql"

//DbConn is used for initialized connection
type DbConn struct {
	DB *sql.DB
}
