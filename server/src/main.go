package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Quote struct {
	ObjectId primitive.ObjectID `bson:"_id"`
	Id       int                `json:"id"`
	Quote    string             `json:"quote"`
}
type Client struct {
	client *redis.Client
}

func main() {

	port := os.Getenv("PORT")
	fmt.Println("Starting HTTP server on PORT", port)

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/quote", getQuote)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Failed.")
	} else {
		fmt.Println("Success.")
	}
}

func (q Quote) MarshalBinary() ([]byte, error) {
	return []byte(fmt.Sprintf("%v-%v", q.Id, q.Quote)), nil
}

// Helper function to enable cors
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Get root (/)
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /")
	enableCors(&w)
	io.WriteString(w, "OK\n")
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

// Get a random Quote
func getQuote(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GET /quote")
	enableCors(&w)

	redis, err := newRedis()
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	}
	fmt.Println("Try to get quotes from cache")
	val := redis.getQuotes()
	if val != nil {
		fmt.Println("Get quote from Cache")
		rand.Seed(time.Now().Unix())
		result := val[rand.Int()%len(val)]
		// cursor.Decode(&result)

		resp, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("Failed to convert quote to json: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Failed to convert quote to json")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(resp))

		return
	}
	// step 1: connect to the database
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "You must set your 'MONGODB_URI' environmental variable.")
		return
	}

	fmt.Printf("Connecting to database at: %s... ", uri)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Failed to connect to the database")
		return
	}
	fmt.Println("Success")

	// step 2: fetch the quotes collection
	fmt.Printf("Fetching quotes... ")

	db := client.Database("quotesdb")
	coll := db.Collection("quotes")

	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Printf("Failed to open cursor: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Failed to open cursor")
		return
	}

	var results []Quote
	if err = cursor.All(context.TODO(), &results); err != nil {
		fmt.Printf("Failed to fetch quotes: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Failed to fetch quotes")
		return
	}
	fmt.Println("Done. Got", len(results), "quotes")

	// step 3: return a random quote
	rand.Seed(time.Now().Unix())
	result := results[rand.Int()%len(results)]
	fmt.Println("Set quotes to redis")
	redis.setQuotes(results)
	cursor.Decode(&result)

	resp, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Failed to convert quote to json: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Failed to convert quote to json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(resp))
}
