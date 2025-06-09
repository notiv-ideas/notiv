package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func initCfg() {
	if err := loadCfg(); err != nil {
		if os.IsNotExist(err) {
			createCfg()
		} else {
			panic(err)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Println("notiv configs already exist", filepath.Join(dir, notivDir))
	}
}

func readCfg() {
	// Open the config file
	file, err := os.Open(confPath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer file.Close()

	// Create a map to hold the config values
	var cfg map[string]string

	// Unmarshal the JSON data into the map
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		fmt.Println("Error unmarshaling config data:", err)
		return
	}
}

func listCfg() {
	readCfg()
	// Range over the map and print key-value pairs
	for k, v := range cfg {
		fmt.Println(k, v)
	}
}

func loadCfg() error {

	if share_flag != nil && *share_flag != "" {
		projectDir = *share_flag
	} else if projectDir != "" {
		projectDir = "."
	}
	notivDir = filepath.Join(projectDir, ".notiv")
	confDir = filepath.Join(projectDir, ".notiv", "config")
	confPath = filepath.Join(confDir, "config.json")
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	return nil
}

func createCfg() {
	// Do not overwrite
	if _, err := os.Stat(confPath); err == nil {
		fmt.Println("config file is already created", confPath)
		return
	}
	if err := os.MkdirAll(confDir, 0755); err != nil {
		if !strings.Contains(err.Error(), "file exists") {
			panic(err)
		}
	}
	conf, err := os.Create(confPath)
	if err != nil {
		panic(err)
	}
	defer conf.Close()

	newConfig := map[string]string{
		"defaultDecryptPolicy":     "false",
		"defaultDiskStore":         "true",
		"defaultEncryptPolicy":     "true",
		"defaultEncryptedCHECKSUM": encrypt(checksum, readPassword()),
		"defaultnotivDir":          notivDir,
		"defaultProjectDir":        ".",
		"defaultShareDir":          filepath.Join(notivDir, "share"),
		"defaultDataDir":           filepath.Join(notivDir, "share", "data"),
	}

	encoder := json.NewEncoder(conf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newConfig); err != nil {
		panic(err)
	}
	fmt.Println("config.json created at", confPath)
	// Load the config after creation
	if err := loadCfg(); err != nil {
		panic(err)
	}
}

func editCfg(key, value string) {
	if _, ok := cfg[key]; !ok {
		fmt.Println("cannot add adhoc configs")
		return
	}

	switch key {
	case "defaultEncryptedCHECKSUM":
		fmt.Printf("Cannot modify %s; use 'config password new'\n", key)
		return
	case "defaultDiskStore", "defaultEncryptPolicy", "defaultDecryptPolicy":
		value = strings.ToLower(value)
		if value != "true" && value != "false" {
			fmt.Printf("notiv config: %s , must be either true or false\n", key)
			return
		}
	}

	cfg[key] = value
	conf, err := os.OpenFile(confPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer conf.Close()
	encoder := json.NewEncoder(conf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		panic(err)
	}
	readCfg() // paranoia
	fmt.Println("notiv config: updated", key, cfg[key])
}
