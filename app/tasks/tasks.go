package tasks

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"todoo/app"
	"todoo/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/grift/grift"
	"github.com/wawandco/fako"
)

type Fako struct {
	FirstName   string `fako:"first_name"`
	LastName    string `fako:"last_name"`
	Email       string `fako:"email_address"`
	Description string `fako:"paragraph"`
	Title       string `fako:"title"`
	Priority    string `fako:"a_gen"`
}

// Init the tasks with some common tasks that come from
// grift
func init() {
	buffalo.Grifts(app.New())
}

func priorityTask(x int) string {
	var y string
	if x == 1 {
		y = "a"
	} else if x == 2 {
		y = "b"
	} else if x == 3 {
		y = "c"
	}
	return y
}
func checkComplet(x int) bool {
	var y bool
	if x == 1 {
		y = true
	} else if x == 2 {
		y = false
	}
	return y
}
func date() time.Time {

	x := rand.Intn(30-10) + 10
	dateString := fmt.Sprintf("2022-09-%vT10:12:11+06:00", x)
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		log.Fatal(err)
	}
	return date
}

var _ = grift.Add("admins", func(c *grift.Context) error {
	db := models.DB()
	for i := 0; i < 5; i++ {
		var f Fako
		fako.Fill(&f)
		admin := &models.User{
			FirstName:            f.FirstName, //fmt.Sprintf("Admin %v", i+1),
			LastName:             f.LastName,
			Email:                f.Email,
			Rol:                  "admin",
			StatusUser:           "activated",
			Password:             "javier",
			PasswordConfirmation: "javier"}
		admin.CreateByAdmin(db)
	}
	for i := 0; i < 200; i++ {
		var f Fako
		fako.Fill(&f)
		u := &models.User{
			FirstName:            f.FirstName,
			LastName:             f.LastName,
			Email:                f.Email,
			Rol:                  "user",
			StatusUser:           "activated",
			Password:             "javier",
			PasswordConfirmation: "javier"}
		u.CreateByAdmin(db)
		for j := 0; j < 5; j++ {
			x := rand.Intn(3-0) + 1
			y := rand.Intn(3-0) + 1
			t := &models.Task{
				Title:        f.Title, //fmt.Sprintf("task %v", j+1),
				LimitData:    date(),
				Description:  f.Description,
				CheckComplet: checkComplet(y),
				Priority:     fmt.Sprint(priorityTask(x)), //Priority["a","b","c"]
				UserID:       u.ID}
			db.Create(t)
		}
	}

	return nil
})

var _ = grift.Add("deleteUsers", func(c *grift.Context) error {
	db := models.DB()
	users := models.Users{}
	q := db.Q()
	q.Where("rol = ? ", "user")
	if err := q.All(&users); err != nil {
		return err
	}
	for _, user := range users {
		if err := db.Destroy(&user); err != nil {
			return err
		}
	}
	return nil
})
var _ = grift.Add("deleteAdmins", func(c *grift.Context) error {
	db := models.DB()
	users := models.Users{}
	q := db.Q()
	q.Where("email <> ? ", "jhernandez@wawand.co")
	if err := q.All(&users); err != nil {
		return err
	}
	for _, user := range users {
		if err := db.Destroy(&user); err != nil {
			return err
		}
	}
	return nil
})
