package client

import (
    "errors"
    "fmt"
    "io"
    "strings"
    "time"
    "net/http"
    "net/url"
    "encoding/json"
)

type Client struct {
    Endpoint     string
    AuthToken    string
    CommonName   string
    SerialNumber string
}
func NewClient(endpoint, authToken, commonName string, serialNumber ...string) (*Client, error) {
    if endpoint == "" || authToken == "" {
        return nil, errors.New("endpoint and token must be provided")
    }
    var sn string
    if len(serialNumber) > 0 {
        sn = serialNumber[0]
    }
    return &Client{Endpoint: endpoint, AuthToken: authToken, CommonName: commonName, SerialNumber: sn}, nil
}

// GetSerialNumber retrieves the serial number for the newest certificate with given common name from the Key Manager REST API.
func (c *Client) GetSerialNumber(commonName string) (string, error) {
    apiurl := c.Endpoint + "/api/pki/restapi/getAllSSLCertificates"

    req, err := http.NewRequest("GET", apiurl, nil)
    if err != nil {
        return "", err
    }

    req.Header.Set("AUTHTOKEN", c.AuthToken)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("non-200 response: %s", resp.Status)
    }

    var result map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return "", err
    }

    details, ok := result["details"].([]interface{})
    if !ok {
        return "", errors.New("details not found in response")
    }

    var latestSerialNumber string
    var latestExpiryDate time.Time

    for _, detail := range details {
        detailMap, ok := detail.(map[string]interface{})
        if !ok {
            continue
        }
        if detailMap["Common Name"] == commonName {
            serialNumber, ok := detailMap["serialNumber"].(string)
            if !ok {
                continue
            }
            expiryDateStr, ok := detailMap["ExpiryDate"].(string)
            if !ok {
                continue
            }
            layout := "Jan 2, 2006"
            expiryDate, err := time.Parse(layout, expiryDateStr)
            if err != nil {
                continue
            }
            if latestSerialNumber == "" || expiryDate.After(latestExpiryDate) {
                latestSerialNumber = serialNumber
                latestExpiryDate = expiryDate
            }
        }
    }

    if latestSerialNumber == "" {
        return "", errors.New("common name not found in response")
    }

    return latestSerialNumber, nil
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