package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Lgdev07/deskify/twitch"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

func InitTwitchCmd(db *gorm.DB) {

	var cmdTwitch = &cobra.Command{
		Use:   "twitch [action]",
		Short: "Do an action with the command",
		Long:  "twitch command preceed an action.",
		Args:  cobra.MinimumNArgs(1),
	}

	var cmdTwitchAdd = &cobra.Command{
		Use:   "add [channel]",
		Short: "Add an channel to be notified when goes live",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			channel := fmt.Sprint(strings.Join(args, " "))
			TwitchAddChannel(db, channel)
		},
	}

	var cmdTwitchRem = &cobra.Command{
		Use:   "rem [channel]",
		Short: "Remove an channel",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			channel := fmt.Sprint(strings.Join(args, " "))
			TwitchRemoveChannel(db, channel)
		},
	}

	rootCmd.AddCommand(cmdTwitch)

	cmdTwitch.AddCommand(cmdTwitchAdd)
	cmdTwitch.AddCommand(cmdTwitchRem)

}

func TwitchAddChannel(db *gorm.DB, channelName string) {
	channel := twitch.Twitch{}

	db.Model(&twitch.Twitch{}).Where("name = ?", channelName).First(&channel)

	if channel.Name != "" {
		fmt.Println("There is already a channel with the same name")
		return
	}

	newChannel := &twitch.Twitch{
		Name:   channelName,
		IsLive: false,
	}

	db.Create(newChannel)
	fmt.Printf("Channel %s Added with success\n", channelName)

}

func TwitchRemoveChannel(db *gorm.DB, channelName string) {

	channel := twitch.Twitch{}

	db.Model(&twitch.Twitch{}).Where("name = ?", channelName).First(&channel)

	if channel.Name == "" {
		fmt.Println("We did not find a channel with that name")
		return
	}

	err := db.Model(&twitch.Twitch{}).Where("name = ?", channelName).Delete(&twitch.Twitch{}).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Channel %s Deleted with success\n", channelName)
}
