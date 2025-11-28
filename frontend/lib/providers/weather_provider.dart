import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class WeatherProvider with ChangeNotifier {
  String? _city;
  String? _imageBase64;
  bool _isLoading = false;
  String? _error;

  String? get city => _city;
  String? get imageBase64 => _imageBase64;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> fetchWeather({String? city, double? lat, double? lng}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final Uri uri;
      if (lat != null && lng != null) {
        uri = Uri.parse('http://localhost:8080/api/weather?lat=$lat&lng=$lng');
      } else if (city != null && city.isNotEmpty) {
        uri = Uri.parse('http://localhost:8080/api/weather?city=$city');
      } else {
        // Default to SF if nothing provided
         uri = Uri.parse('http://localhost:8080/api/weather?city=San Francisco');
      }
      
      final response = await http.get(uri);

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _city = data['city'];
        _imageBase64 = data['image_base64'];
      } else {
        _error = 'Failed to load data: ${response.statusCode}';
      }
    } catch (e) {
      _error = 'Error: $e';
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
