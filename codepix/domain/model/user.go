package model

import (
  "time"
  uuid "github.com/satori/go.uuid"
  "github.com/asaskevich/govalidator"
)

type User struct {
  Base          `valid:"required"`
  Name  string  `json: "name" valid:"notnull"`
  Email string  `json: "email" valid:"notnull"`
}

func (user *User) isValid() error {
  _, err := govalidator.ValidateStruct(user)

  if err != nil {
    return err
  }

  return nil
}

func NewUser(name string, email string) (*User, error) {
  user := User {Name: name, Email: email}

  user.ID = uuid.NewV4().String()
  user.CreatedAt = time.Now()
  user.UpdateAt = time.Now()

  err := user.isValid()

  if err != nil {
    return nil, err
  }

  return &user, nil
}