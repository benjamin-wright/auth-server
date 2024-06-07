package sut

import (
	"context"
	"fmt"
	"time"

	rds "github.com/benjamin-wright/db-operator/v2/pkg/redis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Client struct {
	redis *redis.Client
}

func New() (*Client, error) {
	cfg, err := rds.ConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("error getting redis config: %v", err)
	}

	client, err := rds.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %v", err)
	}

	return &Client{
		redis: client,
	}, nil
}

func (c *Client) Get() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate sut: %v", err)
	}

	cmd := c.redis.Set(context.Background(), uuid.String(), true, 5*time.Minute)
	if cmd.Err() != nil {
		return "", fmt.Errorf("failed to set sut: %v", cmd.Err())
	}

	return uuid.String(), nil
}

func (c *Client) Check(value string) (bool, error) {
	cmd := c.redis.Get(context.Background(), value)
	if cmd.Err() != nil {
		return false, fmt.Errorf("failed to get value: %v", cmd.Err())
	}

	if err := c.delete(value); err != nil {
		return false, fmt.Errorf("failed to delete sut: %v", err)
	}

	return cmd.Val() != "", nil
}

func (c *Client) delete(sut string) error {
	cmd := c.redis.Del(context.Background(), sut)
	if cmd.Err() != nil {
		return fmt.Errorf("failed to delete sut: %v", cmd.Err())
	}

	return nil
}
