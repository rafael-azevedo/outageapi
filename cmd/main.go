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

	fmt.Println(`
	 .d8888b.  888                      888                   
	d88P  Y88b 888                      888                   
	Y88b.      888                      888                   
	 "Y888b.   888888  8888b.  88888b.  888  .d88b.  .d8888b  
	    "Y88b. 888        "88b 888 "88b 888 d8P  Y8b 88K      
	      "888 888    .d888888 888  888 888 88888888 "Y8888b. 
	Y88b  d88P Y88b.  888  888 888 d88P 888 Y8b.          X88 
	 "Y8888P"   "Y888 "Y888888 88888P"  888  "Y8888   88888P' 
	                           888                            
	                           888                            
	                           888   
    
	  01010011 01110100 01100001 01110000 01101100 01100101 
	  01110011 00100000 01001111 01110101 01110100 01100001 
	  01100111 01100101 00100000 01000001 01010000 01001001 

	`)
	log.Printf("Outage Api Started at http://%s\n", srv.Addr)
	fmt.Println(router.GetActiveRoutes(srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
