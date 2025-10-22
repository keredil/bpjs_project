class Program {
  final int id;
  final String name;
  final String deskripsi;
  final String bentukManfaat;
  final String manfaatLengkap;

  Program({
    required this.id,
    required this.name,
    required this.deskripsi,
    required this.bentukManfaat,
    required this.manfaatLengkap,
  });

  factory Program.fromJson(Map<String, dynamic> json) {
    return Program(
      id: json['id'],
      name: json['name'],
      deskripsi: json['deskripsi'],
      bentukManfaat: json['bentuk_manfaat'],
      manfaatLengkap: json['manfaat_lengkap'],
    );
  }
}
