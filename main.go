package main

import (
	"bda/api"
	crawlerCollector "bda/crawler_collector"
	pingerCollector "bda/pinger_collector"
	"bda/types"
	"flag"
	"net"
)

func main_() {
	seedIpRaw := flag.String("seedip", "", "Seed node IP to begin crawl")
	seedPortRaw := flag.Uint("seedport", 99999, "Seed node PORT to begin crawl")
	flag.Parse()

	if seedIpRaw == nil || *seedIpRaw == "" {
		panic("Invalid seed node IP")
	}
	seedIp := net.ParseIP(*seedIpRaw)
	if seedIp == nil {
		panic("Invalid seed node IP")
	}

	if seedPortRaw == nil || *seedPortRaw > 65535 {
		panic("Invalid seed node IP")
	}
	seedPort := uint16(*seedPortRaw)
	addrsees := make(chan types.AddrChanMsg)
	versions := make(chan types.VersionChanMsg)

	crawlerCollector.Collect(addrsees, versions, seedIp, seedPort)
}

func main__() {
	pingerCollector.Collect()
}

func main() {
	api.Start()
}
