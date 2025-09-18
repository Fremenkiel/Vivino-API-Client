package main

import (
	"compress/gzip"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
	"vivino/models"
)

const baseURL string = "https://www.vivino.com/api"

type VivinoClient struct {
	client    *http.Client
	userAgent string
	token     string
}

type VivinoResponse struct {
	Reviews []models.VivinoRating `json:"reviews"`
}

var vivinoRating VivinoResponse

func NewVivinoClient() *VivinoClient {
	return &VivinoClient{
		client:    &http.Client{Timeout: 30 * time.Second},
		userAgent: getRandomUserAgent(),
		token:     generateToken(),
	}
}

func (vc *VivinoClient) refreshToken() {
	vc.token = generateToken()
	vc.userAgent = getRandomUserAgent()
}

func generateToken() string {
	timestamp := time.Now().Unix()
	nonce, _ := rand.Int(rand.Reader, big.NewInt(1000000))

	data := fmt.Sprintf("%d:%d", timestamp, nonce.Int64())
	hash := md5.Sum([]byte(data))

	return fmt.Sprintf("%x:%d:%s", hash, timestamp, hex.EncodeToString(hash[:8]))
}

func getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	return userAgents[n.Int64()]
}

func (vc *VivinoClient) makeRequest(endpoint string, params map[string]string) (VivinoResponse, error) {
	fullURL := baseURL + endpoint

	if len(params) > 0 {
		u, _ := url.Parse(fullURL)
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return vivinoRating, err
	}

	vc.setHeaders(req)

	resp, err := vc.client.Do(req)
	if err != nil {
		return vivinoRating, err
	}
	defer resp.Body.Close()

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))

	if resp.StatusCode == 403 || resp.StatusCode == 429 {
		vc.refreshToken()
		time.Sleep(2 * time.Second)
		return vc.makeRequest(endpoint, params)
	}

	var reader io.Reader = resp.Body

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return vivinoRating, fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return vivinoRating, err
	}

	if err = json.Unmarshal(body, &vivinoRating); err != nil {
		return vivinoRating, err
	}

	return vivinoRating, nil
}

func (vc *VivinoClient) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", vc.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Referer", "https://www.vivino.com/")
	req.Header.Set("Origin", "https://www.vivino.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
}

func main() {
	args := os.Args[1:]

	wineIdIndex := slices.IndexFunc(args, func(s string) bool { return strings.Split(s, "=")[0] == "--id" })
	if wineIdIndex == -1 {
		fmt.Println("No wine ID provided")
		return
	}
	wineID := strings.Split(args[wineIdIndex], "=")[1]
	fmt.Println("Wine ID provided")

	yearIdIndex := slices.IndexFunc(args, func(s string) bool { return strings.Split(s, "=")[0] == "--year" })
	if yearIdIndex == -1 {
		fmt.Println("No year provided")
		return
	}
	year := strings.Split(args[yearIdIndex], "=")[1]

	client := NewVivinoClient()

	endpoint := fmt.Sprintf("/wines/%s/reviews", wineID)

	fmt.Printf("Fetching wine data for ID: %s\n", wineID)

	response, err := client.makeRequest(endpoint, map[string]string{
		"year":     year,
		"per_page": "1",
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Response:")
	fmt.Println(response.Reviews[0].Vintage.Statistics.RatingsAverage)
}
