package main

import (
	"time"
)

type Account struct {
	ID        int 	 	`json:"id"`
	FirstName string 	`json:"firstName"`
	LastName  string 	`json:"lastName"`
	Username  string 	`json:"username"`
	Email	  string 	`json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string
	LastName  string
	Username  string
	Email 	  string
}

func NewAccount(firstName, lastName, username, email string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email: 	   email,
		CreatedAt: time.Now().UTC(),
	}
}