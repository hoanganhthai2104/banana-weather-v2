import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:geolocator/geolocator.dart';
import '../providers/weather_provider.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final TextEditingController _controller = TextEditingController();
  bool _isInitializing = true;

  @override
  void initState() {
    super.initState();
    _determinePosition();
  }

  Future<void> _determinePosition() async {
    final weatherProvider = Provider.of<WeatherProvider>(context, listen: false);
    
    try {
      bool serviceEnabled;
      LocationPermission permission;

      // Test if location services are enabled.
      serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        await weatherProvider.fetchWeather(city: 'San Francisco');
        return;
      }

      permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
        if (permission == LocationPermission.denied) {
          await weatherProvider.fetchWeather(city: 'San Francisco');
          return;
        }
      }
      
      if (permission == LocationPermission.deniedForever) {
        await weatherProvider.fetchWeather(city: 'San Francisco');
        return;
      } 

      Position position = await Geolocator.getCurrentPosition();
      await weatherProvider.fetchWeather(lat: position.latitude, lng: position.longitude);
    } catch (e) {
       await weatherProvider.fetchWeather(city: 'San Francisco');
    } finally {
      if (mounted) {
        setState(() {
          _isInitializing = false;
        });
      }
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final weatherProvider = Provider.of<WeatherProvider>(context);

    return Scaffold(
      backgroundColor: Colors.black,
      body: Center(
        child: AspectRatio(
          aspectRatio: 9 / 16,
          child: Container(
            decoration: BoxDecoration(
              color: Colors.grey[900],
              borderRadius: BorderRadius.circular(16),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.5),
                  blurRadius: 20,
                  spreadRadius: 5,
                ),
              ],
            ),
            clipBehavior: Clip.hardEdge,
            child: Stack(
              children: [
                // Image Layer
                if (weatherProvider.imageBase64 != null)
                  Positioned.fill(
                    child: Image.memory(
                      base64Decode(weatherProvider.imageBase64!),
                      fit: BoxFit.cover,
                    ),
                  )
                else
                  Positioned.fill(
                    child: Image.asset(
                      'assets/images/placeholder_vertical.png',
                      fit: BoxFit.cover,
                    ),
                  ),

                // Loading Indicator
                if (weatherProvider.isLoading || _isInitializing)
                  const Positioned.fill(
                    child: Center(
                      child: CircularProgressIndicator(
                        color: Colors.yellowAccent,
                      ),
                    ),
                  ),

                // Gradient Overlay for Text Readability
                Positioned(
                  bottom: 0,
                  left: 0,
                  right: 0,
                  height: 200,
                  child: Container(
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        begin: Alignment.topCenter,
                        end: Alignment.bottomCenter,
                        colors: [
                          Colors.transparent,
                          Colors.black.withOpacity(0.8),
                        ],
                      ),
                    ),
                  ),
                ),



                // Error Message
                if (weatherProvider.error != null)
                  Positioned(
                    top: 100,
                    left: 20,
                    right: 20,
                    child: Container(
                      padding: const EdgeInsets.all(8),
                      color: Colors.red.withOpacity(0.8),
                      child: Text(
                        weatherProvider.error!,
                        style: const TextStyle(color: Colors.white),
                        textAlign: TextAlign.center,
                      ),
                    ),
                  ),

                // Input Field
                Positioned(
                  bottom: 40,
                  left: 20,
                  right: 20,
                  child: Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: _controller,
                          style: const TextStyle(color: Colors.white),
                          decoration: InputDecoration(
                            hintText: 'Enter City...',
                            hintStyle: TextStyle(color: Colors.white.withOpacity(0.5)),
                            filled: true,
                            fillColor: Colors.white.withOpacity(0.1),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(30),
                              borderSide: BorderSide.none,
                            ),
                            contentPadding: const EdgeInsets.symmetric(
                              horizontal: 20,
                              vertical: 16,
                            ),
                          ),
                          onSubmitted: (value) {
                            if (value.isNotEmpty) {
                              weatherProvider.fetchWeather(city: value);
                            }
                          },
                        ),
                      ),
                      const SizedBox(width: 10),
                      FloatingActionButton(
                        onPressed: () {
                          if (_controller.text.isNotEmpty) {
                            weatherProvider.fetchWeather(city: _controller.text);
                          }
                        },
                        backgroundColor: Colors.yellowAccent,
                        foregroundColor: Colors.black,
                        child: const Icon(Icons.arrow_forward),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
