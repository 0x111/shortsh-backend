package models

import (
	"time"
)

type Visitors struct {
	Id        int64     `xorm:"'id' index pk autoincr" json:"id"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`

	Ip       string `xorm:"ip" json:"url,required"`
	Url      int64  `xorm:"index" json:"url_id"` // relation to the Shortened URL
	Referrer string `xorm:"referrer" json:"referrer"`
}
