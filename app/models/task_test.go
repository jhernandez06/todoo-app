package models_test

import (
	"time"
	"todoo/app/models"

	"github.com/gofrs/uuid"
)

func (ms *ModelSuite) Test_Task_Validate() {
	id := uuid.Must(uuid.DefaultGenerator.NewV4())
	u := ms.CreateUser()
	task := &models.Task{
		Title:        "Testing",
		LimitData:    time.Now(),
		Description:  "test",
		CheckComplet: false,
		Priority:     "j", //Priority["a","b","c"]
		UserID:       u.ID}
	verrs := task.Validate(ms.DB)
	ms.True(verrs.HasAny())

	task = &models.Task{
		Title:        "Testing",
		LimitData:    time.Now(),
		Description:  "test",
		CheckComplet: false,
		Priority:     "a",
		UserID:       id} //id no Found
	verrs = task.Validate(ms.DB)
	ms.True(verrs.HasAny())

	task = &models.Task{
		Title:        "", // can not be blank
		LimitData:    time.Now(),
		Description:  "test",
		CheckComplet: false,
		Priority:     "a",
		UserID:       u.ID}
	verrs = task.Validate(ms.DB)
	ms.True(verrs.HasAny())

}
