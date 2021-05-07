package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/tjarratt/babble"
)

func main() {
	mmURL := os.Getenv("MM_URL")
	if mmURL == "" {
		fmt.Println("env var MM_URL must be set")
		return
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		fmt.Println("env var WEBHOOK_URL must be set")
		return
	}

	postRateMilliseconds, err := strconv.Atoi(os.Getenv("POST_RATE_MS"))
	if postRateMilliseconds == 0 || err != nil {
		fmt.Println("env var POST_RATE_MS must be set")
		return
	}

	editRateMilliseconds, err := strconv.Atoi(os.Getenv("EDIT_RATE_MS"))
	if postRateMilliseconds == 0 || err != nil {
		fmt.Println("env var EDIT_RATE_MS must be set")
		return
	}

	botToken := os.Getenv("BOT_TOKEN")
	if mmURL == "" {
		fmt.Println("env var BOT_TOKEN must be set")
		return
	}

	postID := os.Getenv("POST_ID_TO_EDIT")
	if mmURL == "" {
		fmt.Println("env var POST_ID_TO_EDIT must be set")
		return
	}

	babbler := babble.NewBabbler()
	babbler.Separator = " "
	babbler.Count = 10

	client := model.NewAPIv4Client(mmURL)
	client.SetToken(botToken)

	edit := time.NewTicker(time.Duration(math.Max(float64(editRateMilliseconds), 250) * 1000 * 1000))
	post := time.NewTicker(time.Duration(math.Max(float64(postRateMilliseconds), 250) * 1000 * 1000))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			close(stop)
			return
		case <-edit.C:
			if editRateMilliseconds > 0 {
				editPost(client, postID, babbler)
			}
		case <-post.C:
			if postRateMilliseconds > 0 {
				postToWebhook(webhookURL, babbler)
			}
		}
	}

}

func editPost(client *model.Client4, postID string, babbler babble.Babbler) {
	message := babbler.Babble()
	fmt.Println("Editing post with ID " + postID + " with message \"" + message + "\"")

	_, resp := client.PatchPost(postID, &model.PostPatch{Message: &message})
	if resp.Error != nil {
		fmt.Println("error: " + resp.Error.Error())
	}
	fmt.Println(resp.StatusCode)
}

func postToWebhook(webhookURL string, babbler babble.Babbler) {
	message := babbler.Babble()
	fmt.Println("Posting \"" + message + "\" to " + webhookURL)

	resp, err := http.Post(webhookURL, "application/json", strings.NewReader("{\"text\":\""+message+"\"}"))
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println(resp.Status)
}
