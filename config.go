package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	ProxyPort  uint16   `yaml:"proxy-port"`
	TargetHost string   `yaml:"target-host"`
	Tls        Tls      `yaml:"tls"`
	Headers    Headers  `yaml:"headers"`
	UrlBlocked []string `yaml:"url-blocking"`
	// Settings   Settings `yaml:"settings"`
}

type Tls struct {
	Active   bool   `yaml:"active"`
	CertFile string `yaml:"cert-file"`
	KeyFile  string `yaml:"key-file"`
}

type Headers struct {
	Active      bool
	AllowOnly   []string          `yaml:"allow-only"`
	BlockOnly   []string          `yaml:"block-only"`
	Replace     map[string]string `yaml:"replace"`
	Set         map[string]string `yaml:"set"`
	CloneFilter func(header http.Header) http.Header
}

type Settings struct {
	Active          bool
	CaseSentisitive bool `yaml:"case-sensitive"`
}

func parseConfigFile(configFile string) Config {
	if !fileExists(configFile) {
		exitWithError(fmt.Errorf("given config file '%s' does not exist", configFile))
	}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		exitWithError(fmt.Errorf("failed to open config file: %v", err))
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		exitWithError(fmt.Errorf("failed to parse config file: %v", err))
	}
	if !connIsOk(config.TargetHost) {
		exitWithError(fmt.Errorf("can't connect to target host '%s'", config.TargetHost))
	}
	if config.Tls.Active {
		if !fileExists(config.Tls.CertFile) {
			exitWithError(fmt.Errorf("certification file '%s' does not exist", config.Tls.CertFile))
		}
		if !fileExists(config.Tls.KeyFile) {
			exitWithError(fmt.Errorf("certification key file '%s' does not exist", config.Tls.KeyFile))
		}
	}
	config.Headers.Active = true
	config.Headers.CloneFilter = genHeaderFilter(config.Headers)
	return config
}

func parseConfigCli(target string, port string, tlsOpts Tls) Config {
	proxyPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln("Failed to parse CLI options:", err)
	}
	return Config{
		ProxyPort:  uint16(proxyPort),
		TargetHost: target,
		Tls:        tlsOpts,
		Headers:    Headers{Active: false},
		// Settings:   Settings{Active: false},
	}
}

func genHeaderFilter(headers Headers) func(header http.Header) http.Header {
	if headers.Active {
		var allowAll, blockAll bool
		sets := map[string][]string{}
		for h, v := range headers.Set {
			sets[http.CanonicalHeaderKey(h)] = []string{v}
		}
		replaces := map[string][]string{}
		for h, v := range headers.Replace {
			replaces[http.CanonicalHeaderKey(h)] = []string{v}
		}
		allows := map[string]struct{}{}
		for _, h := range headers.AllowOnly {
			if h == "*" {
				allowAll = true
				break
			}
			allows[http.CanonicalHeaderKey(h)] = struct{}{}
		}
		blocks := map[string]struct{}{}
		for _, h := range headers.BlockOnly {
			if h == "*" {
				blockAll = true
				break
			}
			blocks[http.CanonicalHeaderKey(h)] = struct{}{}
		}
		if blockAll {
			return func(header http.Header) http.Header {
				newHeader := http.Header{}
				for h, v := range header {
					newHeader[h] = v
				}
				for h, v := range sets {
					newHeader[h] = v
				}
				return newHeader
			}
		}
		return func(header http.Header) http.Header {
			newHeader := http.Header{}
			for h, v := range header {
				h := http.CanonicalHeaderKey(h)
				if _, isBlocked := blocks[h]; isBlocked {
					continue
				}
				if !allowAll {
					if _, isAllowed := allows[h]; !isAllowed {
						continue
					}
				}
				if rplc, mustBeReplace := replaces[h]; mustBeReplace {
					newHeader[h] = rplc
				} else {
					newHeader[h] = v
				}
			}
			for h, v := range sets {
				newHeader[h] = v
			}
			return newHeader
		}
	} else {
		return func(header http.Header) http.Header {
			/*newHeader := http.Header{}
			for h, v := range header {
				newHeader[h] = v
			}
			return newHeader*/
			return header.Clone()
		}
	}
}
