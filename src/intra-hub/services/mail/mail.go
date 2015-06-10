package mail

import (
	"net/smtp"

	"intra-hub/confperso"
	"intra-hub/models"
    "github.com/astaxie/beego"
    "text/template"
    "bytes"
    "strconv"
)

var (
	auth      = smtp.PlainAuth("", confperso.EmailUsername, confperso.EmailPassword, confperso.EmailHost)
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
    b := bytes.NewBufferString("")
    if err := tmpl.Execute(b, data); err != nil {
        beego.Error(err)
        return err
    }
    sendMail(user.Email, b.String())
	return nil
}

func sendMail(to string, body string) {
    if err := smtp.SendMail("", auth, confperso.EmailUsername, []string{to}, []byte(body)); err != nil {
        beego.Warn("MAIL ERROR", err)
    }

}