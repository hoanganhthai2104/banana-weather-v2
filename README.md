# Banana Weather: AI-Powered Atmospheric Art

![Banana Weather Banner](frontend/assets/images/placeholder_vertical.png)

**Banana Weather** is a "GenMedia" web application that visualizes the current "vibe" and atmospheric essence of a location using Generative AI. 

It combines precise Geolocation with the creative power of **Google Gemini 3 Pro Image** (Nano Banana Pro) to generate high-fidelity, vertical (9:16) 3D isometric art representing the weather, architecture, and mood of your city in real-time.

## Features

*   **AI-Generated Atmospheric Art:** unique, non-deterministic visuals for every request.
*   **Smart Geolocation:** Automatically detects your location or resolves "City, State" from manual input using Google Maps.
*   **Google Search Grounding:** The AI retrieves real-time weather data (temperature, conditions) during image generation.
*   **Responsive Flutter Web UI:** Designed for a mobile-first, full-screen experience.

## Architecture

The system follows a Client-Server architecture:

1.  **The Backend (Go 1.25):** 
    *   **Orchestrator:** Handles API requests, validating inputs.
    *   **Geocoding:** Interacts with Google Maps Platform to resolve coordinates and clean city names.
    *   **Generative AI:** Constructs the prompt and invokes Vertex AI (`gemini-3-pro-image-preview`) with Google Search grounding.
    *   **Server:** Serves the compiled Flutter Web application.

2.  **The Frontend (Flutter Web):**
    *   **UI:** A single-page application using `Provider` for state management.
    *   **Visuals:** Displays the generated art in a 9:16 aspect ratio with a manual override for city selection.

## Setup & Deployment

### 1. Prerequisites
*   **Google Cloud Project** with the following APIs enabled:
    *   Vertex AI API
    *   Google Maps Geocoding API
    *   Cloud Run Admin API
    *   Cloud Logging API
    *   Cloud Storage API
*   **Google Cloud CLI (`gcloud`)** installed and authenticated.

### 2. Environment Configuration
Create a `.env` file in the project root:

```bash
GOOGLE_CLOUD_PROJECT="your-gcp-project-id"
GOOGLE_CLOUD_LOCATION="global" 
GOOGLE_MAPS_API_KEY="your-maps-api-key"
PORT=8080
GENMEDIA_BUCKET="your-gcs-bucket-name" # Optional (for future history features)
```

### 3. Infrastructure Setup

#### Service Account
Create a dedicated Service Account with Least Privilege (Vertex AI User, Logging Writer).
```bash
./setup_sa.sh
```

#### GCS Bucket (Optional/Future)
If using storage features, configure your bucket for CORS (to allow the Web App to display images) and Permissions.

```bash
# 1. Create Bucket
gcloud storage buckets create gs://your-bucket-name --project=your-project --location=us-central1

# 2. Configure CORS (using provided cors.json)
gcloud storage buckets update gs://your-bucket-name --cors-file=cors.json

# 3. Grant Permissions (Public Read, Service Account Write)
gcloud storage buckets add-iam-policy-binding gs://your-bucket-name --member=allUsers --role=roles/storage.objectViewer
gcloud storage buckets add-iam-policy-binding gs://your-bucket-name --member="serviceAccount:your-sa-email" --role=roles/storage.objectAdmin
```

### 4. Local Development
Use the developer helper script to build the frontend and run the backend.

```bash
./dev.sh
# Use --quick to skip rebuilding the frontend if only backend changed
./dev.sh --quick
```

### 5. Deployment
Deploy to **Google Cloud Run**.

1.  Ensure you have run `./setup_sa.sh`.
2.  Open `deploy.sh` and uncomment the `ARGS+=( "--service-account" ... )` line to use your secured Service Account.
3.  Run the deployment:

```bash
./deploy.sh
```

## Technology Stack

| Component | Tech | Key Libraries |
| :--- | :--- | :--- |
| **Backend** | Go | `chi`, `google.golang.org/genai`, `googlemaps.github.io/maps` |
| **Frontend** | Flutter | `provider`, `geolocator`, `google_fonts` |
| **AI Model** | Vertex AI | `gemini-3-pro-image-preview` (Nano Banana Pro) |
| **Infra** | GCP | Cloud Run, Cloud Storage |

---
*Note: This application is a demo of the "GenMedia" pattern.*
