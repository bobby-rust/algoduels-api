package main

import "math/rand"

type Account struct {
	ID        int 	 `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	username  string `json:"username"`
	email	  string `json:"email"`
}

func NewAccount(firstName, lastName, username, email string) *Account {
	return &Account{
		ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		username:  username,
		email: 	   email,
	}
}