package main

import (
	"time"
)

type Account struct {
	UserID        		int   	  `json:"user_id"`
	FirstName 			string 	  `json:"firstName"`
	LastName  			string 	  `json:"lastName"`
	Username  			string 	  `json:"username"`
	Email	  			string 	  `json:"email"`
	EncryptedPassword	string	  `json:"encryptedPassword"`
	CreatedAt 			time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string
	LastName  string
	Username  string
	Email 	  string
	Password  string
}

func NewAccount(firstName, lastName, username, email, password string) *Account {
	return &Account{
		FirstName: 		   firstName,
		LastName:  		   lastName,
		Username:  		   username,
		Email: 	   		   email,
		EncryptedPassword:  password,
		CreatedAt: 		   time.Now().UTC(),
	}
}