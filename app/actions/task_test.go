package actions_test

import (
	"fmt"
	"time"
	"todoo/app/models"
)

func (as *ActionSuite) Test_Task_New() {
	as.Login()
	res := as.HTML("/tasks/new/").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New Task")
}

func (as *ActionSuite) Test_Users_Create() {
	count, err := as.DB.Count("tasks")
	as.NoError(err)
	as.Equal(0, count)
	user := as.Login()
	tasks := models.Tasks{
		{Title: "test create 1",
			LimitData:   time.Now(),
			Description: "Testing 1",
			UserID:      user.ID,
			Priority:    "a"},
		{Title: "test create 2",
			LimitData:   time.Now(),
			Description: "Testing 2",
			UserID:      user.ID,
			Priority:    "b"},
		{Title: "test create 3",
			LimitData:   time.Now(),
			Description: "Testing 3",
			UserID:      user.ID,
			Priority:    "c"},
	}
	for _, task := range tasks {
		resp := as.HTML("/tasks/create").Post(task)
		as.Equal(303, resp.Code)
	}
	count, err = as.DB.Count("tasks")
	as.NoError(err)
	as.Equal(3, count)
}

func (as *ActionSuite) Test_Task_List() {
	count, err := as.DB.Count("tasks")
	as.NoError(err)
	as.Equal(0, count)
	user := as.Login()
	tasks := models.Tasks{
		{Title: "test create 1",
			LimitData:    time.Now(),
			Description:  "Testing 1",
			CheckComplet: true,
			Priority:     "a",
			UserID:       user.ID},
		{Title: "test create 2",
			LimitData:    time.Now(),
			Description:  "Testing 2",
			CheckComplet: false,
			Priority:     "b",
			UserID:       user.ID,
		},
		{Title: "test create 3",
			LimitData:    time.Now(),
			Description:  "Testing 3",
			CheckComplet: true,
			Priority:     "c",
			UserID:       user.ID,
		},
	}
	for _, task := range tasks {
		resp := as.HTML("/tasks/create").Post(task)
		as.Equal(303, resp.Code)
	}
	res := as.HTML("/tasks").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	for _, task := range tasks {
		as.Contains(body, fmt.Sprintf("%s", task.Title))
	}
	resCompleted := as.HTML("/tasks/?check_complet=true").Get()
	as.Equal(200, resCompleted.Code)
	body = resCompleted.Body.String()
	as.NotContains(body, "task create 2")
	resincomplet := as.HTML("/tasks/?check_complet=false").Get()
	as.Equal(200, resincomplet.Code)
	body = resincomplet.Body.String()
	as.NotContains(body, "task create 1")
	as.NotContains(body, "task create 3")
	q := as.DB.Q()
	q.Where("check_complet = ?", "true").Count(&models.Task{})
	countComplet, err := q.Count("tasks")
	as.NoError(err)
	as.Equal(2, countComplet)
}

func (as *ActionSuite) Test_Task_Show() {
	count, err := as.DB.Count("tasks")
	as.NoError(err)
	as.Equal(0, count)
	user := as.Login()
	task_test := &models.Task{
		Title:       "test Show 1 ",
		LimitData:   time.Now(),
		Description: "Testing",
		Priority:    "a",
		UserID:      user.ID,
	}
	verrs1, err := as.DB.ValidateAndCreate(task_test)
	as.NoError(err)
	as.False(verrs1.HasAny())
	res := as.HTML("/tasks/show/{%s}", task_test.ID).Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, task_test.Title)
	res2 := as.HTML("/show/{javier}").Get()
	as.Equal(404, res2.Code)
}

