package utils

import (
	"github.com/go-xorm/xorm"
	"github.com/short-sh/shortsh-backend/models"
)

func UrlExists(engine *xorm.Engine, url string) (shortshurl *models.Url, exists bool) {
	var urlMeta = &models.Url{Url: url}
	has, err := engine.Get(urlMeta)

	if !has || err != nil {
		return nil, false
	}

	return urlMeta, true
}

func GetShortDomain(engine *xorm.Engine, urlMeta *models.Url) *models.ShortDomains {
	shortDomain := models.ShortDomains{Id: urlMeta.ShortDomain}
	has, err := engine.Get(&shortDomain)

	if !has || err != nil {
		return nil
	}

	return &shortDomain
}
