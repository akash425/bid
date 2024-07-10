package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

type Banner struct {
	Type int `json:"type"`
}

type Bid struct {
	ID     string `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Banner Banner `json:"banner"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)
}

func validateRequest(bid Bid) error {
	if bid.ID == "" {
		return fmt.Errorf("Missing 'id'")
	}

	if bid.Banner.Type != 1 && bid.Banner.Type != 2 {
		return fmt.Errorf("Invalid 'banner.type'")
	}
	if bid.Width <= 0 || bid.Height <= 0 {
		return fmt.Errorf("Invalid 'width' or 'height'")
	}

	return nil
}

func bidHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Missing authorization header")
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token != "test1234" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Invalid authorization header")
		return
	}

	var bid Bid
	json.NewDecoder(r.Body).Decode(&bid)
	fmt.Print(bid)

	validateRequest(bid)

	if err := validateRequest(bid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	id := bid.ID
	if rdb.Exists(ctx, id).Val() == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ID not found in redis"})
		return
	}
	impression := rdb.HGet(ctx, id, "impression").Val()
	click := rdb.HGet(ctx, id, "click").Val()
	videoURL := rdb.HGet(ctx, id, "video_url").Val()
	videoStart := rdb.HGet(ctx, id, "video_start").Val()
	videoEnd := rdb.HGet(ctx, id, "video_end").Val()

	var response string

	bannerType := bid.Banner.Type

	if bannerType == 1 {
		jsTemplate := rdb.Get(ctx, fmt.Sprintf("%s_js", id)).Val()
		response = strings.Replace(jsTemplate, "{impression}", impression, -1)
		response = strings.Replace(response, "{click}", click, -1)
		w.Header().Set("Content-Type", "text/javascript")
	} else if bannerType == 2 {
		xmlTemplate := rdb.Get(ctx, fmt.Sprintf("%s_xml", id)).Val()
		response = strings.Replace(xmlTemplate, "{impression}", impression, -1)
		response = strings.Replace(response, "{click}", click, -1)
		response = strings.Replace(response, "{video_url}", videoURL, -1)
		response = strings.Replace(response, "{video_start}", videoStart, -1)
		response = strings.Replace(response, "{video_end}", videoEnd, -1)
		w.Header().Set("Content-Type", "text/xml")
	}
	fmt.Fprintf(w, ":%v", response)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/bid", bidHandler).Methods("POST")
	r.HandleFunc("/", healthCheck).Methods("GET")
	http.ListenAndServe(":5004", r)
	fmt.Println("Server listening on 5000")
}
