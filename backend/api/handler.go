package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"banana-weather/pkg/genai"
	"banana-weather/pkg/maps"
	"banana-weather/pkg/storage"
)

type Handler struct {
	Maps    *maps.Service
	GenAI   *genai.Service
	Storage *storage.Service
}

type WeatherResponse struct {
	City        string `json:"city"`
	ImageBase64 string `json:"image_base64"`
}

type Preset struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
	VideoURL string `json:"video_url"`
}

func (h *Handler) HandleGetPresets(w http.ResponseWriter, r *http.Request) {
	// Try to read presets.json from GCS
	data, err := h.Storage.ReadObject(r.Context(), "presets.json")
	if err != nil {
		// If error (e.g. not found), return mock/empty for now
		log.Printf("Failed to read presets.json (using mock): %v", err)
		mock := []Preset{
			{
				ID:       "ft_collins",
				Name:     "Fort Collins, CO",
				ImageURL: "https://storage.googleapis.com/generative-bazaar-001-banana-weather/image_1764375772967339950.png", // Example from logs
				VideoURL: "https://storage.googleapis.com/generative-bazaar-001-banana-weather/videos/535901273950979597/sample_0.mp4",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mock)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) HandleGetWeather(w http.ResponseWriter, r *http.Request) {
	// Check for SSE support
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Helper to send SSE events
	sendEvent := func(event string, data string) {
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
		flusher.Flush()
	}

	city := r.URL.Query().Get("city")
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	var formattedCity string
	var err error

	log.Printf("Received weather request. City: %s, Lat: %s, Lng: %s", city, latStr, lngStr)

	sendEvent("status", "Identifying location...")

	if latStr != "" && lngStr != "" {
		// Handle Coordinates
		var lat, lng float64
		fmt.Sscanf(latStr, "%f", &lat)
		fmt.Sscanf(lngStr, "%f", &lng)
		
		formattedCity, err = h.Maps.GetReverseGeocoding(r.Context(), lat, lng)
		if err != nil {
			log.Printf("Error reverse geocoding: %v", err)
			sendEvent("error", "Failed to resolve location: "+err.Error())
			return
		}
	} else {
		// Handle City Name (or default)
		if city == "" {
			city = "San Francisco"
		}

		// 1. Resolve City
		formattedCity, _, _, err = h.Maps.GetCityLocation(r.Context(), city)
		if err != nil {
			log.Printf("Error resolving location for city '%s': %v", city, err)
			sendEvent("error", "Failed to find city: "+err.Error())
			return
		}
	}
	
	log.Printf("Resolved location to: %s", formattedCity)
	sendEvent("status", "Found location: "+formattedCity)

	// 2. Generate Image
	sendEvent("status", fmt.Sprintf("Getting a banana image of the weather for %s...", formattedCity))
	
	// Use formattedCity to ensure the AI gets the full context (e.g. "Paris, TX" vs "Paris, France")
	imgBase64, err := h.GenAI.GenerateImage(r.Context(), formattedCity, "")
	if err != nil {
		log.Printf("Error generating image for '%s': %v", formattedCity, err)
		sendEvent("error", "Failed to generate image: "+err.Error())
		return
	}
	log.Printf("Successfully generated image for: %s", formattedCity)

	resp := WeatherResponse{
		City:        formattedCity,
		ImageBase64: imgBase64,
	}

	jsonData, _ := json.Marshal(resp)
	sendEvent("result", string(jsonData))

	// 3. Generate Video (If Storage is available)
	if h.Storage == nil {
		log.Printf("Storage service not available, skipping video generation.")
		return
	}

	sendEvent("status", "Preparing for animation...")

	// Upload Image
	fileName := fmt.Sprintf("image_%d.png", time.Now().UnixNano())
	gsURI, _, err := h.Storage.UploadImage(r.Context(), imgBase64, fileName)
	if err != nil {
		log.Printf("Failed to upload image for video gen: %v", err)
		// Don't send error event, just stop, as user already has image
		return
	}

	sendEvent("status", "Animating (Veo 3.1)... this may take a minute.")

	// Call Veo (Returns GS URI now)
	prompt := "The camera moves in parallax as the elements in the image move naturally, while the forecast data—the bold title remain fixed."
	videoGsURI, err := h.GenAI.GenerateVideo(r.Context(), gsURI, prompt)
	if err != nil {
		log.Printf("Veo generation failed: %v", err)
		sendEvent("error", "Video generation failed (Beta). Enjoy the image!")
		return
	}

	sendEvent("status", "Finalizing video...")

	// Convert gs://bucket/path to https://storage.googleapis.com/bucket/path
	// Simple replacement works because the bucket is configured as public.
	// Format: gs://bucket/videos/xyz.mp4 -> https://storage.googleapis.com/bucket/videos/xyz.mp4
	// Note: If the output URI format is different, we might need parsing. 
	// Assuming standard OutputGCSURI behavior.
	
	// Safer: Use storage client helper? Or just string replace.
	// String replace is robust for standard GCS.
	publicVideoURL := "https://storage.googleapis.com/" + videoGsURI[5:] // Strip gs://

	log.Printf("Video available at: %s", publicVideoURL)
	sendEvent("video", publicVideoURL)
}
