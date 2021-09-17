package mailers

import (
	"todoo/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
)

func SendWelcomeEmails(u *models.User) error {
	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = "Welcome Email"
	m.From = "testmailersjavier@gmail.com"
	m.Subject = "New Contact"
	m.To = []string{u.Email}
	err := m.AddBody(r.HTML("welcome_email.html"), render.Data{})
	if err != nil {
		return err
	}

	return smtp.Send(m)
}

func SendMail(c buffalo.Context, u *models.User) error {
	m := mail.New(c)

	// fill in with your stuff:
	m.Subject = "Welcome Email"
	m.From = "testmailersjavier@gmail.com"
	m.Subject = "New Contact"
	m.To = []string{u.Email}
	err := m.AddBody(r.HTML("welcome_email.plush.html"), render.Data{})
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
