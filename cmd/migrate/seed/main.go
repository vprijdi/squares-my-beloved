package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/muzhiknastya/squares-my-beloved/internal/db"
	"github.com/muzhiknastya/squares-my-beloved/internal/env"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

func main() {
	godotenv.Load(".env")

	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/squares?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store)
}
