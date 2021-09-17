// middleware package is intended to host the middlewares used
// across the app.
package middleware

import (
	"net/http"
	"time"
	"todoo/app/models"

	tx "github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	paramlogger "github.com/gobuffalo/mw-paramlogger"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

var (
	// Transaction middleware wraps the request with a pop
	// transaction that is committed on success and rolled
	// back when errors happen.
	Transaction = tx.Transaction(models.DB())

	// ParameterLogger logs out parameters that the app received
	// taking care of sensitive data.
	ParameterLogger = paramlogger.ParameterLogger

	// CSRF middleware protects from CSRF attacks.
	CSRF = csrf.New
)

func NTasksIncomplet(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tx := models.DB()
		q := tx.Q()
		tasks := models.Tasks{}
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			if u.Rol == "admin" {
				q.Where("check_complet = false")
				if err := q.All(&tasks); err != nil {
					return err
				}
			} else {
				q.Where("check_complet = false").Where("user_id = ?", uid)
				if err := q.All(&tasks); err != nil {
					return err
				}
			}
			c.Set("current_user", u)
			c.Set("ntasks", len(tasks))
		}
		return next(c)
	}
}

func TimeMW(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		t := time.Now()
		c.Set("Date", t.Format("Monday 02, Jan 2006"))
		c.Set("Date1", t.Format("2006-01-02T15:04"))
		return next(c)
	}
}

func EditTaskMW(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tx := models.DB()
		currentUser := c.Value("current_user").(*models.User)
		task := models.Task{}
		taskID := c.Param("task_id")
		tx.Find(&task, taskID)
		if currentUser.Rol == "user" && currentUser.ID != task.UserID {
			c.Flash().Add("danger", "You must be authorized")
			c.Redirect(http.StatusSeeOther, "/tasks")
		}
		if task.CheckComplet {
			c.Flash().Add("danger", "cannot edit a completed task")
			c.Redirect(http.StatusSeeOther, "/tasks")
		}
		return next(c)
	}
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				c.Session().Clear()
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())
			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}

// AdminRequired requires a user to be logged in and to be an admin before accessing a route.
func Admin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user, ok := c.Value("current_user").(*models.User)
		if ok && user.Rol == "admin" {
			return next(c)
		}
		c.Flash().Add("danger", "You are not authorized.")
		return c.Redirect(302, "/tasks")
	}
}

func Active(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user, ok := c.Value("current_user").(*models.User)
		if ok && user.StatusUser == "activated" {
			return next(c)
		}
		c.Session().Clear()
		c.Flash().Add("danger", "inactive user")
		return c.Redirect(302, "/")
	}
}
func Invited(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user, ok := c.Value("current_user").(*models.User)
		if ok && user.StatusUser == "invited" {
			return next(c)
		} else if user.StatusUser == "activated" || user.StatusUser == "disabled" {
			c.Flash().Add("danger", "You are not authorized to view that page.")
			return c.Redirect(302, "/tasks")
		}
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/")
	}
}
