package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Host Struct

type Host struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Config struct {
	Hosts []Host `yaml:"hosts"`
}

func readConfig(filename string) (Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// Define the gauge metric for network route status.
var (
	networkRouteUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_route_up",
			Help: "Indicates whether the network route is up. 1 = up, 0 = down",
		},
		[]string{"name", "endpoint"},
	)
)

// init function to register metrics
func init() {
	prometheus.MustRegister(networkRouteUp)
}

// Checks whether a network route is up by attempting to establish a TCP connection.
// If the connection succeeds, the route is up. If it fails, the route is down.
func checkNetworkRoute(name string, endpoint string) {
	_, err := net.DialTimeout("tcp", endpoint, 1*time.Second)
	if err != nil {
		networkRouteUp.WithLabelValues(name, endpoint).Set(0)
	} else {
		networkRouteUp.WithLabelValues(name, endpoint).Set(1)
	}
}

// Loop all over endpoints and check each one.
func checkAllNetworkRoutes(hosts []Host) {
	for _, host := range hosts {
		endpoint := fmt.Sprintf("%s:%d", host.Address, host.Port)
		checkNetworkRoute(host.Name, endpoint)
	}
}

func main() {

	// Define a string flag with a default value and a description.
	configPath := flag.String("config", "hosts.yaml", "path to the YAML configuration file")

	// Parse the flags.
	flag.Parse()

	// Reading config yaml file
	config, err := readConfig(*configPath)

	if err != nil {
		log.Fatalf("Error reading hosts yaml file: %v", err)
	}

	// Creating a list of endpoints from config.
	// endpoints := make([]string, len(config.Hosts))

	// for i, host := range config.Hosts {
	// 	endpoints[i] = fmt.Sprintf("%s:%d", host.Address, host.Port)
	// }

	// check the endpoints every 5 seconds
	go func() {
		for range time.Tick(5 * time.Second) {
			checkAllNetworkRoutes(config.Hosts)
		}
	}()

	log.Println("Starting server on port 2112...")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
