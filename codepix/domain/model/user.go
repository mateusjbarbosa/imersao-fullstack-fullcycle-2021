package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Base  `valid:"required"`
	Name  string `json:"name" gorm:"type:varchar(50)" valid:"notnull"`
	Email string `json:"email" gorm:"type:varchar(255)" valid:"notnull"`
}

func (user *User) isValid() error {
	_, err := govalidator.ValidateStruct(user)

	if err != nil {
		return err
	}

	return nil
}

func NewUser(name string, email string) (*User, error) {
	user := User{Name: name, Email: email}

	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := user.isValid()

	if err != nil {
		return nil, err
	}

	return &user, nil
}
