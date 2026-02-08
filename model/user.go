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
	Id       int64
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

	res, err := stmt.Exec(&obj.Name, &obj.Password, &obj.Role, &obj.Bio)
	if err != nil {
		return err
	}
	obj.Id, _ = res.LastInsertId()
	return nil
}

func FetchUserByName(db *sql.DB, name string) (*User, error) {
	row := db.QueryRow("SELECT `ID`, `NAME`, `PASSWORD`, `ROLE`, `BIO` FROM `USERS` WHERE `NAME` = ?", name)
	user := &User{}
	err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Role, &user.Bio)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FetchUserById(db *sql.DB, id int64) (*User, error) {
	row := db.QueryRow("SELECT `ID`, `NAME`, `PASSWORD`, `ROLE`, `BIO` FROM `USERS` WHERE `ID` = ?", id)
	user := &User{}
	err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Role, &user.Bio)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UserPublic is a safe representation of a user without the password hash.
type UserPublic struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Bio  string `json:"bio"`
}

func (u *User) Public() UserPublic {
	return UserPublic{
		Id:   u.Id,
		Name: u.Name,
		Role: u.Role,
		Bio:  u.Bio,
	}
}

func FetchAllUsers(db *sql.DB) ([]UserPublic, error) {
	rows, err := db.Query("SELECT `ID`, `NAME`, `ROLE`, `BIO` FROM `USERS` ORDER BY `ID`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]UserPublic, 0)
	for rows.Next() {
		var u UserPublic
		err := rows.Scan(&u.Id, &u.Name, &u.Role, &u.Bio)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func UpdateUser(db *sql.DB, id int64, role string, bio string) error {
	_, err := db.Exec("UPDATE `USERS` SET `ROLE` = ?, `BIO` = ? WHERE `ID` = ?", role, bio, id)
	return err
}

func UpdateUserPassword(db *sql.DB, id int64, password string) error {
	hash := calcHash(password)
	_, err := db.Exec("UPDATE `USERS` SET `PASSWORD` = ? WHERE `ID` = ?", hash, id)
	return err
}

func DeleteUser(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM `USERS` WHERE `ID` = ?", id)
	return err
}

func IsValidRole(role string) bool {
	return role == ROLE_ADMIN || role == ROLE_MONITOR || role == ROLE_OPERATOR
}
