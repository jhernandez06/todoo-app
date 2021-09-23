package models_test

import (
	"testing"
	"todoo/app/models"

	"github.com/gobuffalo/suite/v3"
)

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	suite.Run(t, &ModelSuite{
		Model: suite.NewModel(),
	})
}
func (ms *ModelSuite) CreateUser() *models.User {
	user := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandez",
		Email:                "javier@wawand.co",
		StatusUser:           "activated",
		Rol:                  "admin",
		Password:             "javier",
		PasswordConfirmation: "javier",
	}
	verrs, err := user.CreateByAdmin(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	return user
}
