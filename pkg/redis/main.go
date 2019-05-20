package redis

import (
	"fmt"

	"github.com/mediocregopher/radix/v3"
)

const (
	busVehicleDataChannel = "busVehicleDataChannel"
)

// Client represents layer between writer and Redis
type Client struct {
	pool *radix.Pool
	conn radix.Conn
	ps   radix.PubSubConn
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
// and sends it to PubSub subscribers.
// Returns the number of clients that received the message (and error)
func (c *Client) SaveBusVehicleDataMessage(vehicleID string, message []byte) (int, error) {
	err := c.pool.Do(
		radix.Cmd(nil, "SET", fmt.Sprintf("busVehicleData:%s", vehicleID), string(message)),
	)
	if err != nil {
		return 0, err
	}

	var reply int
	err = c.pool.Do(
		radix.Cmd(&reply, "PUBLISH", busVehicleDataChannel, string(message)),
	)

	return reply, err
}
