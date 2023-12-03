package store

import (
	"example.com/go_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetPgConn() *gorm.DB {
	// 实际应用中dsn从配置文件中读取
	dsn := "host=localhost user=tyler password=123 dbname=playground port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: false},
	})
	if err != nil {
		panic("failed to connect postgres database")
	}
	for _, s := range models.ModelList {
		if !db.Migrator().HasTable(s) {
			db.Migrator().CreateTable(s)
		}
	}
	return db
}

// db.Model(&user).Where("name = ?", "Race").Preload("UserInfo.MemberRight").First(&user)
// db.Model(&user).Association("UserInfo").Find(&user_info)
