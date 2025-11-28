package genai

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

type Service struct {
	client *genai.Client
}

func NewService(ctx context.Context) (*Service, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = os.Getenv("PROJECT_ID")
	}
	if projectID == "" {
		return nil, fmt.Errorf("PROJECT_ID or GOOGLE_CLOUD_PROJECT not set")
	}

	location := os.Getenv("GOOGLE_CLOUD_LOCATION")
	if location == "" {
		location = "us-central1"
	}

	log.Printf("Initializing GenAI Service. Project: %s, Location: %s", projectID, location)

	c, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend:  genai.BackendVertexAI,
		Project:  projectID,
		Location: location,
	})
	if err != nil {
		return nil, err
	}
	return &Service{client: c}, nil
}

// GenerateImage generates a 9:16 image for the given city.
func (s *Service) GenerateImage(ctx context.Context, city string) (string, error) {
	prompt := fmt.Sprintf(`Present a clear, 45° top-down view of a vertical (9:16) isometric miniature 3D cartoon scene, highlighting iconic landmarks centered in the composition to showcase precise and delicate modeling.

The scene features soft, refined textures with realistic PBR materials and gentle, lifelike lighting and shadow effects. Weather elements are creatively integrated into the urban architecture, establishing a dynamic interaction between the city's landscape and atmospheric conditions, creating an immersive weather ambiance.

Use a clean, unified composition with minimalistic aesthetics and a soft, solid-colored background that highlights the main content. The overall visual style is fresh and soothing.

Display a prominent weather icon at the top-center, with the date (x-small text) and temperature range (medium text) beneath it. The city name (large text) is positioned directly above the weather icon. The weather information has no background and can subtly overlap with the buildings.

The text should match the input city's native language.
Please retrieve current weather conditions for the specified city before rendering.

City name: %s`, city)

	// Nano Banana Pro corresponds to 'gemini-3-pro-image-preview' or similar.
	model := "gemini-3-pro-image-preview"

	log.Printf("Generating image for city: %s using model: %s (GenerateContent)", city, model)

	resp, err := s.client.Models.GenerateContent(ctx, model, genai.Text(prompt), &genai.GenerateContentConfig{
		ResponseModalities: []string{"IMAGE"},
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		},
	})
	if err != nil {
		log.Printf("GenAI GenerateContent failed: %v", err)
		return "", fmt.Errorf("genai error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Printf("GenAI returned no candidates or parts")
		return "", fmt.Errorf("no content generated")
	}

	// Iterate through parts to find the image
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			log.Printf("Image generated successfully. Bytes: %d", len(part.InlineData.Data))
			return base64.StdEncoding.EncodeToString(part.InlineData.Data), nil
		}
	}
	
	log.Printf("No inline image data found in response")
	return "", fmt.Errorf("no image data found in response")
}

func ptr[T any](v T) *T {
	return &v
}
