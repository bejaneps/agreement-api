package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bejaneps/docs-agreement-api/backend/routers"
	"github.com/gin-gonic/gin"
)

var (
	logger   *os.File
	recovery *os.File
)

var port = flag.String("port", ":5050", "-port :portnumber")

func init() {
	var err error

	logger, err = os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}

	recovery, err = os.Create("recovery.txt")
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode("release")
}

func main() {
	flag.Parse()

	router := gin.New()
	router.Use(gin.LoggerWithWriter(logger), gin.RecoveryWithWriter(recovery))

	router.POST("/document/create", routers.DocCreateHandler)
	router.POST("/document/perm", routers.DocPermHandler)
	router.POST("/document/sign", routers.DocSignHandler)
	router.POST("/document/list", routers.DocListHandler)

	var server = &http.Server{
		Addr:    *port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()

	logger.Close()
	recovery.Close()
}
