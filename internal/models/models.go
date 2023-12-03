package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(50);not null"`
	Phone    uint   `gorm:"not null;unique"`
	Email    string
	PassWord string   `gorm:"type:char(64);"` // primitive pwd + createdAt HASH256
	UserInfo UserInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// keep the hash sum of primitive pwd as the user's password
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now().Local()
	h := sha256.New()
	_, err := h.Write([]byte(u.PassWord + u.CreatedAt.String()))
	if err != nil {
		return err
	}
	u.PassWord = hex.EncodeToString(h.Sum(nil))
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	h := sha256.New()
	_, err := h.Write([]byte(u.PassWord + u.CreatedAt.String()))
	if err != nil {
		return err
	}
	u.PassWord = hex.EncodeToString(h.Sum(nil))
	return nil
}

type UserInfo struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	UserID      uint
	MemberLevel uint8
	MemberRight MemberRight `gorm:"foreignKey:MemberLevel;references:Level;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// MemberRight 是UserInfo的主表，应在UserInfo之前创建，不然会报错
type MemberRight struct {
	ID       uint    `gorm:"primaryKey;autoIncrement"`
	Level    uint8   `gorm:"unique;not null"`
	Discount float32 `gorm:"not null"`
}

type Item struct {
	gorm.Model
	BarCode     int         `gorm:"unique;not null"`
	Name        string      `gorm:"type:varchar(100);not null"`
	ItemPicture ItemPicture `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ItemPrice   ItemPrice   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Products    []Product   `gorm:"many2many:item_products;"`
}

type ItemPicture struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	PicPath string
	ItemID  uint
}

type ItemPrice struct {
	gorm.Model
	Price  float32
	ItemID uint
}

type Product struct {
	gorm.Model
	BarCode     int         `gorm:"unique;not null"`
	Name        string      `gorm:"type:varchar(100);not null"`
	ProductInfo ProductInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ProductInfo struct {
	gorm.Model
	Price     float32
	Weight    int16
	Volume    float32
	ProductID uint
}

// ordered by master follow relationship
var ModelList []interface{} = []interface{}{
	&User{}, &MemberRight{}, &UserInfo{}, &Item{}, &ItemPicture{}, &ItemPrice{},
	&Product{}, &ProductInfo{},
}

func (u *Item) Get(name string) any {
	return reflect.ValueOf(u).FieldByName(name).Interface()
}

func (u *Item) Set(name string, val any) error {

	s := reflect.ValueOf(u).FieldByName(name)
	if s.Kind() == reflect.ValueOf(val).Kind() {
		switch s.Kind() {
		case reflect.Int:
			v := val.(int)
			s.SetInt(int64(v))
		case reflect.String:
			v := val.(string)
			s.SetString(v)
		}
		return nil
	} else {
		return errors.New("x's type is not the same type as v's element type")
	}
}

// // TableName 会将 User 的表名重写为 `profiles`
// func (User) TableName() string {
// 	return "profiles"
// }
