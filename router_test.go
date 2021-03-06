package rest

import (
	"net/url"
	"testing"
)

func TestFindRouteAPI(t *testing.T) {

	r := router{
		routes: []Route{
			Route{
				HttpMethod: "GET",
				PathExp:    "/",
			},
		},
	}

	err := r.start()
	if err != nil {
		t.Fatal()
	}

	// full url string
	input := "http://example.org/"
	route, params, pathMatched, err := r.findRoute("GET", input)
	if err != nil {
		t.Fatal()
	}
	if route.PathExp != "/" {
		t.Error()
	}
	if len(params) != 0 {
		t.Error()
	}
	if pathMatched != true {
		t.Error()
	}

	// part of the url string
	input = "/"
	route, params, pathMatched, err = r.findRoute("GET", input)
	if err != nil {
		t.Fatal()
	}
	if route.PathExp != "/" {
		t.Error()
	}
	if len(params) != 0 {
		t.Error()
	}
	if pathMatched != true {
		t.Error()
	}

	// url object
	urlObj, err := url.Parse("http://example.org/")
	if err != nil {
		t.Fatal()
	}
	route, params, pathMatched = r.findRouteFromURL("GET", urlObj)
	if route.PathExp != "/" {
		t.Error()
	}
	if len(params) != 0 {
		t.Error()
	}
	if pathMatched != true {
		t.Error()
	}
}

func TestNoRoute(t *testing.T) {

	r := router{
		routes: []Route{},
	}

	err := r.start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/notfound"
	route, params, pathMatched, err := r.findRoute("GET", input)
	if err != nil {
		t.Fatal()
	}

	if route != nil {
		t.Error("should not be able to find a route")
	}
	if params != nil {
		t.Error("params must be nil too")
	}
	if pathMatched != false {
		t.Error()
	}
}

func TestDuplicatedRoute(t *testing.T) {

	r := router{
		routes: []Route{
			Route{
				HttpMethod: "GET",
				PathExp:    "/",
			},
			Route{
				HttpMethod: "GET",
				PathExp:    "/",
			},
		},
	}

	err := r.start()
	if err == nil {
		t.Error("expected the duplicated route error")
	}
}

func TestRouteOrder(t *testing.T) {

	r := router{
		routes: []Route{
			Route{
				HttpMethod: "GET",
				PathExp:    "/r/:id",
			},
			Route{
				HttpMethod: "GET",
				PathExp:    "/r/*rest",
			},
		},
	}

	err := r.start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/r/123"
	route, params, pathMatched, err := r.findRoute("GET", input)
	if err != nil {
		t.Fatal()
	}

	if route.PathExp != "/r/:id" {
		t.Errorf("both match, expected the first defined, got %s", route.PathExp)
	}
	if params["id"] != "123" {
		t.Error()
	}
	if pathMatched != true {
		t.Error()
	}
}

func TestSimpleExample(t *testing.T) {

	r := router{
		routes: []Route{
			Route{
				HttpMethod: "GET",
				PathExp:    "/resources/:id",
			},
			Route{
				HttpMethod: "GET",
				PathExp:    "/resources",
			},
		},
	}

	err := r.start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/resources/123"
	route, params, pathMatched, err := r.findRoute("GET", input)
	if err != nil {
		t.Fatal()
	}

	if route.PathExp != "/resources/:id" {
		t.Error()
	}
	if params["id"] != "123" {
		t.Error()
	}
	if pathMatched != true {
		t.Error()
	}
}
