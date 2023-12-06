package models

import (
	"fmt"
	"net/http"
)

type User struct {
	Id         int    `json:"id" xml:"id" toml:"id"`
	Name       string `json:"name" xml:"name" toml:"name"`
	Age        int    `json:"age" xml:"age" toml:"age"`
	Salary     int    `json:"salary" xml:"salary" toml:"salary"`
	Occupation string `json:"occupation" xml:"occupation" toml:"occupation"`
}

type NewUser struct {
	Name       string `json:"name" xml:"name" toml:"name"`
	Age        int    `json:"age" xml:"age" toml:"age"`
	Salary     int    `json:"salary" xml:"salary" toml:"salary"`
	Occupation string `json:"occupation" xml:"occupation" toml:"occupation"`
}

func (i *User) Bind(r *http.Request) error {
	if i.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}

func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
