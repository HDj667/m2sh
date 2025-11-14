package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/term"

	"cert.at/m2sh/internal/config"
)

func handleAPIError(err error, operation string) {
	if err != nil {
		var appError *model.AppError
		if errors.As(err, &appError) {
			log.Fatalf("%s failed (Status: %d, msg: %s (ID: %s))", operation, appError.StatusCode, appError.Message, appError.Id)
		} else {
			log.Fatalf("%s failed: %v", operation, err)
		}
	}
}

func main() {
	// Load configuration from INI file and environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	mattermostURL := cfg.MattermostURL
	username := cfg.Username
	password := cfg.Password

	// If password is not set in config or environment, prompt for it
	if password == "" {
		fmt.Print("Mattermost Password: ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("Error reading password: %v", err)
		}
		fmt.Println()
		password = string(bytePassword)
	}

	client := model.NewAPIv4Client(mattermostURL)
	fmt.Printf("connecting to %s as %s...\n", mattermostURL, username)

	var totpToken string
	fmt.Print("TOTP Token: ")
	_, err = fmt.Scanln(&totpToken)
	if err != nil {
		log.Fatalf("Error reading TOTP token: %v", err)
	}

	ctx := context.Background()

	user, resp, err := client.LoginWithMFA(ctx, username, password, totpToken)
	handleAPIError(err, "MFA login")
	//fmt.Printf("MFA login %s|%s ==> successfull", username, user.Username)

	authToken := resp.Header.Get(model.HeaderToken)
	if authToken == "" {
		log.Fatal("login ok (Status 200), but no auth token returned from server")
	}
	client.SetToken(authToken)
	fmt.Println("--- Authentication successful  ---")
	fmt.Printf("logged in as %s (ID: %s)\n", user.Username, user.Id)

	fmt.Println("\n--- Listing Teams ---")
	teams, _, err := client.GetTeamsForUser(ctx, "me", "")
	handleAPIError(err, "GetTeamsForUser")

	if len(teams) == 0 {
		fmt.Println("User has no teams")
		return
	}

	fmt.Println("\n--- Listing Channels ---")
	for _, team := range teams {
		fmt.Printf("== Team: %s\n", team.DisplayName)

		channels, _, err := client.GetChannelsForTeamForUser(ctx, team.Id, "me", false, "")
		handleAPIError(err, fmt.Sprintf("Channels from team %s", team.DisplayName))
		if len(channels) == 0 {
			fmt.Println("   (Team has no channels)")
		} else {
			for _, channel := range channels {
				channelName := channel.DisplayName
				if channelName == "" {
					channelName = channel.Name
				}
				if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
					//channelName = fmt.Sprintf("%s (%s)", channelName, channel.Type)
					continue
				}
				fmt.Printf("   - Channel: %s (Typ: %s, ID: %s)\n", channelName, channel.Type, channel.Id)
			}
		}
	}
}
