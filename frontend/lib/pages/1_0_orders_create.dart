import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:frontend/core/constants.dart';
import 'package:frontend/models/client.dart';
import 'package:frontend/services/auth_service.dart';

class OrdersCreatePage extends StatefulWidget {
  final String title;
  const OrdersCreatePage({super.key, required this.title});

  @override
  State<OrdersCreatePage> createState() => _OrdersCreatePageState();
}

class _OrdersCreatePageState extends State<OrdersCreatePage> {
  // 1. VARIABILE DE STARE (Bordul de control)
  Client? _selectedClient;
  String? _selectedPaymentType;
  final ValueNotifier<bool> _loadingNotifier = ValueNotifier<bool>(false);

  // 2. LOGICA DE CĂUTARE CLIENȚI (GO BACKEND)
  Future<List<Client>> _searchClientsFromGo(String query) async {
    if (query.length < 3) return [];
    _loadingNotifier.value = true;

    
    try {
      debugPrint("HEADERS TRIMISE: ${AuthService.getHeaders()}");
      final response = await http.get(
        Uri.parse('${AppConfig.clientsSearchEndpoint}?q=$query'),
        headers: AuthService.getHeaders(),
      ).timeout(  
        const Duration(milliseconds: AppConfig.apiTimeout),
        onTimeout: () {
          debugPrint("REQUEST TIMEOUT!");
          throw Exception("Timeout la căutare clienți");
        },
      );  
      if (response.statusCode == 200) {
        debugPrint("DATE PRIMITE DE LA SERVER: ${response.body}");
        List<dynamic> data = json.decode(response.body);
        return data.map((json) => Client.fromJson(json)).toList();
      }
    } catch (e) {
      debugPrint("Eroare la căutare: $e");
    } finally {
      _loadingNotifier.value = false;
    }
    return [];
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: SingleChildScrollView(
        // Ca să nu avem erori de spațiu când apare tastatura
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // --- CÂMPUL 1: TIP PLATĂ (Нал / Бнал) ---
            const Text(
              "Выберите тип оплаты:",
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            DropdownButtonFormField<String>(
              initialValue: _selectedPaymentType,
              decoration: InputDecoration(
                prefixIcon: const Icon(Icons.payments_outlined),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              hint: const Text("Нал / Бнал"),
              items: const [
                DropdownMenuItem(value: "Нал", child: Text("Наличные (Нал)")),
                DropdownMenuItem(
                  value: "Бнал",
                  child: Text("Безналичные (Бнал)"),
                ),
              ],
              onChanged: (val) => setState(() => _selectedPaymentType = val),
            ),

            const SizedBox(height: 24),

            // --- CÂMPUL 2: CĂUTARE CLIENT ---
            const Text(
              "Выберите клиента:",
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            Autocomplete<Client>(
              displayStringForOption: (Client c) => c.name,
              optionsBuilder: (textValue) =>
                  _searchClientsFromGo(textValue.text),
              onSelected: (selection) =>
                  setState(() => _selectedClient = selection),
              fieldViewBuilder: (context, controller, focusNode, onSubmitted) {
                return ValueListenableBuilder<bool>(
                  valueListenable: _loadingNotifier,
                  builder: (context, isLoading, _) {
                    return TextFormField(
                      controller: controller,
                      focusNode: focusNode,
                      decoration: InputDecoration(
                        hintText: "Введите min 3 символа...",
                        prefixIcon: const Icon(Icons.person_search),
                        suffixIcon: isLoading
                            ? const Padding(
                                padding: EdgeInsets.all(12),
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                ),
                              )
                            : const Icon(Icons.arrow_drop_down),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                    );
                  },
                );
              },
            ),

            const SizedBox(height: 40),

            // --- BUTONUL DE SALVARE ---
            ElevatedButton(
              onPressed:
                  (_selectedClient != null && _selectedPaymentType != null)
                  ? () {
                      /* Aici va veni funcția de Submit */
                    }
                  : null,
              style: ElevatedButton.styleFrom(
                minimumSize: const Size.fromHeight(50),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              child: const Text("СОЗДАТЬ ЗАКАЗ"),
            ),
          ],
        ),
      ),
    );
  }
}
