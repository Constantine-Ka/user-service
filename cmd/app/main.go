package main

import (
	"context"
	"fmt"
	"github.com/Constantine-Ka/user-service/model"
	"github.com/Constantine-Ka/user-service/pkg/handler"
	"github.com/Constantine-Ka/user-service/pkg/repository"
	"github.com/Constantine-Ka/user-service/pkg/service"
	"github.com/Constantine-Ka/user-service/tools/logging"
	"github.com/gookit/ini/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logging.GetLogger()
	fmt.Println("start")
	errEnv := godotenv.Load()
	if errEnv != nil {
		logger.Info(errEnv)
	}
	ini.LoadFiles("configs/config.ini")
	logger.Info(ini.Int("port"))
	dbConfig := ini.StringMap("db")
	//db, errDB := repository.NewPostgresDB(repository.Config{
	//	Host:     viper.GetString("db.host1"),
	//	Port:     viper.GetString("db.port1"),
	//	Username: viper.GetString("db.username"),
	//	Password: os.Getenv("DB_PASSWORD"),
	//	DBName:   viper.GetString("db.dbname"),
	//	SSLMode:  viper.GetString("db.sslmode"),
	//})
	db, errDB := repository.NewPostgresDB(repository.Config{
		Host:     dbConfig["host1"],
		Port:     dbConfig["port1"],
		Username: dbConfig["username"],
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   dbConfig["dbname"],
		SSLMode:  dbConfig["sslmode"],
	})
	if errDB != nil {
		logrus.Fatalf("failed to initialize db:%s", errDB.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(model.Server)
	go func() {
		if err := srv.Run(ini.String("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server:%s", err.Error())
		}
	}()
	logrus.Print("github.com/Constantine-Ka/user-service App is was Started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("github.com/Constantine-Ka/user-service App is Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down:%s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db-connection shutting down:%s", err.Error())
	}
}
