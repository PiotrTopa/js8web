package model

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
)

var (
	DEFAULT_ADMIN_USERNAME = "admin"
	DEFAULT_ADMIN_PASSWORD = "admin"

	ROLE_ADMIN    = "admin"
	ROLE_MONITOR  = "monitor"
	ROLE_OPERATOR = "operator"
)

var DefaultAdminUser = User{
	Name:     DEFAULT_ADMIN_USERNAME,
	Password: calcHash(DEFAULT_ADMIN_PASSWORD),
	Role:     ROLE_ADMIN,
}

type User struct {
	Id       int
	Name     string
	Password string
	Role     string
	Bio      string
}

func calcHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (obj *User) SetPassword(password string) {
	obj.Password = calcHash(password)
}

func (obj *User) CheckPassword(password string) bool {
	return calcHash(password) == obj.Password
}

func (obj *User) Insert(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO `USERS`(`NAME`, `PASSWORD`, `ROLE`, `BIO`) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&obj.Name, &obj.Password, &obj.Role, &obj.Bio)
	if err != nil {
		return err
	}
	return nil
}
