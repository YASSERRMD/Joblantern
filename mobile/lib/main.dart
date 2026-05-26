// Joblantern mobile — minimal Flutter client over the Joblantern HTTP API.
//
// This Phase 23 scaffold ships the screens + the API client; the
// drift-based offline cache, biometric lock, deep links, and CI build
// pipeline are follow-up commits documented in docs/MOBILE.md.

import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

const _defaultApiBase = 'https://joblantern.example';

void main() => runApp(const JoblanternApp());

class JoblanternApp extends StatelessWidget {
  const JoblanternApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Joblantern',
      theme: ThemeData.dark(useMaterial3: true),
      home: const HomeScreen(),
    );
  }
}

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});
  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final _text = TextEditingController();
  final _company = TextEditingController();
  final _country = TextEditingController();
  bool _busy = false;

  Future<void> _submit() async {
    setState(() => _busy = true);
    try {
      final prefs = await SharedPreferences.getInstance();
      final base = prefs.getString('api_base') ?? _defaultApiBase;
      final api = JoblanternApi(base);
      final id = await api.verify(
        listingText: _text.text,
        companyName: _company.text,
        jurisdiction: _country.text.toUpperCase(),
      );
      if (!mounted) return;
      Navigator.push(context, MaterialPageRoute(
        builder: (_) => ResultScreen(apiBase: base, id: id),
      ));
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    } finally {
      setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Joblantern')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
          TextField(
            controller: _text, maxLines: 8,
            decoration: const InputDecoration(
              labelText: 'Paste the recruiter message',
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _company,
            decoration: const InputDecoration(labelText: 'Company name', border: OutlineInputBorder()),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _country, maxLength: 2,
            decoration: const InputDecoration(labelText: 'Country (ISO-2)', border: OutlineInputBorder()),
          ),
          const SizedBox(height: 16),
          FilledButton(
            onPressed: _busy ? null : _submit,
            child: _busy ? const CircularProgressIndicator() : const Text('Verify'),
          ),
        ]),
      ),
    );
  }
}

class ResultScreen extends StatefulWidget {
  final String apiBase;
  final String id;
  const ResultScreen({super.key, required this.apiBase, required this.id});
  @override
  State<ResultScreen> createState() => _ResultScreenState();
}

class _ResultScreenState extends State<ResultScreen> {
  Map<String, dynamic>? _rec;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _poll();
    _timer = Timer.periodic(const Duration(seconds: 2), (_) => _poll());
  }

  Future<void> _poll() async {
    try {
      final api = JoblanternApi(widget.apiBase);
      final rec = await api.get(widget.id);
      setState(() => _rec = rec);
      final status = rec['status'];
      if (status == 'completed' || status == 'failed') _timer?.cancel();
    } catch (_) {}
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final verdict = _rec?['verdict'] as Map<String, dynamic>?;
    final risk = verdict?['overall_risk']?.toString() ?? 'pending';
    final color = switch (risk) {
      'red' => Colors.red,
      'yellow' => Colors.orange,
      'green' => Colors.green,
      _ => Colors.grey,
    };
    return Scaffold(
      appBar: AppBar(title: Text('Verification ${widget.id.substring(0, 8)}')),
      body: Center(child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(12)),
            child: Text(risk.toUpperCase(),
              style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: Colors.white)),
          ),
          const SizedBox(height: 16),
          if (verdict != null) ...[
            Text('Confidence ${(((verdict['confidence'] as num?)?.toDouble() ?? 0) * 100).toStringAsFixed(0)}%'),
            const SizedBox(height: 12),
            ...((verdict['reasons'] as List?)?.take(5).map((r) => Text('• $r')) ?? const []),
          ] else
            const Text('Working…'),
        ]),
      )),
    );
  }
}

class JoblanternApi {
  final String base;
  JoblanternApi(this.base);

  Future<String> verify({
    required String listingText,
    required String companyName,
    required String jurisdiction,
  }) async {
    final r = await http.post(
      Uri.parse('$base/api/v1/verify'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'listing_text': listingText,
        'company_name': companyName,
        'jurisdiction': jurisdiction,
      }),
    );
    if (r.statusCode != 202) {
      throw Exception('verify ${r.statusCode}: ${r.body}');
    }
    return (jsonDecode(r.body) as Map<String, dynamic>)['verification_id'] as String;
  }

  Future<Map<String, dynamic>> get(String id) async {
    final r = await http.get(Uri.parse('$base/api/v1/verifications/$id'));
    if (r.statusCode != 200) throw Exception('get ${r.statusCode}');
    return jsonDecode(r.body) as Map<String, dynamic>;
  }
}
