package main

import (
	"context"
	"course/redis/repository"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

var client *redis.Client
var db *pgx.Conn

func main() {
	ctx := context.Background()
	var err error

	db, err = repository.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	client = repository.ConnectCache(ctx)

	defer db.Close(ctx)
	defer client.Close()

	http.HandleFunc("/product", ProductHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ProductHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	product, err := repository.GetProduct(context.Background(), db, id, client)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
	}

	resp, err := json.Marshal(product)
	if err != nil {
		log.Fatal("Can't parse json")
	}

	//returning response to client
	_, err = writer.Write(resp)
	if err != nil {
		log.Fatal("Can't write")
	}
}
