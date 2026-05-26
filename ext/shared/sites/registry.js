// Joblantern WebExtension — site registry.
//
// Each supported job board registers an extractor that, given the current
// document, returns a Joblantern Submission object or null when the page
// is not actually a job listing.
//
// Extractors are intentionally small and defensive: DOMs change often, so
// every selector is wrapped in a try/catch and missing fields are simply
// left blank. The Joblantern server gracefully tolerates partial input.

window.Joblantern = window.Joblantern || {};
window.Joblantern.sites = [];

window.Joblantern.registerSite = function (site) {
  if (!site || typeof site.match !== "function" || typeof site.extract !== "function") {
    return;
  }
  window.Joblantern.sites.push(site);
};

window.Joblantern.detectSite = function () {
  for (const s of window.Joblantern.sites) {
    try {
      if (s.match(window.location)) return s;
    } catch (_) {
      /* ignore */
    }
  }
  return null;
};

// Utility used by every extractor.
window.Joblantern.textOf = function (selector) {
  try {
    const el = document.querySelector(selector);
    return el ? el.textContent.trim() : "";
  } catch (_) {
    return "";
  }
};
