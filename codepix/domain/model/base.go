package model

import (
  "time"
  "github.com/asaskevich/govalidator"
)

type Base struct {
	ID        string    `json: "id" valid:"uuid"`
	CreatedAt time.Time `json: "created_at" valid:"-"`
	UpdateAt  time.Time `json: "update_at" valid:"-"`
}

func init() {
  govalidator.SetFieldsRequiredByDefault(true)
}