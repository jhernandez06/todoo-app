package render

import (
	"fmt"
	"time"
	base "todoo"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
)

// Engine for rendering across the app, it provides
// the base for rendering HTML, JSON, XML and other formats
// while also defining thing like the base layout.
var Engine = render.New(render.Options{
	HTMLLayout:   "application.plush.html",
	TemplatesBox: base.Templates,
	AssetsBox:    base.Assets,
	Helpers:      Helpers,
})

// Helpers available for the plush templates, there are
// some helpers that are injected by Buffalo but this is
// the list of custom Helpers.
var Helpers = map[string]interface{}{
	// partialFeeder is the helper used by the render engine
	// to find the partials that will be used, this is important
	"partialFeeder": base.Templates.FindString,
	"FormatDate":    Status,
	"Icon":          Icon,
	"status":        CheckStatus,
	"addTask":       AddTask,
	"byCompleted":   Completed,
	"priority":      Priority,
	"equalsID":      EqualsID,
}

func Status(completed bool, date time.Time, dateUpdate time.Time) string {
	var status string
	dateCurrent := time.Now()
	if completed {
		status = fmt.Sprintf("Was completed on %v", dateUpdate.Format("Monday 02, Jan 2006 at 15:04"))
	} else {
		if dateCurrent.After(date) {
			status = fmt.Sprintf("was to be completed on %v", date.Format("Monday 02, Jan 2006 at 15:04"))
		} else {
			status = fmt.Sprintf("needs to be completed on %v", date.Format("Monday 02, Jan 2006 at 15:04"))
		}
	}
	return status
}
func Icon(k string) string {
	var icon string
	if k == "info" {
		icon = "#info-fill"
	} else if k == "danger" {
		icon = "#exclamation-triangle-fill"
	} else if k == "success" {
		icon = "#check-circle-fill"
	}
	return icon
}
func CheckStatus(x string, y string) string {
	if x == y {
		y = "font-weight-lighter"
	}
	return y
}
func AddTask(x string) string {
	var y string
	if x == "true" {
		y = "d-none"
	}
	return y
}
func Completed(x string) string {
	var y string
	if x == "true" {
		y = "d"
	}
	return y
}
func Priority(priority string) string {
	var flag string
	if priority == "a" {
		flag = "danger"
	} else if priority == "b" {
		flag = "warning"
	} else if priority == "c" {
		flag = "success"
	}
	return flag
}

func EqualsID(u1 uuid.UUID, u2 uuid.UUID) bool {

	if u1 == u2 {

		return true
	}

	return false

}
