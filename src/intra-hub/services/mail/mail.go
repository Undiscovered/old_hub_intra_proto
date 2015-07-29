package mail

import (
	"bytes"
	"strconv"
	"text/template"

	"intra-hub/confperso"
	"intra-hub/models"

	"archive/zip"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const (
	activationURL = confperso.Protocol + "://" + confperso.Domain + "/users/activate/"
)

func SendUserCreated(user *models.User) error {
	data := make(map[string]interface{})
	data["User"] = user
	data["Link"] = activationURL + strconv.FormatInt(int64(user.Id), 10) + "/" + user.Token
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

func SendForgotPassword(user *models.User) error {
	data := make(map[string]interface{})
	data["User"] = user
	data["Link"] = activationURL + strconv.FormatInt(int64(user.Id), 10) + "/" + user.Token
	tmpl, err := template.New("mail").Parse(templateForgotPassword)
	if err != nil {
		beego.Error(err)
		return err
	}
	tmplSubject, err := template.New("subject").Parse(subjectForgotPassword)
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
	email.Subject = subject
	email.To = []string{to}
	email.HTML = body
	email.From = confperso.EmailUsername
	if err := email.Send(); err != nil {
		beego.Warn("MAIL ERROR", err)
	}
}

func SendBackupEmail(filepath string) error {
	config := `{"username":"` + confperso.EmailUsername + `","password":"` + confperso.EmailPassword + `","host":"` +
		confperso.EmailHost + `","port":` + confperso.EmailHostPort + `}`
	email := utils.NewEMail(config)
	email.Subject = "Backup - " + time.Now().Format(time.RFC3339)
	email.To = []string{confperso.EmailUsername}
	email.From = confperso.EmailUsername
	fileName := path.Base(filepath)
	zipName := path.Dir(filepath) + "/" + strings.TrimSuffix(fileName, path.Ext(fileName)) + ".zip"
	zipFile, err := os.Create(zipName)
	if err != nil {
		return err
	}
	w := zip.NewWriter(zipFile)
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := w.Create(fileName)
	if err != nil {
		return err
	}
	sqlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	if _, err := f.Write(sqlFile); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	zipFile.Close()
	if zipFile, err = os.Open(zipName); err != nil {
		return err
	}
	if _, err := email.Attach(zipFile, path.Base(zipName), "application/zip"); err != nil {
		return err
	}
	if err := email.Send(); err != nil {
		return err
	}
	beego.Warn(filepath, fileName, zipName)
	return nil
}
