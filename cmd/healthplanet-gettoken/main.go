package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const redirectURI = "https://www.healthplanet.jp/success.html"

type AuthorizeResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	healthPlanetClientId := os.Getenv("HEALTHPLANET_CLIENT_ID")
	healthPlanetClientSecret := os.Getenv("HEALTHPLANET_CLIENT_SECRET")

	values := url.Values{}
	values.Add("client_id", healthPlanetClientId)
	values.Add("redirect_uri", redirectURI)
	values.Add("scope", "innerscan")
	values.Add("response_type", "code")

	fmt.Printf("Authorize URL: %s\n", fmt.Sprintf("https://www.healthplanet.jp/oauth/auth?%s", values.Encode()))
	fmt.Println("")

	fmt.Print("Input code: ")
	var code string
	fmt.Scan(&code)

	values = url.Values{}
	values.Add("client_id", healthPlanetClientId)
	values.Add("client_secret", healthPlanetClientSecret)
	values.Add("redirect_uri", redirectURI)
	values.Add("code", code)
	values.Add("grant_type", "authorization_code")

	res, err := http.Post(fmt.Sprintf("https://www.healthplanet.jp/oauth/token?%s", values.Encode()), "application/json", nil)
	if err != nil {
		fmt.Printf("failed to get token: %+v", err)
		os.Exit(1)
	}
	if res.StatusCode < 200 || 400 <= res.StatusCode {
		fmt.Printf("failed to get token(invalid status code): %d", res.StatusCode)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("failed to get token(read body error): %+v", err)
		os.Exit(1)
	}

	fmt.Println("")
	var resData AuthorizeResponse
	if err = json.Unmarshal(resBody, &resData); err != nil {
		fmt.Printf("failed to parse response: %v", err)
		os.Exit(1)
	}
	fmt.Printf("AccessToken: %s\n", resData.AccessToken)
}
