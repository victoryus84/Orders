import 'dart:convert';
import 'package:http/http.dart' as http;
import 'auth_service.dart'; 
import '../core/constants.dart';
import '../core/logger.dart'; 
import '../models/client.dart';
import '../models/contract.dart';


class ApiService {
  Future<List<Client>> searchClients(String query) async {
    myLog(
      "1. Încep căutarea pentru: $query",
    ); // Să vedem dacă funcția e apelată

    if (query.length < 3) {
      myLog("2. Query prea scurt");
      return [];
    }
    
    try {
      myLog("3. Trimit cerere la: ${AppConfig.clientsSearchEndpoint}?q=$query");

      final response = await http.get(
        Uri.parse('${AppConfig.clientsSearchEndpoint}?q=$query'),
        headers: AuthService.getHeaders(),
      );

      myLog("4. Status Code: ${response.statusCode}");
      myLog("5. Body: ${response.body}"); // Aici vedem dacă vin datele

      if (response.statusCode == 200) {
        List<dynamic> data = json.decode(response.body);
        return data.map((json) => Client.fromJson(json)).toList();
      }
    } catch (e) {
      myLog(
        "EROARE GRAVĂ:",
        error: e,
      ); // Aici ne va zice dacă e problemă de rețea sau SSL
    }
    return [];
  }

  Future<List<Contract>> fetchContracts(String clientId) async {
    try {
      final response = await http.get(
        Uri.parse('${AppConfig.baseUrl}/clients/$clientId/contracts'),
        headers: AuthService.getHeaders(),
      );
      if (response.statusCode == 200) {
        List<dynamic> data = json.decode(response.body);
        return data.map((json) => Contract.fromJson(json)).toList();
      }
    } catch (e) {
      myLog("Eroare API Contracte", error: e);
    }
    return [];
  }
}
