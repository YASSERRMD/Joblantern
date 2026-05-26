const browserAPI = typeof browser !== "undefined" ? browser : chrome;

const DEFAULTS = {
  apiBase: "http://localhost:8080",
  apiKey: "",
  cacheDays: 7,
  telemetry: false,
};

function $(id) { return document.getElementById(id); }

function load() {
  browserAPI.storage.local.get(DEFAULTS, (items) => {
    $("apiBase").value = items.apiBase;
    $("apiKey").value = items.apiKey;
    $("cacheDays").value = items.cacheDays;
    $("telemetry").checked = !!items.telemetry;
  });
}

document.getElementById("settings").addEventListener("submit", (e) => {
  e.preventDefault();
  const values = {
    apiBase: $("apiBase").value.trim() || DEFAULTS.apiBase,
    apiKey: $("apiKey").value.trim(),
    cacheDays: Math.max(1, Math.min(30, parseInt($("cacheDays").value, 10) || DEFAULTS.cacheDays)),
    telemetry: $("telemetry").checked,
  };
  browserAPI.storage.local.set(values, () => {
    $("status").textContent = "Saved.";
    setTimeout(() => ($("status").textContent = ""), 1500);
  });
});

load();
