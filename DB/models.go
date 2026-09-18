package DB

import (
	"gorm.io/gorm"
	"time"
)

const (
	ByID = iota
	ByEmail
	ByName
	ByCID
	ALL
)

type DetailInfo struct {
	User
	UserInfo
}
type User struct {
	ID         int64  `gorm:"primaryKey"`
	Name       string `gorm:"size:64;not null"`
	Online     int    `gorm:"-"`
	LastOnline *time.Time
	Email      string         `gorm:"size:128;unique"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Password   string         `gorm:"size:255;not null" json:"password,omitempty"`
}

type UserInfo struct {
	UserID    int64      `json:"-" gorm:"primaryKey"`
	Addr      string     `json:"addr" binding:"omitempty,min=1,max=255"`
	Age       int        `json:"-" binding:"omitempty,gte=0,lte=150"`
	Birthday  *time.Time `json:"birthday" binding:"omitempty,datetime=2006:01:02"`
	Gender    string     `json:"gender" binding:"omitempty,oneof=male female unknown"`
	Signature string     `json:"signature" binding:"omitempty,max=255"`
	CreatedAt time.Time  `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"-" gorm:"autoUpdateTime"`

	User *User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type Group struct {
	ID             int64     `gorm:"primarykey"`
	WriteSeq       uint64    `gorm:"default:1"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	ConversationId string
	MemberNum      uint16
	GroupName      string `gorm:"size:100;not null"`
	OwnerID        int64  `gorm:"primaryKey;comment:群主ID"`
	OwnerName      string `gorm:"size:100;not null"`
	Owner          *User  `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserUser struct {
	ActiveID       int64     `gorm:"primaryKey;comment:主动添加方"`
	PassiveID      int64     `gorm:"primaryKey;comment:被动添加方"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	ID             int64
	ConversationId string
	GlobalWriteSeq uint64 `gorm:"default:1"`
	ActiveReadSeq  uint64 `gorm:"default:1"`
	PassiveReadSeq uint64 `gorm:"default:1"`
	Status         string `gorm:""` // apply | reject | accept
	Active         *User  `gorm:"foreignKey:ActiveID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Passive        *User  `gorm:"foreignKey:PassiveID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type GroupUser struct {
	GroupID        int64  `gorm:"primaryKey"`
	UserID         int64  `gorm:"primaryKey"`
	ReadSeq        uint64 `gorm:"default:1"`
	ConversationId string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	Status         string    `gorm:""` // apply | reject | accept
	Group          *Group    `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	User           *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
type GroupMessage struct {
	ID             int64 `gorm:"primaryKey"`
	Seq            uint64
	ToID           int64
	FromID         int64
	ConversationId string
	Content        string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	To   *Group `gorm:"foreignKey:ToID;references:ID;constraint:OnDelete:CASCADE"  json:"-"`
	From *User  `gorm:"foreignKey:FromID;references:ID;constraint:OnDelete:CASCADE"  json:"-"`
}
type PrivateMessage struct {
	ID             int64 `gorm:"primaryKey"`
	ToID           int64
	FromID         int64
	ConversationId string
	Content        string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	Seq            uint64
	To             *User `gorm:"foreignKey:ToID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	From           *User `gorm:"foreignKey:FromID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type NoticeMessage struct {
	ID      int64 `gorm:"primaryKey"`
	ToID    int64
	Read    bool
	Content string `gorm:"type:text"`
}
