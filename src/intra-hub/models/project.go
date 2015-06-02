package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/docker/docker/pkg/stringutils"
	"strconv"
	"strings"
	"time"
)

func init() {
	orm.RegisterModel(&Project{})
}

const (
	ProjectStatusWaitingManager              = "WAITING_MANAGER"
	ProjectStatusToFormalize                 = "TO_FORMALIZE"
	ProjectStatusPedagogicalValidationNeeded = "PEDAGOGICAL_VALIDATION_NEEDED"
	ProjectStatusPitchRequired               = "PITCH_REQUIRED"
	ProjectStatusInProgress                  = "IN_PROGRESS"
	ProjectStatusWaitingForOralDefense       = "WAITING_FOR_DEFENSE"
	ProjectStatusCreditsGiven                = "CREDITS_GIVEN"
	ProjectStatusBlockedByPedagogy           = "BLOCKED_BY_PEDAGOGY"
	ProjectStatusAbandoned                   = "ABANDONED"
)

var (
	EveryProjectStatus = []string{ProjectStatusWaitingManager, ProjectStatusToFormalize,
		ProjectStatusPedagogicalValidationNeeded, ProjectStatusPitchRequired,
		ProjectStatusInProgress, ProjectStatusWaitingForOralDefense, ProjectStatusCreditsGiven,
		ProjectStatusBlockedByPedagogy, ProjectStatusAbandoned}
)

type Project struct {
	Id               int
	Name             string         `orm:"unique;size(128)" form:"name"`
	ShortDescription string         `orm:"size(128)" form:"shortDescription"`
	Status           string         `orm"size(128)" form:"status"`
	History          []*HistoryItem `orm:"null;rel(m2m)"`
	Members          []*User        `orm:"null;reverse(many)"`
	Manager          *User          `orm:"null;rel(fk)"`
	Themes           []*Theme       `orm:"null;rel(m2m)"`
	Technos          []*Techno      `orm:"null;rel(m2m)"`
	Created          time.Time      `orm:"auto_now_add;type(datetime)"`
	Updated          time.Time      `orm:"auto_now;type(datetime)"`

	// Non Persistent fields

	ManagerLogin string `orm:"-" form:"managerLogin"`
	MembersID    string `orm:"-" form:"membersId"`
	MemberCount  int    `orm:"-"`
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
		if member.City.Name == "" {
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
	if !stringutils.InSlice(EveryProjectStatus, p.Status) {
		v.SetError("Status", "unknown or empty status")
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
	members := strings.Split(p.MembersID, ",")
Loop:
	for _, memberId := range members {
		id, err := strconv.ParseInt(memberId, 10, 64)
		if err != nil {
			v.SetError("MembersID", err.Error())
			return
		}
		for _, member := range p.Members {
			if int(id) == member.Id {
				continue Loop
			}
		}
		p.Members = append(p.Members, &User{Id: int(id)})
	}
}
