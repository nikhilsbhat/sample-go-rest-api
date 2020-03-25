// Package api has various components to handle the request that reaches this sample app.
// API is split to various segments to handle the request better.
// These components are handlers, middlewares, router and routes.
package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	decode "github.com/nikhilsbhat/config/decode"
)

var (
	//port   = "8080"
	config = `{
		"port": "80",
		"logpath": "neuron.log"
    }`
)

// Config required by the sample api app.
type Config struct {
	AppPort          string     `json:"port"`
	LogPath          string     `json:"logpath"`
	DefaultIdentites []Identity `json:"defaults"`
}

// API enables api for this app.
func API() {

	//Initializing router to prepare neuron to serve endpoints
	rout := new(MuxIn)

	// setting configurations
	conf, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stdout, err.(error).Error())
		os.Exit(1)
	}
	if len(conf.LogPath) != 0 {
		logp, err := conf.configLog()
		if err != nil {
			fmt.Fprintf(os.Stdout, err.(error).Error())
			rout.Apilog = os.Stdout
		} else {
			rout.Apilog = logp
		}
	} else {
		rout.Apilog = os.Stdout
	}

	router := rout.NewRouter()

	type r struct {
		router *mux.Router
		port   string
	}
	neurouter := r{router: router, port: conf.AppPort}
	fmt.Fprintf(os.Stdout, fmt.Sprintf("App is runnig on port: %s", conf.AppPort))
	errCh := make(chan error, 1)
	//starting the neuron on specified port
	var wg sync.WaitGroup
	wg.Add(1)
	go func(neurouter r) {
		starterr := http.ListenAndServe(":"+neurouter.port, neurouter.router)
		if starterr != nil {
			errCh <- starterr
		}
	}(neurouter)
	httperr := <-errCh
	if httperr != nil {
		log.Fatal(httperr)
	}
}

func (c *Config) logstat() bool {
	return statfile(c.LogPath)
}

func (c *Config) configLog() (io.Writer, error) {
	if c.logstat() != true {
		//fmt.Fprintf(os.Stdout, "unable to locate the log file specified, switching to STDOUT\n")
		return nil, fmt.Errorf("unable to locate the log file specified, switching to STDOUT\n")
	}
	path, err := os.OpenFile(c.LogPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		//fmt.Fprintf(os.Stdout, "unable open the log file specified, switching to STDOUT\n")
		return nil, fmt.Errorf("unable open the log file specified, switching to STDOUT\n")
	}
	return path, nil
}

func getConfig() (Config, error) {

	// fetching the configurations from env variable.
	if len(os.Getenv("API_CONFIG_PATH")) == 0 {
		fmt.Fprintf(os.Stdout, "no config file was specified, switching to default config\n")
		conf, err := decodeConfig([]byte(config))
		if err != nil {
			return Config{}, err
		}
		return conf, nil
	}

	if statfile(os.Getenv("API_CONFIG_PATH")) != true {
		fmt.Fprintf(os.Stdout, "unable to locate the config file specified, switching to default config\n")
		conf, err := decodeConfig([]byte(config))
		if err != nil {
			return Config{}, err
		}
		return conf, nil
	}
	// if len(os.Getenv("API_CONFIG4")) == 0 {
	// 	if err := decode.JsonDecode([]byte(config), &cnf); err != nil {
	// 		fmt.Println("Error Decoding JSON to gcpSVCred")
	// 		return Config{}, nil
	// 	}
	// 	return cnf, nil
	// }

	// fetching the configurations from the config file.
	jsonCont, err := decode.ReadFile(os.Getenv("API_CONFIG_PATH"))
	if err != nil {
		return Config{}, err
	}
	conf, err := decodeConfig(jsonCont)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
}

func decodeConfig(cont []byte) (Config, error) {
	var cnf Config
	if err := decode.JsonDecode(cont, &cnf); err != nil {
		return Config{}, fmt.Errorf(("Oops...! an error occured while decoding config file"))
	}
	return cnf, nil
}

func statfile(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
