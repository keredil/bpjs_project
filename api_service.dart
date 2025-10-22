import 'package:http/http.dart' as http;
import 'dart:convert';
import '../models/program.dart';

class ApiService {
  static const String baseUrl = "http://192.168.1.10:8080"; // ganti IP kamu

  static Future<List<Program>> getPrograms() async {
    final response = await http.get(Uri.parse("$baseUrl/programs"));
    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body)['data'];
      return data.map((e) => Program.fromJson(e)).toList();
    } else {
      throw Exception("Gagal mengambil data dari server");
    }
  }
}
