package main

import (
	"context"
	"fmt"
	"github.com/garv2003/go-load-balancer/internals/algo"
	"github.com/garv2003/go-load-balancer/internals/models"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ServerManager interface {
	GetServer(ctx context.Context, servers []*models.Server) (url.URL, error)
}

type UpdateServerStats interface {
	DecrementConnection(servers []*models.Server, server url.URL)
}

type UpdateAvgTime interface {
	UpdateServerAvgTime(servers []*models.Server, serverUrl url.URL, t time.Duration)
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
	case "weightedResponseTime":
		return &algo.WeightedResponseTime{}
	case "IpHash":
		return &algo.IpHash{}
	case "random":
		return algo.NewRandom()
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

		clientIP := request.RemoteAddr
		if host, _, err := net.SplitHostPort(clientIP); err == nil {
			clientIP = host
		}

		ctx := context.WithValue(request.Context(), "client-ip", clientIP)

		serverUrl, err := serverManager.GetServer(ctx, serverPool.Servers)

		if err != nil {
			fmt.Println(err)
			return
		}

		update, ok := serverManager.(UpdateServerStats)

		if !ok {
			defer update.DecrementConnection(serverPool.Servers, serverUrl)
		}

		updateTime, err1 := serverManager.(UpdateAvgTime)

		if !err1 {
			updateTime.UpdateServerAvgTime(serverPool.Servers, serverUrl, time.Since(now))
		}

		proxyHttp := httputil.NewSingleHostReverseProxy(&serverUrl)
		proxyHttp.ServeHTTP(writer, request)
	})

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		return
	}
}
