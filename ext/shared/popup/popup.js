const browserAPI = typeof browser !== "undefined" ? browser : chrome;

document.getElementById("openSettings").addEventListener("click", (e) => {
  e.preventDefault();
  if (browserAPI.runtime.openOptionsPage) {
    browserAPI.runtime.openOptionsPage();
  } else {
    window.open(browserAPI.runtime.getURL("options/options.html"));
  }
});

// Show a short status hint when the active tab is one of the supported sites.
const SUPPORTED = [
  /linkedin\.com\/jobs\//,
  /indeed\.[a-z.]+\/viewjob/,
  /bayt\.com\/(en|ar)\//,
  /naukrigulf\.com\/job-listing-/,
  /gulftalent\.com\/[a-z]{2}\/jobs\//,
  /jobstreet\.[a-z.]+\/job\//,
];

browserAPI.tabs.query({ active: true, currentWindow: true }, (tabs) => {
  if (!tabs || !tabs[0]) return;
  const url = tabs[0].url || "";
  const match = SUPPORTED.find((re) => re.test(url));
  const state = document.getElementById("state");
  if (match) {
    state.textContent = "Active on this page — look for the inline badge near the job title.";
  } else {
    state.textContent = "Open a supported job board to see the inline verdict.";
  }
});
