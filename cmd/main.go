package main

import (
	"fmt"
	"github.com/garv2003/go-load-balancer/internals/algo"
	"github.com/garv2003/go-load-balancer/internals/models"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ServerManager interface {
	GetServer(servers []*models.Server) (url.URL, error)
}

type UpdateServerStats interface {
	IncrementConnection(serverUrl url.URL)
	DecrementConnection(serverUrl url.URL)
}

type UpdateAvgTime interface {
	UpdateServerAvgTime(serverUrl url.URL, t time.Duration)
}

var serverPool models.ServerPool

func GetServerManger(strategies string) ServerManager {
	switch strategies {
	case "roundRobin":
		return &algo.RoundRobin{}
	case "weightedRoundRobin":
		return &algo.WeightedRoundRobin{}
	case "leastConnection":
		return &algo.LeastConnection{}
	case "weightedLeastConnection":
		return &algo.WeightedLeastConnection{}
	case "leastResponseTime":
		return &algo.LeastResponseTime{}
	case "weightedResponseTime:":
		return &algo.WeightedResponseTime{}
	case "IpHash":
		return &algo.IpHash{}
	case "random":
		return &algo.Random{}
	default:
		return nil
	}
}

func healthCheck() {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

func main() {
	serverPool.AddServer("http://localhost:7891")
	serverPool.AddServer("http://localhost:7892")
	serverPool.AddServer("http://localhost:7893")
	serverPool.AddServer("http://localhost:7894")

	mux := http.NewServeMux()
	serverManager := GetServerManger("roundRobin")

	go healthCheck()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		serverUrl, err := serverManager.GetServer(serverPool.Servers)

		if err != nil {
			fmt.Println(err)
			return
		}

		update, ok := serverManager.(UpdateServerStats)

		if !ok {
			update.IncrementConnection(serverUrl)
			defer update.DecrementConnection(serverUrl)
		}

		updateTime, err1 := serverManager.(UpdateAvgTime)

		if !err1 {
			updateTime.UpdateServerAvgTime(serverUrl, time.Since(now))
		}

		proxyHttp := httputil.NewSingleHostReverseProxy(&serverUrl)
		proxyHttp.ServeHTTP(writer, request)
	})

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		return
	}
}
