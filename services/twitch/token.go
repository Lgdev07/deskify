package twitch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Lgdev07/deskify/utils"
	"github.com/joho/godotenv"
)

var authToken string = getAuthToken()

func updateBearerTokenDotEnv(newAuthToken string) {
	dotEnvPath := utils.DotEnvPath()
	file, err := os.Open(dotEnvPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var text string
	var updated bool = false

	for scanner.Scan() {
		array := strings.Split(scanner.Text(), "=")
		variableName := array[0]
		variableValue := array[1]

		if variableName == "BEARER_TOKEN" {
			variableValue = newAuthToken
			updated = true
		}

		text += fmt.Sprintf("%s=%s\n", variableName, variableValue)
	}

	if !updated {
		text += fmt.Sprintf("BEARER_TOKEN=%s\n", newAuthToken)
	}

	authToken = newAuthToken

	err = ioutil.WriteFile(dotEnvPath, []byte(text), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func getAuthToken() string {
	dotEnvPath := utils.DotEnvPath()
	err := godotenv.Load(dotEnvPath)
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	return os.Getenv("BEARER_TOKEN")
}

func createNewToken() string {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%v&client_secret=%v&grant_type=client_credentials", clientID, clientSecret)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseInterface := make(map[string]interface{})

	err = json.Unmarshal(body, &responseInterface)
	if err != nil {
		log.Fatal(err)
	}

	return responseInterface["access_token"].(string)
}
