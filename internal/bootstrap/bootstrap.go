package bootstrap

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"

	migrateps "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"user_service/internal/config"
	"user_service/internal/repository"
	"user_service/pkg/logging"
)

type Container struct {
	DB           *gorm.DB
	Config       *config.Config
	Repositories map[string]interface{}
}

func Init() (*Container, error) {
	logger := logging.Instance

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading configuration:", err)
		return nil, err
	}

	fmt.Println("╭╮╱╭┳━━━┳━━━┳━━━╮╭━━━┳━━━┳━━━┳╮╱╱╭┳━━┳━━━┳━━━╮\n┃┃╱┃┃╭━╮┃╭━━┫╭━╮┃┃╭━╮┃╭━━┫╭━╮┃╰╮╭╯┣┫┣┫╭━╮┃╭━━╯\n┃┃╱┃┃╰━━┫╰━━┫╰━╯┃┃╰━━┫╰━━┫╰━╯┣╮┃┃╭╯┃┃┃┃╱╰┫╰━━╮\n┃┃╱┃┣━━╮┃╭━━┫╭╮╭╯╰━━╮┃╭━━┫╭╮╭╯┃╰╯┃╱┃┃┃┃╱╭┫╭━━╯\n┃╰━╯┃╰━╯┃╰━━┫┃┃╰╮┃╰━╯┃╰━━┫┃┃╰╮╰╮╭╯╭┫┣┫╰━╯┃╰━━╮\n╰━━━┻━━━┻━━━┻╯╰━╯╰━━━┻━━━┻╯╰━╯╱╰╯╱╰━━┻━━━┻━━━╯")
	db, err := initDatabase(cfg)
	if err != nil {
		return nil, err
	}

	repositories := initRepositories(db)

	logger.Info("✅ Dependencies initialized successfully")

	return &Container{
		DB:           db,
		Config:       cfg,
		Repositories: repositories,
	}, nil
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger := logging.Instance

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Db,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Error connecting to the database:", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Error getting generic DB object:", err)
		return nil, err
	}

	g, err2, done := initMigrations(sqlDB, logger)
	if done {
		return g, err2
	}

	return db, nil
}

func initMigrations(sqlDB *sql.DB, logger *logrus.Logger) (*gorm.DB, error, bool) {
	driver, err := migrateps.WithInstance(sqlDB, &migrateps.Config{})
	if err != nil {
		logger.Fatal("Error initializing migration driver:", err)
		return nil, err, true
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		logger.Fatal("Error creating migration instance:", err)
		return nil, err, true
	}

	version, _, err := m.Version()
	if err != nil {
		logger.Fatal("Error checking migration version:", err)
		return nil, err, true
	}

	if version == 0 {
		logger.Info("⚡ Applying migrations for the first time")
	} else {
		logger.Info("⚡ Database already migrated (version:", version, ")")
	}

	if version == 1 {
		logger.Warn("⚠️ Database is in dirty state, forcing migration to version 1")
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("Migration failed:", err)
		return nil, err, true
	}
	logger.Info("✅ Migrations applied successfully")
	return nil, nil, false
}

func initRepositories(db *gorm.DB) map[string]interface{} {
	return map[string]interface{}{
		"user":     repository.NewUserRepository(db),
		"follower": repository.NewFollowerRelationRepository(db),
	}
}

func (b *Container) GetRepository(name string) (interface{}, error) {
	repo, exists := b.Repositories[name]
	if !exists {
		return nil, fmt.Errorf("repository %s not found", name)
	}
	return repo, nil
}
