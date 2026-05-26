// Joblantern Service Worker
//
// Strategy:
//   * Precache the small set of shell assets at install.
//   * Static assets under /static/: cache-first.
//   * GET /verifications/<id>: network-first with offline fallback.
//   * Everything else: pass through.
//
// Background sync queues failed POST /verify submissions for retry when
// the browser is back online.

const CACHE = "joblantern-shell-v1";
const SHELL = [
  "/",
  "/static/manifest.webmanifest",
  "/static/offline.html",
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE).then((c) => c.addAll(SHELL)).then(() => self.skipWaiting())
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)))
    ).then(() => self.clients.claim())
  );
});

self.addEventListener("fetch", (event) => {
  const req = event.request;
  if (req.method !== "GET") return;

  const url = new URL(req.url);

  // Static assets: cache-first.
  if (url.pathname.startsWith("/static/")) {
    event.respondWith(
      caches.match(req).then((hit) =>
        hit ||
        fetch(req).then((resp) => {
          const copy = resp.clone();
          caches.open(CACHE).then((c) => c.put(req, copy)).catch(() => {});
          return resp;
        })
      )
    );
    return;
  }

  // Verification result pages: network-first, fall back to cached or offline page.
  if (url.pathname.startsWith("/verifications/")) {
    event.respondWith(
      fetch(req)
        .then((resp) => {
          const copy = resp.clone();
          caches.open(CACHE).then((c) => c.put(req, copy)).catch(() => {});
          return resp;
        })
        .catch(() =>
          caches.match(req).then((hit) => hit || caches.match("/static/offline.html"))
        )
    );
    return;
  }

  // Home page: stale-while-revalidate.
  if (url.pathname === "/") {
    event.respondWith(
      caches.match(req).then((hit) => {
        const refresh = fetch(req).then((resp) => {
          const copy = resp.clone();
          caches.open(CACHE).then((c) => c.put(req, copy)).catch(() => {});
          return resp;
        }).catch(() => hit || caches.match("/static/offline.html"));
        return hit || refresh;
      })
    );
  }
});

// Background sync — replay queued submissions.
self.addEventListener("sync", (event) => {
  if (event.tag !== "joblantern-verify-queue") return;
  event.waitUntil(flushQueue());
});

async function flushQueue() {
  try {
    const all = await idbAll();
    for (const item of all) {
      try {
        const resp = await fetch("/api/v1/verify", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(item.payload),
        });
        if (resp.ok) await idbDelete(item.id);
      } catch (_) { /* keep queued */ }
    }
  } catch (_) { /* ignore */ }
}

// Tiny IndexedDB helpers for the offline queue.
function idbOpen() {
  return new Promise((res, rej) => {
    const r = indexedDB.open("joblantern-pwa", 1);
    r.onupgradeneeded = () => r.result.createObjectStore("queue", { keyPath: "id", autoIncrement: true });
    r.onsuccess = () => res(r.result);
    r.onerror = () => rej(r.error);
  });
}
async function idbAll() {
  const db = await idbOpen();
  return new Promise((res, rej) => {
    const tx = db.transaction("queue", "readonly");
    const req = tx.objectStore("queue").getAll();
    req.onsuccess = () => res(req.result);
    req.onerror = () => rej(req.error);
  });
}
async function idbDelete(id) {
  const db = await idbOpen();
  return new Promise((res, rej) => {
    const tx = db.transaction("queue", "readwrite");
    tx.objectStore("queue").delete(id);
    tx.oncomplete = () => res();
    tx.onerror = () => rej(tx.error);
  });
}

// Push notification handler.
self.addEventListener("push", (event) => {
  let data = {};
  try { data = event.data && event.data.json(); } catch (_) {}
  const title = (data && data.title) || "Joblantern verdict ready";
  const body = (data && data.body) || "Your verification has completed.";
  const url = (data && data.url) || "/";
  event.waitUntil(
    self.registration.showNotification(title, {
      body,
      icon: "/static/icons/icon-192.png",
      badge: "/static/icons/icon-192.png",
      data: { url },
    })
  );
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();
  const url = (event.notification.data && event.notification.data.url) || "/";
  event.waitUntil(self.clients.openWindow(url));
});
