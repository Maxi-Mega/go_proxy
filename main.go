package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("")

	for _, arg := range os.Args[1:] {
		a := strings.ToLower(arg)
		if a == "-h" || a == "--help" {
			displayHelp()
			os.Exit(0)
		}
	}
	var config Config
	var configFile string
	setStringFlag(&configFile, "--conf", "", false, "The path to the configuration file")
	if len(configFile) > 0 {
		config = parseConfigFile(configFile)
	} else {
		var redirectHost string
		var port string
		var cert, key string
		setStringFlag(&redirectHost, "--host", "", true, "The host the proxy will redirect to")
		setStringFlag(&port, "--port", "", true, "The port the proxy will listen to")
		setStringFlag(&cert, "--cert", "", false, "The path to the certification PEM file")
		setStringFlag(&key, "--key", "", false, "The path to the certification key PEM file")

		tlsOpts := checkTls(cert, key)
		config = parseConfigCli(redirectHost, port, tlsOpts)
	}
	proxy := Proxy{config: config}

	fmt.Println("Redirection set to", proxy.GetTargetHost())
	if proxy.IsTlsActive() {
		fmt.Println("Starting proxy on port", proxy.GetProxyPort(), "with TLS\n")
		log.Fatal(http.ListenAndServeTLS(proxy.PortToString(), proxy.GetCertFile(), proxy.GetKeyFile(), &proxy))
	} else {
		fmt.Println("Starting proxy on port", proxy.GetProxyPort(), "...\n")
		log.Fatal(http.ListenAndServe(proxy.PortToString(), &proxy))
	}
}
