package main

import (
	"context"
	"fmt"
	"log"
	"os"

	htf "healthplanet-to-fitbit"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	healthPlanetAccessToken := os.Getenv("HEALTHPLANET_ACCESS_TOKEN")
	fitbitClientId := os.Getenv("FITBIT_CLIENT_ID")
	fitbitClientSecret := os.Getenv("FITBIT_CLIENT_SECRET")
	fitbitAccessToken := os.Getenv("FITBIT_ACCESS_TOKEN")
	fitbitRefreshToken := os.Getenv("FITBIT_REFRESH_TOKEN")

	// Initialize API clients
	healthPlanetAPI := htf.HealthPlanetAPI{
		AccessToken: healthPlanetAccessToken,
	}
	fitbitApi := htf.NewFitbitAPI(fitbitClientId, fitbitClientSecret, fitbitAccessToken, fitbitRefreshToken)

	// Initialize Context
	ctx := context.Background()

	// Get data from HealthPlanet
	scanData, err := healthPlanetAPI.AggregateInnerScanData(ctx)
	if err != nil {
		log.Fatalf("failed to aggregate inner scan data: %+v", err)
	}

	// Save data to Fitbit
	for t, data := range scanData {
		weightLog, err := fitbitApi.GetBodyWeightLog(t)
		if err != nil {
			log.Fatalf("failed to get weight log from fitbit: %+v", err)
		}

		if len(weightLog.Weight) > 0 {
			log.Printf("%s: record is found", t)
			continue
		}

		if data.Weight != nil {
			if err := fitbitApi.CreateWeightLog(*data.Weight, t); err != nil {
				log.Fatalf("failed to create weight log: time: %s, err: %+v", t, err)
			}
		}

		if data.Fat != nil {
			if err := fitbitApi.CreateBodyFatLog(*data.Fat, t); err != nil {
				log.Fatalf("failed to create fat log: time: %s, err: %+v", t, err)
			}
		}

		printFloat := func(f *float64) string {
			if f == nil {
				return "nil"
			}
			return fmt.Sprintf("%.2f", *f)
		}

		log.Printf("%s: saved, weight: %s, fat: %s", t, printFloat(data.Weight), printFloat(data.Fat))
	}

	log.Printf("done")
}
