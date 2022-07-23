package credentials

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"runtime"
)

const credentialsFile = ".credentials.json"
const configFolder = ".pocket-cli"

type Credentials struct {
	RequestToken string `json:"request_token"`
	AccessToken  string `json:"access_token"`
}

func ReadStoredCredentials() (Credentials, error) {
	credFilepath := getFullCredentialFilePath()
	if _, err := os.Stat(credFilepath); errors.Is(err, fs.ErrNotExist) {
		log.Println("No stored credentials exist yet")
		return Credentials{}, err
	}

	f, err := os.ReadFile(credFilepath)
	if err != nil {
		log.Printf("Failed to read stored credentials: %v\n", err)
		return Credentials{}, err
	}

	var c Credentials

	err = json.Unmarshal(f, &c)
	if err != nil {
		log.Fatal("Failed to parse credentials file")
	}
	log.Println("Read stored credentials from disk")
	return c, nil
}

func StoreCredentials(c Credentials) {
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Failed to parse credentials to store")
	}

	err = os.Mkdir(getAppConfigFolderPath(), 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = os.WriteFile(getFullCredentialFilePath(), bytes, 0600)
	if err != nil {
		log.Fatal("Failed to store credentials file")
	}
}

func RemoveStoredCredentials() {
	err := os.Remove(getFullCredentialFilePath())
	if err != nil {
		log.Fatal("Failed to delete credentials file: %w", err)
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
