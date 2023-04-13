package main

import (
	"bufio"
	"fmt"
	"github.com/gofantasy"
	"os"
)

func main() {
	redirectURL := "https://www.eoslaomao.com"
	clientID := os.Getenv("YAHOO_CLIENT_ID")
	//client secret is no need if app is registered as public client
	clientSecret := os.Getenv("YAHOO_CLIENT_SECRET")
	yo2 := gofantasy.NewYahooOAuth2(clientID, clientSecret, redirectURL)
	authCodeUrl, err := yo2.GetAuthCodeUrl()
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
	token, err := yo2.GetAccessToken(code)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token.AccessToken)
}
