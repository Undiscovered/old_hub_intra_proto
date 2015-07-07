package db

import (
	"sync"

	"intra-hub/models"

	"fmt"
	"github.com/astaxie/beego"
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
	if err := clearUserRelation(user); err != nil {
		return err
	}
	if err := setUserRelation(user); err != nil {
		return err
	}
	_, err := orm.NewOrm().Update(user)
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
	err := LoadUserInfo(userDb)
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

func SetManagerProjects(user *models.User) error {
	var projects []*models.Project
	_, err := QueryProjects().Filter("Manager", user).All(&projects)
	user.ProjectsManaged = projects
	return err
}

func GetUsersPaginated(page, limit int, queryFilter map[string]interface{}) (*models.ItemPaginated, error) {
	queryHelper := make(map[string]string)
	values := make(map[string]interface{}, 0)
	for key, value := range queryFilter {
		beego.Warn(key, value)
		switch key {
		case "promotions":
			s := value.([]string)
			if s[0] == "" {
				continue
			}
			questionMarks := ""
			for range s {
				questionMarks += "?, "
			}
			if questionMarks != "" {
				questionMarks = questionMarks[:len(questionMarks)-2]
			}
			queryHelper[key] = fmt.Sprintf("AND user.promotion_id IN (%s)", questionMarks)
			values[key] = value.([]string)
		case "cities":
			s := value.([]string)
			if s[0] == "" {
				continue
			}
			questionMarks := ""
			for range s {
				questionMarks += "?, "
			}
			if questionMarks != "" {
				questionMarks = questionMarks[:len(questionMarks)-2]
			}
			queryHelper[key] = fmt.Sprintf("AND user.city_id IN (%s)", questionMarks)
			values[key] = value.([]string)
		case "skills":
			s := value.([]string)
			if s[0] == "" {
				continue
			}
			questionMarks := ""
			for range s {
				questionMarks += "?, "
			}
			if questionMarks != "" {
				questionMarks = questionMarks[:len(questionMarks)-2]
			}
			queryHelper[key] = fmt.Sprintf("AND user_skills.skill_id IN (%s)", questionMarks)
			values[key] = value.([]string)
		case "themes":
			s := value.([]string)
			if s[0] == "" {
				continue
			}
			questionMarks := ""
			for range s {
				questionMarks += "?, "
			}
			if questionMarks != "" {
				questionMarks = questionMarks[:len(questionMarks)-2]
			}
			queryHelper[key] = fmt.Sprintf("AND user_themes.theme_id IN (%s)", questionMarks)
			values[key] = value.([]string)
		case "name":
			if value.(string) == "" {
				continue
			}
			queryHelper[key] = `AND (user.first_name LIKE ? OR
user.last_name LIKE ? OR
user.login LIKE ? OR
CONCAT(user.first_name, ' ', user.last_name) LIKE ?)`
			values[key] = []string{value.(string), value.(string), value.(string), value.(string)}
		}
	}
	beego.Warn(queryHelper)
	o := orm.NewOrm()
	raw := fmt.Sprintf(`SELECT
user.id, user.login, user.first_name, user.last_name, user.email, user.picture, user.promotion_id, user.city_id,
user_skills.user_id AS skills_user_id,
user_themes.user_id AS themes_user_id,
user_projects.user_id
FROM user
INNER JOIN user_skills, user_themes, user_projects
WHERE (user_skills.user_id = user.id
OR user_themes.user_id = user.id
OR user_projects.user_id = user.id)
%s
%s
%s
%s
%s
GROUP BY user.id`, queryHelper["promotions"], queryHelper["cities"], queryHelper["skills"], queryHelper["themes"], queryHelper["name"])
	finalValues := make([]string, 0)
	for _, v := range values {
		switch v.(type) {
		case []string:
			finalValues = append(finalValues, v.([]string)...)
		case string:
			finalValues = append(finalValues, v.(string))
		}
	}
	users := make([]*models.User, 0)
	if _, err := o.Raw(raw, finalValues).QueryRows(&users); err != nil {
		return nil, err
	}
	for _, u := range users {
		if err := LoadUserInfo(u); err != nil {
			return nil, err
		}
	}
	itemPaginated := &models.ItemPaginated{
		Items:          users,
		ItemCount:      len(users),
		TotalItemCount: len(users),
		CurrentPage:    page + 1,
		TotalPageCount: len(users)/limit + 1,
	}
	return itemPaginated, nil
}

func LoadUserInfo(user *models.User) error {
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
	themes := make([]*models.Theme, 0)
	if _, err := o.Raw(`SELECT theme.id, theme.name, user_themes.level, user_themes.user_id, user_themes.theme_id
            FROM user_themes INNER JOIN theme
            WHERE user_themes.theme_id = theme.id
            AND user_id = ?`, mapslice.MapSliceToIntUnsafe(user.Themes, "Id"), user.Id).QueryRows(&themes); err != nil {
		return err
	}
	user.Themes = themes
	return nil
}

func loadEveryInfoOfUsers(users []*models.User) error {
	wg := sync.WaitGroup{}
	errorChan := make(chan error, 1)
	for _, u := range users {
		wg.Add(1)
		go func(w *sync.WaitGroup, user *models.User) {
			defer wg.Done()
			if err := LoadUserInfo(user); err != nil {
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
