package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"todoo/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthLanding shows a landing page to login
func AuthLanding(c buffalo.Context) error {
	return c.Render(200, r.HTML("auth/landing.plush.html"))
}

// AuthNew loads the signin page
func AuthNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("auth/new.plush.html"))
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}
	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(u.Email))).First(u)

	// helper function to handle bad attempts
	bad := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password")
		c.Set("errors", verrs)
		c.Set("user", u)
		return c.Render(http.StatusUnauthorized, r.HTML("user/index.plush.html"))
	}
	inactived := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "inactive user")
		c.Set("errors", verrs)
		c.Set("user", u)
		return c.Render(http.StatusUnauthorized, r.HTML("user/index.plush.html"))
	}
	invited := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "invited, add password")
		c.Set("errors", verrs)
		c.Set("user", u)
		return c.Render(http.StatusUnauthorized, r.HTML("user/password.plush.html"))
	}
	if u.StatusUser == "invited" {
		return invited()
	}
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			//couldn't find an user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}
	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}
	if u.StatusUser == "disabled" {
		return inactived()
	}
	msg := fmt.Sprintf("Welcome %s!!", u.FirstName)
	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", msg)
	return c.Redirect(302, "/tasks")
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out!")
	return c.Redirect(302, "/")
}
