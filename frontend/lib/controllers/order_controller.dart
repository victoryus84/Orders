import 'package:flutter/material.dart';
import '../models/client.dart';
import '../models/contract.dart';
import '../services/api_service.dart';

class OrderCreateController extends ChangeNotifier {
  final ApiService _api = ApiService();

  Client? selectedClient;
  Contract? selectedContract;
  String? paymentType;

  List<Contract> availableContracts = [];
  bool isLoadingContracts = false;

  // Verifică dacă butonul de salvare poate fi apăsat
  bool get isValid =>
      selectedClient != null && selectedContract != null && paymentType != null;

  void setPaymentType(String? type) {
    paymentType = type;
    notifyListeners();
  }

  void selectContract(Contract? contract) {
    selectedContract = contract;
    notifyListeners();
  }

  // LOGICA DE DEPENDENȚĂ: Client -> Contract
  Future<void> selectClient(Client client) async {
    selectedClient = client;
    selectedContract = null; // Resetăm contractul vechi!
    availableContracts = [];
    isLoadingContracts = true;
    notifyListeners();

    availableContracts = await _api.fetchContracts(client.id);

    isLoadingContracts = false;
    notifyListeners();
  }
}
