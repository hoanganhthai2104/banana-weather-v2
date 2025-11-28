# Flutter UI Hierarchy

## Main Structure

*   **MaterialApp** (Root)
    *   **MultiProvider** (State Management)
        *   `ChangeNotifierProvider<WeatherProvider>`
    *   **HomeScreen** (Main View)

## HomeScreen Widget Tree

*   **Scaffold**
    *   `backgroundColor`: Black
    *   **Center**
        *   **AspectRatio** (9:16) - Enforces the vertical "Phone Wallpaper" ratio.
            *   **Container** (Main Card)
                *   `decoration`: Grey bg, rounded corners, shadow.
                *   `clipBehavior`: HardEdge (clips content to rounded corners).
                *   **Stack** (Layers)
                    1.  **Image Layer** (`Positioned.fill`)
                        *   *Condition:* `weatherProvider.imageBase64 != null`
                        *   **Image.memory** (Base64 decoded image)
                        *   *Else:* **Icon** (Image Not Supported placeholder)
                    2.  **Loading Layer** (`Positioned.fill`)
                        *   *Condition:* `weatherProvider.isLoading || _isInitializing`
                        *   **Center** -> **CircularProgressIndicator**
                    3.  **Gradient Overlay** (`Positioned` bottom)
                        *   **Container** with `LinearGradient` (Transparent -> Black) for text readability.
                    4.  **City Name** (`Positioned` top)
                        *   *Condition:* `weatherProvider.city != null`
                        *   **Text** (Custom Font: `Cinzel`)
                    5.  **Error Message** (`Positioned` top-middle)
                        *   *Condition:* `weatherProvider.error != null`
                        *   **Container** (Red bg) -> **Text**
                    6.  **Input Controls** (`Positioned` bottom)
                        *   **Row**
                            *   **Expanded** -> **TextField** (City Input)
                                *   `controller`: `_controller`
                                *   `onSubmitted`: Triggers fetch.
                            *   **SizedBox** (Spacer)
                            *   **FloatingActionButton**
                                *   `onPressed`: Triggers fetch.
                                *   `child`: Arrow Icon.

## State Management (WeatherProvider)

*   `city` (String?): Current city name.
*   `imageBase64` (String?): Generated image data.
*   `isLoading` (bool): Loading state.
*   `error` (String?): Error message.
*   `fetchWeather({city, lat, lng})`: Async method to update state.
