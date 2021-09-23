package actions_test

import (
	"fmt"
	"todoo/app/models"
)

func (as *ActionSuite) Test_User_Index() {
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Wawandco 2021 training app")
}

func (as *ActionSuite) Test_User_New() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	res := as.HTML("/user/new").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New User")
}

func (as *ActionSuite) Test_User_NewByAdmin() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	as.Login()
	res := as.HTML("/user/newByAdmin").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New User")
}
func (as *ActionSuite) Test_Create_UserByAdmin() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	as.Login()
	u := &models.User{
		FirstName:  "Ares",
		LastName:   "Dog",
		Email:      "ares@gmail.com",
		Rol:        "user",
		StatusUser: "activated"}
	resp := as.HTML("/user/createByAdmin").Post(u)
	as.Equal(303, resp.Code)
	count2, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count2)
}

func (as *ActionSuite) Test_Create_User() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	users := models.Users{
		{FirstName: "Javier",
			LastName:             "Hernandez",
			Email:                "javier@gmail.com",
			Rol:                  "admin",
			StatusUser:           "activated",
			Password:             "javier",
			PasswordConfirmation: "javier"},
		{FirstName: "Eduardo",
			LastName:             "Gomez",
			Email:                "eduardo@gmail.com",
			Rol:                  "user",
			StatusUser:           "activated",
			Password:             "eduardo",
			PasswordConfirmation: "eduardo"},
	}
	for _, user := range users {
		resp := as.HTML("/user/create").Post(user)
		as.Equal(303, resp.Code)
	}
	count2, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count2)
}

func (as *ActionSuite) Test_User_List() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	users := models.Users{
		{FirstName: "Javier",
			LastName:             "Hernandez",
			Email:                "javier@gmail.com",
			Rol:                  "admin",
			StatusUser:           "activated",
			Password:             "javier",
			PasswordConfirmation: "javier"},
		{FirstName: "Eduardo",
			LastName:             "Gomez",
			Email:                "eduardo@gmail.com",
			Rol:                  "user",
			StatusUser:           "activated",
			Password:             "eduardo",
			PasswordConfirmation: "eduardo"},
	}
	for _, user := range users {
		resp := as.HTML("/user/create").Post(user)
		as.Equal(303, resp.Code)
	}
	count2, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count2)
	as.Login()
	res := as.HTML("/user/list").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	for _, user := range users {
		as.Contains(body, fmt.Sprintf("%s", user.Email))
	}
	count1, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(3, count1)
}

func (as *ActionSuite) Test_Show_User() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	user := as.Login()
	resp404 := as.HTML("/user/show/javier").Get()
	as.Equal(404, resp404.Code)
	respShow := as.HTML("/user/show/{%s}", user.ID).Get()
	as.Equal(200, respShow.Code)
	body := respShow.Body.String()
	as.Contains(body, fmt.Sprintf("%s", user.Email))

	count1, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, count1)
}

func (as *ActionSuite) Test_Delete_User() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)

	user := &models.User{
		FirstName:            "Ares",
		LastName:             "Dog",
		Email:                "ares@gmail.com",
		Rol:                  "admin",
		StatusUser:           "activated",
		Password:             "javier",
		PasswordConfirmation: "javier"}
	verrs, err := user.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Login()
	count1, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count1)
	resp404 := as.HTML("/user/delete/javier").Delete()
	as.Equal(404, resp404.Code)
	respDelete := as.HTML("/user/delete/{%s}", user.ID).Delete()
	as.Equal(303, respDelete.Code)
	count2, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, count2)
}

func (as *ActionSuite) Test_Edit_User() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	user := &models.User{
		FirstName:            "Ares",
		LastName:             "Dog",
		Email:                "ares@gmail.com",
		Rol:                  "admin",
		StatusUser:           "activated",
		Password:             "javier",
		PasswordConfirmation: "javier"}
	verrs, err := user.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Login()
	resp404 := as.HTML("/user/edit/javier").Get()
	as.Equal(404, resp404.Code)
	respEdit := as.HTML("/user/edit/{%s}", user.ID).Get()
	as.Equal(200, respEdit.Code)
	body := respEdit.Body.String()
	as.Contains(body, fmt.Sprintf("%s", user.Email))
	as.Contains(body, "Edit User")
	count1, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count1)
}

func (as *ActionSuite) Test_Update_User() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	user := &models.User{
		FirstName:            "Ares",
		LastName:             "Dog",
		Email:                "ares@gmail.com",
		Rol:                  "admin",
		StatusUser:           "activated",
		Password:             "javier",
		PasswordConfirmation: "javier"}
	verrs, err := user.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Login()
	respUpdate := as.HTML("/user/update/{%s}", user.ID).Put(&models.User{ID: user.ID,
		FirstName:  "Ares",
		LastName:   "Update",
		Email:      "ares@gmail.com",
		Rol:        "admin",
		StatusUser: "activated"})
	as.Equal(303, respUpdate.Code)
	err = as.DB.Reload(user)
	as.NoError(err)
	as.Equal("Update", user.LastName)
	resp404 := as.HTML("/user/update/javier").Put(&models.User{ID: user.ID, FirstName: "Test", LastName: "Update", Email: "javier@gmail.com"})
	as.Equal(404, resp404.Code)
	respVacio := as.HTML("/user/update/{%s}", user.ID).Put(&models.User{ID: user.ID, FirstName: "", LastName: "Update", Email: "javier@gmail.com"})
	as.Equal(303, respVacio.Code)
	count1, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count1)
}

func (as *ActionSuite) Test_Active_User() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	user := &models.User{
		FirstName:            "Ares",
		LastName:             "Dog",
		Email:                "ares@gmail.com",
		Rol:                  "admin",
		StatusUser:           "disabled",
		Password:             "javier",
		PasswordConfirmation: "javier"}
	verrs, err := user.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Login()
	respActive := as.HTML("/user/active/{%s}", user.ID).Put(&models.User{ID: user.ID,
		FirstName:            "Ares",
		LastName:             "Dog",
		Email:                "ares@gmail.com",
		Rol:                  "admin",
		StatusUser:           "disabled",
		Password:             "javier",
		PasswordConfirmation: "javier"})
	as.Equal(303, respActive.Code)
	err = as.DB.Reload(user)
	as.NoError(err)
	as.Equal("activated", user.StatusUser)
	resp404 := as.HTML("/user/update/javier").Put(&models.User{ID: user.ID, FirstName: "Test", LastName: "Update", Email: "javier@gmail.com"})
	as.Equal(404, resp404.Code)
	count1, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count1)
}

func (as *ActionSuite) Test_User_Password() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	as.Login()
	u := models.User{
		FirstName: "Abby",
		LastName:  "Dog",
		Rol:       "user",
		Email:     "abby@gmail.com"}
	resp := as.HTML("/user/createByAdmin").Post(u)
	as.Equal(303, resp.Code)
	resp = as.HTML("/signin").Post(u)
	body := resp.Body.String()
	as.Contains(body, "Add Password")
}

func (as *ActionSuite) Test_Update_Password_User() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	u := as.Invited()
	resp2 := as.HTML("/user/updatePassword/{%s}", u.ID).Put(&models.User{ID: u.ID,
		FirstName:            "Abby",
		LastName:             "Dog",
		Email:                "abby@gmail.com",
		Rol:                  "user",
		StatusUser:           "activated",
		Password:             "javier",
		PasswordConfirmation: "javier"})
	as.Equal(302, resp2.Code)
	// as.DB.Reload(u)
	// as.Equal("activated", u.StatusUser)
	count1, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, count1)
}
