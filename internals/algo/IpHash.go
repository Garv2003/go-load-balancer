package algo

import (
	"errors"
	"fmt"
	"github.com/garv2003/go-load-balancer/internals/models"
	"net/url"
	"sync/atomic"
)

type IpHash struct {
	count atomic.Int64
}

func (rr *IpHash) GetCount() int {
	return int(rr.count.Load())
}

func (rr *IpHash) IncrementCount() {
	rr.count.Add(1)
}

func (rr *IpHash) GetServer(servers []*models.Server) (url.URL, error) {
	if len(servers) == 0 {
		fmt.Println("there is no server in servers list")
		return url.URL{}, errors.New("there is no server in servers list")
	}

	server := servers[rr.GetCount()%len(servers)]
	rr.IncrementCount()

	return server.ServerUrl, nil
}
