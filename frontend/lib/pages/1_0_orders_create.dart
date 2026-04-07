import 'package:flutter/material.dart';
import '../models/client.dart';
import '../models/contract.dart';
import '../services/api_service.dart';
import '../controllers/order_controller.dart';
import '../core/logger.dart';

class OrdersCreatePage extends StatefulWidget {
  final String title;
  const OrdersCreatePage({super.key, required this.title});

  @override
  State<OrdersCreatePage> createState() => _OrdersCreatePageState();
}

class _OrdersCreatePageState extends State<OrdersCreatePage> {
  final OrderCreateController _controller = OrderCreateController();
  final ApiService _api = ApiService();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: ListenableBuilder(
        listenable: _controller,
        builder: (context, _) {
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildPaymentDropdown(),
                const SizedBox(height: 24),
                _buildClientAutocomplete(),
                const SizedBox(height: 24),
                // ADAUGĂM UN KEY: Acest lucru forțează resetarea dropdown-ului
                // de contracte când se schimbă clientul selectat.
                _buildContractDropdown(
                  key: ValueKey(_controller.selectedClient?.id ?? 'none'),
                ),
                const SizedBox(height: 40),
                _buildSubmitButton(),
              ],
            ),
          );
        },
      ),
    );
  }

  // --- COMPONENTELE EXTRASE ---

  Widget _buildPaymentDropdown() {
    return DropdownButtonFormField<String>(
      // REZOLVARE DEPRECATION: Folosim initialValue în loc de value
      initialValue: _controller.paymentType,
      decoration: const InputDecoration(
        labelText: "Tip Plată",
        border: OutlineInputBorder(),
      ),
      items: const [
        DropdownMenuItem(value: "Нал", child: Text("Наличные (Нал)")),
        DropdownMenuItem(value: "Бнал", child: Text("Безналичные (Бнал)")),
      ],
      onChanged: (val) => _controller.setPaymentType(val),
    );
  }

Widget _buildClientAutocomplete() {
    return Autocomplete<Client>(
      // 1. Aici transformăm obiectul în text pentru listă
      displayStringForOption: (Client c) => c.name.toString(),

      // 2. Aici se întâmplă "magia" căutării
      optionsBuilder: (TextEditingValue textValue) async {
        if (textValue.text.length < 3) return const Iterable<Client>.empty();

        // 1. Luăm Messenger-ul ÎNAINTE de await (aici context e sigur valid)
        final messenger = ScaffoldMessenger.of(context);

        // 2. Facem cererea la server
        final rezultate = await _api.searchClients(textValue.text);

        // 3. Verificăm dacă pagina mai e pe ecran (ca să nu facem prostii)
        if (!mounted) return rezultate;

        // 4. Folosim 'messenger' (variabila salvată), NU mai scriem 'context' aici!
        messenger.clearSnackBars();
        messenger.showSnackBar(
          SnackBar(
            content: Text(
              rezultate.isEmpty
                  ? "Nimic găsit pentru '${textValue.text}'"
                  : "Am găsit ${rezultate.length} clienți",
            ),
            backgroundColor: rezultate.isEmpty ? Colors.orange : Colors.blue,
            duration: const Duration(milliseconds: 800),
          ),
        );

        return rezultate;
      },

      // 3. Ce se întâmplă când dai click pe un client din listă
      onSelected: (Client selection) {
        _controller.selectClient(selection);
        myLog("Client selectat: ${selection.name}");
      },

      // 4. Cum arată câmpul unde scrii
      fieldViewBuilder: (context, textController, focusNode, onFieldSubmitted) {
        return TextFormField(
          controller: textController,
          focusNode: focusNode,
          decoration: const InputDecoration(
            labelText: "Selectați Clientul",
            prefixIcon: Icon(Icons.person_search),
            border: OutlineInputBorder(),
          ),
        );
      },
    );
  }

  Widget _buildContractDropdown({Key? key}) {
    return DropdownButtonFormField<Contract>(
      key: key, // Folosim cheia pentru resetare
      // REZOLVARE DEPRECATION: initialValue în loc de value
      initialValue: _controller.selectedContract,
      decoration: InputDecoration(
        labelText: "Contract",
        border: const OutlineInputBorder(),
        suffixIcon: _controller.isLoadingContracts
            ? const SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(strokeWidth: 2),
              )
            : null,
      ),
      hint: Text(
        _controller.selectedClient == null
            ? "Alegeți clientul mai întâi"
            : "Selectați contractul",
      ),
      items: _controller.availableContracts
          .map(
            (c) => DropdownMenuItem(
              value: c,
              // REZOLVARE INT/STRING: toString() forțează conversia
              child: Text(c.name.toString()),
            ),
          )
          .toList(),
      onChanged: _controller.isLoadingContracts
          ? null
          : (val) => _controller.selectContract(val),
    );
  }

  Widget _buildSubmitButton() {
    return ElevatedButton(
      onPressed: _controller.isValid ? () => _handleSave() : null,
      style: ElevatedButton.styleFrom(
        minimumSize: const Size.fromHeight(50),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      ),
      child: const Text("СОЗДАТЬ ЗАКАЗ"),
    );
  }

  void _handleSave() {
    // Folosim string interpolation "${...}" pentru a evita erorile de tip la print/log
    debugPrint("Salvare comandă pentru: ${_controller.selectedClient?.name}");
  }
}
