package db

import (
	"sync"

	"intra-hub/models"
	"intra-hub/services/cache"

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

func getQuestionMarksFromSlice(s []string) (questionMarks string) {
	for range s {
		questionMarks += "?, "
	}
	if questionMarks != "" {
		return questionMarks[:len(questionMarks)-2]
	}
	return
}

func getQuestionMarks(count int) (questionMarks string) {
	for i := 0; i < count; i++ {
		questionMarks += "?, "
	}
	return
}

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
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, err
	}
	qb.Select("id", "picture", "first_name", "last_name", "login").
		From("user").
		Where("first_name LIKE ?").
		Or("last_name LIKE ?").
		Or("login LIKE ?").
		Or("CONCAT(first_name, ' ', last_name) LIKE ?").
		Limit(15)
	_, err = orm.NewOrm().Raw(qb.String(), query, query, query, query).QueryRows(&usersFound)
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

func GetUserByEmail(email string) (*models.User, error) {
	userDb := &models.User{}
	if err := QueryUser().Filter("Email", email).RelatedSel().One(userDb); err != nil {
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
		if _, err := o.LoadRelated(u.User, "Promotion"); err != nil {
			return nil, err
		}
	}
	return
}

// Example requests constructed:
// UPDATE `user_projects` SET `pedagogically_validated` = CASE
// WHEN `user_id` = 16475 AND `project_id` = 3 THEN 2
// WHEN `user_id` = 20438 AND `project_id` = 1 THEN 2
// END
// WHERE (`user_id` = 16475 AND `project_id` = 3) OR (`user_id` = 20438 AND `project_id` = 1)
func PedagogicallyValidation(userIDs, projectIDs, validation []interface{}) error {
	rawSQL := "UPDATE `user_projects` SET `pedagogically_validated` = CASE\n"
	valueArray := make([]interface{}, 0, len(userIDs)+len(projectIDs)+len(validation))
	for i := range userIDs {
		rawSQL += fmt.Sprintf("WHEN `user_id` = ? AND `project_id` = ? THEN ?\n")
		valueArray = append(valueArray, userIDs[i], projectIDs[i], validation[i])
	}
	rawSQL += "END\nWHERE "
	for i := range userIDs {
		rawSQL += "(`user_id` = ? AND `project_id` = ?)"
		valueArray = append(valueArray, userIDs[i], projectIDs[i])
		if i+1 < len(userIDs) {
			rawSQL += " OR\n"
		}
	}
	_, err := orm.NewOrm().Raw(rawSQL, valueArray...).Exec()
	return err
}

func SetManagerProjects(user *models.User) error {
	var projects []*models.Project
	_, err := QueryProjects().Filter("Manager", user).All(&projects)
	user.ProjectsManaged = projects
	return err
}

