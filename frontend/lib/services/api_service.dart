import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'auth_service.dart';
import '../core/constants.dart';
import '../core/logger.dart';
import '../models/client.dart';
import '../models/contract.dart';

class ApiService {
  // 1. CĂUTARE CLIENȚI
  Future<List<Client>> searchClients(String query) async {
    myLog("🔍 Căutăm clienți: $query");

    if (query.length < 3) return [];

    try {
      final url = Uri.parse('${AppConfig.clientsSearchEndpoint}?q=$query');
      final response = await http.get(url, headers: AuthService.getHeaders());
      debugPrint("🚀 JSON RAW DE LA SERVER: ${response.body}");
      if (response.statusCode == 200) {
        List<dynamic> data = json.decode(response.body);
        return data.map((json) => Client.fromJson(json)).toList();
      }
    } catch (e) {
      myLog("❌ Eroare API Căutare: $e");
    }
    return [];
  }

  // 2. FETCH CONTRACTE (Piesa care lipsea!)
  Future<List<Contract>> fetchContracts(String clientId) async {
    myLog("📄 Cerem contractele pentru clientul: $clientId");

    try {
      final url = Uri.parse('${AppConfig.baseUrl}/contracts/client/$clientId');
      final response = await http.get(url, headers: AuthService.getHeaders());

      if (response.statusCode == 200) {
        List<dynamic> data = json.decode(response.body);
        return data.map((json) => Contract.fromJson(json)).toList();
      }
    } catch (e) {
      myLog("❌ Eroare API Contracte: $e");
    }
    return [];
  }
}
