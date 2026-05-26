import 'package:flutter_test/flutter_test.dart';
import 'package:joblantern_mobile/main.dart';

void main() {
  testWidgets('home renders with Verify button', (tester) async {
    await tester.pumpWidget(const JoblanternApp());
    expect(find.text('Joblantern'), findsWidgets);
    expect(find.text('Verify'), findsOneWidget);
  });
}
