import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Program BPJS Ketenagakerjaan',
      theme: ThemeData(
        primarySwatch: Colors.green,
        scaffoldBackgroundColor: Colors.grey[100],
      ),
      home: const ProgramBPJSPage(),
    );
  }
}

class ProgramBPJSPage extends StatefulWidget {
  const ProgramBPJSPage({super.key});

  @override
  State<ProgramBPJSPage> createState() => _ProgramBPJSPageState();
}

class _ProgramBPJSPageState extends State<ProgramBPJSPage> {
  List<dynamic> programs = [];
  bool isLoading = true;

  @override
  void initState() {
    super.initState();
    fetchPrograms();
  }

  Future<void> fetchPrograms() async {
    try {
      final response = await http.get(Uri.parse("http://localhost:8080/programs"));

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        setState(() {
          programs = data['data'];
          isLoading = false;
        });
      } else {
        throw Exception("Gagal memuat data (${response.statusCode})");
      }
    } catch (e) {
      setState(() => isLoading = false);
      print("Error: $e");
    }
  }

  void showProgramDetail(Map<String, dynamic> program) {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) {
        return Padding(
          padding: const EdgeInsets.all(20),
          child: SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: Container(
                    width: 40,
                    height: 4,
                    margin: const EdgeInsets.only(bottom: 16),
                    decoration: BoxDecoration(
                      color: Colors.grey[400],
                      borderRadius: BorderRadius.circular(10),
                    ),
                  ),
                ),
                Text(
                  program['nama'] ?? 'Tanpa nama',
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Colors.green,
                  ),
                ),
                const SizedBox(height: 10),
                Text(
                  program['deskripsi'] ?? 'Tidak ada deskripsi',
                  style: const TextStyle(fontSize: 15, color: Colors.black87),
                ),
                const Divider(height: 30, thickness: 1),
                const Text(
                  "Bentuk Manfaat:",
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                ),
                const SizedBox(height: 6),
                Text(
                  program['bentuk_manfaat'] ?? 'Tidak tersedia',
                  style: const TextStyle(fontSize: 14),
                ),
                const SizedBox(height: 12),
                const Text(
                  "Manfaat Lengkap:",
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                ),
                const SizedBox(height: 6),
                Text(
                  program['manfaat_lengkap'] ?? 'Tidak tersedia',
                  style: const TextStyle(fontSize: 14),
                ),
                const SizedBox(height: 20),
              ],
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          "Program BPJS Ketenagakerjaan",
          style: TextStyle(
            fontWeight: FontWeight.bold,
            fontSize: 18,
          ),
        ),
        centerTitle: true,
        backgroundColor: Colors.green[700],
        elevation: 0,
      ),
      body: isLoading
          ? const Center(child: CircularProgressIndicator(color: Colors.green))
          : programs.isEmpty
              ? const Center(child: Text("Tidak ada data program"))
              : Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: ListView.builder(
                    itemCount: programs.length,
                    itemBuilder: (context, index) {
                      final program = programs[index];
                      return GestureDetector(
                        onTap: () => showProgramDetail(program),
                        child: Container(
                          margin: const EdgeInsets.only(bottom: 12),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(16),
                            gradient: const LinearGradient(
                              colors: [Color(0xFF66BB6A), Color(0xFF43A047)],
                              begin: Alignment.topLeft,
                              end: Alignment.bottomRight,
                            ),
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withOpacity(0.1),
                                blurRadius: 6,
                                offset: const Offset(0, 3),
                              ),
                            ],
                          ),
                          child: ListTile(
                            leading: const Icon(Icons.card_membership,
                                color: Colors.white, size: 36),
                            title: Text(
                              program['nama'] ?? 'Tanpa nama',
                              style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 16,
                                  fontWeight: FontWeight.bold),
                            ),
                            subtitle: Text(
                              program['deskripsi'] ?? 'Tanpa deskripsi',
                              style: const TextStyle(
                                color: Colors.white70,
                                fontSize: 13,
                              ),
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}
