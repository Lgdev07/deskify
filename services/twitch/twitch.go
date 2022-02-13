package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Twitch is the struct associated with database model.
type Twitch struct {
	gorm.Model
	Name   string `gorm:"size:100;not null" json:"name"`
	IsLive bool   `json:"is_live"`
}

// Initialize runs the functionality.
func Initialize(wg *sync.WaitGroup, db *gorm.DB) {
	wg.Add(1)

	go func() {
		for {
			makeRequest(db)
			time.Sleep(50 * time.Second)
		}
	}()
}

func makeRequest(db *gorm.DB) {
	channels := &[]Twitch{}

	err := db.Model(&Twitch{}).Where("is_live = 0").Find(&channels).Error
	if err != nil {
		log.Fatal(err)
	}

	for _, channel := range *channels {
		fmt.Printf("Looking for channel: %s\n", channel.Name)
		isLive := request("https://api.twitch.tv/helix/streams?user_login=" + channel.Name)

		if isLive {
			notify(channel.Name)
			updateChannelSetIsLiveTrue(db, channel.Name)
		}
	}
	VerifyWitchChannelsAreActive(db)
}

func notify(channelName string) {
	err := beeep.Notify(channelName+" Is Live!", "Go Watch", "assets/twitch.png")
	if err != nil {
		panic(err)
	}
}

func updateChannelSetIsLiveTrue(db *gorm.DB, channelName string) {
	err := db.Model(&Twitch{}).Where("name = ?", channelName).Update("is_live", true).Error
	if err != nil {
		log.Fatal(err)
	}
}

func request(url string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	bearerToken := fmt.Sprintf("Bearer %s", authToken)

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

		newToken := createNewToken()
		updateBearerTokenDotEnv(newToken)

		return request(url)
	}

	dataResponse := responseInterface["data"].([]interface{})

	return len(dataResponse) > 0
}

func VerifyWitchChannelsAreActive(db *gorm.DB) {
	channels := &[]Twitch{}

	err := db.Model(&Twitch{}).Where("is_live = 1").Find(&channels).Error
	if err != nil {
		log.Fatal(err)
	}

	for _, channel := range *channels {
		isLive := request("https://api.twitch.tv/helix/streams?user_login=" + channel.Name)

		if !isLive {
			UpdateChannelSetIsLiveFalse(db, channel.Name)
		}
	}
}

func UpdateChannelSetIsLiveFalse(db *gorm.DB, channelName string) {
	err := db.Model(&Twitch{}).Where("name = ?", channelName).Update("is_live", false).Error
	if err != nil {
		log.Fatal(err)
	}
}
