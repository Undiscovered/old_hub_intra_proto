package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/docker/docker/pkg/stringutils"
)

func init() {
	orm.RegisterModel(&Project{})
}

type Project struct {
	Id                     int
	Name                   string         `json:"name" orm:"unique;size(128)" form:"name"`
	ShortDescription       string         `json:"shortDescription" orm:"size(128)" form:"shortDescription"`
	Status                 *ProjectStatus `json:"status" orm:"rel(fk)"`
	Manager                *User          `json:"manager" orm:"null;rel(fk)"`
	History                []*HistoryItem `json:"history" orm:"null;rel(m2m)"`
	Members                []*User        `json:"members" orm:"null;reverse(many)"`
	Themes                 []*Theme       `json:"themes" orm:"null;rel(m2m)"`
	Created                time.Time      `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated                time.Time      `json:"updated" orm:"auto_now;type(datetime)"`
	PedagogicallyValidated bool           `json:"pedagogicallyValidated" form:pedagogicallyValidated"`

	// Non Persistent fields

	ManagerLogin string `orm:"-" form:"managerLogin"`
	MembersID    string `orm:"-" form:"membersId"`
	ThemesID     string `orm:"-" form:"themesId"`
	MemberCount  int    `json:"memberCount" orm:"-"`
	StatusName   string `orm:"-" form:"status"`
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
	if p.MembersID == "" {
		return
	}
	// Convert the string MembersID to an array of User.
	// MembersId has the format 1,2,3,4 etc.
	members := strings.Split(p.MembersID, ",")
LoopMembers:
	for _, memberId := range members {
		id, err := strconv.ParseInt(memberId, 10, 64)
		if err != nil {
			v.SetError("MembersID", err.Error())
			return
		}
		for _, member := range p.Members {
			if int(id) == member.Id {
				continue LoopMembers
			}
		}
		p.Members = append(p.Members, &User{Id: int(id)})
	}
	// Convert the string ThemeID to an array of Theme.
	// ThemeID has the format 1,2,3,4 etc.
	themes := strings.Split(p.ThemesID, ",")
LoopTheme:
	for _, themeId := range themes {
		id, err := strconv.ParseInt(themeId, 10, 64)
		if err != nil {
			v.SetError("ThemesID", err.Error())
			return
		}
		for _, theme := range p.Themes {
			if int(id) == theme.Id {
				continue LoopTheme
			}
		}
		p.Themes = append(p.Themes, &Theme{Id: int(id)})
	}
}
