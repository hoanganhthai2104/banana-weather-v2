# Banana Weather: AI-Powered Atmospheric Art

![Banana Weather Banner](frontend/assets/images/placeholder_vertical.png)

**Banana Weather** is a "GenMedia" web application that visualizes the current "vibe" and atmospheric essence of a location using Generative AI. 

It combines precise Geolocation with the creative power of **Google Gemini 3 Pro Image** (Nano Banana Pro) and **Google Veo 3.1** to generate high-fidelity, vertical (9:16) 3D isometric art and looping videos representing the weather, architecture, and mood of your city in real-time.

## Features

*   **AI-Generated Atmospheric Art:** Unique, non-deterministic visuals for every request.
*   **Cinematic Video Loops:** Transitions from static image to a "Parallax" animation using **Veo 3.1 Fast**.
*   **Smart Geolocation:** Automatically detects your location or resolves "City, State" from manual input using Google Maps.
*   **Fictional Locations:** Supports generating scenes for fictional worlds (e.g., Arrakis, Middle-earth) via the Presets system.
*   **Presets Gallery:** A curated list of pre-generated scenes categorized by theme (e.g., Star Wars, Landmarks), accessible via the side menu.
*   **Responsive Flutter Web UI:** Mobile-first design with a clean, "Digital Picture Frame" aesthetic.

## Architecture

The system follows a Client-Server architecture:

1.  **The Backend (Go 1.25):** 
    *   **Orchestrator:** Handles API requests, validating inputs.
    *   **Geocoding:** Interacts with Google Maps Platform.
    *   **Generative AI:**
        *   **Image:** Vertex AI `gemini-3-pro-image-preview` (GenerateContent).
        *   **Video:** Vertex AI `veo-3.1-fast-generate-preview` (GenerateVideos + Polling).
    *   **Storage:** Manages GCS uploads for images, videos, and `presets.json`.
    *   **Server:** Serves the compiled Flutter Web application.

2.  **The Frontend (Flutter Web):**
    *   **UI:** Single-page app using `Provider` and `ChangeNotifier`.
    *   **Video:** `video_player` integration with seamless transition.
    *   **Theme:** Light/Dark mode with dynamic "Glassmorphism" overlays.
    *   **Navigation:** Drawer-based Preset Gallery with categories.

## Setup & Deployment

### 1. Prerequisites
*   **Google Cloud Project** with APIs enabled: Vertex AI, Maps, Cloud Run, GCS.
*   **GCS Bucket:** Publicly readable with CORS configured.

### 2. Environment Configuration
Create a `.env` file:

```bash
GOOGLE_CLOUD_PROJECT="your-gcp-project-id"
GOOGLE_CLOUD_LOCATION="global" 
GOOGLE_MAPS_API_KEY="your-maps-api-key"
GENMEDIA_BUCKET="your-gcs-bucket-name"
PORT=8080
```

### 3. Development
*   **Run Local:** `./dev.sh`
*   **Deploy:** `./deploy.sh`

### 4. Tools
**Preset Generator:**
Populate the gallery with pre-generated content.

```bash
cd backend
go run cmd/generate_preset/main.go -csv presets.csv
```

## Coding Conventions
*   **Frontend:** Use `kDebugMode` for API URLs. Use `AnimatedOpacity` for non-blocking UI transitions.
*   **Backend:** Use `aiplatform` SDK for LRO polling. Store assets in GCS.
