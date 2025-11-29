package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"banana-weather/pkg/genai"
	"banana-weather/pkg/storage"
	"github.com/joho/godotenv"
)

type Preset struct {

	ID       string `json:"id"`

	Name     string `json:"name"`

	Category string `json:"category"`

	ImageURL string `json:"image_url"`

	VideoURL string `json:"video_url"`

}



func main() {

	// Load .env

	_ = godotenv.Load("../../.env") 

	_ = godotenv.Load("../.env")

	_ = godotenv.Load(".env")



	csvPath := flag.String("csv", "", "Path to CSV file (format: id,name,city,category,context)")

	force := flag.Bool("force", false, "Force overwrite existing presets")

	

	// Single mode flags

	city := flag.String("city", "", "City name")

	ctxPrompt := flag.String("context", "", "Extra prompt context")

	name := flag.String("name", "", "Display name")

	category := flag.String("category", "General", "Category name")

	id := flag.String("id", "", "Unique ID")

	

	flag.Parse()



	ctx := context.Background()



	// Init Services

	genaiService, err := genai.NewService(ctx)

	if err != nil {

		log.Fatalf("Failed to init GenAI: %v", err)

	}

	storageService, err := storage.NewService(ctx)

	if err != nil {

		log.Fatalf("Failed to init Storage: %v", err)

	}



		// Pre-load existing registry



		existingPresets := loadRegistry(ctx, storageService)



		existingMap := make(map[string]Preset)



		for _, p := range existingPresets {



			existingMap[p.ID] = p



		}



	



		var newPresets []Preset



	



		if *csvPath != "" {



			// Batch Mode



			log.Printf("Running in Batch Mode from %s (Force: %v)", *csvPath, *force)



			f, err := os.Open(*csvPath)



			if err != nil {



				log.Fatalf("Failed to open CSV: %v", err)



			}



			defer f.Close()



	



			reader := csv.NewReader(f)



			records, err := reader.ReadAll()



			if err != nil {



				log.Fatalf("Failed to read CSV: %v", err)



			}



	



			for i, row := range records {



				if i == 0 { continue } // Skip Header



				if len(row) < 4 { continue }



				



				pID := row[0]



				pName := row[1]



				pCity := row[2]



				pCat := row[3]



				pCtx := ""



				if len(row) > 4 { pCtx = row[4] }



	



				if existing, ok := existingMap[pID]; ok && !*force {



					log.Printf("Skipping generation for [%s], updating metadata only.", pID)



					// Patch metadata, preserve URLs



					patched := existing



					patched.Name = pName



					patched.Category = pCat



					newPresets = append(newPresets, patched)



					continue



				}



	



				log.Printf("Processing [%d/%d]: %s (%s)", i, len(records)-1, pName, pID)



				p, err := processPreset(ctx, genaiService, storageService, pID, pName, pCity, pCat, pCtx)



				if err != nil {



					log.Printf("Error processing %s: %v", pID, err)



					continue



				}



				newPresets = append(newPresets, *p)



			}



	



		} else {



			// Single Mode



			if *city == "" || *name == "" || *id == "" {



				log.Fatal("Missing required flags: -city, -name, -id (or -csv)")



			}



			



			if existing, ok := existingMap[*id]; ok && !*force {



				log.Printf("Skipping generation for [%s], updating metadata only.", *id)



				patched := existing



				patched.Name = *name



				patched.Category = *category



				newPresets = append(newPresets, patched)



			} else {



				p, err := processPreset(ctx, genaiService, storageService, *id, *name, *city, *category, *ctxPrompt)



				if err != nil {



					log.Fatalf("Error: %v", err)



				}



				newPresets = append(newPresets, *p)



			}



		}



	// Update JSON

	if len(newPresets) > 0 {

		updateRegistry(ctx, storageService, newPresets, *force)

	}

}



func processPreset(ctx context.Context, gs *genai.Service, ss *storage.Service, id, name, city, category, promptCtx string) (*Preset, error) {

	// 1. Generate Image

	log.Printf("Generating image for '%s'...", city)

	imgBase64, err := gs.GenerateImage(ctx, city, promptCtx)

	if err != nil {

		return nil, fmt.Errorf("image gen failed: %w", err)

	}



	// 2. Upload Image

	imgFileName := fmt.Sprintf("preset_%s_image_%d.png", id, time.Now().Unix())

	gsImageURI, publicImageURL, err := ss.UploadImage(ctx, imgBase64, imgFileName)

	if err != nil {

		return nil, fmt.Errorf("image upload failed: %w", err)

	}

	log.Printf("Image uploaded: %s", publicImageURL)



	// 3. Generate Video

	log.Printf("Generating video (Veo)...")

	videoPrompt := "The camera moves in parallax as the elements in the image move naturally, while the forecast data—the bold title remain fixed."

	videoGsURI, err := gs.GenerateVideo(ctx, gsImageURI, videoPrompt)

	if err != nil {

		return nil, fmt.Errorf("video gen failed: %w", err)

	}

	

	bucketName := os.Getenv("GENMEDIA_BUCKET")

	publicVideoURL := strings.Replace(videoGsURI, "gs://"+bucketName, "https://storage.googleapis.com/"+bucketName, 1)

	log.Printf("Video generated: %s", publicVideoURL)



	return &Preset{

		ID:       id,

		Name:     name,

		Category: category,

		ImageURL: publicImageURL,

		VideoURL: publicVideoURL,

	},

	nil

}



func loadRegistry(ctx context.Context, ss *storage.Service) []Preset {

	var presets []Preset

	data, err := ss.ReadObject(ctx, "presets.json")

	if err == nil {

		json.Unmarshal(data, &presets)

	}

	return presets

}



func updateRegistry(ctx context.Context, ss *storage.Service, newItems []Preset, force bool) {

	log.Printf("Updating presets.json...")

	existing := loadRegistry(ctx, ss)

	

	// Merge Strategy

	// If force is true, newItems replace existing ones with same ID.

	// Since we already filtered based on force in main loop, newItems contains ONLY things we want to write.

	// BUT we need to merge them into the existing list (replacing if match).

	

	finalMap := make(map[string]Preset)

	for _, p := range existing {

		finalMap[p.ID] = p

	}

	for _, p := range newItems {

		finalMap[p.ID] = p

	}

	

	var final []Preset

	for _, p := range finalMap {

		final = append(final, p)

	}

	

	// Sort? Map iteration is random. 

	// Ideally keep order or sort by Name/ID.

	// Let's just write.



	newData, _ := json.MarshalIndent(final, "", "  ")

	_, err := ss.UploadBytes(ctx, newData, "presets.json", "application/json")

	if err != nil {

		log.Fatalf("Failed to save presets.json: %v", err)

	}

	log.Printf("Registry updated successfully.")

}
