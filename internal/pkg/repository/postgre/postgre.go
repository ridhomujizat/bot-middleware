package postgre

import (
	"bot-middleware/internal/pkg/util"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func getDSN() string {

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		util.GodotEnv("DB_HOST"),
		util.GodotEnv("DB_USERNAME"),
		util.GodotEnv("DB_PASSWORD"),
		util.GodotEnv("DB_NAME"),
		util.GodotEnv("DB_PORT"),
	)
}

func initDBMaster() (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		util.HandleAppError(err, "initDBMaster", "connectionMasterDB", true)
		return nil, err
	}

	fmt.Println("Database connection established")

	return gormDB, nil
}

func GetDB() (*gorm.DB, error) {
	db, err := initDBMaster()
	if err != nil {
		return db, err
	}

	return db, nil
}
