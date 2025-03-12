package main

import (
    "flag"
    "fmt"
    "os"

    "gopkg.in/yaml.v2"
    "key-manager-client/pkg/client"
)

type Config struct {
    Endpoint     string `yaml:"endpoint"`
    Token        string `yaml:"token"`
    CommonName   string `yaml:"common_name"`
    SerialNumber string `yaml:"serial_number"`
    KeyFile      string `yaml:"key_file"`
    CerFile      string `yaml:"cer_file"`
}

func main() {
    configPath := flag.String("config", "/etc/keymanager-client.conf", "Path to the YAML configuration file")
    flag.Parse()

    if *configPath == "" {
        fmt.Println("Missing required flag --config")
        os.Exit(1)
    }

    data, err := os.ReadFile(*configPath)
    if err != nil {
        fmt.Printf("Failed to read config file: %v\n", err)
        os.Exit(1)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        fmt.Printf("Failed to parse config file: %v\n", err)
        os.Exit(1)
    }

    if cfg.Token == "" || cfg.Endpoint == "" || cfg.CommonName == "" || cfg.SerialNumber == "" || cfg.KeyFile == "" || cfg.CerFile == "" {
        fmt.Println("Config file must contain 'endpoint', 'token', 'common_name', 'serial_number', 'key_file', and 'cer_file' values")
        os.Exit(1)
    }

    c, err := client.NewClient(cfg.Endpoint, cfg.Token, cfg.CommonName, cfg.SerialNumber)
    if err != nil {
        fmt.Printf("Failed to create API client: %v\n", err)
        os.Exit(1)
    }

    // Download Key
    keyData, err := c.DownloadKey()
    if err != nil {
        fmt.Printf("Failed to download key: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Key downloaded successfully")
    err = os.WriteFile(cfg.KeyFile, keyData, 0644)
    if err != nil {
        fmt.Printf("Failed to write key to file: %v\n", err)
        os.Exit(1)
    }

    // Download CER
    cerData, err := c.DownloadCER()
    if err != nil {
        fmt.Printf("Failed to download CER: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("CER downloaded successfully")
    err = os.WriteFile(cfg.CerFile, cerData, 0644)
    if err != nil {
        fmt.Printf("Failed to write CER to file: %v\n", err)
        os.Exit(1)
    }
}