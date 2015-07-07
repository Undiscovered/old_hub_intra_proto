package db

import (
	"sync"

	"intra-hub/models"

	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/pikanezi/mapslice"
	"github.com/saschpe/tribool"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserTable        = "user"
	UserProjectTable = "user_projects"
)

func QueryUser() orm.QuerySeter {
	return orm.NewOrm().QueryTable(UserTable)
}

func QueryUserProject() orm.QuerySeter {
	return orm.NewOrm().QueryTable(UserProjectTable)
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
	userDB.City = user.City
	userDB.Promotion = user.Promotion
	userDB.Group = user.Group
	userDB.Email = user.Email
	userDB.Skills = user.Skills
	userDB.Themes = user.Themes
	userDB.PhoneNumber = user.PhoneNumber
	userDB.LastName = user.LastName
	userDB.FirstName = user.FirstName
	userDB.Picture = user.Picture
	if err := clearUserRelation(userDB); err != nil {
		return err
	}
	if err := setUserRelation(userDB); err != nil {
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
	err := loadUserInfo(userDb)
	return userDb, err
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

func GetEveryUserProjectsByValidation(pedagogicallyValidation tribool.Tribool) (userProjects []*models.UserProjects, err error) {
	userProjects = make([]*models.UserProjects, 0)
	o := orm.NewOrm()
	_, err = o.QueryTable(UserProjectTable).Filter("PedagogicallyValidated", pedagogicallyValidation).All(&userProjects)
	for _, u := range userProjects {
		if _, err := o.LoadRelated(u, "Project"); err != nil {
			return nil, err
		}
		if _, err := o.LoadRelated(u, "User"); err != nil {
			return nil, err
		}
	}
	return
}

func ValidatePedagogicallyUser(userID int, projectID int, pedagogicallyValidation tribool.Tribool) error {
	_, err := QueryUserProject().Filter("user_id", userID).Filter("project_id", projectID).
		Update(orm.Params{"pedagogically_validated": pedagogicallyValidation})
	return err
}

func loadEveryInfoOfUsers(users []*models.User) error {
	wg := sync.WaitGroup{}
	errorChan := make(chan error, 1)
	for _, u := range users {
		wg.Add(1)
		go func(w *sync.WaitGroup, user *models.User) {
			defer wg.Done()
			if err := loadUserInfo(user); err != nil {
				errorChan <- err
			}
		}(&wg, u)
	}
	wg.Wait()
	if len(errorChan) > 0 {
		select {
		case err := <-errorChan:
			return err
		}
	}
	return nil
}

func loadUserInfo(user *models.User) error {
	o := orm.NewOrm()
	if _, err := o.LoadRelated(user, "City"); err != nil {
		return err
	}
	if _, err := o.LoadRelated(user, "Promotion"); err != nil {
		return err
	}
	skills := make([]*models.Skill, 0)
	if _, err := o.Raw(`SELECT skill.id, skill.name, user_skills.level, user_skills.user_id, user_skills.skill_id
            FROM user_skills INNER JOIN skill
            WHERE user_skills.skill_id = skill.id
            AND user_id = ?`, mapslice.MapSliceToIntUnsafe(user.Skills, "Id"), user.Id).QueryRows(&skills); err != nil {
		return err
	}
	user.Skills = skills
	return nil
}

func setUserRelation(user *models.User) error {
	o := orm.NewOrm()
	if len(user.Skills) != 0 {
		for _, sk := range user.Skills {
			uk := &models.UserSkills{
				User:  user,
				Skill: sk,
				Level: sk.Level,
			}
			o.Insert(uk)
		}
	}
	if len(user.Themes) != 0 {
		for _, th := range user.Themes {
			ut := &models.UserThemes{
				User:  user,
				Theme: th,
				Level: th.Level,
			}
			o.Insert(ut)
		}
	}

	return nil
}

func clearUserRelation(user *models.User) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	o := orm.NewOrm()
	o.QueryM2M(user, "Skills").Clear()
	o.QueryM2M(user, "Themes").Clear()
	return err
}
