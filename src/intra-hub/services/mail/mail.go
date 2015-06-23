package mail

import (
	"bytes"
	"strconv"
	"text/template"

	"intra-hub/confperso"
	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/go-gomail/gomail"
)

var (
	auth = gomail.LoginAuth(confperso.EmailUsername, confperso.EmailPassword, confperso.EmailHost)
)

func SendUserCreated(user *models.User) error {
	data := make(map[string]interface{})
	data["User"] = user
	data["Link"] = "http://localhost:8080/users/activate/" + strconv.FormatInt(int64(user.Id), 10) + "/" + user.Token
	tmpl, err := template.New("mail").Parse(templateUserCreated)
	if err != nil {
		beego.Error(err)
		return err
	}
	tmplSubject, err := template.New("subject").Parse(subjectUserCreated)
	if err != nil {
		beego.Error(err)
		return err
	}
	b := bytes.NewBufferString("")
	if err := tmpl.Execute(b, data); err != nil {
		beego.Error(err)
		return err
	}
	bsubject := bytes.NewBufferString("")
	if err := tmplSubject.Execute(bsubject, data); err != nil {
		beego.Error(err)
		return err
	}
	sendMail(user.Email, bsubject.String(), b.String())
	return nil
}

func sendMail(to string, subject, body string) {
	config := `{"username":"` + confperso.EmailUsername + `","password":"` + confperso.EmailPassword + `","host":"` +
		confperso.EmailHost + `","port":` + confperso.EmailHostPort + `}`
	email := utils.NewEMail(config)
	beego.Warn(subject)
	email.Subject = subject
	email.Auth = auth
	email.To = []string{to}
	email.HTML = body
	email.From = confperso.EmailUsername
	if err := email.Send(); err != nil {
		beego.Warn("MAIL ERROR", err)
	}
}
