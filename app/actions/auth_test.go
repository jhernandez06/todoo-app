package actions_test

func (as *ActionSuite) Test_Auth_Create() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	u1 := as.Activated()
	resp1 := as.HTML("/signin").Post(u1)
	as.Equal(302, resp1.Code)
	u2 := as.Invited()
	resp2 := as.HTML("/signin").Post(u2)
	as.Equal(401, resp2.Code)
	body := resp2.Body.String()
	as.Contains(body, "Add Password")
	u3 := as.disabled()
	resp3 := as.HTML("/signin").Post(u3)
	as.Equal(401, resp3.Code)
	count, err = as.DB.Count("users")
	as.NoError(err)
	as.Equal(3, count)
}

func (as *ActionSuite) Test_Auth_Destroy() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)
	user := as.Login()
	respShow := as.HTML("/user/show/{%s}", user.ID).Get()
	as.Equal(200, respShow.Code)
	resp := as.HTML("/signout").Delete()
	as.Equal(302, resp.Code)
	respShow = as.HTML("/user/show/{%s}", user.ID).Get()
	as.Equal(302, respShow.Code)
}
