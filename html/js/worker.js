const staticSite = "dev-shopping-v1"
const assets = [
    "/",
    "/index.html",
    "/js/app.js",
    "/css/main.css",
    "/css/normalize.css"
]

self.addEventListener("install", installEvent => {
    installEvent.waitUntil(
        caches.open(staticSite).then(cache => {
            cache.addAll(assets)
        })
    )
})

self.addEventListener("fetch", fetchEvent => {
    fetchEvent.respondWith(
        caches.match(fetchEvent.request).then(res => {
            return res || fetch(fetchEvent.request)
        })
    )
})