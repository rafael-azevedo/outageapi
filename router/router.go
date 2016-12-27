package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafael-azevedo/outageapi/logger"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"OutageList",
		"GET",
		"/outage",
		OutageList,
	},
	Route{
		"AssignNode",
		"POST",
		"/assign",
		AssignNode,
	},
	Route{
		"DeassignNode",
		"POST",
		"/deassign",
		DeassignNode,
	},
	Route{
		"GetRequest",
		"Get",
		"/status/{id}",
		GetRequest,
	},
}

func BuildRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func GetActiveRoutes(address string) []string {
	var routeString []string
	routeString = append(routeString, fmt.Sprintf("%s\n", "Active Endpoints :"))
	for _, route := range routes {
		rs := fmt.Sprintf("%s%s : %s\n", address, route.Pattern, route.Method)
		routeString = append(routeString, rs)
	}
	return routeString
}
