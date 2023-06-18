package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var (
	loaders = map[string]func([]byte, any) error{
		".json": LoadFromJsonBytes,
		".yaml": LoadFromYamlBytes,
		".yml":  LoadFromYamlBytes,
	}
)

// Load loads config into v object from .json, .yaml, .yml file
// Note: Load will return error, you need handle the error on your own
// param file: file path
// param v: convert into v object
// param opts: customize Option, eg.  Load(file, v, UseEnv())
func Load(file string, v any, opts ...Option) error {
	// read file content
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Obtain the loader based on the file suffix
	loader, ok := loaders[strings.ToLower(path.Ext(file))]
	if !ok {
		return fmt.Errorf("unrecognized file type: %s", file)
	}

	// Use all custom configuration options
	// Add the option you want to modify to customOpts object
	var customOpts options
	for _, opt := range opts {
		opt(&customOpts)
	}

	// use environment variables option
	// replaces ${var} or $var in the string according to values of environment variables
	if customOpts.env {
		// parse content to v object
		return loader([]byte(os.ExpandEnv(string(content))), v)
	}

	// parse content to v object
	return loader(content, v)
}

// MustLoad loads config into v object from .json, .yaml, .yml file
// Note: exit on error
func MustLoad(file string, v any, opts ...Option) {
	if err := Load(file, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", file, err.Error())
	}
}
