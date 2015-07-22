package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/docker/docker/pkg/stringutils"
)

func init() {
	orm.RegisterModel(&Project{})
}

type Project struct {
	Id                  int            `json:"id" form:"id" csv:"id"`
	Name                string         `json:"name" orm:"unique;size(128)" form:"name" csv:"nom du projet"`
	ShortDescription    string         `json:"shortDescription" orm:"size(256)" form:"shortDescription" csv:"description"`
	CompleteDescription string         `json:"completeDescription" orm:"null;type(text)" form:"completeDescription" csv:"-"`
	Repository          string         `json:"repository" form:"repository" csv:"repository"`
	Website             string         `json:"website" form:"website" csv:"website"`
	Image               string         `json:"image" form:"image" csv:"-"`
	Status              *ProjectStatus `json:"status" orm:"rel(fk)"`
	Manager             *User          `json:"manager" orm:"null;rel(fk)"`
	History             []*HistoryItem `json:"history" orm:"null;rel(m2m)" csv:"-"`
	Members             []*User        `json:"members" orm:"null;reverse(many);rel_through(intra-hub/models.UserProjects)" csv:"-"`
	Themes              []*Theme       `json:"themes" orm:"null;rel(m2m)" csv:"-"`
	Technos             []*Skill       `json:"themes" orm:"null;rel(m2m)" csv:"-"`
	Comments            []*Comment     `json:"comments" orm:"null;rel(m2m)" csv:"-"`
	Created             time.Time      `json:"created" orm:"auto_now_add;type(datetime)" csv:"creation"`
	Updated             time.Time      `json:"updated" orm:"auto_now;type(datetime)" csv:"mise a jour"`

	// Non Persistent fields

	ManagerLogin   string `orm:"-" form:"managerLogin" csv:"-"`
	MembersID      string `orm:"-" form:"membersId" csv:"-"`
	ThemesID       string `orm:"-" form:"themesId" csv:"-"`
	TechnosID      string `orm:"-" form:"technosId" csv:"-"`
	MemberCount    int    `json:"memberCount" orm:"-" csv:"nombre de membres"`
	StatusName     string `orm:"-" form:"status" csv:"-"`
	MembersName    string `orm:"-" form:'-" csv:"membres"`
	SkillsName     string `orm:"-" form:"-" csv:"skills"`
	ThemesName     string `orm:"-" form:"-" csv:"themes"`
	CitiesName     string `orm:"-" form:"-" csv:"villes"`
	PromotionsName string `orm:"-" form:"-" csv:"promotions"`
}

func (p *Project) Cities() (s string) {
	m := make(map[string]bool)
	for _, member := range p.Members {
		if member.City.Name == "" {
			continue
		}
		m[member.City.Name] = true
	}
	for city := range m {
		s += city + ", "
	}
	if len(s) > 2 {
		s = s[:len(s)-2]
	}
	return
}

func (p *Project) Promotions() (s string) {
	m := make(map[string]bool)
	for _, member := range p.Members {
		if member.Promotion.Name == "" {
			continue
		}
		m[member.Promotion.Name] = true
	}
	for promo := range m {
		s += promo + ", "
	}
	if len(s) > 2 {
		s = s[:len(s)-2]
	}
	return
}

func (p *Project) IsManager(login string) bool {
	return p.Manager != nil && p.Manager.Login == login
}

func (p *Project) HasTechno(technoID int) bool {
	for _, t := range p.Technos {
		if t.Id == technoID {
			return true
		}
	}
	return false
}

func (p *Project) HasTheme(themeID int) bool {
	for _, t := range p.Themes {
		if t.Id == themeID {
			return true
		}
	}
	return false
}

func (p *Project) Valid(v *validation.Validation) {
	if !stringutils.InSlice(EveryProjectStatus, p.StatusName) {
		v.SetError("Status", "unknown or empty status: "+p.StatusName)
	}
	if p.Name == "" {
		v.SetError("Name", "empty name")
	}
	if p.ShortDescription == "" {
		v.SetError("ShortDescription", "short description empty")
	}
}
