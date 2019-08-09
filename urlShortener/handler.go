package main

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL string  `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, val := range pathUrls {
		pathsToUrls[val.Path] = val.URL
	}
	return pathsToUrls
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yml)
	if err != nil {
		panic(err)
	}

	pathsToUrls := buildMap(pathUrls)

	return MapHandler(pathsToUrls, fallback), nil
}
