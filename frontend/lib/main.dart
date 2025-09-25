import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/services/auth_service.dart';
import 'pages/0_0_login.dart';
import 'pages/0_1_menu.dart';
import 'pages/1_0_orders_create.dart';
import 'pages/1_1_orders_list.dart';
import 'pages/2_0_equipment.dart';
import 'pages/3_0_calculations.dart';
import 'pages/4_0_shipment.dart';

void main() {
  runApp(const MyApp());
}

final _router = GoRouter(
  initialLocation: '/0_0_login',
  routes: [
    GoRoute(path: '/0_0_login', builder: (context, state) => const LoginPage()),
    GoRoute(path: '/0_1_menu', builder: (context, state) => const MenuPage()),
    GoRoute(path: '/1_0_orders_create', builder: (context, state) => OrdersCreatePage(title: "Создать заявку"),),
    GoRoute(path: '/1_1_orders_list', builder: (context, state) => const OrdersListPage(title: "Список заявок"),),
    GoRoute(path: '/2_0_equipment', builder: (context, state) => const EquipmentPage(title: "Оборудование"),),
    GoRoute(path: '/3_0_calculations', builder: (context, state) => const CalculationsPage(title: "Взаиморасчеты"),),
    GoRoute(path: '/4_0_shipment', builder: (context, state) => const ShipmentPage(title: "Отгрузка 4 мес."),),
],

/// 🔹 Асинхронный redirect (при старте приложения загружаем токен из памяти)
  redirect: (context, state) {
    final token = AuthService.getToken();
    final loggingIn = state.matchedLocation == '/0_0_login';

    if (token == null && !loggingIn) return '/0_0_login';
    if (loggingIn) return '/0_1_menu';
    return null;
  },
);

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'Orders App',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.lightGreen),
      ),
      routerConfig: _router,
    );
  }
}
