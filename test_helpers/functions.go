package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	goutils "github.com/axdbertuol/goutils/functions"
	"gorm.io/gorm"

	"github.com/spf13/viper"
)

type E2ESuitCase struct {
	PgContainer *PgContainer
	Ctx         context.Context
	Config      *viper.Viper
	ScriptPath  string
	DB          *gorm.DB
}

func (esc *E2ESuitCase) MustInitEnvironment() {
	var (
		dsn    string
		config = esc.Config
		ctx    = esc.Ctx
	)

	_, ok := os.LookupEnv("CI")
	if ok {
		dbConnStr, err := goutils.GetConnection(config)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to get connection: %v", err))
		}
		dsn = *dbConnStr
	} else {
		postgresContainer := NewPgContainer(config, esc.ScriptPath)
		if err := postgresContainer.StartPostgresContainer(ctx); err != nil {
			log.Fatal(err.Error())
		}
		esc.PgContainer = postgresContainer
		dsn = postgresContainer.DSN
		// Clean up the container

	}
	db, err := goutils.ConnectToDb(
		dsn,
		goutils.OpenDatabaseConnection,
		0,
		1,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	esc.DB = db
}

func (esc *E2ESuitCase) Setup(modelsList ...interface{}) error {

	// Run migrations
	if modelsList == nil {
		return errors.New("modelsList is nil")
	}
	if err := esc.DB.AutoMigrate(modelsList...); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	return nil
}

func (esc *E2ESuitCase) CreateEntities(entities ...interface{}) error {
	for _, entity := range entities {
		if err := esc.DB.Create(entity).Error; err != nil {
			return err
		}
	}
	return nil
}

func (esc *E2ESuitCase) Cleanup(targets ...interface{}) error {

	if targets == nil || len(targets) == 0 {
		return errors.New("no targets to clean")
	}
	for _, target := range targets {
		if err := esc.DB.Migrator().DropTable(target); err != nil {
			log.Fatalf("failed to drop table " + err.Error())
		}
	}
	return nil
}
