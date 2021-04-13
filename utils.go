package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func setStringFlag(strVar *string, name, defaultValue string, required bool, help string) {
	args := os.Args[1:]
	if len(args) >= 2 {
		for ai := 0; ai < len(args)-1; ai++ {
			if args[ai] == name {
				*strVar = args[ai+1]
				return
			}
		}
	}
	if required {
		fmt.Println("[Error] Missing mandatory parameter:", name)
		if len(help) > 0 {
			fmt.Printf("%s: %s\n", name, help)
		}
		os.Exit(1)
	}
	*strVar = defaultValue
}

func checkTls(cert, key string) Tls {
	if len(cert) == 0 && len(key) == 0 {
		return Tls{Active: false}
	}
	if !fileExists(cert) {
		log.Fatalf("Invalid cert file: '%s' was not found.\n", cert)
	}
	if !fileExists(key) {
		log.Fatalf("Invalid key file: '%s' was not found.\n", key)
	}
	return Tls{
		Active:   true,
		CertFile: cert,
		KeyFile:  key,
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func connIsOk(host string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	req, err := http.NewRequestWithContext(ctx, "GET", host, nil)
	if err != nil {
		cancel()
		return false
	}
	_, err = http.DefaultClient.Do(req)
	cancel()
	return err == nil
}

func exitWithError(err error) {
	fmt.Println("/!\\ Fatal Error /!\\")
	fmt.Println(err)
	os.Exit(1)
}

func displayHelp() {
	fmt.Println(`Proxy - Help

	The proxy can run with TLS, as long as both certification and key files are provided.

	Arguments:
		-h, --help
			Displays this help

		--conf <conf-file.yml>
			The location of the configuration file
			If a config file is provided, none of the following parameters are required, even the mandatory ones
		
		--host <host>
		--> [Mandatory]
			The target host the proxy will redirect to
		
		--port <port>
		--> [Mandatory]
			The port the proxy will listen to
		
		--cert <cert-file.pem>
			The location of the certification file

		--key <key-file.pem>
			The location of the key file`)
}
