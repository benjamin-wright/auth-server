package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/benjamin-wright/auth-server/internal/api/request"
)

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{url: url}
}

type GetLoginTokenResponse struct {
	Token  string `json:"token"`
	MaxAge int    `json:"maxAge"`
}

func (c *Client) GetLoginToken(subject string) (*GetLoginTokenResponse, error) {
	var response GetLoginTokenResponse
	status, err := request.Get(context.TODO(), c.url+"/"+subject+"/login", &response)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", status)
	}

	return &response, nil
}

type GetAdminTokenResponse struct {
	Token  string `json:"token"`
	MaxAge int    `json:"maxAge"`
}

func (c *Client) GetAdminToken(subject string) (*GetAdminTokenResponse, error) {
	var response GetAdminTokenResponse
	status, err := request.Get(context.TODO(), c.url+"/"+subject+"/admin", &response)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", status)
	}

	return &response, nil
}
