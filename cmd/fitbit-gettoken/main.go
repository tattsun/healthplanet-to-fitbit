package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	htf "healthplanet-to-fitbit"
)

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func genCodeChallenge() (verifier string, challenge string) {
	verifier = randomString(128)
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sum[:])
	return
}

func main() {
	godotenv.Load(".env")

	clientID := os.Getenv("FITBIT_CLIENT_ID")
	clientSecret := os.Getenv("FITBIT_CLIENT_SECRET")

	conf := htf.GetFitbitConfig(clientID, clientSecret)

	verifier, challenge := genCodeChallenge()

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		ctx := context.Background()
		token, err := conf.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", verifier))
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "err: %v", err)
			return
		}

		fmt.Fprintf(w, "AccessToken: %s\n", token.AccessToken)
		fmt.Fprintf(w, "RefreshToken: %s", token.RefreshToken)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := conf.AuthCodeURL("state",
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			oauth2.SetAuthURLParam("code_challenge", challenge),
			oauth2.AccessTypeOffline,
		)

		http.Redirect(w, r, url, http.StatusFound)
	})

	fmt.Println("Open: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
