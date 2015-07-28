package controllers

import (
	"fmt"
	"strings"

	"intra-hub/db"
	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"strconv"
)

type ProjectController struct {
	BaseController
}

func (c *ProjectController) NestedPrepare() {
	c.RequireLogin()
}

func (c *ProjectController) IntroView() {
	c.TplNames = "project/intro.html"
}

func (c *ProjectController) ListView() {
	c.TplNames = "project/list.html"
	handleError := func(err error) {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
	}
	queryFilter := make(map[string]interface{})
	queryFilter["promotions"] = strings.Split(c.GetString("promotions", ""), ",")
	queryFilter["cities"] = strings.Split(c.GetString("cities", ""), ",")
	queryFilter["managers"] = strings.Split(c.GetString("managers", ""), ",")
	queryFilter["status"] = strings.Split(c.GetString("status", ""), ",")
	queryFilter["student"] = c.GetString("student", "")
	queryFilter["name"] = c.GetString("name", "")
	page, err := c.GetInt("page", 1)
	limit, err := c.GetInt("limit", 25)
	if page <= 0 {
		c.Redirect(fmt.Sprintf("/projects?page=1&limit=25"), 301)
		return
	}
	if limit == 0 {
		limit = 25
	}
	paginatedItems, err := db.GetProjectsPaginated(page, limit, queryFilter)
	if err != nil {
		handleError(err)
		return
	}
	paginatedItems.SetPagesToShow()
	promotions, err := db.GetEveryPromotion()
	if err != nil {
		handleError(err)
		return
	}
	cities, err := db.GetEveryCities()
	if err != nil {
		handleError(err)
		return
	}
	managers, err := db.GetManagersOrAdmin()
	if err != nil {
		handleError(err)
		return
	}
	c.Data["Status"] = models.EveryProjectStatus
	c.Data["Managers"] = managers
	c.Data["Cities"] = cities
	c.Data["Promotions"] = promotions
	c.Data["Limit"] = limit
	c.Data["PaginatedItems"] = paginatedItems
	c.Data["HasNextPage"] = paginatedItems.CurrentPage+1 <= paginatedItems.TotalPageCount
	c.Data["HasPreviousPage"] = paginatedItems.CurrentPage != 1
	c.Data["ShowGoToFirst"] = paginatedItems.PagesToShow[0] != 1
	c.Data["ShowGoToLast"] = paginatedItems.PagesToShow[len(paginatedItems.PagesToShow)-1] != paginatedItems.TotalPageCount
}

