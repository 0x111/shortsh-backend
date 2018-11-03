package models

import (
	"github.com/0x111/shortsh-backend/randomizer"
	"time"
)

type ShortShUrl struct {
	Id        int64     `xorm:"'id' index pk autoincr" json:"id"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`

	Url         string `xorm:"url" json:"url,required"`
	ShortDomain string `xorm:"short_domain" json:"short_domain"`
	ShortId     string `xorm:"short_id" json:"short_id"`
}

func (t *ShortShUrl) SetRandomID(length int) {
	if length == 0 {
		length = 3
	}

	t.ShortId = randomizer.String(length)
}

func (t *ShortShUrl) SetDomain() {
	t.ShortDomain = "a2.to"
}
