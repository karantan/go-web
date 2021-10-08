package user

import (
	"goweb/db"
	"testing"
)

func TestFindOneUser(t *testing.T) {
	db, _ := db.GetDB("postgres://postgres:secret@localhost:5432/goweb")
	FindOneUser(db, &User{Username: "bar"})
}
