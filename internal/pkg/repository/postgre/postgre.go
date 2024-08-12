package postgre

import (
	"bot-middleware/internal/pkg/util"
	"fmt"

	appAccount "bot-middleware/internal/application/account"
	appBot "bot-middleware/internal/application/bot"
	appSession "bot-middleware/internal/application/session"

	"strconv"

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
		util.HandleAppError(err, "initDBMaster", "Gorm Open", true)
	}

	sync, err := strconv.ParseBool(util.GodotEnv("DB_SYNC"))
	if err != nil {
		util.HandleAppError(err, "initDBMaster", "ParseBool", true)
	}

	fmt.Println("sync", sync)

	if sync {
		err := gormDB.AutoMigrate(&appAccount.AccountSetting{})
		if err != nil {
			return nil, err
		}
		err = gormDB.AutoMigrate(&appBot.ServerBot{})
		if err != nil {
			return nil, err
		}
		err = gormDB.AutoMigrate(&appSession.Session{})
		if err != nil {
			return nil, err
		}
		err = gormDB.AutoMigrate(&appSession.SessionHistory{})
		if err != nil {
			return nil, err
		}
	}

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
