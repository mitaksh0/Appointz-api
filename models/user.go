package models

import (
	"database/sql"
	"strconv"

	"github.com/appointments_api/db"
)

func GetUser(username, email, role string) ([]Credentials, error) {

	sqlQuery := `SELECT username, email FROM users WHERE username = $1 OR (email = $2 AND role = $3)`
	rows, err := db.Db.Query(sqlQuery, username, email, role)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var dbUsers []Credentials
	for rows.Next() {
		var dbUser Credentials
		var err error

		err = rows.Scan(&dbUser.Username, &dbUser.Email)
		if err != nil {
			return nil, err
		}

		dbUsers = append(dbUsers, dbUser)
	}

	return dbUsers, nil
}

func SetUser(user Credentials) error {
	sqlStr := `INSERT INTO users(name, username, password_hash, role, email) VALUES($1, $2, $3, $4, $5)`
	_, err := db.Db.Query(sqlStr, user.Name, user.Username, user.Password, user.Role, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func GetAuthUser(username, role string) (Credentials, error) {

	var dbUser Credentials

	sqlQuery := `SELECT id, password_hash FROM users WHERE username = $1 AND role = $2`

	err := db.Db.QueryRow(sqlQuery, username, role).Scan(&dbUser.ID, &dbUser.Password)
	if err != nil && err != sql.ErrNoRows {
		return dbUser, err
	}

	return dbUser, nil
}

// get user's role. If id is passed returns single map row with role as key and info as value array first index, otherwise lists all users with rols array
func GetRoles(id int) (map[string][]Role, error) {
	var userRoles = map[string][]Role{}

	sqlStr := `SELECT id, name, role FROM users`

	if id > 0 {
		sqlStr += ` WHERE id = ` + strconv.Itoa(id)
	}

	sqlStr += ` ORDER BY name, created_at`

	rows, err := db.Db.Query(sqlStr)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var role Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Type,
		)

		if err != nil {
			return nil, err
		}

		userRoles[role.Type] = append(userRoles[role.Type], role)
	}

	return userRoles, nil
}

func GetUsers(role string) ([]Role, error) {
	sqlStr := `SELECT id, name FROM users
	WHERE role = $1 AND name IS NOT NULL
	ORDER BY name, username`
	rows, err := db.Db.Query(sqlStr, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []Role

	for rows.Next() {
		var user Role
		err := rows.Scan(&user.ID, &user.Name)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
