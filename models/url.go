package models

import (
	"github.com/0x111/shortsh-backend/randomizer"
	"time"
)

type Url struct {
	Id        int64     `xorm:"'id' index pk autoincr" json:"id"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`

	Url string `xorm:"url" json:"url,required"`
	//ShortDomain string `xorm:"short_domain" json:"short_domain"`
	ShortId     string `xorm:"short_id" json:"short_id"`
	ShortDomain int64  `xorm:"index" json:"-"`
}

func (t *Url) SetRandomID(length int) {
	if length == 0 {
		length = 3
	}

	t.ShortId = randomizer.String(length)
}

func (t *Url) SetDomain() {
	// 1 is now hardcoded, it is the index of the first short domain which is reserved by the system
	t.ShortDomain = 1
}
