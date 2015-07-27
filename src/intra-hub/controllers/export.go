package controllers

import (
	"encoding/csv"
	"github.com/astaxie/beego"
	"github.com/gocarina/gocsv"
	"github.com/pikanezi/mapslice"
	"intra-hub/db"
	"intra-hub/models"
	"io"
	"strings"
)

type ExportController struct {
	BaseController
}

func (c *ExportController) Projects() {
	c.EnableRender = false
	c.RequireAdmin()
	pagination, err := db.GetProjectsPaginated(1, 10000, make(map[string]interface{}))
	if err != nil {
		beego.Error(err)
		return
	}
	gocsv.SetCSVWriter(func(writer io.Writer) *csv.Writer {
		myWriter := csv.NewWriter(writer)
		myWriter.Comma = ';'
		return myWriter
	})
	for _, p := range pagination.Items.([]*models.Project) {
		p.MembersName = strings.Join(mapslice.MapSliceToStringUnsafe(p.Members, "Login"), ",")
		p.SkillsName = strings.Join(mapslice.MapSliceToStringUnsafe(p.Technos, "Name"), ",")
		p.ThemesName = strings.Join(mapslice.MapSliceToStringUnsafe(p.Themes, "Name"), ",")
		p.CitiesName = p.Cities()
		p.PromotionsName = p.Promotions()
		p.MemberCount = len(p.Members)
	}
	raw, err := gocsv.MarshalBytes(pagination.Items)
	if err != nil {
		return
	}
	c.Ctx.Output.Header("Content-Type", "text/csv")
	c.Ctx.WriteString(string(raw))
}
