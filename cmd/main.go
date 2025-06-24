package main

import (
	"context"
	"fmt"
	"github.com/garv2003/go-load-balancer/internals/algo"
	"github.com/garv2003/go-load-balancer/internals/config"
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

func healthCheck(healthInterval int) {
	t := time.NewTicker(time.Second * time.Duration(healthInterval))
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
	configPath := "config.yaml"
	cfg := &config.Config{}
	if err := cfg.ReloadFromFile(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	for _, server := range cfg.SafeGet().Servers {
		serverPool.AddServer(server)
	}

	go config.WatchConfig(configPath, cfg, &serverPool)

	for _, v := range cfg.Servers {
		serverPool.AddServer(v)
	}

	mux := http.NewServeMux()
	serverManager := GetServerManger(cfg.Strategy)

	go healthCheck(cfg.HealthCheckInterval)

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

	err := http.ListenAndServe(":"+cfg.Port, mux)

	if err != nil {
		return
	}
}
