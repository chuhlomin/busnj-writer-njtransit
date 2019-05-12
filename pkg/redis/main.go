package redis

import (
	"fmt"

	"github.com/mediocregopher/radix/v3"
)

// Client represents layer between writer and Redis
type Client struct {
	pool *radix.Pool
}

// NewClient creates new Client
func NewClient(network string, addr string, size int) (*Client, error) {
	pool, err := radix.NewPool(network, addr, size)
	if err != nil {
		return nil, err
	}

	return &Client{
		pool: pool,
	}, nil
}

// SaveBusVehicleDataMessage saves individual BusVehicleData message
func (c *Client) SaveBusVehicleDataMessage(vehicleID string, message []byte) error {
	err := c.pool.Do(
		radix.Cmd(nil, "SET", fmt.Sprintf("busVehicleData:%s", vehicleID), string(message)),
	)
	return err
}
