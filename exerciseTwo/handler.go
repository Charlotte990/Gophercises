package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler takes a map of path/urls and if the requesting URL path contains the map path, it will redirect to the url in the map
// if no map path matches the request url path, it will deliver the fallback handler (hello() called via defaultMux())
// which displays Hello, World!
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, redirectURL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler takes a yaml of path/url is mapped to the format required by MapHandler, then passed into MapHandler
// Or if the path is not in the yaml, make MapHandler call itself
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	type config []struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}
	var c config
	err := yaml.Unmarshal(yml, &c)

	if err != nil {
		return nil, err
	}

	pathMap := map[string]string{}

	for _, i := range c {
		pathMap[i.Path] = i.URL
	}

	return MapHandler(pathMap, fallback), nil

}
