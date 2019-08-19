package datamodels

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64  `json:"id" form:"id"`
	Firstname string `json:"firstname" form:"firstname"`
}
