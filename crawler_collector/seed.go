package crawler_collector

import (
	"bda/types"
	"context"
	"crypto/rand"
	"math/big"
	"net"
)

func PickRandomDnsSeed() (net.IP, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort("8.8.8.8", "53"))
		},
	}

	ips, err := resolver.LookupIP(context.Background(), "ip", types.MainnetDnsSeed)
	if err != nil {
		return net.IP{}, err
	}

	iBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(ips))))
	if err != nil {
		return net.IP{}, err
	}

	return ips[iBig.Int64()], nil
}
