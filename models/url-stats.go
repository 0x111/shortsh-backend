package models

import "time"

// this will be returned on the stats path
type UrlStats struct {
	Url string `xorm:"url" json:"url"`

	//ShortDomain string `xorm:"short_domain" json:"short_domain"`
	ShortId     string `xorm:"short_id" json:"short_id"`
	ShortDomain int64  `xorm:"index" json:"-"`
	Count       int64
}

type UrlStatRet struct {
	//+----+---------------------+-----------------------------------------------------+--------------+---------------+
	//| id | created_at          | url                                                 | short_domain | visitorscount |
	//+----+---------------------+-----------------------------------------------------+--------------+---------------+
	//|  5 | 2018-11-15 17:31:07 | https://github.com/short-sh/shortsh-backend/issues/1#d |            1 |             2 |
	//+----+---------------------+-----------------------------------------------------+--------------+---------------+
	Id          int64     `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	Url         string    `json:"url"`
	ShortDomain string    `json:"short_domain"`
	Secure      int64     `json:"-"`
	Count       int64     `json:"count"`
}

type UrlStatDaily struct {
	Day    string `json:"day"`
	Clicks int64  `json:"clicks"`
}
