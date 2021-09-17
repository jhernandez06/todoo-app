package actions_test

import (
	"testing"
	"todoo/app"
	"todoo/app/models"

	"github.com/gobuffalo/suite/v3"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	bapp := app.New()

	as := &ActionSuite{suite.NewAction(bapp)}
	suite.Run(t, as)

}
func (as *ActionSuite) CreateUser() *models.User {
	user := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandez",
		Email:                "javier@wawand.co",
		StatusUser:           "activated",
		Rol:                  "admin",
		Password:             "javier",
		PasswordConfirmation: "javier",
	}

	verrs, err := user.CreateByAdmin(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	return user
}

func (as *ActionSuite) Login() *models.User {
	user := as.CreateUser()
	as.Session.Set("current_user_id", user.ID)
	return user
}

func (as *ActionSuite) Invited() *models.User {
	u := &models.User{
		FirstName:  "Abby",
		LastName:   "Dog",
		Rol:        "user",
		Email:      "abby@gmail.com",
		StatusUser: "invited"}
	if err := as.DB.Create(u); err != nil {
		as.Fail("Fail")
	}
	return u
}

func (as *ActionSuite) disabled() *models.User {
	user := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandez",
		Email:                "javierdisabled@wawand.co",
		StatusUser:           "disabled",
		Rol:                  "admin",
		Password:             "javier",
		PasswordConfirmation: "javier",
	}
	verrs, err := user.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	return user
}
func (as *ActionSuite) Activated() *models.User {
	user := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandez",
		Email:                "javieractivated@wawand.co",
		StatusUser:           "activated",
		Rol:                  "admin",
		Password:             "javier",
		PasswordConfirmation: "javier",
	}
	verrs, err := user.CreateByAdmin(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	return user
}
