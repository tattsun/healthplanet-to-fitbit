package htf

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var tz *time.Location

func init() {
	t, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("failed to load location: %v", err)
	}
	tz = t
}

type InnerScanTag int64

const (
	InnerScanTagWeight     InnerScanTag = 6021
	InnerScanTagBodyFatPct InnerScanTag = 6022
)

type InnerScanData struct {
	Date    string `json:"date"`
	KeyData string `json:"keydata"`
	Model   string `json:"model"`
	Tag     string `json:"tag"`
}

type AggregatedInnerScanData struct {
	Weight *float64
	Fat    *float64
}

type AggregatedInnerScanDataMap map[time.Time]*AggregatedInnerScanData

func (d *InnerScanData) Time() (time.Time, error) {
	layout := "200601021504"
	t, err := time.ParseInLocation(layout, d.Date, tz)
	if err != nil {
		return time.Time{}, err
	}

	return t.UTC(), nil
}

type InnerScanResponse struct {
	BirthDate string          `json:"birth_date"`
	Height    string          `json:"height"`
	Sex       string          `json:"sex"`
	Data      []InnerScanData `json:"data"`
}

type HealthPlanetAPI struct {
	AccessToken string
}

func (api *HealthPlanetAPI) AggregateInnerScanData(ctx context.Context) (AggregatedInnerScanDataMap, error) {
	weights, err := api.GetInnerScan(ctx, InnerScanTagWeight)
	if err != nil {
		return nil, err
	}

	fats, err := api.GetInnerScan(ctx, InnerScanTagBodyFatPct)
	if err != nil {
		return nil, err
	}

	m := make(AggregatedInnerScanDataMap, len(weights.Data))

	for _, weight := range weights.Data {
		t, err := weight.Time()
		if err != nil {
			log.Printf("invalid time: %+v", err)
			continue
		}

		data, err := strconv.ParseFloat(weight.KeyData, 64)
		if err != nil {
			log.Printf("invalid weight: %+v", err)
			continue
		}

		m[t] = &AggregatedInnerScanData{
			Weight: &data,
		}
	}

	for _, fat := range fats.Data {
		t, err := fat.Time()
		if err != nil {
			log.Printf("invalid time: %+v", err)
			continue
		}

		data, err := strconv.ParseFloat(fat.KeyData, 64)
		if err != nil {
			log.Printf("invalid fat: %+v", err)
			continue
		}

		if d, ok := m[t]; ok {
			d.Fat = &data
		} else {
			log.Printf("weight data not found: %+v", fat)
		}
	}

	return m, nil
}

func (api *HealthPlanetAPI) GetInnerScan(ctx context.Context, tag InnerScanTag) (InnerScanResponse, error) {
	values := url.Values{}
	values.Add("access_token", api.AccessToken)
	values.Add("date", "1")
	values.Add("tag", strconv.Itoa(int(tag)))

	url := fmt.Sprintf("https://www.healthplanet.jp/status/innerscan.json?%s", values.Encode())
	res, err := http.Get(url)
	if err != nil {
		return InnerScanResponse{}, errors.Wrap(err, "failed to fetch response")
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || 400 <= res.StatusCode {
		return InnerScanResponse{}, errors.Errorf("failed to get inner scan(invalid status code): %d", res.StatusCode)
	}

	dec := json.NewDecoder(res.Body)
	var resData InnerScanResponse
	if err = dec.Decode(&resData); err != nil {
		return InnerScanResponse{}, errors.Wrap(err, "failed to parse response")
	}

	return resData, nil
}
