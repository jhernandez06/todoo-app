package actions

import (
	"net/http"
	"todoo/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func TasksList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := models.Tasks{}
	currentUser := c.Value("current_user").(*models.User)
	status := c.Param("check_complet")
	q := tx.PaginateFromParams(c.Params())

	if status == "true" || status == "false" {
		q.Where("check_complet = ?", status)
	}
	if currentUser.Rol == "user" {
		q.Where("user_id = ?", currentUser.ID)
	}

	if err := q.Order("priority, check_complet, limit_data").All(&tasks); err != nil {
		return err
	}

	c.Set("tasks", tasks)
	c.Set("paginationTasks", q.Paginator)
	return c.Render(http.StatusOK, r.HTML("task/index.plush.html"))
}

func NewTask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	users := models.Users{}
	user := models.User{}
	q := tx.Q()
	if currentUser.Rol == "admin" {
		q.Where("status_user = ?", "activated")
		if err := q.Order("first_name asc").All(&users); err != nil {
			return err
		}
		UsersList := []map[string]interface{}{}
		for _, user := range users {
			oneUser := map[string]interface{}{
				user.FirstName + " " + user.LastName: uuid.FromStringOrNil(user.ID.String()),
			}
			UsersList = append(UsersList, oneUser)
		}

		c.Set("usersList", UsersList)
		c.Set("user", user)
		c.Set("users", users)
		c.Set("task", models.Task{})
		return c.Render(http.StatusOK, r.HTML("task/new.plush.html"))
	}
	c.Set("user", user)
	c.Set("task", models.Task{})
	return c.Render(http.StatusOK, r.HTML("task/new.plush.html"))
}

func CreateTask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	task := models.Task{}
	users := models.Users{}
	user := models.User{}
	q := tx.Q()
	q.Where("status_user = ?", "activated")
	if err := q.Order("first_name asc").All(&users); err != nil {
		return err
	}
	if currentUser.Rol == "admin" {

		UsersList := []map[string]interface{}{}
		for _, user := range users {
			oneUser := map[string]interface{}{
				user.FirstName + " " + user.LastName: uuid.FromStringOrNil(user.ID.String()),
			}
			UsersList = append(UsersList, oneUser)
		}

		if err := c.Bind(&task); err != nil {
			return errors.WithStack(err)
		}
		verrs := task.Validate(tx)
		if verrs.HasAny() {
			c.Set("errors", verrs)
			c.Set("task", task)
			c.Set("usersList", UsersList)
			c.Set("user", user)
			c.Set("users", users)
			return c.Render(http.StatusOK, r.HTML("task/new.plush.html"))
		}
		if err := tx.Create(&task); err != nil {
			return err
		}
		c.Flash().Add("success", "task created success")
		return c.Redirect(http.StatusSeeOther, "/tasks")
	}

	if err := c.Bind(&task); err != nil {
		return err
	}
	task.UserID = currentUser.ID
	verrs := task.Validate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("task", task)
		c.Set("user", user)
		return c.Render(http.StatusOK, r.HTML("task/new.plush.html"))
	}
	if err := tx.Create(&task); err != nil {
		return err
	}
	c.Flash().Add("success", "task created success")
	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func ShowTask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	task := models.Task{}
	user := models.User{}
	taskID, _ := uuid.FromString(c.Param("task_id"))
	if err := tx.Find(&task, taskID); err != nil {
		c.Flash().Add("danger", "a task with that ID was not found")
		return c.Redirect(http.StatusNotFound, "/tasks")
	}
	if err := tx.Find(&user, task.UserID); err != nil {
		return err
	}
	if currentUser.Rol == "user" && currentUser.ID != task.UserID {
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/tasks")
	}
	c.Set("user", user)
	c.Set("task", task)
	return c.Render(http.StatusOK, r.HTML("task/show.plush.html"))
}

func EditTask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	task := models.Task{}
	users := models.Users{}
	user := models.User{}
	taskID := c.Param("task_id")
	if err := tx.Find(&task, taskID); err != nil {
		c.Flash().Add("danger", "a task with that ID was not found")
		return c.Redirect(http.StatusNotFound, "tasks/")
	}
	q := tx.Q()
	q.Where("status_user = ?", "activated")
	if err := q.Order("first_name asc").All(&users); err != nil {
		return err
	}

	UsersList := []map[string]interface{}{}
	for _, user := range users {
		oneUser := map[string]interface{}{
			user.FirstName + " " + user.LastName: user.ID,
		}
		UsersList = append(UsersList, oneUser)
	}
	if currentUser.Rol == "user" && currentUser.ID != task.UserID {
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/tasks")
	}
	c.Set("usersList", UsersList)
	c.Set("user", user)
	c.Set("users", users)
	c.Set("task", task)
	return c.Render(http.StatusOK, r.HTML("task/edit.plush.html"))
}

func UpdateTask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	task := models.Task{}
	taskID := c.Param("task_id")
	users := models.Users{}
	user := models.User{}
	q := tx.Q()
	q.Where("status_user = ?", "activated")
	if err := q.Order("first_name asc").All(&users); err != nil {
		return err
	}
	UsersList := []map[string]interface{}{}
	for _, user := range users {
		oneUser := map[string]interface{}{
			user.FirstName + " " + user.LastName: user.ID,
		}
		UsersList = append(UsersList, oneUser)
	}
	if err := tx.Find(&task, taskID); err != nil {
		c.Flash().Add("danger", "a task with that ID was not found")
		return c.Redirect(404, "/tasks")
	}
	if err := c.Bind(&task); err != nil {
		return err
	}
	verrs := task.Validate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("task", task)
		c.Set("usersList", UsersList)
		c.Set("user", user)
		c.Set("users", users)
		return c.Render(http.StatusSeeOther, r.HTML("task/edit.plush.html"))
	}
	if currentUser.Rol == "user" && currentUser.ID != task.UserID {
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/tasks")
	}
	if err := tx.Update(&task); err != nil {
		return err
	}
	c.Flash().Add("success", "task updated success")
	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func DestroyTask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	task := models.Task{}
	taskID, _ := uuid.FromString(c.Param("task_id"))
	if err := tx.Find(&task, taskID); err != nil {
		c.Flash().Add("danger", "no task found with that ID")
		return c.Redirect(404, "/tasks")
	}
	if currentUser.Rol == "user" && currentUser.ID != task.UserID {
		c.Flash().Add("danger", "You are not authorized")
		return c.Redirect(302, "/tasks")
	}
	if err := tx.Destroy(&task); err != nil {
		return err
	}
	c.Flash().Add("success", "task destroyed success")
	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func UpdateTaskCheck(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(*models.User)
	task := models.Task{}
	taskID := c.Param("task_id")
	if err := tx.Find(&task, taskID); err != nil {
		return err
	}
	if err := c.Bind(&task); err != nil {
		return err
	}
	if currentUser.Rol == "user" && currentUser.ID != task.UserID {
		c.Flash().Add("danger", "You are not authorized")
		return c.Redirect(302, "/tasks")
	}
	if !(task.CheckComplet) {
		task.CheckComplet = true
		c.Flash().Add("info", "task completed success, Congratulations")
	} else if task.CheckComplet {
		task.CheckComplet = false
		c.Flash().Add("info", "the task returned to incomplete tasks")
	}
	if err := tx.Update(&task); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/tasks")
}
