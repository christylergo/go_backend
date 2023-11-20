package authen

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func wirteUserInfoToPg(user *UserInfo) error {
	dsn := "host=localhost user=tyler password=123 dbname=playground port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		panic("failed to connect database")
	}
	result := db.Create(user) // 通过数据的指针来创建
	return result.Error
}
