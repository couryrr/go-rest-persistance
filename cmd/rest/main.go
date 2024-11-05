package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/couryrr/go-rest-persistance/internal"
	_ "github.com/mattn/go-sqlite3"
)

var handler *internal.UserHandler

func init() {
    db, err := sql.Open("sqlite3", "users.db")

    if err != nil {
        log.Fatal("database failed to open")
    }
    
	repo := internal.NewSQLiteRepository(db)
	err = repo.Migrate()
	if err != nil {
		log.Fatal("unable to run database migration: ", err)
	}
    handler = internal.NewUserHandler(repo)
    
    if err != nil {
        log.Fatalf("article_handler failed to create: %s", err) 
    }
}

func main(){
	router := http.NewServeMux()
	router.Handle("/users/", http.StripPrefix("/users", handler.GetHandler()))

	server := &http.Server{
		Addr: ":8080",
		Handler: http.TimeoutHandler(router, 1*time.Second, "Timed out!\n"),
	}

	go func(){
		log.Printf("starting server on: %s \n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed){
				log.Fatalf("HTTP server error: %v", err)
			}
		}
		log.Println("stopping server...")
	}()

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)
	<- sChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	log.Println("called shutdown with ctx")
	err := server.Shutdown(shutdownCtx)
	log.Println("shutdown done check for error")

	if err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}

	log.Println("gracefully shutdown complete...")
}

