package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Task is used by pop to map your tasks database table to your go code.
type Task struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	LimitData    time.Time `json:"limit_data" db:"limit_data"`
	Description  string    `json:"description" db:"description"`
	CheckComplet bool      `json:"check_complet" db:"check_complet"`
	Priority     string    `json:"priority" db:"priority"`
	UserID       uuid.UUID `json:"user_id" pg:"type:uuid" db:"user_id" `
	User         User      `belongs_to:"user"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (t Task) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tasks is not required by pop and may be deleted
type Tasks []Task

// String is not required by pop and may be deleted
func (t Tasks) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Task) Validate(tx *pop.Connection) *validate.Errors {

	return validate.Validate(
		&validators.StringIsPresent{Field: t.Title, Name: "Title"},
		&validators.TimeIsPresent{Field: t.LimitData, Name: "Limit Data"},
		&validators.StringIsPresent{Field: t.Description, Name: "Description"},
		&validators.UUIDIsPresent{Name: "UserID", Field: t.UserID, Message: "UserID"},
		&UserIDNotFound{Name: "UserID", Field: t.UserID, tx: tx},
		&UserIDNotValid{Name: "UserID", Field: t.UserID, tx: tx},
		&validators.FuncValidator{
			Field:   t.Priority,
			Name:    "Priority",
			Message: "%s is an invalid priority",
			Fn: func() bool {
				priorities := [3]string{"a", "b", "c"}
				for _, p := range priorities {
					if t.Priority == p {
						return true
					}
				}
				return false
			},
		},
	)
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Task) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Task) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type UserIDNotFound struct {
	Name  string
	Field uuid.UUID
	tx    *pop.Connection
}

func (v *UserIDNotFound) IsValid(errors *validate.Errors) {
	query := v.tx.Where("id = ?", v.Field).Where("status_user = ?", "activated")
	queryUser := User{}
	err := query.First(&queryUser)
	if err != nil {
		errors.Add(validators.GenerateKey(v.Name), "invalid user")
	}
}

type UserIDNotValid struct {
	Name  string
	Field uuid.UUID
	tx    *pop.Connection
}

func (v *UserIDNotValid) IsValid(errors *validate.Errors) {
	//id := IsValidUUID(v.Field.String())
	fmt.Println(len(v.Field.String()))
	// if id {
	// 	errors.Add(validators.GenerateKey(v.Name), "ID not valid!")
	// } else
	if len(v.Field.String()) != 36 {
		errors.Add(validators.GenerateKey(v.Name), "ID not valid!")
	}
}
