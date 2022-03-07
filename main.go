package main

import (
	"bda/api"
	crawlerCollector "bda/crawler_collector"
	pingerCollector "bda/pinger_collector"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "pinger":
			pingerCollector.Start()
		case "crawler":
			crawlerCollector.Start()
		case "api":
			api.Start()
		}
	}
}
