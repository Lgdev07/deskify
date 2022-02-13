package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Lgdev07/deskify/services/twitch"
	"github.com/Lgdev07/deskify/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var avaliableConfigs = []string{"CLIENT_ID", "CLIENT_SECRET"}

func InitTwitchCmd(db *gorm.DB) {
	var cmdTwitch = &cobra.Command{
		Use:   "twitch [action]",
		Short: "Add Twitch channels and be notified when it goes live",
		Long:  "Add Twitch channels and be notified when it goes live",
		Args:  cobra.MinimumNArgs(1),
	}

	var cmdTwitchRun = &cobra.Command{
		Use:   "run",
		Short: "Check if a channel is live and be notified",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			twitch.Initialize(&wg, db)
			wg.Wait()
		},
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

	var cmdTwitchConfig = &cobra.Command{
		Use:   "config [name] [value]",
		Short: "Add values to your configs",
		Long: ("two avaliable configs\n" +
			"client_id and client_secret\n" +
			"example: deskify twitch client_id example_value"),
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			value := args[1]
			TwitchAddConfig(name, value)
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

	var cmdTwitchList = &cobra.Command{
		Use:   "list",
		Short: "List all channels",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			TwitchListChannels(db)
		},
	}

	rootCmd.AddCommand(cmdTwitch)

	cmdTwitch.AddCommand(cmdTwitchRun)
	cmdTwitch.AddCommand(cmdTwitchAdd)
	cmdTwitch.AddCommand(cmdTwitchRem)
	cmdTwitch.AddCommand(cmdTwitchList)

	cmdTwitch.AddCommand(cmdTwitchConfig)
}

func TwitchAddConfig(name, value string) {
	if !validateConfigName(name) {
		configs := strings.Join(avaliableConfigs, ", ")
		fmt.Printf("Avaliable configs: %s\n", configs)
		return
	}

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

		if variableName == name {
			variableValue = value
			updated = true
		}

		text += fmt.Sprintf("%s=%s\n", variableName, variableValue)
	}

	if !updated {
		text += fmt.Sprintf("%s=%s\n", name, value)
	}

	err = ioutil.WriteFile(dotEnvPath, []byte(text), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func validateConfigName(inputName string) bool {
	for _, configName := range avaliableConfigs {
		if configName == inputName {
			return true
		}
	}
	return false
}

func TwitchAddChannel(db *gorm.DB, channelName string) {
	channel := twitch.Twitch{}

	db.Model(&twitch.Twitch{}).Where("name = ?", channelName).First(&channel)

	if channel.Name != "" {
		fmt.Println("A channel with the same name already exists")
		return
	}

	newChannel := &twitch.Twitch{
		Name:   channelName,
		IsLive: false,
	}

	db.Create(newChannel)
	fmt.Printf("Channel %s added successfully\n", channelName)
}

func TwitchRemoveChannel(db *gorm.DB, channelName string) {
	channel := twitch.Twitch{}

	db.Model(&twitch.Twitch{}).Where("name = ?", channelName).First(&channel)

	if channel.Name == "" {
		fmt.Println("We couldn't find a channel with that name")
		return
	}

	err := db.Model(&twitch.Twitch{}).Where("name = ?", channelName).Delete(&twitch.Twitch{}).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Channel %s successfully deleted\n", channelName)
}

func TwitchListChannels(db *gorm.DB) {
	twitchList := []twitch.Twitch{}

	db.Model(&twitch.Twitch{}).Find(&twitchList)

	if len(twitchList) == 0 {
		fmt.Println("No channels found")
		return
	}

	for _, value := range twitchList {
		fmt.Printf("Channel: %s\n", value.Name)
	}
}
