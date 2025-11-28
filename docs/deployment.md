# Deployment Guide

## Prerequisites

1.  **Google Cloud Project**
2.  **APIs Enabled:**
    *   Cloud Run API
    *   Vertex AI API
    *   Google Maps Geocoding API
    *   Artifact Registry API
3.  **Google Cloud CLI (`gcloud`)** installed and authenticated.

## Configuration

Ensure you have a `.env` file in the project root:

```bash
GEMINI_API_KEY="your-vertex-ai-key"
GOOGLE_MAPS_API_KEY="your-maps-key"
PROJECT_ID="your-gcp-project-id"
```

## Deployment Steps

1.  **Run the Deployment Script:**
    The `deploy.sh` script handles building the frontend and deploying the container to Cloud Run.

    ```bash
    ./deploy.sh
    ```

2.  **Access the Application:**
    Once deployment is complete, the script will output the Service URL (e.g., `https://banana-weather-xyz.a.run.app`).

## Notes

*   **Service Account:** By default, the script uses the Compute Engine default service account. For better security, create a dedicated service account with specific permissions (Vertex AI User, Maps API User) and update `deploy.sh` to use it.
*   **Region:** Defaults to `us-central1`. Change via `REGION` env var if needed.
