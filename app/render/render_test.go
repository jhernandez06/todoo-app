package render_test

import (
	"fmt"
	"testing"
	"time"
	"todoo/app/render"
)

func TestStatus(t *testing.T) {
	dateUpdate := time.Now()
	dateLimit := dateUpdate.Add(time.Hour * 24 * 7)
	dateCurrent := time.Now()
	statusTrue := render.Status(true, dateLimit, dateUpdate)
	if statusTrue != fmt.Sprintf("Was completed on %v", dateUpdate.Format("Monday 02, Jan 2006 at 15:04")) {
		t.Error("Error")
		t.Fail()
	}
	statusFalse := render.Status(false, dateLimit, dateUpdate)
	if dateCurrent.After(dateLimit) {
		if statusFalse != fmt.Sprintf("was to be completed on %v", dateLimit.Format("Monday 02, Jan 2006 at 15:04")) {
			t.Error("Error")
			t.Fail()
		}
	} else {
		if statusFalse != fmt.Sprintf("needs to be completed on %v", dateLimit.Format("Monday 02, Jan 2006 at 15:04")) {
			t.Error("Error")
			t.Fail()
		}
	}
}

func TestIcon(t *testing.T) {
	var icon string
	if icon = render.Icon("info"); icon != "#info-fill" {
		t.Error("Error")
		t.Fail()
	}
	if icon = render.Icon("danger"); icon != "#exclamation-triangle-fill" {
		t.Error("Error")
		t.Fail()
	}
	if icon = render.Icon("success"); icon != "#check-circle-fill" {
		t.Error("Error")
		t.Fail()
	}

}
func TestCheck(t *testing.T) {
	if check := render.CheckStatus("status", "status"); check != "font-weight-lighter" {
		t.Error("Error")
		t.Fail()
	}
}
func TestAddTask(t *testing.T) {
	if add := render.AddTask("true"); add != "d-none" {
		t.Error("Error")
		t.Fail()
	}
}

func TestCompleted(t *testing.T) {
	if completed := render.Completed("true"); completed != "d" {
		t.Error("Error")
		t.Fail()
	}
}
