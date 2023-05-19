package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ignoxx/podara/poc3/api"
	"github.com/ignoxx/podara/poc3/storage"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address")
	dbFile := flag.String("dbFile", "podara.db", "database file")
	flag.Parse()

	store := storage.NewSqliteStorage(*dbFile)
	server := api.NewServer(*listenAddr, store)
	fmt.Println("Server started on port", *listenAddr)
	log.Fatal(server.Start())
}