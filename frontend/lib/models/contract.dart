class Contract {
  final String id;
  final String name;

  Contract({required this.id, required this.name});

  factory Contract.fromJson(Map<String, dynamic> json) {
    return Contract(
      id: (json['id'] ?? '').toString(),
      name: (json['name'] ?? 'Contract fără nume').toString(),
    );
  }
}
