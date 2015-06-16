package models

import "github.com/astaxie/beego/orm"

const (
	ProjectStatusWaitingManager        = "WAITING_MANAGER"
	ProjectStatusToFormalize           = "TO_FORMALIZE"
	ProjectStatusPitchRequired         = "PITCH_REQUIRED"
	ProjectStatusInProgress            = "IN_PROGRESS"
	ProjectStatusWaitingForOralDefense = "WAITING_FOR_DEFENSE"
	ProjectStatusCreditsGiven          = "CREDITS_GIVEN"
	ProjectStatusAbandoned             = "ABANDONED"
)

var (
	EveryProjectStatus = []string{ProjectStatusWaitingManager, ProjectStatusToFormalize,
		ProjectStatusPitchRequired, ProjectStatusInProgress,
        ProjectStatusWaitingForOralDefense, ProjectStatusCreditsGiven,
		ProjectStatusAbandoned}
)

func init() {
	orm.RegisterModel(&ProjectStatus{})
}

type ProjectStatus struct {
	Id   int
	Name string `orm:"size(128);unique"`
}
