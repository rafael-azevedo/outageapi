package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rafael-azevedo/outageapi/router"
)

var failCount int

func init() {
	//if failCount is higher than 0 panic for each Enviromental Variable that is not present add 1 to failCount
	failCount = 0
	config := os.Getenv("OUTAGECONF")
	logDir := os.Getenv("OUTAGELOGDIR")
	if config == "" {
		failCount++
		log.Println("Enviromental Variable(OUTAGECONF) for the config file location is not set")
	}

	if logDir == "" {
		failCount++
		log.Println("Enviromental Variable(OUTAGELOGDIR) for the config file location is not set")
	}
}
func main() {
	if failCount > 0 {
		log.Fatalln("Enviroment not set correctly")
	}

	routerMux := router.BuildRouter()
	srv := &http.Server{
		Handler:      routerMux,
		Addr:         "127.0.0.1:1234",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Printf("Application Started at http://%s\n", srv.Addr)
	fmt.Println(router.GetActiveRoutes(srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
