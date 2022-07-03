package login

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
)

const credentialsFile = ".credentials.json"
const configFolder = ".pocket-cli"

type credentials struct {
	RequestToken string `json:"request_token"`
	AccessToken string `json:"access_token"`
}

func ReadStoredCredentials() (credentials, error) {
	f, err := os.ReadFile(getFullCredentialFilePath())
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

	err = os.Mkdir(getAppConfigFolderPath(), 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = os.WriteFile(getFullCredentialFilePath(), bytes, 0644)
	if err != nil {
		log.Println("Failed to store credentials file")
		panic(err)
	}
}

func RemoveStoredCredentials() {
	err := os.Remove(getFullCredentialFilePath())
	if err != nil {
		log.Println("Failed to delete credentials file")
		return
	}
	log.Println("Removed stored credentials")
}

func getFullCredentialFilePath() string {
	return path.Join(getAppConfigFolderPath(), credentialsFile)
}

func getAppConfigFolderPath() string {
	usrConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Can't find user config dir for platform %v: storing config in current dir", runtime.GOOS)
		usrConfigDir = "."
	}
	return path.Join(usrConfigDir, configFolder)
}
