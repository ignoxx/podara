package main

import (
	"flag"
	"os"

	"github.com/ignoxx/podara/poc3/api"
	"github.com/ignoxx/podara/poc3/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address")
	baseStorageDir := flag.String("baseStorageDir", "./.db", "directory for storing files like images or the sqlite database")
	dbFileName := flag.String("dbFile", "podara.db", "database file name")
	flag.Parse()

	imageDir := *baseStorageDir + "/images"
	audioDir := *baseStorageDir + "/audio"
	dbDir := *baseStorageDir + "/" + *dbFileName

	// Create the directories if they don't exist
	os.MkdirAll(imageDir, os.ModePerm)
	os.MkdirAll(audioDir, os.ModePerm)

	store := storage.NewSqliteStorage(dbDir)
	server := api.NewServer(*listenAddr, store, imageDir, audioDir)

	log.WithFields(log.Fields{
		"listenAddr":     *listenAddr,
		"baseStorageDir": *baseStorageDir,
		"dbFileName":     *dbFileName,
		"imageDir":       imageDir,
		"audioDir":       audioDir,
		"dbDir":          dbDir,
	}).Info("Starting server")

	log.Fatal(server.Start())
}
