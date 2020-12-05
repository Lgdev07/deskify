package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var updatedBearerToken string

type Twitch struct {
	gorm.Model
	Name   string `gorm:"size:100;not null" json:"name"`
	IsLive bool   `json:"is_live"`
}

func Initialize(wg *sync.WaitGroup, db *gorm.DB) {
	wg.Add(1)

	go func() {
		for {
			MakeRequest(db)
			time.Sleep(50 * time.Second)
		}
	}()

}

func MakeRequest(db *gorm.DB) {
	channels := &[]Twitch{}

	err := db.Model(&Twitch{}).Where("is_live = 0").Find(&channels).Error
	if err != nil {
		log.Fatal(err)
	}

	for _, channel := range *channels {
		fmt.Printf("Looking for channel: %s\n", channel.Name)
		isLive := Request("https://api.twitch.tv/helix/streams?user_login=" + channel.Name)

		if isLive {
			CallNotify(channel.Name)
			SetTrue(db, channel.Name)
		}
	}
	VerifyWitchAreActive(db)
}

func CallNotify(channelName string) {
	err := beeep.Notify(channelName+" Is Live!", "Go Watch", "assets/twitch.png")
	if err != nil {
		panic(err)
	}
}

func SetTrue(db *gorm.DB, channelName string) {
	err := db.Model(&Twitch{}).Where("name = ?", channelName).Update("is_live", true).Error
	if err != nil {
		log.Fatal(err)
	}
}

func Request(url string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	bearerToken := getBearerToken()

	req.Header.Set("client-id", os.Getenv("CLIENT_ID"))
	req.Header.Set("Authorization", bearerToken)

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

	if responseInterface["data"] == nil {
		fmt.Println("Invalid Token, creating a new one...")

		input, err := ioutil.ReadFile(".env")
		if err != nil {
			log.Fatalln(err)
		}

		lines := strings.Split(string(input), " ")
		token := lines[len(lines)-1]
		newToken := createNewToken()
		newOutput := strings.Replace(string(input), token, newToken, 1)

		err = ioutil.WriteFile(".env", []byte(newOutput), 0644)
		if err != nil {
			log.Fatalln(err)
		}

		updatedBearerToken = fmt.Sprint("Bearer ", newToken)

		return Request(url)
	}

	dataResponse := responseInterface["data"].([]interface{})

	if len(dataResponse) > 0 {
		return true
	}

	return false

}

func getBearerToken() string {
	token := os.Getenv("BEARER_TOKEN")
	if updatedBearerToken == "" {
		return token
	}

	if token == updatedBearerToken {
		return token
	}

	return updatedBearerToken
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

func VerifyWitchAreActive(db *gorm.DB) {
	channels := &[]Twitch{}

	err := db.Model(&Twitch{}).Where("is_live = 1").Find(&channels).Error
	if err != nil {
		log.Fatal(err)
	}

	for _, channel := range *channels {
		isLive := Request("https://api.twitch.tv/helix/streams?user_login=" + channel.Name)

		if !isLive {
			SetFalse(db, channel.Name)
		}
	}
}

func SetFalse(db *gorm.DB, channelName string) {
	err := db.Model(&Twitch{}).Where("name = ?", channelName).Update("is_live", false).Error
	if err != nil {
		log.Fatal(err)
	}
}
