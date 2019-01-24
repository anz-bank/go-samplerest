package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.service.anz/go/samplerest/pkg/pet"

	"github.com/go-chi/chi"
	mw "github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	storeimpl = kingpin.Flag("datastore", "Storage used, one of {mem, pq}").Short('d').Default("mem").Enum("mem", "pq")
	port      = kingpin.Flag("port", "port").Short('p').Default("4852").Int()
)

func createStore() (pet.Storer, error) {
	switch *storeimpl {
	case "mem":
		return &pet.MemStore{}, nil
	case "pq":
		return nil, errors.New("postgres store not yet implemented")
	}
	return nil, errors.New("Unknown store implementation, must be either 'mem' or 'pq'")
}

func main() {
	kingpin.Parse()
	store, err := createStore()
	if err != nil {
		log.Fatalf("Could not connect data storage. %v", err)
	}
	router := chi.NewRouter()
	router.Use(mw.Logger)
	service := pet.NewPetService(store)
	pet.SetupRoutes(router, service)
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", *port),
	}
	log.Infoln("Server listening on port", *port)
	log.Fatal(server.ListenAndServe())
}
