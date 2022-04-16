package main

import (
	"bda/api"
	crawlerCollector "bda/crawler_collector"
	"bda/logger"
	"bda/models"
	pingerCollector "bda/pinger_collector"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func WaitForDb(log logger.Logger) {
	log.Info(logrus.Fields{}, "Waiting for DB")
	for {
		isLive := true
		db, err := models.GetDb()
		if err != nil {
			isLive = false
		}
		rawDb, err := db.DB()
		if err != nil {
			isLive = false
		}
		err = rawDb.Ping()
		if err != nil {
			isLive = false
		}
		if isLive {
			break
		} else {
			log.Info(logrus.Fields{}, "DB unavaiable, waiting")
			time.Sleep(5 * time.Second)
			continue
		}
	}
	log.Info(logrus.Fields{}, "DB active, continuing")
}

func main() {
	log := logger.Logger{
		Prefix: "MAIN",
	}
	WaitForDb(log)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "pinger":
			pingerCollector.Start()
		case "crawler":
			crawlerCollector.Start()
		case "api":
			api.Start()
		}
	} else {
		fmt.Println("Select entrypoint parameter: pinger, crawler or api")
	}
}
