package orm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitTestSuite() {
	dsn := ""
	err := NewOrmClient(&Config{
		Config: &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		},
		Dial: postgres.Open(dsn),
	})
	if err != nil {
		panic(err)
	}
}
