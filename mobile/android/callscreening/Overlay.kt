// Overlay renders a translucent SYSTEM_ALERT_WINDOW over the
// incoming-call screen with a band-coloured banner and a short
// reason. Tapping the banner opens the in-app report flow.
package org.joblantern.app.callscreening

import android.content.Context
import android.graphics.Color
import android.graphics.PixelFormat
import android.view.Gravity
import android.view.WindowManager
import android.widget.LinearLayout
import android.widget.TextView

object Overlay {
    fun show(ctx: Context, number: String, headline: String, reason: String) {
        val wm = ctx.getSystemService(Context.WINDOW_SERVICE) as WindowManager
        val view = LinearLayout(ctx).apply {
            orientation = LinearLayout.VERTICAL
            setBackgroundColor(bandColor(headline))
            setPadding(32, 32, 32, 32)
            addView(TextView(ctx).apply {
                text = headline
                textSize = 22f
                setTextColor(Color.BLACK)
            })
            addView(TextView(ctx).apply {
                text = "$number — $reason"
                textSize = 16f
                setTextColor(Color.BLACK)
            })
        }
        val params = WindowManager.LayoutParams(
            WindowManager.LayoutParams.MATCH_PARENT,
            WindowManager.LayoutParams.WRAP_CONTENT,
            WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,
            WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE,
            PixelFormat.TRANSLUCENT
        ).apply { gravity = Gravity.TOP }
        wm.addView(view, params)
    }

    private fun bandColor(headline: String): Int = when {
        headline.contains("scam", true) -> Color.parseColor("#ff6666")
        headline.contains("unverified", true) -> Color.parseColor("#ffd966")
        else -> Color.parseColor("#a6d49f")
    }
}
