package login

import (
	"encoding/json"
	"log"
	"os"
)

const credentialsFile = ".credentials.json"

type credentials struct {
	RequestToken string `json:"request_token"`
	AccessToken string `json:"access_token"`
}

func ReadStoredCredentials() (credentials, error) {
	f, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Println("No stored credentials found")
		return credentials{}, err
	}

	var c credentials 

	err = json.Unmarshal(f, &c)
	if err != nil {
		log.Println("Failed to parse credentials file")
		panic(err)
	}
	log.Println("Read stored credentials from disk")
	return c, nil
}

func StoreCredentials(c credentials) {
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Println("Failed to parse credentials to store")
		panic(err)
	}
	err = os.WriteFile(credentialsFile, bytes, 0644)
	if err != nil {
		log.Println("Failed to store credentials file")
		panic(err)
	}
}

func RemoveStoredCredentials() {
	err := os.Remove(credentialsFile)
	if err != nil {
		log.Println("Failed to delete credentials file")
		return
	}
	log.Println("Removed stored credentials")
}