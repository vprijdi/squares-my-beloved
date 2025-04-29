package main

import (
	"log"

	"github.com/muzhiknastya/squares-my-beloved/internal/env"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