func (as *ActionSuite) Test_Task_Edit() {
	user := as.Login()
	task_test1 := &models.Task{
		Title:       "test Edit 1 ",
		LimitData:   time.Now(),
		Description: "Testing",
		UserID:      user.ID,
		Priority:    "a"}
	task_test2 := &models.Task{
		Title:        "test Edit 2 ",
		LimitData:    time.Now(),
		Description:  "Testing",
		CheckComplet: true,
		UserID:       user.ID,
		Priority:     "b"}
	verrs1, err := as.DB.ValidateAndCreate(task_test1)
	as.NoError(err)
	as.False(verrs1.HasAny())
	verrs2, err := as.DB.ValidateAndCreate(task_test2)
	as.NoError(err)
	as.False(verrs2.HasAny())
	res := as.HTML("/tasks/edit/{%s}", task_test1.ID).Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Edit Task", task_test1.Title)
	res2 := as.HTML("/tasks/edit/{%s}", task_test2.ID).Get()
	as.Equal(303, res2.Code)
	res3 := as.HTML("/tasks/edit/{javier}").Get()
	as.Equal(404, res3.Code)
}

func (as *ActionSuite) Test_Task_Update() {
	user := as.Login()
	task_test := &models.Task{
		Title:       "test Update 1",
		LimitData:   time.Now(),
		Description: "Testing",
		Priority:    "a",
		UserID:      user.ID,
	}
	verrs, err := as.DB.ValidateAndCreate(task_test)
	as.NoError(err)
	as.False(verrs.HasAny())
	res := as.HTML("/tasks/update/{%s}", task_test.ID).Put(&models.Task{ID: task_test.ID, Title: "Learn Go", LimitData: time.Now(), Description: "Testing Update", Priority: "a", UserID: user.ID})
	as.Equal(303, res.Code)
	err = as.DB.Reload(task_test)
	as.NoError(err)
	as.Equal("Learn Go", task_test.Title)
	res2 := as.HTML("/tasks/update/{%s}", task_test.ID).Put(&models.Task{ID: task_test.ID, Title: "", LimitData: time.Now(), Description: "Testing Update", UserID: user.ID})
	as.Equal(303, res2.Code)
	res3 := as.HTML("/tasks/update/{javier}").Put(&models.Task{Title: "Learn Go", LimitData: time.Now(), Description: "Testing Update", UserID: user.ID})
	as.Equal(404, res3.Code)
	err = as.DB.Reload(task_test)
	as.NoError(err)
}

func (as *ActionSuite) Test_Task_Delete() {
	user := as.Login()
	task_test1 := &models.Task{
		Title:       "test delete 1 ",
		LimitData:   time.Now(),
		Description: "Testing",
		UserID:      user.ID,
		Priority:    "a"}
	task_test2 := &models.Task{
		Title:        "test delete 2 ",
		LimitData:    time.Now(),
		Description:  "Testing",
		CheckComplet: true,
		UserID:       user.ID,
		Priority:     "b"}
	verrs1, err := as.DB.ValidateAndCreate(task_test1)
	as.NoError(err)
	as.False(verrs1.HasAny())
	verrs2, err := as.DB.ValidateAndCreate(task_test2)
	as.NoError(err)
	as.False(verrs2.HasAny())
	res := as.HTML("/tasks/delete/{%s}", task_test1.ID).Delete()
	as.Equal(303, res.Code)
	res3 := as.HTML("/tasks/delete/{javier}").Delete()
	as.Equal(404, res3.Code)
}

func (as *ActionSuite) Test_Task_CheckUpdate() {
	user := as.Login()
	task_test := &models.Task{
		Title:       "test CheckUpdate ",
		LimitData:   time.Now(),
		Description: "Testing",
		UserID:      user.ID,
		Priority:    "b"}
	verrs1, err := as.DB.ValidateAndCreate(task_test)
	as.NoError(err)
	as.False(verrs1.HasAny())
	res := as.HTML("/tasks/updateCheck/{%s}", task_test.ID).Put(&models.Task{ID: task_test.ID, Title: "test CheckUpdate", LimitData: time.Now(), Description: "Testing", CheckComplet: false, UserID: user.ID, Priority: "a"})
	as.Equal(303, res.Code)
	err = as.DB.Reload(task_test)
	as.NoError(err)
	as.Equal(true, task_test.CheckComplet)
	res3 := as.HTML("/tasks/updatecheck/{javier}").Put(&models.Task{ID: task_test.ID, CheckComplet: true})
	as.Equal(404, res3.Code)
}
