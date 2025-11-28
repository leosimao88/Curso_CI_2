package database

import (
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/guilhermeonrails/api-go-gin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)
func ConectaComBancoDeDados() error {
	// Use explicit DB_* env variables to avoid clobbering system vars (e.g. USER)
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "root"
	}
	if password == "" {
		password = "root"
	}
	if dbname == "" {
		dbname = "root"
	}

	stringDeConexao := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

	// retry logic: try for up to ~90 seconds (30 tries * 3s)
	for i := 0; i < 30; i++ {
		DB, err = gorm.Open(postgres.Open(stringDeConexao))
		if err == nil {
			// try to ping the underlying sql DB
			sqlDB, derr := DB.DB()
			if derr == nil {
				if perr := sqlDB.Ping(); perr != nil {
					err = perr
				} else {
					err = nil
					break
				}
			} else {
				err = derr
			}
		}
		// sleep a bit and try again
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		// If the error is "database does not exist", try to create it using the default 'postgres' DB
		if strings.Contains(strings.ToLower(err.Error()), "does not exist") || strings.Contains(err.Error(), "3D000") {
			// connect to default postgres db
			adminDSN := "host=" + host + " user=" + user + " password=" + password + " dbname=postgres port=" + port + " sslmode=disable"
			adminDB, aerr := gorm.Open(postgres.Open(adminDSN))
			if aerr != nil {
				return fmt.Errorf("failed to connect to admin DB to create database %s: %v", dbname, aerr)
			}
			// create DB
			createQuery := fmt.Sprintf("CREATE DATABASE %s", dbname)
			if execErr := adminDB.Exec(createQuery).Error; execErr != nil {
				return fmt.Errorf("failed to create database %s: %v", dbname, execErr)
			}
			// close adminDB and try again once
			sqlAdmin, derr := adminDB.DB()
			if derr == nil {
				sqlAdmin.Close()
			}
			// try to connect again
			DB, err = gorm.Open(postgres.Open(stringDeConexao))
			if err != nil {
				return fmt.Errorf("failed to connect to newly created database %s: %v", dbname, err)
			}
		} else {
			return fmt.Errorf("erro ao conectar com banco de dados: %v", err)
		}
	}

	if amErr := DB.AutoMigrate(&models.Aluno{}); amErr != nil {
		return fmt.Errorf("failed to run automigrate: %v", amErr)
	}
	return nil
}
