package rest

import (
	"github.com/HagerHu/go-json-rest/trie"
	"net/url"
	"strings"
)

// TODO
// support for #param placeholder ?

type router struct {
	routes                 []Route
	disableTrieCompression bool
	index                  map[*Route]int
	trie                   *trie.Trie
}

// This validates the Routes and prepares the Trie data structure.
// It must be called once the Routes are defined and before trying to find Routes.
// The order matters, if multiple Routes match, the first defined will be used.
func (self *router) start() error {

	self.trie = trie.New()
	self.index = map[*Route]int{}

	for i, _ := range self.routes {

		// pointer to the Route
		route := &self.routes[i]

		// insert in the Trie
		err := self.trie.AddRoute(
			strings.ToUpper(route.HttpMethod),
			route.PathExp,
			route,
		)
		if err != nil {
			return err
		}

		// index
		self.index[route] = i
	}

	if self.disableTrieCompression == false {
		self.trie.Compress()
	}

	// TODO validation of the PathExp ? start with a /
	// TODO url encoding

	return nil
}

// return the result that has the route defined the earliest
func (self *router) ofFirstDefinedRoute(matches []*trie.Match) *trie.Match {
	minIndex := -1
	matchesByIndex := map[int]*trie.Match{}

	for _, result := range matches {
		route := result.Route.(*Route)
		routeIndex := self.index[route]
		matchesByIndex[routeIndex] = result
		if minIndex == -1 || routeIndex < minIndex {
			minIndex = routeIndex
		}
	}

	return matchesByIndex[minIndex]
}

// Return the first matching Route and the corresponding parameters for a given URL object.
func (self *router) findRouteFromURL(httpMethod string, urlObj *url.URL) (*Route, map[string]string, bool) {

	// lookup the routes in the Trie
	// TODO verify url encoding
	matches, pathMatched := self.trie.FindRoutesAndPathMatched(
		strings.ToUpper(httpMethod),
		urlObj.Path,
	)

	// short cuts
	if len(matches) == 0 {
		// no route found
		return nil, nil, pathMatched
	}

	if len(matches) == 1 {
		// one route found
		return matches[0].Route.(*Route), matches[0].Params, pathMatched
	}

	// multiple routes found, pick the first defined
	result := self.ofFirstDefinedRoute(matches)
	return result.Route.(*Route), result.Params, pathMatched
}

// Parse the url string (complete or just the path) and return the first matching Route and the corresponding parameters.
func (self *router) findRoute(httpMethod, urlStr string) (*Route, map[string]string, bool, error) {

	// parse the url
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, false, err
	}

	route, params, pathMatched := self.findRouteFromURL(httpMethod, urlObj)
	return route, params, pathMatched, nil
}
