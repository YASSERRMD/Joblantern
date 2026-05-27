// Package ui renders the rental submission flow. Layout mirrors the
// existing job flow so users can move between verticals without
// re-learning the form.
package ui

const SubmitForm = `<!doctype html>
<html><head><meta charset="utf-8"><title>Rental check</title>
<link rel="stylesheet" href="/static/joblantern.css"></head>
<body>
<main class="container">
<h1>Check a rental listing</h1>
<form action="/rental/submit" method="post">
  <label>Listing URL <input name="listing_url" type="url" required></label>
  <label>City <input name="city" required></label>
  <label>Country (ISO-2) <input name="country" maxlength="2" required></label>
  <label>Monthly rent <input name="monthly_rent" type="number" step="0.01"></label>
  <label>Currency <input name="currency" maxlength="3"></label>
  <label>Deposit method <input name="deposit_method"></label>
  <label>Landlord contact phone <input name="contact_phone" type="tel"></label>
  <label>Landlord contact email <input name="contact_email" type="email"></label>
  <fieldset><legend>Listing photos (URLs, one per line)</legend>
    <textarea name="image_urls" rows="4"></textarea>
  </fieldset>
  <button type="submit">Check</button>
</form>
</main>
</body></html>`
