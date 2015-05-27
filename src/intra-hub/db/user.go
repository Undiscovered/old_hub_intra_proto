package db

import (
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"intra-hub/models"
)

const (
	UserTable = "user"
)

func QueryUser() orm.QuerySeter {
	return orm.NewOrm().QueryTable(UserTable)
}

func CheckUserCredentials(user *models.User) (*models.User, error) {
	userDb, err := GetUserByLogin(user.Login)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(user.Password)); err != nil {
		return nil, err
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
	_, err = QueryUser().Filter("Groups__Group__Name", models.UserGroupManager).All(&managers)
	return
}

func GetUserByLogin(login string) (*models.User, error) {
	userDb := &models.User{}
	if err := QueryUser().Filter("Login", login).One(userDb); err != nil {
		return nil, err
	}
	return userDb, nil
}