func GetUsersPaginated(page, limit int, queryFilter map[string]interface{}) (*models.ItemPaginated, error) {
	queryHelper := make([]interface{}, 0)
	values := make([]interface{}, 0)
	for key, value := range queryFilter {
		switch key {
		case "promotions":
			s := value.([]string)
			if s[0] == "" {
				queryHelper = append(queryHelper, "")
				continue
			}
			queryHelper = append(queryHelper, fmt.Sprintf("OR (user.promotion_id IN (%s))", getQuestionMarksFromSlice(s)))
			promotionsIDs := make([]int, len(s))
			for i, ss := range s {
				promotionsIDs[i] = cache.Promotions[ss].Id
			}
			values = append(values, promotionsIDs)
		case "cities":
			s := value.([]string)
			if s[0] == "" {
				queryHelper = append(queryHelper, "")
				continue
			}
			queryHelper = append(queryHelper, fmt.Sprintf("OR (user.city_id IN (%s))", getQuestionMarksFromSlice(s)))
			citiesIDs := make([]int, len(s))
			for i, ss := range s {
				citiesIDs[i] = cache.Cities[ss].Id
			}
			values = append(values, citiesIDs)
		case "skills":
			s := value.([]string)
			if s[0] == "" {
				queryHelper = append(queryHelper, "")
				continue
			}
			queryHelper = append(queryHelper, fmt.Sprintf("OR (user_skills.user_id = user.id AND user_skills.skill_id IN (%s))", getQuestionMarksFromSlice(s)))
			skillIDs := make([]int, len(s))
			for i, ss := range s {
				beego.Warn(ss)
				skillIDs[i] = cache.Skills[ss].Id
			}
			values = append(values, skillIDs)
		case "themes":
			s := value.([]string)
			if s[0] == "" {
				queryHelper = append(queryHelper, "")
				continue
			}
			queryHelper = append(queryHelper, fmt.Sprintf("OR (user_themes.user_id = user.id AND user_themes.theme_id IN (%s))", getQuestionMarksFromSlice(s)))
			themeIDs := make([]int, len(s))
			for i, ss := range s {
				themeIDs[i] = cache.Themes[ss].Id
			}
			values = append(values, themeIDs)
		case "name":
			if value.(string) == "" {
				queryHelper = append(queryHelper, "")
				continue
			}
			queryHelper = append(queryHelper, `OR user.first_name LIKE ? OR
user.last_name LIKE ? OR
CONCAT(user.first_name, ' ', user.last_name) LIKE ?`)
			values = append(values, []string{"%" + value.(string) + "%", "%" + value.(string) + "%", "%" + value.(string) + "%"})
		case "login":
			if value.(string) == "" {
				queryHelper = append(queryHelper, "")
				continue
			}
			queryHelper = append(queryHelper, `OR (user.login LIKE ?)`)
			values = append(values, "%"+value.(string)+"%")
		case "email":
			if value.(string) == "" {
				queryHelper = append(queryHelper, "")
				continue
			}
			queryHelper = append(queryHelper, `OR (user.email LIKE ?)`)
			values = append(values, "%"+value.(string)+"%")
		}
	}
	for i, v := range queryHelper {
		s := v.(string)
		if s != "" {
			queryHelper[i] = "AND (" + queryHelper[i].(string)[2:]
			queryHelper[len(queryHelper)-1] = queryHelper[len(queryHelper)-1].(string) + ")"
			break
		}
	}
	if queryHelper[0] != "" {
	}
	queryHelper = append(queryHelper, limit, (page-1)*limit)
	o := orm.NewOrm()
	raw := fmt.Sprintf(`SELECT
user.id, user.login, user.first_name, user.last_name, user.email, user.picture, user.promotion_id, user.city_id,
user_skills.user_id,
user_themes.user_id,
user_projects.user_id
FROM user
LEFT JOIN user_skills ON user_skills.user_id = user.id
LEFT JOIN user_themes ON user_themes.user_id = user.id
LEFT JOIN user_projects ON user_projects.user_id = user.id
WHERE (user_skills.user_id = user.id
OR user_themes.user_id = user.id
OR user_projects.user_id = user.id)
%s %s %s %s %s %s %s
GROUP BY user.id
LIMIT %d OFFSET %d`, queryHelper...)
	users := make([]*models.User, 0)
	if _, err := o.Raw(raw, values).QueryRows(&users); err != nil {
		return nil, err
	}
	for _, u := range users {
		if err := LoadUserInfo(u); err != nil {
			return nil, err
		}
	}
	rawCount := fmt.Sprintf(`SELECT COUNT(DISTINCT user.id) AS count
FROM user
LEFT JOIN user_skills ON user_skills.user_id = user.id
LEFT JOIN user_themes ON user_themes.user_id = user.id
LEFT JOIN user_projects ON user_projects.user_id = user.id
WHERE (user_skills.user_id = user.id
OR user_themes.user_id = user.id
OR user_projects.user_id = user.id)
%s %s %s %s %s %s %s
GROUP BY user.id`, queryHelper[:len(queryHelper)-2]...)
	var res []orm.Params
	if _, err := o.Raw(rawCount, values).Values(&res); err != nil {
		return nil, err
	}
	count := len(res)
	itemPaginated := &models.ItemPaginated{
		Items:          users,
		ItemCount:      len(users),
		TotalItemCount: count,
		CurrentPage:    page,
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

func LoadUserProjects(user *models.User) error {
	_, err := QueryProjects().Filter("Members__User__Id", user.Id).All(&user.Projects)
	return err
}

func LoadUserProjectsManaged(user *models.User) error {
	_, err := QueryProjects().Filter("Manager", user).All(&user.ProjectsManaged)
	return err
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
