package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/benjamin-wright/auth-server/internal/api/request"
)

var ErrUserExists = errors.New("user already exists")

type Client struct {
	url string
}

func New(URL string) *Client {
	return &Client{
		url: URL,
	}
}

type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

type AddUserResponse struct {
	ID string `json:"id"`
}

func (c *Client) AddUser(ctx context.Context, username string, password string, admin bool) (*AddUserResponse, error) {
	var response AddUserResponse
	status, err := request.Post(ctx, c.url, AddUserRequest{
		Username: username,
		Password: password,
		Admin:    admin,
	}, &response)
	if err != nil {
		return nil, err
	}

	if status == http.StatusConflict {
		return nil, ErrUserExists
	}

	if status != http.StatusCreated {
		return nil, fmt.Errorf("failed with status code %d", status)
	}

	return &response, nil
}

type GetUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

func (c *Client) GetUser(ctx context.Context, username string) (*GetUserResponse, error) {
	var response GetUserResponse
	status, err := request.Get(ctx, c.url+"/"+url.PathEscape(username), &response)
	if err != nil {
		return nil, err
	}

	if status == http.StatusNotFound {
		return nil, nil
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("failed with status code %d", status)
	}

	return &response, nil
}

type ListUsersResponse struct {
	Users []GetUserResponse `json:"users"`
}

func (c *Client) ListUsers(ctx context.Context) (*ListUsersResponse, error) {
	var response ListUsersResponse
	status, err := request.Get(ctx, c.url+"/", &response)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("failed with status code %d", status)
	}

	return &response, nil
}

type CheckPasswordRequest struct {
	Password string `json:"password"`
}

type CheckPasswordResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

func (c *Client) CheckPassword(ctx context.Context, username string, password string) (*CheckPasswordResponse, bool, error) {
	var response CheckPasswordResponse
	status, err := request.Put(
		ctx,
		fmt.Sprintf("%s/%s/password", c.url, url.PathEscape(username)),
		CheckPasswordRequest{
			Password: password,
		},
		&response,
	)

	if status == http.StatusUnauthorized {
		return nil, false, nil
	}

	if status > 299 {
		return nil, false, fmt.Errorf("failed with status %d", status)
	}

	if err != nil {
		return nil, false, err
	}

	return &response, true, nil
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	status, err := request.Delete(ctx, fmt.Sprintf("%s/id/%s", c.url, id))
	if err != nil {
		return err
	}

	if status != http.StatusNoContent {
		return fmt.Errorf("failed with status code %d", status)
	}

	return nil
}
