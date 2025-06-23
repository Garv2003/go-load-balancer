package cmd

import (
	"fmt"
	"github.com/garv2003/go-load-balancer/internals/algo"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServerManager interface {
	GetServer(servers []*models.Server) (url.URL, error)
}

type UpdateServerStats interface {
	IncrementConnection(serverUrl url.URL)
	DecrementConnection(serverUrl url.URL)
}

type UpdateAvgTime interface {
	UpdateServerAvgTime(serverUrl url.URL)
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

func main() {
	serverPool.AddServer("http://localhost:7891")
	serverPool.AddServer("http://localhost:7892")
	serverPool.AddServer("http://localhost:7893")
	serverPool.AddServer("http://localhost:7894")

	mux := http.NewServeMux()
	serverManager := GetServerManger("roundRobin")

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		serverUrl, err := serverManager.GetServer(serverPool.Servers)

		if err != nil {
			fmt.Println(err)
			return
		}

		proxyHttp := httputil.NewSingleHostReverseProxy(&serverUrl)
		proxyHttp.ServeHTTP(writer, request)
	})

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		return
	}
}
