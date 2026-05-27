// Package ui renders the marketplace submission form.
package ui

const SubmitForm = `<!doctype html>
<html><head><meta charset="utf-8"><title>Marketplace check</title>
<link rel="stylesheet" href="/static/joblantern.css"></head>
<body>
<main class="container">
<h1>Check a marketplace listing</h1>
<form action="/marketplace/submit" method="post">
  <label>Listing URL <input name="listing_url" type="url" required></label>
  <label>Platform <input name="platform" required></label>
  <label>Category <input name="category"></label>
  <label>Country (ISO-2) <input name="country" maxlength="2" required></label>
  <label>Price <input name="price" type="number" step="0.01"></label>
  <label>Currency <input name="currency" maxlength="3"></label>
  <label>Payment method <input name="payment_method"></label>
  <label>Shipping method <input name="shipping_method"></label>
  <label>Seller phone <input name="contact_phone" type="tel"></label>
  <label>Seller email <input name="contact_email" type="email"></label>
  <button type="submit">Check</button>
</form>
</main>
</body></html>`
