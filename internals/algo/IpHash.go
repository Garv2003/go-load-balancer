package algo

import (
	"context"
	"errors"
	"hash/fnv"
	"net/url"

	"github.com/garv2003/go-load-balancer/internals/models"
)

type IpHash struct{}

func hashIP(ip string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(ip))
	return h.Sum32()
}

func (ip *IpHash) GetServer(ctx context.Context, servers []*models.Server) (url.URL, error) {
	clientIPRaw := ctx.Value("client-ip")
	clientIP, ok := clientIPRaw.(string)
	if !ok || clientIP == "" {
		return url.URL{}, errors.New("client IP not found in context")
	}

	var aliveServers []*models.Server
	for _, s := range servers {
		if s.IsAlive {
			aliveServers = append(aliveServers, s)
		}
	}
	
	if len(aliveServers) == 0 {
		return url.URL{}, errors.New("no alive servers")
	}

	index := int(hashIP(clientIP)) % len(aliveServers)
	return aliveServers[index].ServerUrl, nil
}
