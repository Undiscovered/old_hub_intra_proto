package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
	"intra-hub/services/cache"
)

const (
	PromotionTable = "promotion"
)

func QueryPromotions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(PromotionTable)
}

func GetEveryPromotion() (promotions []*models.Promotion, err error) {
	_, err = QueryPromotions().All(&promotions)
	cache.SetPromotions(promotions)
	return
}

func GetPromotionByName(name string) (*models.Promotion, error) {
	promo := &models.Promotion{}
	err := QueryPromotions().Filter("Name", name).One(promo)
	return promo, err
}
