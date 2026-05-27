// SPDX-License-Identifier: Apache-2.0
//
// JoblanternCallScreeningService extends the Phase 23 mobile app with
// real-time scam-number lookups using Android's CallScreeningService.
package org.joblantern.app.callscreening

import android.telecom.Call
import android.telecom.CallScreeningService
import android.util.Log

class JoblanternCallScreeningService : CallScreeningService() {

    override fun onScreenCall(details: Call.Details) {
        val rawNumber = details.handle?.schemeSpecificPart
        if (rawNumber.isNullOrBlank()) {
            respond(details, allow = true)
            return
        }
        val verdict = Lookup.shared.classify(rawNumber)
        when (verdict.band) {
            "red" -> Overlay.show(this, rawNumber, "Possible scam recruiter", verdict.reason)
            "yellow" -> Overlay.show(this, rawNumber, "Unverified recruiter", verdict.reason)
            "green" -> Overlay.show(this, rawNumber, "Known legitimate recruiter", verdict.reason)
        }
        Log.i("Joblantern", "screened number=${rawNumber.takeLast(4)} band=${verdict.band}")
        respond(details, allow = verdict.band != "red")
    }

    private fun respond(details: Call.Details, allow: Boolean) {
        val response = CallResponse.Builder()
            .setDisallowCall(!allow)
            .setRejectCall(false)
            .setSkipCallLog(false)
            .setSkipNotification(false)
            .build()
        respondToCall(details, response)
    }
}
