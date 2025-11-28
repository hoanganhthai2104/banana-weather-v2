# Banana Weather: AI-Powered Atmospheric Art

*A GenMedia application that visualizes the current "vibe" of a city using Generative AI.*

## Project Overview
Banana Weather is a simple yet powerful application that demonstrates the "Oracle & Temple" architecture. It allows users to select a city and generates a high-fidelity, 9:16 artistic representation of that location's current essence using **Google Gemini 3 Pro Image (Nano Banana Pro)** via Vertex AI.

## Architecture

The system follows a client-server model:

1.  **The Backend (Go):** The orchestration layer.
    *   **Role:** API Server, Static File Host, and GenAI Orchestrator.
    *   **Responsibilities:**
        *   **Geocoding:** Google Maps API for location resolution and reverse geocoding.
        *   **Image Generation:** Vertex AI (`gemini-3-pro-image-preview`) using `GenerateContent`.
        *   **Serving:** Hosts the compiled Flutter Web application.

2.  **The Frontend (Flutter):** The user interface.
    *   **Role:** Web-based client (Single Page App).
    *   **Stack:** Flutter Web, `provider` for state management.
    *   **UI:**
        *   Full-screen 9:16 image display.
        *   Geolocation-based initial fetch.
        *   Manual city search overlay.

## Technology Stack

| Component | Technology | Key Libraries |
| :--- | :--- | :--- |
| **Backend** | Go 1.25+ | `chi`, `google.golang.org/genai`, `googlemaps.github.io/maps` |
| **Frontend** | Flutter (Dart) | `provider`, `google_fonts`, `geolocator` |
| **Cloud** | Google Cloud Run | Vertex AI, Maps Platform |

## Development Guide

### 1. Environment Variables
Create a `.env` file:

```bash
GOOGLE_CLOUD_PROJECT="<your-gcp-project>"
GOOGLE_CLOUD_LOCATION="global" # Required for Gemini 3 Pro Image
GOOGLE_MAPS_API_KEY="<your-maps-key>"
PORT=8080
```

### 2. Running Local (Dev)
Use the helper script to build frontend and run backend:
```bash
./dev.sh
```
*   `./dev.sh --quick` skips the Flutter build if only backend changes were made.

### 3. Deployment
Deploy to Google Cloud Run:
```bash
./deploy.sh
```

## Coding Conventions

### General
*   **Task Management:** Use the `bd` tool for all work items.
    *   **Workflow:** Create (`bd create`), Implement, Close (`bd close`).
    *   **Priorities:** P0 (Critical) to P3 (Nice to have).
*   **Documentation:** Maintain `docs/` and `GEMINI.md` with architectural decisions.

### Go (Backend)
*   **Structure:** `pkg/` for business logic (GenAI, Maps), `api/` for HTTP handlers.
*   **GenAI:** Use `GenerateContent` for multimodal models (even for image generation).
*   **Logging:** Use `log` for tracking request flows and errors.

### Flutter (Frontend)
*   **State:** Use `Provider` pattern.
*   **Location:** Handle permissions gracefully with `geolocator`.
*   **Style:** Dark mode aesthetic, focused on the image.

## Troubleshooting
*   **500 Errors:** Check backend logs.
*   **"Model not found":** Ensure `GOOGLE_CLOUD_LOCATION` is set to `global` for `gemini-3-pro-image-preview`.
