package client

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

type Client struct {
    Endpoint     string
    AuthToken    string
    CommonName   string
    SerialNumber string
}

func NewClient(endpoint, authToken, commonName, serialNumber string) (*Client, error) {
    if endpoint == "" || authToken == "" {
        return nil, errors.New("endpoint and token must be provided")
    }
    return &Client{Endpoint: endpoint, AuthToken: authToken, CommonName: commonName, SerialNumber: serialNumber}, nil
}

// DownloadKey downloads the certificate key using the Key Manager REST API.
func (c *Client) DownloadKey() ([]byte, error) {
    return c.downloadCertificate("KEY")
}

// DownloadPEM downloads the certificate PEM using the Key Manager REST API.
func (c *Client) DownloadCER() ([]byte, error) {
    return c.downloadCertificate("CER")
}

func (c *Client) downloadCertificate(fileType string) ([]byte, error) {
    apiurl := c.Endpoint + "/api/pki/restapi/getCertificate"

    inputData := `{"operation":{"Details":{"common_name":"___CommonName___","serial_number":"___SerialNumber___","fileType":"___FieldType___"}}}`

    // Create form values and set the raw JSON string
    form := url.Values{}
    form.Set("INPUT_DATA", inputData)

    // Replace placeholders with real values
    requestUrl := apiurl + "?" + form.Encode()
    requestUrl = strings.Replace(requestUrl, "___CommonName___", c.CommonName, -1)
    requestUrl = strings.Replace(requestUrl, "___SerialNumber___", c.SerialNumber, -1)
    requestUrl = strings.Replace(requestUrl, "___FieldType___", fileType, -1)


    // Create the HTTP request
    req, err := http.NewRequest("GET", requestUrl, nil)
    if err != nil {
        return nil, err
    }

    // Set authorization and content type headers
    req.Header.Set("AUTHTOKEN", c.AuthToken)
    req.Header.Set("Content-Type", "application/json")


    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("non-200 response: %s", resp.Status)
    }

    // Read and return the certificate data.
    responseData, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return responseData, nil
}