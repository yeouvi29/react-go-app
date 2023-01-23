package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type Client struct {
	client *redis.Client
}

func newRedis() (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        "redis:6379",
		DB:          0, // use default DB
		DialTimeout: 100 * time.Millisecond,
		ReadTimeout: 100 * time.Millisecond,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) getQuotes() (quotes []Quote) {
	val, err := c.client.Get("quotes").Result()
	if err != nil {
		return nil
	}

	resp := []Quote{}
	err = json.Unmarshal([]byte(val), &resp)
	if err != nil {
		log.Fatal(err)
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	// Publish using Redis PubSub
	if err := c.client.Publish("send-user-name", payload).Err(); err != nil {
		log.Fatal(err)
	}

	return resp
}

func (c *Client) setQuotes(quotes []Quote) {
	json, err := json.Marshal(quotes)
	if err != nil {
		log.Fatal(err)
	}

	c.client.Set("quotes", json, 20*time.Second)
}
