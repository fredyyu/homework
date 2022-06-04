package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"homework/conf"
	"homework/database"
	"homework/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	conf.InitConfig()
	database.MongoInit()

}

func main() {

	routers := routers.InitRouter()
	readTimeOut := time.Duration(60) * time.Second
	writeTimeOut := time.Duration(60) * time.Second
	port := fmt.Sprintf(":%s", viper.GetString("SERVER_PORT"))
	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Addr:           port,
		Handler:        routers,
		ReadTimeout:    readTimeOut,
		WriteTimeout:   writeTimeOut,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s\n", port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server Closing in 5 sec")
	select {
	case <-ctx.Done():
		log.Println("Server Closed")
	}
}
