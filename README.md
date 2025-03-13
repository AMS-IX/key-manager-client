# key-manager-client



## Purpose

The `keymanager-client` is a command-line tool designed to interact with a Key Manager API to download certificate keys and certificates (CER files). It retrieves the necessary information from a YAML configuration file.

## How does it work

1. Download the binary from [Gihub Releases page](https://github.com/AMS-IX/key-manager-client)
2. The `keymanager-client` requires a configuration file to run.
The configuration file is in YAML format and contains the following parameters:

*endpoint*: The base URL of the Key Manager API.
*token*: The authentication token for accessing the API.
*common_name*: The common name of the certificate you want to download.
*serial_number* (string, optional): The serial number of the certificate. If not provided, the client will attempt to retrieve the serial number of the newest certificate with the given common name from the API.
*key_file* (string, required): The path to the file where the downloaded certificate key will be saved.
*cer_file* (string, required): The path to the file where the downloaded certificate (CER) will be saved.
Example Configuration File (with serial_number):

``` 
endpoint: https://key-manager.is.ams-ix.net
token: your_api_token_here
common_name: "*.example.ams-ix.net"
serial_number: 1234567890abcdef
key_file: /path/to/save/certificate.key
cer_file: /path/to/save/certificate.cer
```
Example Configuration File (without serial_number):

``` 
endpoint: https://key-manager.is.ams-ix.net
token: your_api_token_here
common_name: "*.example.ams-ix.net"
key_file: /path/to/save/certificate.key
cer_file: /path/to/save/certificate.cer
```
3. You specify the path to this file using the `--config` flag. Run the client:
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

## Security Considerations
- Protect the configuration file: The configuration file contains sensitive information (API token). Ensure that it is stored securely and only accessible to authorized users.
- API Token: Treat the API token as a secret. Do not hardcode it directly into scripts or commit it to version control.
- File Permissions: Make sure the key and cer files are stored with the correct permissions.