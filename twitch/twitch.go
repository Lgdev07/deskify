package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Twitch struct {
	gorm.Model
	Name   string `gorm:"size:100;not null" json:"name"`
	IsLive bool   `json:"is_live"`
}

func Initialize(db *gorm.DB) {
	for {
		time.Sleep(2 * time.Second)
		go MakeRequest(db)

	}
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
	req.Header.Set("client-id", os.Getenv("CLIENT_ID"))
	req.Header.Set("Authorization", os.Getenv("BEARER_TOKEN"))

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

	dataResponse := responseInterface["data"].([]interface{})

	if len(dataResponse) > 0 {
		return true
	}

	return false

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
