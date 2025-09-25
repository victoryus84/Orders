import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  String _email = '';
  String _password = '';
  String? _error;

  Future<void> _login() async {
    final url = Uri.parse(
      'http://localhost:8080/login',
    ); // Замените на ваш адрес
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': _email, 'password': _password}),
    );
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final token = data['token'];
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('auth_token', token); // сохраняем токен
      // Сохраните токен, если нужно
      setState(() {
        _error = null;
      });
      // Перенаправление на меню
      if (context.mounted) {
        GoRouter.of(context).go('/0_1_menu');
      }
    } else {
      setState(() {
        _error = 'Неверный логин или пароль';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Вход')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextFormField(
                decoration: const InputDecoration(labelText: 'Email'),
                onChanged: (val) => _email = val,
                validator: (val) =>
                    val == null || val.isEmpty ? 'Введите email' : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                decoration: const InputDecoration(labelText: 'Пароль'),
                obscureText: true,
                onChanged: (val) => _password = val,
                validator: (val) =>
                    val == null || val.isEmpty ? 'Введите пароль' : null,
              ),
              const SizedBox(height: 24),
              if (_error != null)
                Text(_error!, style: const TextStyle(color: Colors.red)),
              ElevatedButton(
                onPressed: () {
                  if (_formKey.currentState!.validate()) {
                    _login();
                  }
                },
                child: const Text('Войти'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
