// Joblantern verdict badge — pure DOM, no framework dependency.

window.Joblantern = window.Joblantern || {};

window.Joblantern.renderBadge = function (anchorEl, state, payload) {
  if (!anchorEl) return null;
  let badge = anchorEl.parentElement.querySelector(".joblantern-badge");
  if (!badge) {
    badge = document.createElement("span");
    badge.className = "joblantern-badge pending";
    anchorEl.parentElement.appendChild(badge);
  }
  badge.className = "joblantern-badge " + state;
  badge.textContent = "";
  const dot = document.createElement("span");
  dot.className = "dot";
  badge.appendChild(dot);
  const label = document.createElement("span");
  label.textContent = " Joblantern: " + state;
  badge.appendChild(label);

  let tooltip = null;
  const showTip = () => {
    if (tooltip) return;
    tooltip = document.createElement("div");
    tooltip.className = "joblantern-tooltip";
    const rect = badge.getBoundingClientRect();
    tooltip.style.left = `${window.scrollX + rect.left}px`;
    tooltip.style.top = `${window.scrollY + rect.bottom + 6}px`;
    tooltip.innerHTML = buildTooltipHTML(state, payload);
    document.body.appendChild(tooltip);
  };
  const hideTip = () => {
    if (tooltip) {
      tooltip.remove();
      tooltip = null;
    }
  };
  badge.onmouseenter = showTip;
  badge.onmouseleave = hideTip;
  badge.onclick = (e) => {
    e.preventDefault();
    if (payload && payload.verification_id && payload.api_base) {
      window.open(payload.api_base + "/verifications/" + payload.verification_id, "_blank");
    }
  };
  return badge;
};

function buildTooltipHTML(state, payload) {
  if (state === "pending") return "<h4>Verifying…</h4>";
  if (state === "error") return `<h4>Error</h4><p>${escape(payload && payload.error)}</p>`;
  if (!payload || !payload.verdict) return `<h4>${state}</h4>`;
  const v = payload.verdict;
  const reasons = (v.reasons || []).slice(0, 5).map((r) => `<li>${escape(r)}</li>`).join("");
  return `
    <h4>${escape(v.overall_risk)} · confidence ${(v.confidence * 100).toFixed(0)}%</h4>
    ${reasons ? "<ul>" + reasons + "</ul>" : ""}
    <p style="margin-top:6px"><a href="#" data-jbl-open>Open full report</a></p>
  `;
}

function escape(s) {
  return String(s || "").replace(/[&<>"']/g, (c) =>
    ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c])
  );
}
