package main

import (
	"flag"

	"github.com/ignoxx/podara/poc3/api"
	"github.com/ignoxx/podara/poc3/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address")
	dbFile := flag.String("dbFile", "podara.db", "database file")
	flag.Parse()

	store := storage.NewSqliteStorage(*dbFile)
	server := api.NewServer(*listenAddr, store)

	log.WithFields(log.Fields{
		"listenAddr": *listenAddr,
		"dbFile":     *dbFile,
	}).Info("Starting server")

	log.Fatal(server.Start())
}
