package utils

import (
	"github.com/0x111/shortsh-backend/models"
	"github.com/go-xorm/xorm"
)

func UrlExists(engine *xorm.Engine, url string) (shortshurl *models.ShortShUrl, exists bool) {
	var urlMeta = &models.ShortShUrl{Url: url}
	has, err := engine.Get(urlMeta)

	if !has || err != nil {
		return nil, false
	}

	return urlMeta, true
}
