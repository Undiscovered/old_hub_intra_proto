package db

import (
	"sync"

	"intra-hub/models"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserTable = "user"
)

func QueryUser() orm.QuerySeter {
	return orm.NewOrm().QueryTable(UserTable)
}

func AddUser(user *models.User) error {
	_, err := orm.NewOrm().Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func EditUserByLogin(login string, user *models.User) error {
	userDB, err := GetUserByLogin(login)
	if err != nil {
		return err
	}
	beego.Warn(user.Skills)
	beego.Warn(userDB)
	userDB.City = user.City
	userDB.Promotion = user.Promotion
	userDB.Group = user.Group
	userDB.Email = user.Email
	userDB.Skills = user.Skills
	o := orm.NewOrm()
	m2m := o.QueryM2M(userDB, "Skills")
	if _, err := m2m.Clear(); err != nil {
		return err
	}
	if _, err := m2m.Add(userDB.Skills); err != nil {
		return err
	}
	_, err = orm.NewOrm().Update(userDB)
	return err
}

func CheckUserCredentials(user *models.User) (*models.User, error) {
	userDb, err := GetUserByLogin(user.Login)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(user.Password)); err != nil {
		return nil, err
	} else if userDb.Password == "" {
		return nil, fmt.Errorf("password not set")
	}
	return userDb, nil
}

func SearchUsers(query string) (usersFound []*models.User, err error) {
	query = "%" + query + "%"
	rawSql := `SELECT id, picture, first_name, last_name, login FROM user WHERE
    first_name LIKE ? OR
    last_name LIKE ? OR
    login LIKE ? OR
    CONCAT(first_name, ' ', last_name) LIKE ? LIMIT 15`
	_, err = orm.NewOrm().Raw(rawSql, query, query, query, query).QueryRows(&usersFound)
	return
}

func GetManagers() (managers []*models.User, err error) {
	_, err = QueryUser().Filter("Group__Name", models.UserGroupManager).All(&managers)
	return
}

func GetManagersOrAdmin() (managers []*models.User, err error) {
    _, err = QueryUser().Filter("Group__Name__in", []string{models.UserGroupManager, models.UserGroupAdmin}).All(&managers)
    return
}

func GetUserByLogin(login string) (*models.User, error) {
	userDb := &models.User{}
	if err := QueryUser().Filter("Login", login).RelatedSel().One(userDb); err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	if _, err := o.LoadRelated(userDb, "Skills"); err != nil {
		return nil, err
	}
	return userDb, nil
}

func ActivateUser(id int, token, password string) (*models.User, error) {
	userDb := &models.User{}
	o := orm.NewOrm()
	if err := o.QueryTable(UserTable).Filter("Id", id).Filter("Token", token).RelatedSel().One(userDb); err != nil {
		return nil, err
	}
	userDb.Token = ""
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userDb.Password = string(hash)
	if _, err := o.Update(userDb); err != nil {
		return nil, err
	}
	return userDb, nil
}

func CheckUserExists(id int, token string) error {
	userDb := &models.User{}
	if err := QueryUser().Filter("Id", id).Filter("Token", token).One(userDb); err != nil {
		return err
	}
	return nil
}

func loadEveryInfoOfUsers(users []*models.User) error {
	wg := sync.WaitGroup{}
	errorChan := make(chan error)
	for _, u := range users {
		wg.Add(1)
		go func(w *sync.WaitGroup, user *models.User) {
			defer wg.Done()
			o := orm.NewOrm()
			if _, err := o.LoadRelated(user, "City"); err != nil {
				errorChan <- err
			}
			if _, err := o.LoadRelated(user, "Promotion"); err != nil {
				errorChan <- err
			}
		}(&wg, u)
	}
	wg.Wait()
	if len(errorChan) > 0 {
		beego.Error("ERROR", errorChan)
		select {
		case err := <-errorChan:
			return err
		}
	}
	return nil
}
