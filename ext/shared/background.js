// Joblantern WebExtension — background script.
//
// Responsibilities:
//   * receive `verify` messages from content scripts
//   * look up an existing verdict in IndexedDB by listing URL (7-day TTL)
//   * otherwise POST /api/v1/verify, then poll /api/v1/verifications/{id}
//     until status="completed" or "failed", with a hard timeout
//   * persist the verdict and reply to the content script
//
// Works on Chrome MV3 (service worker) and Firefox MV2 (event page).

const browserAPI = typeof browser !== "undefined" ? browser : chrome;

const DEFAULTS = {
  apiBase: "http://localhost:8080",
  apiKey: "",
  cacheDays: 7,
  telemetry: false,
};

const DB_NAME = "joblantern";
const STORE = "verdicts";
const POLL_INTERVAL_MS = 1500;
const POLL_TIMEOUT_MS = 60000;

// ---- settings ----

function getSettings() {
  return new Promise((resolve) => {
    browserAPI.storage.local.get(DEFAULTS, (items) => resolve(items));
  });
}

// ---- IndexedDB cache (very small, no library) ----

function openDB() {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, 1);
    req.onupgradeneeded = () => {
      const db = req.result;
      if (!db.objectStoreNames.contains(STORE)) {
        db.createObjectStore(STORE, { keyPath: "url" });
      }
    };
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

async function cacheGet(url) {
  try {
    const db = await openDB();
    return await new Promise((resolve, reject) => {
      const tx = db.transaction(STORE, "readonly");
      const req = tx.objectStore(STORE).get(url);
      req.onsuccess = () => resolve(req.result || null);
      req.onerror = () => reject(req.error);
    });
  } catch (_) {
    return null;
  }
}

async function cachePut(url, value, ttlDays) {
  try {
    const db = await openDB();
    return await new Promise((resolve, reject) => {
      const tx = db.transaction(STORE, "readwrite");
      tx.objectStore(STORE).put({
        url,
        value,
        savedAt: Date.now(),
        expiresAt: Date.now() + ttlDays * 24 * 60 * 60 * 1000,
      });
      tx.oncomplete = () => resolve();
      tx.onerror = () => reject(tx.error);
    });
  } catch (_) {
    /* ignore */
  }
}

function isFresh(entry) {
  return entry && entry.expiresAt && entry.expiresAt > Date.now();
}

// ---- API ----

async function apiVerify(settings, submission) {
  const headers = { "Content-Type": "application/json" };
  if (settings.apiKey) headers["X-Joblantern-API-Key"] = settings.apiKey;
  const resp = await fetch(settings.apiBase + "/api/v1/verify", {
    method: "POST",
    headers,
    body: JSON.stringify(submission),
  });
  if (!resp.ok) throw new Error("verify HTTP " + resp.status);
  return resp.json();
}

async function apiGet(settings, id) {
  const resp = await fetch(settings.apiBase + "/api/v1/verifications/" + id);
  if (!resp.ok) throw new Error("get HTTP " + resp.status);
  return resp.json();
}

async function pollVerdict(settings, id) {
  const deadline = Date.now() + POLL_TIMEOUT_MS;
  while (Date.now() < deadline) {
    const rec = await apiGet(settings, id);
    if (rec.status === "completed" || rec.status === "failed") {
      return rec;
    }
    await new Promise((r) => setTimeout(r, POLL_INTERVAL_MS));
  }
  throw new Error("verdict timeout");
}

// ---- message handler ----

browserAPI.runtime.onMessage.addListener((msg, _sender, sendResponse) => {
  if (!msg || msg.type !== "verify") return false;
  (async () => {
    try {
      const settings = await getSettings();
      const submission = msg.submission || {};
      const url = submission.listing_url || "";
      const cached = url && (await cacheGet(url));
      if (isFresh(cached)) {
        sendResponse({ ...cached.value, api_base: settings.apiBase, cached: true });
        return;
      }
      const created = await apiVerify(settings, submission);
      const rec = await pollVerdict(settings, created.verification_id);
      const result = {
        verification_id: created.verification_id,
        verdict: rec.verdict,
        api_base: settings.apiBase,
      };
      if (url) await cachePut(url, result, settings.cacheDays);
      sendResponse(result);
    } catch (err) {
      sendResponse({ error: String(err && err.message || err) });
    }
  })();
  // Returning true keeps the message channel open for an async sendResponse.
  return true;
});