func (c *ProjectController) EditView() {
	c.RequireManager()
	c.TplNames = "project/edit.html"
	project, err := db.GetProjectByIDOrName(c.GetString(":nameOrId"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/projects?page=1&limit=15", 301)
		return
	}
	managers, err := db.GetManagersOrAdmin()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	themes, err := db.GetEveryThemes()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	skills, err := db.GetEverySkills()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	c.Data["Themes"] = themes
	c.Data["Technos"] = skills
	c.Data["Status"] = models.EveryProjectStatus
	c.Data["Managers"] = managers
	c.Data["Project"] = project
}

func (c *ProjectController) SingleView() {
	c.TplNames = "project/single.html"
	project, err := db.GetProjectByIDOrName(c.GetString(":nameOrId"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/projects?page=1&limit=15", 301)
		return
	}
	c.Data["Project"] = project
}

func (c *ProjectController) AddView() {
	c.RequireManager()
	c.TplNames = "project/add.html"
	managers, err := db.GetManagersOrAdmin()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	themes, err := db.GetEveryThemes()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	skills, err := db.GetEverySkills()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	c.Data["Themes"] = themes
	c.Data["Technos"] = skills
	c.Data["Status"] = models.EveryProjectStatus
	c.Data["Managers"] = managers
}

func (c *ProjectController) Add() {
	c.RequireManager()
	project := &models.Project{}
	if err := c.ParseForm(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	valid := validation.Validation{}
	if b, err := valid.Valid(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	} else if !b {
		beego.Error(valid.Errors[0])
		c.SetErrorAndRedirect(fmt.Errorf(valid.Errors[0].String()))
		return
	}
	if err := setProjectData(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	projectAdded, err := db.AddAndGetProject(project)
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	c.Redirect("/projects/"+strconv.FormatInt(int64(projectAdded.Id), 10), 301)
}

func (c *ProjectController) Edit() {
	c.RequireManager()
	project := &models.Project{}
	if err := c.ParseForm(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	valid := validation.Validation{}
	if b, err := valid.Valid(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	} else if !b {
		beego.Error(valid.Errors[0])
		c.SetErrorAndRedirect(fmt.Errorf(valid.Errors[0].String()))
		return
	}
	if err := setProjectData(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	if _, err := db.EditAndGetProject(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	c.Redirect("/projects/"+strconv.FormatInt(int64(project.Id), 10), 301)
}

func (c *ProjectController) CommentView() {
	c.RequireManager()
	c.TplNames = "project/comment.html"
	project, err := db.GetProjectByIDOrName(c.GetString(":nameOrId"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/projects?page=1&limit=15", 301)
		return
	}
	c.Data["Project"] = project
}

func (c *ProjectController) AddComment() {
	c.RequireManager()
	project, err := db.GetProjectByIDOrName(c.GetString(":nameOrId"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/projects?page=1&limit=15", 301)
		return
	}
	defer c.Redirect("/projects/"+strconv.FormatInt(int64(project.Id), 10)+"/comments", 301)
	comment := &models.Comment{}
	if err := c.ParseForm(comment); err != nil {
		beego.Error(err)
		return
	}
	comment.Author = c.user
	db.AddCommentToProject(comment, project)
}

func (c *ProjectController) CheckName() {
	c.EnableRender = false
	if ok := db.CheckProjectExists(c.GetString("name")); ok {
		c.Ctx.Output.Body([]byte("false"))
	} else {
		c.Ctx.Output.Body([]byte("true"))
	}
}

func setProjectData(project *models.Project) error {
	// Convert the string MembersID to an array of User.
	// MembersId has the format 1,2,3,4 etc.
	members := strings.Split(project.MembersID, ",")
LoopMembers:
	for _, memberId := range members {
		if memberId == "" {
			continue
		}
		id, err := strconv.ParseInt(memberId, 10, 64)
		if err != nil {
			return err
		}
		for _, member := range project.Members {
			if int(id) == member.Id {
				continue LoopMembers
			}
		}
		project.Members = append(project.Members, &models.User{Id: int(id)})
	}

	// Convert the string ThemeID to an array of Theme.
	// ThemeID has the format 1,2,3,4 etc.
	themes := strings.Split(project.ThemesID, ",")
LoopTheme:
	for _, themeId := range themes {
		if themeId == "" {
			continue
		}
		id, err := strconv.ParseInt(themeId, 10, 64)
		if err != nil {
			err = nil
			// If its not an id, its probably a new theme.
			theme := &models.Theme{Name: themeId}
			if theme, err = db.AddAndGetTheme(theme); err != nil {
				return err
			}
			project.Themes = append(project.Themes, theme)
			continue
		}
		for _, theme := range project.Themes {
			if int(id) == theme.Id {
				continue LoopTheme
			}
		}
		project.Themes = append(project.Themes, &models.Theme{Id: int(id)})
	}

	// Convert the string ThemeID to an array of Theme.
	// ThemeID has the format 1,2,3,4 etc.
	technos := strings.Split(project.TechnosID, ",")
LoopTechno:
	for _, technoID := range technos {
		if technoID == "" {
			continue
		}
		id, err := strconv.ParseInt(technoID, 10, 64)
		if err != nil {
			err = nil
			// If its not an id, its probably a new theme.
			techno := &models.Skill{Name: technoID}
			if techno, err = db.AddAndGetSkill(techno); err != nil {
				return err
			}
			project.Technos = append(project.Technos, techno)
			continue
		}
		// If the techno is already set in the project, we skip it.
		for _, techno := range project.Technos {
			if int(id) == techno.Id {
				continue LoopTechno
			}
		}
		project.Technos = append(project.Technos, &models.Skill{Id: int(id)})
	}
	if project.ManagerLogin != "--" {
		manager, err := db.GetUserByLogin(project.ManagerLogin)
		if err != nil {
			return err
		}
		project.Manager = manager
	}
	return nil
}
