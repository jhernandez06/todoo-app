package actions

import (
	"database/sql"
	"net/http"
	"strings"
	"todoo/app/mailers"
	"todoo/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func EditPassword(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	user := models.User{}
	userID := c.Param("user_id")
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "a user with that ID was not found")
		return c.Redirect(http.StatusNotFound, "/user/list")
	}
	if currentUser.ID != user.ID {
		c.Flash().Add("danger", "You are not authorized.")
		return c.Redirect(302, "/tasks")
	}
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("/user/changePassword.plush.html"))
}
func ForgotPassword(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(http.StatusOK, r.HTML("user/recoverAccount.plush.html"))
}
func FindAccount(c buffalo.Context) error {
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
		verrs.Add("email", "email not found")
		c.Set("errors", verrs)
		c.Set("user", u)
		return c.Render(http.StatusUnauthorized, r.HTML("user/recoverAccount.plush.html"))
	}
	invited := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "user invited or disabled")
		c.Set("errors", verrs)
		c.Set("user", u)
		return c.Render(http.StatusUnauthorized, r.HTML("user/recoverAccount.plush.html"))
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			//couldn't find an user with the supplied email address.
			return bad()
		}

		return errors.WithStack(err)
	}
	if u.StatusUser == "disabled" || u.StatusUser == "invited" {
		return invited()
	}
	c.Set("u", u)
	mailers.SendMail(c, u)
	c.Flash().Add("success", "	An email has been sent for you to add your new password ")
	return c.Redirect(303, "/")
}
func Index(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(http.StatusOK, r.HTML("user/index.plush.html"))
}
func NewUser(c buffalo.Context) error {

	c.Set("user", models.User{})
	return c.Render(http.StatusOK, r.HTML("user/new.plush.html"))
}
func NewUserByAdmin(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(http.StatusOK, r.HTML("user/newByAdmin.plush.html"))
}
func CreateUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	verrs, err := user.Create(tx)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusOK, r.HTML("user/new.plush.html"))
	}
	c.Flash().Add("success", "user registered successfully")
	return c.Redirect(http.StatusSeeOther, "/")
}
func CreateUserByAdmin(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	verrs, err := user.Validate(tx)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusOK, r.HTML("user/newByAdmin.plush.html"))
	}
	user.StatusUser = "invited"
	if err := tx.Create(&user); err != nil {
		return err
	}
	c.Flash().Add("success", "user created successfully")
	return c.Redirect(http.StatusSeeOther, "/user/list")
}
func UsersList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	q := tx.PaginateFromParams(c.Params())
	users := models.Users{}

	if err := q.Order("rol,first_name").All(&users); err != nil {
		return err
	}
	c.Set("users", users)
	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.HTML("user/list.plush.html"))
}
func ShowUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	user := models.User{}
	tasksC := models.Tasks{}
	tasksI := models.Tasks{}

	userID := c.Param("user_id")
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "a task with that ID was not found")
		return c.Redirect(http.StatusNotFound, "/user/list")
	}
	if currentUser.Rol == "user" && currentUser.ID != user.ID {
		c.Flash().Add("danger", "You are not authorized.")
		return c.Redirect(302, "/tasks")
	}
	q := tx.Q()
	p := tx.Q()
	q.Where("check_complet = false").Where("user_id = ?", user.ID)
	if err := q.All(&tasksC); err != nil {
		return err
	}
	p.Where("check_complet = true").Where("user_id = ?", user.ID)
	if err := p.All(&tasksI); err != nil {
		return err
	}
	c.Set("tasksC", len(tasksI))
	c.Set("tasksI", len(tasksC))
	c.Set("user", user)

	return c.Render(http.StatusOK, r.HTML("user/show.plush.html"))
}
func EditUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	user := models.User{}
	userID := c.Param("user_id")
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "a user with that ID was not found")
		return c.Redirect(http.StatusNotFound, "/user/list")
	}
	if currentUser.Rol == "user" && currentUser.ID != user.ID {
		c.Flash().Add("danger", "You are not authorized.")
		return c.Redirect(302, "/tasks")
	}
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("user/edit.plush.html"))
}
func PasswordUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userID := c.Param("user_id")
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "a user with that ID was not found")
		return c.Redirect(http.StatusNotFound, "/")
	}
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("user/edit.plush.html"))
}
func UpdateUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	user := models.User{}
	userID := c.Param("user_id")
	if currentUser.Rol == "user" && currentUser.ID.String() != userID {
		c.Flash().Add("danger", "You are not authorized.")
		return c.Redirect(302, "/tasks")
	}
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "a user with that ID was not found")
		return c.Redirect(404, "/user/list")
	}
	if err := c.Bind(&user); err != nil {
		return err
	}
	verrs, err := user.Validate(tx)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusSeeOther, r.HTML("user/edit.plush.html"))
	}

	if err := tx.Update(&user); err != nil {
		return err
	}

	c.Flash().Add("success", "user updated successfully")
	if currentUser.Rol == "admin" {
		return c.Redirect(http.StatusSeeOther, "/user/list")
	}
	return c.Redirect(http.StatusSeeOther, "/tasks")
}
func AddPassword(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userID := c.Param("user_id")
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "a user with that ID was not found")
		return c.Redirect(404, "/user/list")
	}
	if err := c.Bind(&user); err != nil {
		return err
	}
	if user.StatusUser == "invited" {
		user.StatusUser = "activated"
	}
	verrs, err := user.Update(tx)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusSeeOther, r.HTML("user/password.plush.html"))
	}
	c.Flash().Add("success", "user registration completed successfully")
	return c.Redirect(http.StatusSeeOther, "/")
}
func UpdatePassword(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	user := models.User{}
	userID := c.Param("user_id")
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "a user with that ID was not found")
		return c.Redirect(404, "/user/list")
	}
	if currentUser.ID != user.ID {
		c.Flash().Add("danger", "You are not authorized.")
		return c.Redirect(302, "/tasks")
	}
	equal := func() error {
		verrs := validate.NewErrors()
		verrs.Add("password", "use a different password than the current one")
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusUnauthorized, r.HTML("user/changePassword.plush.html"))
	}
	badPassword := func() error {
		verrs := validate.NewErrors()
		verrs.Add("password_old", "Incorrect password")
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusUnauthorized, r.HTML("user/changePassword.plush.html"))
	}
	if err := c.Bind(&user); err != nil {
		return err
	}
	// confirm that the given password matches the hashed password from the db
	err1 := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user.PasswordOld))
	if err1 != nil {
		return badPassword()
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user.Password))
	if err == nil {
		return equal()
	}

	verrs, err := user.Update(tx)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusSeeOther, r.HTML("user/changePassword.plush.html"))
	}
	c.Flash().Add("success", "the password was changed successfully")
	return c.Redirect(http.StatusSeeOther, "/tasks")
}
func DestroyUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	path := "/user/list"
	user := models.User{}
	userID, _ := uuid.FromString(c.Param("user_id"))
	if err := tx.Find(&user, userID); err != nil {
		c.Flash().Add("danger", "no user found with that ID")
		return c.Redirect(404, "/user/list")
	}
	if user.ID == currentUser.ID {
		//c.Session().Clear()
		c.Flash().Add("info", "an administrator can only be deleted or inactivated by another administrator")
		return c.Redirect(http.StatusSeeOther, "/user/list")
	}
	if err := tx.Destroy(&user); err != nil {
		return err
	}
	c.Flash().Add("success", "user destroyed successfully")
	return c.Redirect(http.StatusSeeOther, path)
}
func UpdateUserActive(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	user := models.User{}
	userID := c.Param("user_id")
	if err := tx.Find(&user, userID); err != nil {
		return err
	}
	if err := c.Bind(&user); err != nil {
		return err
	}
	if user.StatusUser == "invited" {
		c.Flash().Add("info", "the user has not assigned a password")
	}
	if user.ID == currentUser.ID {
		c.Flash().Add("info", "an administrator can only be deleted or inactivated by another administrator")
		return c.Redirect(http.StatusSeeOther, "/user/list")
	}
	if user.StatusUser == "disabled" {
		user.StatusUser = "activated"
		c.Flash().Add("info", "User Activated successfully")
	} else if user.StatusUser == "activated" {
		user.StatusUser = "disabled"
		c.Flash().Add("info", "User disable")
	}
	if err := tx.Update(&user); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/user/list")
}
