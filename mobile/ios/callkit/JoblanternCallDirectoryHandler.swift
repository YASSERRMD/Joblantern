// SPDX-License-Identifier: Apache-2.0
//
// JoblanternCallDirectoryHandler is the iOS-side equivalent. Apple
// restricts how aggressively a third-party app can intervene on an
// incoming call: this extension ships a Call Directory that the OS
// consults to display a caller label and an optional block list.
//
// We ship best-effort. Full pre-pickup overlays are not possible on
// stock iOS; users get a labelled caller line.
import CallKit
import Foundation

class JoblanternCallDirectoryHandler: CXCallDirectoryProvider {

    override func beginRequest(with context: CXCallDirectoryExtensionContext) {
        context.delegate = self
        Snapshot.shared.load { entries in
            for e in entries.sorted(by: { $0.numberE164 < $1.numberE164 }) {
                if e.band == "red" {
                    context.addBlockingEntry(withNextSequentialPhoneNumber: e.cxNumber)
                } else {
                    context.addIdentificationEntry(withNextSequentialPhoneNumber: e.cxNumber, label: e.label)
                }
            }
            context.completeRequest()
        }
    }
}

extension JoblanternCallDirectoryHandler: CXCallDirectoryExtensionContextDelegate {
    func requestFailed(for extensionContext: CXCallDirectoryExtensionContext, withError error: Error) {
        NSLog("Joblantern call directory failed: \(error)")
    }
}
