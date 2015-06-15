package controllers

type PedagoController struct {
    BaseController
}

func (c *PedagoController) ValidateProjectView() {
    c.TplNames = "pedago/project-validation.html"
}