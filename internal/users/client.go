package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/benjamin-wright/db-operator/v2/pkg/postgres/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoUser = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")
var ErrComplexity = errors.New("password didn't meet complexity requirements")

type Client struct {
	conn *pgxpool.Pool
}

func New() (*Client, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed getting config from environment: %+v", err)
	}

	cfg.Retry = false

	log.Info().Interface("config", cfg).Msg("the config")

	conn, err := config.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to postgres: %+v", err)
	}

	if conn == nil {
		return nil, errors.New("failed connecting to postgres: no error")
	}

	return &Client{conn: conn}, nil
}

func (c *Client) DeleteAllUsers() error {
	_, err := c.conn.Exec(context.Background(), `DELETE FROM users`)
	if err != nil {
		return fmt.Errorf("failed to clear existing users; %+v", err)
	}

	return nil
}

func (c *Client) AddUser(user User) (string, error) {
	if !CheckPasswordComplexity(user.Password) {
		return "", ErrComplexity
	}

	log.Info().Interface("user", user).Msg("adding user")

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash: %+v", err)
	}
	hash := string(bytes)

	rows, err := c.conn.Query(context.Background(), `INSERT INTO users("name", "password", "admin") VALUES ($1, $2, $3) RETURNING "id"`, user.Name, hash, user.Admin)
	if err != nil {
		return "", fmt.Errorf("failed to add user to database: %+v", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&user.ID); err != nil {
			return user.ID, fmt.Errorf("failed to parse new user ID: %+v", err)
		}

		return user.ID, nil
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" && pgerr.ConstraintName == "users_name_key" {
			return "", ErrUserExists
		}
		return "", fmt.Errorf("failed to add user to database: %+v", err)
	}

	return "", errors.New("failed to find new user ID")
}

var ErrPasswordMismatch = errors.New("password mismatch")

func (c *Client) CheckPassword(user User) (*User, error) {
	rows, err := c.conn.Query(context.Background(), `SELECT "id", "password", "admin" FROM users WHERE "name" = $1`, user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %+v", err)
	}
	defer rows.Close()

	numUsers := 0
	var passwordHash string

	for rows.Next() {
		numUsers += 1
		if err = rows.Scan(&user.ID, &passwordHash, &user.Admin); err != nil {
			return nil, fmt.Errorf("failed to parse new user ID: %+v", err)
		}
	}

	if numUsers == 0 {
		return nil, ErrNoUser
	}

	if numUsers > 1 {
		return nil, fmt.Errorf("expected 1 user, got %d", numUsers)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	if err != nil {
		return nil, ErrPasswordMismatch
	}

	return &user, nil
}

func (c *Client) GetUser(name string) (*User, error) {
	rows, err := c.conn.Query(context.Background(), `SELECT "id", "admin" FROM users WHERE "name" = $1`, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %+v", err)
	}
	defer rows.Close()

	numUsers := 0
	user := User{
		Name: name,
	}

	for rows.Next() {
		numUsers += 1
		if err = rows.Scan(&user.ID, &user.Admin); err != nil {
			return nil, fmt.Errorf("failed to parse new user ID: %+v", err)
		}
	}

	if numUsers == 0 {
		return nil, ErrNoUser
	}

	if numUsers > 1 {
		return nil, fmt.Errorf("expected 1 user, got %d", numUsers)
	}

	return &user, nil
}

func (c *Client) ListUsers() ([]User, error) {
	rows, err := c.conn.Query(context.Background(), `SELECT "id", "name", "password", "admin" FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %+v", err)
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Password, &user.Admin); err != nil {
			return nil, fmt.Errorf("failed to parse new user ID: %+v", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *Client) DeleteUser(id string) error {
	_, err := c.conn.Exec(context.Background(), `DELETE FROM users WHERE "id" = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user from database: %+v", err)
	}

	return nil
}
