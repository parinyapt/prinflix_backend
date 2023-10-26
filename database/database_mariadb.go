package database

import (
	"os"
	"time"

	originalmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/parinyapt/prinflix_backend/logger"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

var DB *gorm.DB

func initializeConnectMariaDB() {
	host := os.Getenv("DATABASE_MARIADB_HOST")
	if os.Getenv("DATABASE_MARIADB_PORT") != "" {
		host = os.Getenv("DATABASE_MARIADB_HOST") + ":" + os.Getenv("DATABASE_MARIADB_PORT")
	}

	dsn := originalmysql.Config{
		User:      os.Getenv("DATABASE_MARIADB_USERNAME"),
		Passwd:    os.Getenv("DATABASE_MARIADB_PASSWORD"),
		Net:       "tcp",
		Addr:      host,
		DBName:    os.Getenv("DATABASE_MARIADB_DBNAME"),
		AllowNativePasswords: true,
		ParseTime: true,
		Loc:       time.Local,
	}
	database, err := gorm.Open(mysql.Open(dsn.FormatDSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect MariaDB database", logger.Field("error", err))
	}

	// AutoMigrate database
	err = database.AutoMigrate(
		modelDatabase.Account{},
		modelDatabase.AuthSession{},
	)
	if err != nil {
		logger.Fatal("Failed to AutoMigrate database", logger.Field("error", err))
	}

	DB = database

	logger.Info("Initialize MariaDB Database Success")
}