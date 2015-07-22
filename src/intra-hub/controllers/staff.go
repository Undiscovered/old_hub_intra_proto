package controllers

import "intra-hub/db"

type StaffController struct {
	BaseController
}

func (c *StaffController) ListView() {
	c.TplNames = "staff/list.html"
	c.RequireLogin()
	managers, err := db.GetManagersOrAdmin()
	if err != nil {
		c.SetErrorAndRedirect(err)
		return
	}
	for _, m := range managers {
		if err := db.LoadUserInfo(m); err != nil {
			c.SetErrorAndRedirect(err)
			return
		}
		if err := db.SetManagerProjects(m); err != nil {
			c.SetErrorAndRedirect(err)
			return
		}
	}
	c.Data["Managers"] = managers
}
