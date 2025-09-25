import 'package:shared_preferences/shared_preferences.dart';

class AuthService {
  static const _tokenKey = "auth_token";

  // Кеш токена в памяти (чтобы не лазить в SharedPreferences каждый раз)
  static String? _cachedToken;

  /// Сохраняем токен
  static Future<void> saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, token);
    _cachedToken = token;
  }

  /// Получаем токен (сначала из кеша, потом из SharedPreferences)
  static Future<String?> getToken() async {
    if (_cachedToken != null) return _cachedToken;
    final prefs = await SharedPreferences.getInstance();
    _cachedToken = prefs.getString(_tokenKey);
    return _cachedToken;
  }

  /// Удаляем токен (разлогин)
  static Future<void> clearToken() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
    _cachedToken = null;
  }
}
