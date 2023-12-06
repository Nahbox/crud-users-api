package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nahbox/crud-users-api/internal/models"
)

var ErrNoMatch = fmt.Errorf("no matching record")

func (db *Database) GetAllUsers() ([]models.User, error) {
	var users []models.User
	rows, err := db.conn.Query("SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Salary, &user.Occupation)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *Database) AddUser(user *models.User) error {
	var id int
	query := `INSERT INTO users (name, age, salary, occupation) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.conn.QueryRow(query, user.Name, user.Age, user.Salary, user.Occupation).Scan(&id)
	if err != nil {
		return err
	}
	user.Id = id
	return nil
}

func (db *Database) GetUserById(userId int) (models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE id = $1;`
	row := db.conn.QueryRow(query, userId)
	switch err := row.Scan(&user.Id, &user.Name, &user.Age, &user.Salary, &user.Occupation); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (db Database) UpdateUser(userId int, userData *models.User) (models.User, error) {
	user := models.User{}
	query := `UPDATE users SET name=$1, age=$2, salary=$3, occupation=$4 WHERE id=$5 RETURNING id, name, age, salary, occupation;`
	err := db.conn.QueryRow(query, userData.Name, userData.Age, userData.Salary, userData.Occupation, userId).Scan(&user.Id, &user.Name, &user.Age, &user.Salary, &user.Occupation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrNoMatch
		}
		return user, err
	}
	return user, nil
}

func (db *Database) DeleteUser(userId int) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := db.conn.Exec(query, userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}
