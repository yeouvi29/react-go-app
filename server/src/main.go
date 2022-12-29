package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "math/rand"
    "net/http"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Quote struct {
    ObjectId    primitive.ObjectID `bson:"_id"`
    Id          int
    Quote       string
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

//
// Get root (/)
// 
func getRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GET /")
    io.WriteString(w, "OK\n")
}

//
// Get a random Quote
//
func getQuote(w http.ResponseWriter, r *http.Request) {

    fmt.Println("GET /quote")

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
    result := results[rand.Int() % len(results)]
    cursor.Decode(&result)

    resp, err := json.Marshal(result)
    if err != nil {
        fmt.Printf("Failed to convert quote to json: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, "Failed to convert quote to json")
        return
    }

    io.WriteString(w, string(resp))
}
