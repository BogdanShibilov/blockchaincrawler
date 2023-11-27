package crawler

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Protocol string

const (
	WSS   Protocol = "wss"
	HTTPS Protocol = "https"
)

type gethConnection struct {
	protocol Protocol
	url      string
}

func NewGethConnection(protocol Protocol, url string) *gethConnection {
	return &gethConnection{
		protocol: protocol,
		url:      url,
	}
}

func (c *gethConnection) getRawUrl() string {
	return string(c.protocol) + "://" + c.url
}

func (c *gethConnection) Dial() (*ethclient.Client, error) {
	return ethclient.Dial(c.getRawUrl())
}

func (c *gethConnection) DialContext(ctx context.Context) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, c.getRawUrl())
}
