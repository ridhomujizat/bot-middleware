package postgre

import (
	"bot-middleware/internal/pkg/util"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type DB struct {
	Master *gorm.DB
}

func getDSN() string {

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
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
	})
	if err != nil {
		util.HandleAppError(err, "initDBMaster", "connectionMasterDB", true)
		return nil, err
	}

	fmt.Println("Database connection established")

	return gormDB, nil
}

func GetDB() (*DB, error) {
	var connection DB
	master, err := initDBMaster()
	if err != nil {
		return &connection, err
	}
	connection.Master = master

	return &connection, nil
}
