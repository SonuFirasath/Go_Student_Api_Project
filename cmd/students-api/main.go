package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SonuFirasath/Go_Student_Api_Project/internal/config"
	"github.com/SonuFirasath/Go_Student_Api_Project/internal/http/handlers/student"
	"github.com/SonuFirasath/Go_Student_Api_Project/internal/storage/sqlite"
)

func main() {
	//Load config

	cfg := config.MustLoad()

	//Database setup

	storage, err :=sqlite.New(cfg)

	if err != nil{
		log.Fatal(err)
	}

	slog.Info("Storag initialized",slog.String("env",cfg.Env),slog.String("version","1.0.0"),slog.String("storage_path",cfg.StoragePath))

	//Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))

	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))

	// Setup server

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("server started",slog.String("address",cfg.Addr))
	fmt.Printf("Server started %s",cfg.Addr)

	done := make(chan os.Signal,1)

	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe()
		if err!=nil{
			log.Fatal("Failed to start server")
		}
	} ()

	<- done

	slog.Info("shutting down the server...")

	ctx, cancle := context.WithTimeout(context.Background(),5*time.Second)

	defer cancle()

	err =	server.Shutdown(ctx)

	if err!=nil{
		slog.Error("Failed to shutdown server...", slog.String("error",err.Error()))
	}

	slog.Info("server shutdown successfully")

}