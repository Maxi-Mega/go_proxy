package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Proxy struct {
	config       Config
	requestCount int
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy.requestCount++
	fmt.Printf("Request %d: %s\n", proxy.requestCount, r.URL.Path)
	if proxy.IsUrlForbidden(r.URL.Path) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "The proxy has blocked this URL")
		return
	}
	requestedUrl := joinURL(proxy.GetTargetHost(), r.URL.Path)
	req, err := http.NewRequest(r.Method, requestedUrl, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	req.Header = proxy.config.Headers.CloneFilter(r.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (proxy Proxy) GetProxyPort() uint16 {
	return proxy.config.ProxyPort
}

func (proxy Proxy) GetTargetHost() string {
	return proxy.config.TargetHost
}

func (proxy Proxy) GetCertFile() string {
	return proxy.config.Tls.CertFile
}

func (proxy Proxy) GetKeyFile() string {
	return proxy.config.Tls.KeyFile
}

func (proxy Proxy) IsTlsActive() bool {
	return proxy.config.Tls.Active
}

func (proxy Proxy) PortToString() string {
	return ":" + strconv.Itoa(int(proxy.config.ProxyPort))
}

func (proxy *Proxy) IsUrlForbidden(url string) bool {
	for _, blocked := range proxy.config.UrlBlocked {
		if strings.HasSuffix(blocked, "*") {
			if strings.HasPrefix(url, blocked[:len(blocked)-1]) {
				return true
			}
		}
		if strings.HasPrefix(blocked, "*") {
			if strings.HasSuffix(url, blocked[1:]) {
				return true
			}
		}
		if url == blocked {
			return true
		}
	}
	return false
}
