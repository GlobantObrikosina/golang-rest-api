package main

import (
	"context"
	"fmt"
	"github.com/GlobantObrikosina/golang-rest-api/db"
	"github.com/GlobantObrikosina/golang-rest-api/handler"
	"github.com/GlobantObrikosina/golang-rest-api/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	database := db.NewDatabase(dbUser, dbPassword, dbName)
	services := service.NewService(database)
	httpHandler := handler.NewHandler(services)

	defer database.Close()
	server := &http.Server{Handler: httpHandler.InitRoutes()}
	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Printf("Server some how didn't start")
		}
	}()
	defer Stop(server)
	log.Printf("Started server on %s", addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
