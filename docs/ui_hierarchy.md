# Flutter UI Hierarchy

## Main Structure

*   **MaterialApp** (Root)
    *   **Theme:** Light/Dark support (Lato font).
    *   **MultiProvider** (State Management)
        *   `ChangeNotifierProvider<WeatherProvider>`: Handles API logic, polling, and state.
        *   `ChangeNotifierProvider<ThemeProvider>`: Handles theme toggling.
    *   **HomeScreen** (Main View)

## HomeScreen Widget Tree (`lib/screens/home_screen.dart`)

*   **Scaffold**
    *   `backgroundColor`: Black (or Themed).
    *   **AppBar**
        *   **Title**: "Banana Weather 🍌🌤️"
        *   **Actions**: Theme Toggle Icon.
    *   **Body**: `Center` -> `AspectRatio` (9:16)
        *   **Container** (Main Frame)
            *   `color`: Black (Ensures no white flash during transitions).
            *   `clipBehavior`: HardEdge (Rounded corners).
            *   **Stack** (Layers)
                1.  **Image Layer** (`Positioned.fill`)
                    *   *Condition:* `weatherProvider.imageBase64 != null`
                    *   **Image.memory** (Base64 decoded generated image).
                    *   *Else:* **Image.asset** (`placeholder_vertical.png`).
                2.  **Video Layer** (`Positioned.fill`)
                    *   *Condition:* `_videoController != null` AND `initialized`.
                    *   **FittedBox** -> **SizedBox** -> **VideoPlayer** (Looping, Muted).
                3.  **Loading Overlay** (`Center`) - *Blocking*
                    *   *Condition:* `weatherProvider.isLoading` (Initial fetch).
                    *   **Container** (Dark Glass) -> `Column` -> `CircularProgressIndicator` + Text.
                4.  **Status Pill** (`Positioned` Bottom) - *Non-Blocking*
                    *   *Condition:* `!isLoading` AND `statusMessage != null` (Video Generation).
                    *   **Container** (Pill Shape, Semi-transparent).
                    *   **Row** -> `CircularProgressIndicator` (Small) + `Text` (Cycling messages e.g. "Teaching pixels to dance...").
                5.  **Error Banner** (`Positioned` Top)
                    *   *Condition:* `weatherProvider.error != null`.
                    *   **Container** (Red) -> **Row** -> `Text` (Error) + `IconButton` (Close).

## State Management

### `WeatherProvider` (`lib/providers/weather_provider.dart`)
*   **State:**
    *   `city` (String?): Resolved city name.
    *   `imageBase64` (String?): Generated image.
    *   `videoUrl` (String?): Public URL of generated video.
    *   `isLoading` (bool): True only during initial location/image fetch.
    *   `statusMessage` (String?): Streaming updates from backend ("Animating...").
    *   `error` (String?): Error messages.
*   **Methods:**
    *   `fetchWeather({city, lat, lng})`: Connects to Backend SSE stream.
    *   `clearError()`: Dismisses error.

### `ThemeProvider` (`lib/providers/theme_provider.dart`)
*   **State:** `ThemeMode` (Light/Dark).
*   **Methods:** `toggleTheme()`.