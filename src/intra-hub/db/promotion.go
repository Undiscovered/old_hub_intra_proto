package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
)

const (
	PromotionTable = "promotion"
)

func QueryPromotions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(PromotionTable)
}

func GetEveryPromotion() (promotions []*models.Promotion, err error) {
	_, err = QueryPromotions().All(&promotions)
	return
}
