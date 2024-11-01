package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/couryrr/go-rest-persistance/internal"
)

func main(){
	router := http.NewServeMux()

	personHandler := internal.NewPersonHandler()
	path, fn := personHandler.HandleAddPerson()
	router.HandleFunc(path, fn)

	server := &http.Server{
		Addr: ":8080",
		Handler: router,
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
	} else {
		log.Println("if I am here shutdown has returned correct?")
	}

	log.Println("gracefully shutdown complete...")
}

