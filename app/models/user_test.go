package models_test

import (
	"todoo/app/models"
)

func (ms *ModelSuite) Test_User_Create() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandes",
		Email:                "javier@gmail.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	ms.Zero(u.PasswordHash)
	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.NotZero(u.PasswordHash)
	ms.Equal("user", u.Rol)
	ms.Equal("disabled", u.StatusUser)
	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_User_Create_ValidationErrors() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.User{
		FirstName: "Javier",
		LastName:  "Hernandes",
		Password:  "password",
	}

	ms.Zero(u.PasswordHash)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_User_Create_UserExists() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandes",
		Email:                "javier@gmail.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	ms.Zero(u.PasswordHash)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.NotZero(u.PasswordHash)

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, count)

	u = &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandes",
		Email:                "javier@gmail.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	verrs, err = u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, count)
}
func (ms *ModelSuite) Test_Update_Password() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)
	u := &models.User{
		FirstName:  "Abby",
		LastName:   "Dog",
		Rol:        "user",
		Email:      "abby@gmail.com",
		StatusUser: "invited"}
	if err := ms.DB.Create(u); err != nil {
		ms.Fail("Fail")
	}

	ms.Zero(u.PasswordHash)
	u = &models.User{
		ID:                   u.ID,
		FirstName:            "Abby",
		LastName:             "Dog",
		Email:                "abby@gmail.com",
		Rol:                  "user",
		StatusUser:           "activated",
		Password:             "abby",
		PasswordConfirmation: "abby",
	}
	verrsUpdate, err := u.Update(ms.DB)
	ms.NoError(err)
	ms.DB.Reload(u)
	ms.NotZero(u.PasswordHash)
	ms.False(verrsUpdate.HasAny())
	ms.NotZero(u.PasswordHash)
}
func (ms *ModelSuite) Test_User_CreateByAdmin() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)
	u := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandez",
		Email:                "javier@gmail.com",
		Rol:                  "superAdmin",
		StatusUser:           "activated",
		Password:             "javier",
		PasswordConfirmation: "javier"}
	verrs, err := u.CreateByAdmin(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)

	u1 := &models.User{
		FirstName:            "Javier",
		LastName:             "Hernandez",
		Email:                "javier@gmail.com",
		Rol:                  "admin",
		StatusUser:           "activated",
		Password:             "javier",
		PasswordConfirmation: "javier"}

	verrs, err = u1.CreateByAdmin(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	count1, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, count1)
}
