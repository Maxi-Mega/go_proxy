# Go Proxy

This project is a golang written proxy, that redirect received request to a precise host. The proxy is able to listen
https requests, if certifications file are provided.

A Yaml configuration file can be provided to replace command line arguments. It also contains a "headers" part which is
used to allow, block, replace or set headers in each request. The config file also provide a "url-blocking" part for
blocking requests to desired URLs.

## Command Line Startup

### CLI arguments list

```
-h, --help
    Displays a help message

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
    The location of the key file
```

### Minimal arguments

```shell
./go_proxy --host "https://host-target.com" --port 8082
```

### With TLS

```shell
./go_proxy --host "https://host-target.com" --port 8082 --cert "./certs/localhost.pem" --key "./certs/localhost-key.pem"
# Needs both certification and key files
```

### With a config file

```shell
./go_proxy --conf "./config.yml"
```

## Configuration file syntax

```yaml
# The port the proxy will listen to
proxy-port: 8082

# The host the proxy will redirect to
target-host: "https://host-target.com"

tls:
  active: true

  # The path to the certification file
  cert-file: "certs/localhost.pem"

  # The path to the certification key file
  key-file: "certs/localhost-key.pem"

headers:
  # 'allow-only' just allows headers who are in the list to pass (wildcard: '*')
  allow-only: [ "*" ]

  # 'block-only' just blocks headers who are in the list (wildcard: '*')
  block-only: [ "Accept-Encoding" ]

  # 'replace' replaces the key by the corresponding value if the key is already present in the header
  replace:
    "User-Agent": "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0"
    "Accept": "*"

  # 'set' sets a new entry 'key' in the header with the corresponding value if it is not present or overwrite the existing one if there is one
  set:
    "Secret-Key": "0123456789"

url-blocking:
  # Use '/lib/*' to block access to any resource under '/lib/'
  - "/lib/*"
  # Use '*.css' to block access to any 'css' resource
  - "*.css"
```