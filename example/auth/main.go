package main

import (
	"bufio"
	"fmt"
	"github.com/gofantasy"
	"os"
)

func main() {
	redirectURL := os.Getenv("YAHOO_REDIRECT_URL")
	clientID := os.Getenv("YAHOO_CLIENT_ID")

	yc := gofantasy.
		NewClient().
		Yahoo().
		OAuth2(clientID, "", redirectURL)

	authCodeUrl, err := yc.GetAuthCodeUrl()
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

	err = yc.GetAccessToken(code).SaveToken("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("save to path: ~/.config/gofantasy/yahoo_token.json")
}
