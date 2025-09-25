// Создаие заказа

import 'package:flutter/material.dart';
import 'package:frontend/widgets/back_menu.dart';
import 'package:frontend/widgets/dictionary_dropdown.dart';
import 'package:multi_dropdown/multi_dropdown.dart';

class OrdersCreatePage extends StatefulWidget {
  final String title;
  const OrdersCreatePage({super.key, required this.title});

  @override
  _OrdersCreatePageState createState() => _OrdersCreatePageState();
}

class _OrdersCreatePageState extends State<OrdersCreatePage> {
  @override  
  Widget build(BuildContext context) {

    // контроллер для выбора клиентов
    final paymentController = MultiSelectController<String>();
    final clientController = MultiSelectController<String>();

    // список Оплаты
    final paymentsMethod = ["НАЛ", "БНАЛ"];
    // список клиентов
    final clients = ["Рога и копыта", "Сервидар", "ПроИТ"];

    return Scaffold(
      appBar: AppBar(title: Text(widget.title), actions: [const BackToMenuButton()]),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              //Text("Страница: $title", style: const TextStyle(fontSize: 18)),
              const SizedBox(height: 20),
              TextFormField(
                decoration: const InputDecoration(labelText: "Enter something"),
              ),
              const SizedBox(height: 20),

              DictionaryDropdown<String>(
                hintText: "Выберите способ оплаты",
                items: paymentsMethod,
                controller: paymentController,
                labelBuilder: (c) =>
                    c, // для строк просто возвращаем саму строку
                validator: (selectedItems) {
                  if (selectedItems == null || selectedItems.isEmpty) {
                    return "cпособ обязателен";
                  }
                  return null;
                },
              ),             
              
              const SizedBox(height: 20),
              
              DictionaryDropdown<String>(
                hintText: "Выберите клиента",
                items: clients,
                controller: clientController,
                labelBuilder: (c) =>
                    c, // для строк просто возвращаем саму строку
                validator: (selectedItems) {
                  if (selectedItems == null || selectedItems.isEmpty) {
                    return "Клиент обязателен";
                  }
                  return null;
                },
              ),             

              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text("Form submitted!")),
                  );
                },
                child: const Text("Submit"),
              ),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
  }
}
