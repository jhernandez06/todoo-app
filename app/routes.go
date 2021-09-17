package app

import (
	base "todoo"
	"todoo/app/actions"
	"todoo/app/middleware"

	"github.com/gobuffalo/buffalo"
)

// SetRoutes for the application
func setRoutes(root *buffalo.App) {
	root.Use(middleware.Transaction)
	root.Use(middleware.ParameterLogger)
	root.Use(middleware.CSRF)
	root.Use(middleware.NTasksIncomplet)
	root.Use(middleware.TimeMW)

	root.Use(middleware.Authorize)
	root.Use(middleware.SetCurrentUser)
	root.Use(middleware.Active)

	root.GET("/", actions.Index)
	root.POST("/signin", actions.AuthCreate)
	root.DELETE("/signout", actions.AuthDestroy)

	root.GET("/tasks/", actions.TasksList)
	root.GET("/tasks/new", actions.NewTask)
	root.POST("/tasks/create", actions.CreateTask)
	root.DELETE("/tasks/delete/{task_id}", actions.DestroyTask)
	root.GET("/tasks/show/{task_id}", actions.ShowTask)
	root.GET("/tasks/edit/{task_id}", middleware.EditTaskMW(actions.EditTask))
	root.PUT("/tasks/update/{task_id}", actions.UpdateTask)
	root.PUT("/tasks/updateCheck/{task_id}", actions.UpdateTaskCheck)

	root.GET("/user/new", actions.NewUser)
	root.GET("/user/newByAdmin", middleware.Admin(actions.NewUserByAdmin))
	root.POST("/user/create", actions.CreateUser)
	root.POST("/user/createByAdmin", middleware.Admin(actions.CreateUserByAdmin))
	root.GET("/user/list", middleware.Admin(actions.UsersList))
	root.GET("/user/show/{user_id}", actions.ShowUser)
	root.DELETE("/user/delete/{user_id}", middleware.Admin(actions.DestroyUser))
	root.GET("/user/edit/{user_id}", actions.EditUser)
	root.GET("/user/changePassword/{user_id}", actions.EditPassword)
	root.GET("/user/password", middleware.Invited(actions.PasswordUser))
	root.PUT("/user/addPassword/{user_id}", actions.AddPassword)
	root.PUT("/user/updatePassword/{user_id}", actions.UpdatePassword)
	root.PUT("/user/update/{user_id}", actions.UpdateUser)
	root.PUT("/user/active/{user_id}", middleware.Admin(actions.UpdateUserActive))
	root.GET("/user/forgotPassword", actions.ForgotPassword)
	root.POST("/user/findAccount", actions.FindAccount)

	root.Middleware.Skip(middleware.SetCurrentUser, actions.Index, actions.AuthCreate, actions.AuthDestroy, actions.NewUser, actions.CreateUser, actions.FindAccount, actions.AddPassword)
	root.Middleware.Skip(middleware.Authorize, actions.Index, actions.AuthCreate, actions.AuthDestroy, actions.NewUser, actions.CreateUser, actions.PasswordUser, actions.AddPassword, actions.ForgotPassword, actions.FindAccount)
	root.Middleware.Skip(middleware.Active, actions.Index, actions.AuthCreate, actions.AuthDestroy, actions.NewUser, actions.CreateUser, actions.PasswordUser, actions.AddPassword, actions.ForgotPassword, actions.FindAccount)
	root.ServeFiles("/", base.Assets)
}
