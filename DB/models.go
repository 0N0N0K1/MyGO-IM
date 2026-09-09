package DB

import (
	"gorm.io/gorm"
	"time"
)

const (
	ByID = iota
	ByEmail
	ByName
	ALL
)

type DetailInfo struct {
	User
	UserInfo
}
type User struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"size:64;not null"`
	Online     int    `gorm:"-"`
	LastOnline *time.Time
	Email      string         `gorm:"size:128;unique"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Password   string         `gorm:"size:255;not null" json:"password,omitempty"`
}

type UserInfo struct {
	UserID    uint   `gorm:"primaryKey"` // 关联到 User.ID
	Addr      string `gorm:"size:255"`
	Age       int
	Gender    string    `gorm:"size:10"`
	Signature string    `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

type Group struct {
	gorm.Model
	GroupName string `gorm:"size:100;not null"`
	OwnerID   uint   `gorm:"primaryKey;comment:群主ID"`
	OwnerName string `gorm:"size:100;not null"`

	Owner *User `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserUser struct {
	ActiveID  uint      `gorm:"primaryKey;comment:主动添加方"`
	PassiveID uint      `gorm:"primaryKey;comment:被动添加方"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Active  *User `gorm:"foreignKey:ActiveID;references:ID;constraint:OnDelete:CASCADE"`
	Passive *User `gorm:"foreignKey:PassiveID;references:ID;constraint:OnDelete:CASCADE"`
}

type GroupUser struct {
	GroupID   uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Group *Group `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	User  *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
type GroupMessage struct {
	ID       uint `gorm:"primaryKey"`
	ToID     uint
	FromID   uint
	FromName string
	ToName   string
	Content  string `gorm:"type:text"`

	To   *Group `gorm:"foreignKey:ToID;references:ID;constraint:OnDelete:CASCADE"  json:"-"`
	From *User  `gorm:"foreignKey:FromID;references:ID;constraint:OnDelete:CASCADE"  json:"-"`
}
type PrivateMessage struct {
	ID       uint `gorm:"primaryKey"`
	ToID     uint
	FromID   uint
	FromName string
	ToName   string
	Content  string `gorm:"type:text"`

	To   *User `gorm:"foreignKey:ToID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	From *User `gorm:"foreignKey:FromID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
