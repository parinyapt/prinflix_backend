package database

import (
	"os"
	"time"

	originalmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/parinyapt/prinflix_backend/logger"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	utilsDatabase "github.com/parinyapt/prinflix_backend/utils/database"
)

var DB *gorm.DB

func initializeConnectMariaDB() {
	host := os.Getenv("DATABASE_MARIADB_HOST")
	if os.Getenv("DATABASE_MARIADB_PORT") != "" {
		host = os.Getenv("DATABASE_MARIADB_HOST") + ":" + os.Getenv("DATABASE_MARIADB_PORT")
	}

	dsn := originalmysql.Config{
		User:                 os.Getenv("DATABASE_MARIADB_USERNAME"),
		Passwd:               os.Getenv("DATABASE_MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 host,
		DBName:               os.Getenv("DATABASE_MARIADB_DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
	}

	var gormConfig *gorm.Config
	if os.Getenv("DEPLOY_MODE") == "development" {
		gormConfig = &gorm.Config{}
	} else {
		gormConfig = &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		}
	}

	database, err := gorm.Open(mysql.Open(dsn.FormatDSN()), gormConfig)
	if err != nil {
		logger.Fatal("Failed to connect MariaDB database", logger.Field("error", err))
	}

	// AutoMigrate database
	err = database.AutoMigrate(
		modelDatabase.Account{},
		modelDatabase.AuthSession{},
		modelDatabase.AccountOAuth{},
		modelDatabase.TemporaryCode{},
		modelDatabase.MovieCategory{},
		modelDatabase.Movie{},
		modelDatabase.FavoriteMovie{},
		modelDatabase.WatchSession{},
		modelDatabase.WatchHistory{},
		modelDatabase.OauthState{},
		modelDatabase.Review{},
	)
	if err != nil {
		logger.Fatal("Failed to AutoMigrate database", logger.Field("error", err))
	}

	review_stat_query := database.Raw(`
		WITH review_stat AS(
			SELECT
					review_movie_uuid,
					COUNT(review_rating) AS review_total_count,
					COUNT(CASE WHEN review_rating = 3 THEN 1 END) AS review_good_count,
					COUNT(CASE WHEN review_rating = 2 THEN 1 END) AS review_fair_count,
					COUNT(CASE WHEN review_rating = 1 THEN 1 END) AS review_bad_count
			FROM prinflix_review
			GROUP BY review_movie_uuid
		)
		SELECT
			movie_uuid,
			CASE 
				WHEN review_total_count IS NULL THEN 0 
				ELSE review_total_count
			END AS review_total_count,
			CASE 
				WHEN review_good_count IS NULL THEN 0 
				ELSE review_good_count
			END AS review_good_count,
			CASE 
				WHEN review_fair_count IS NULL THEN 0 
				ELSE review_fair_count
			END AS review_fair_count,
			CASE 
				WHEN review_bad_count IS NULL THEN 0 
				ELSE review_bad_count
			END AS review_bad_count
		FROM prinflix_movie
		LEFT JOIN review_stat ON prinflix_movie.movie_uuid = review_stat.review_movie_uuid;
	`)
	database.Migrator().CreateView(utilsDatabase.GenerateTableName("view_review_stat"), gorm.ViewOption{Query: review_stat_query, Replace: true})

	DB = database

	logger.Info("Initialize MariaDB Database Success")
}
