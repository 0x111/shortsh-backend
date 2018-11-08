package models

import "time"

type ShortDomains struct {
	Id        int64     `xorm:"'id' index pk autoincr" json:"id"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`

	ShortDomain string `xorm:"short_domain" json:"short_domain"`
	//Url   int64  `xorm:"index" json:"url_id"` // relation to the Shortened URL
}
