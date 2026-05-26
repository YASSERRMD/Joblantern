// Joblantern content-script entry point.
// Detects the site, extracts a submission, asks the background worker to
// verify, then renders the inline badge with the resulting verdict.

(function () {
  // Cross-browser browser API alias.
  const browserAPI = typeof browser !== "undefined" ? browser : chrome;

  function findAnchor(site) {
    try {
      return site.anchorSelector && site.anchorSelector();
    } catch (_) {
      return null;
    }
  }

  function go() {
    const site = window.Joblantern.detectSite();
    if (!site) return;
    const submission = site.extract();
    if (!submission) return;
    const anchor = findAnchor(site);
    if (!anchor) return;

    window.Joblantern.renderBadge(anchor, "pending", null);

    browserAPI.runtime.sendMessage(
      { type: "verify", submission, site: site.id },
      (resp) => {
        if (!resp) {
          window.Joblantern.renderBadge(anchor, "error", { error: "no response" });
          return;
        }
        if (resp.error) {
          window.Joblantern.renderBadge(anchor, "error", resp);
          return;
        }
        const state = (resp.verdict && resp.verdict.overall_risk) || "pending";
        window.Joblantern.renderBadge(anchor, state, resp);
      }
    );
  }

  // SPA-friendly: many job boards swap content without a full reload.
  let lastURL = window.location.href;
  go();
  setInterval(() => {
    if (window.location.href !== lastURL) {
      lastURL = window.location.href;
      // Remove old badges before re-running.
      document.querySelectorAll(".joblantern-badge").forEach((b) => b.remove());
      document.querySelectorAll(".joblantern-tooltip").forEach((t) => t.remove());
      go();
    }
  }, 1500);
})();
