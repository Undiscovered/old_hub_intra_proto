package models

import "github.com/astaxie/beego/orm"

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

func init() {
	orm.RegisterModel(&ProjectStatus{})
}

type ProjectStatus struct {
	Id   int
	Name string `orm:"size(128);unique"`
}
