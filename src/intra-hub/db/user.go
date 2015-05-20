package db
import (
    "github.com/astaxie/beego/orm"
    "intra-hub/models"
    "golang.org/x/crypto/bcrypt"
)

const (
    UserTable = "user"
)

func QueryUser() orm.QuerySeter {
    return orm.NewOrm().QueryTable(UserTable)
}

func CheckUserCredentials(user *models.User) (*models.User, error) {
    userDb := &models.User{}
    if err := QueryUser().Filter("login", user.Login).One(userDb); err != nil {
        return nil, err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(user.Password)); err != nil {
        return nil, err
    }
    return userDb, nil
}