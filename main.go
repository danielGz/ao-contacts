package main

import (
	"accelone-contacts/api"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// command-line arguments
	var port int
	flag.IntVar(&port, "port", 8000, "HTTP server port") // Default to 8000
	// Parse the flags
	flag.Parse()
	// using gorilla router due to its method based routing capability, subroutes and variables in route-patterns, good for RESTful solutions
	router := mux.NewRouter()
	router.Use(api.JsonContentTypeMiddleware)
	// separated each API into its registrar to keep the endpoint list small and prevent mix-up
	api.RegisterCreateContactsAPI(router)

	log.Info().Msgf("Listening on port %d", port)

	addr := fmt.Sprintf(":%d", port) // Format the api address
	log.Fatal().Err(http.ListenAndServe(addr, router)).Msg("terminated")
}
