# key-manager-client



## Purpose

This repository holds the code and build piepline for a client for our [Key Manager Plus instance](https://key-manager.is.ams-ix.net/apiclient/index.jsp#/Dashboard/All)

## How does it work

1. Download the binary from [Gihub Releases page](https://github.com/AMS-IX/key-manager-client)
2. Prepare client config file *client-conf.yaml*:
```
endpoint: "https://key-manager.xxx.ams-ix.net"
token: "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
common_name: "*.example.ams-ix.net"
serial_number: "1234567890abcdefghijklmnopqrstuvwxyz"
key_file: "/etc/tls/private/certificate.key"
cer_file: "/etc/tls//certificate.cer"
```
3. Run the client:
```
./bin/key-manager-client --config client-conf.yaml
```
4. Restart/reset any service that would use the new ceritificate (depends on the service)


## Building your own

1. Clone this repository
2. Ensure you have the version of go specified in *go.mod* installed - see [official documentation](https://go.dev/doc/install)
3. set environment variables:
```
export GOBIN=/home/go/key-manager-client/bin
export GOPATH=/home/go/key-manager-client
```
4. change directory to the location of the code:
```
/home/go/key-manager-client
```
5. Run the build:
```
go build -o key-manager-client ./cmd/keymanager-client
```