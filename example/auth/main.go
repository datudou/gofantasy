package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/gofantasy"
)

func main() {
	redirectURL := os.Getenv("YAHOO_REDIRECT_URL")
	clientID := os.Getenv("YAHOO_CLIENT_ID")
	ctx := context.Background()

	ya := gofantasy.
		NewYahooClient().OAuth2(clientID, "", redirectURL)

	authCodeUrl, err := ya.GetAuthCodeUrl()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("please copy the below link to the browser:\n %s", authCodeUrl)
	fmt.Println("\nEnter the code: ")
	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return
	}

	err = ya.GetAccessToken(ctx, code)
	if err != nil {
		fmt.Println(err)
	}
	err = ya.SaveToken("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("save to path: ~/.config/gofantasy/yahoo_token.json")
}
