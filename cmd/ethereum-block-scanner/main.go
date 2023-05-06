package main

import (
	"context"
	"log"
	"os"

	"github.com/powerslider/ethereum-block-scanner/pkg/transport/client/jsonrpc"

	"github.com/gorilla/mux"
	"github.com/powerslider/ethereum-block-scanner/pkg/transport/server"

	"github.com/powerslider/ethereum-block-scanner/pkg/handlers"

	"github.com/joho/godotenv"
	"github.com/powerslider/ethereum-block-scanner/pkg/configs"
)

// @title Ethereum Block Scanner API
// @version 1.0
// @description API for exploring Ethereum blocks.
// @termsOfService http://swagger.io/terms/

// @contact.name Tsvetan Dimitrov
// @contact.email tsvetan.dimitrov23@gmail.com

// @license.name MIT
// @license.url https://www.mit.edu/~amini/LICENSE.md

// @host 0.0.0.0:8080
// @BasePath /
func main() {
	ctx := context.Background()

	var err error

	setEnvironment()

	conf := configs.InitializeConfig()
	client := jsonrpc.NewDefaultClient(conf.EthereumHost)

	router := mux.NewRouter()
	router = handlers.InitializeHandlers(conf, router, client)
	s := server.NewServer(conf, router)

	if err = s.Run(ctx); err != nil {
		log.Fatal("error starting the HTTP server: ", err)
	}
}

func setEnvironment() {
	_, foundHost := os.LookupEnv("SERVER_HOST")
	_, foundPort := os.LookupEnv("SERVER_PORT")

	if !foundHost && !foundPort {
		err := godotenv.Load(".env.dist")
		if err != nil {
			panic(err)
		}
	}
}
