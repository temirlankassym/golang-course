package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type Product struct {
	Id          int
	Name        string
	Description string
	Price       float64
}

func Connect(ctx context.Context) (*pgx.Conn, error) {
	config := getConfig()

	db, err := pgx.Connect(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("cannot connect to database")
	}

	return db, nil
}

func ConnectCache(ctx context.Context) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func getConfig() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return fmt.Sprintf("postgresql://%s:%s@localhost:%s?sslmode=disable",
		os.Getenv("DBNAME"), os.Getenv("PASSWORD"), os.Getenv("PORT"))
}

func GetProduct(ctx context.Context, db *pgx.Conn, id string, client *redis.Client) (Product, error) {
	//Retrieving from Redis
	val, err := client.Get(ctx, id).Result()

	product := Product{}

	if errors.Is(err, redis.Nil) {
		//If not found in Redis
		fmt.Println("From Database")

		err := db.QueryRow(ctx, "Select * from products where id = $1", id).Scan(&product.Id, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return product, err
		}

		productJson, err := json.Marshal(product)
		if err != nil {
			fmt.Println(err)
		}

		//Setting TTL to 10 seconds
		err = client.Set(ctx, id, productJson, time.Minute).Err()
		if err != nil {
			return Product{}, fmt.Errorf("Can't Set: %w", err)
		}
	} else {
		//If found in Redis just return
		fmt.Println("From Redis")

		err := json.Unmarshal([]byte(val), &product)
		if err != nil {
			return Product{}, err
		}
	}

	return product, nil
}
