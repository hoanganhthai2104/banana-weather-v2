package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"banana-weather/pkg/genai"
	"banana-weather/pkg/maps"
)

type Handler struct {
	Maps  *maps.Service
	GenAI *genai.Service
}

type WeatherResponse struct {
	City        string `json:"city"`
	ImageBase64 string `json:"image_base64"`
}

func (h *Handler) HandleGetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	var formattedCity string
	var err error

	log.Printf("Received weather request. City: %s, Lat: %s, Lng: %s", city, latStr, lngStr)

	if latStr != "" && lngStr != "" {
		// Handle Coordinates
		var lat, lng float64
		fmt.Sscanf(latStr, "%f", &lat)
		fmt.Sscanf(lngStr, "%f", &lng)
		
		formattedCity, err = h.Maps.GetReverseGeocoding(r.Context(), lat, lng)
		if err != nil {
			log.Printf("Error reverse geocoding: %v", err)
			http.Error(w, "Failed to resolve location: "+err.Error(), http.StatusBadRequest)
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
			http.Error(w, "Failed to find city: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	
	log.Printf("Resolved location to: %s", formattedCity)

	// 2. Generate Image
	// Use formattedCity to ensure the AI gets the full context (e.g. "Paris, TX" vs "Paris, France")
	imgBase64, err := h.GenAI.GenerateImage(r.Context(), formattedCity)
	if err != nil {
		log.Printf("Error generating image for '%s': %v", formattedCity, err)
		http.Error(w, "Failed to generate image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully generated image for: %s", formattedCity)

	resp := WeatherResponse{
		City:        formattedCity,
		ImageBase64: imgBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
