package main

import (
	"io"
	"log"
	"net/http"
)

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Pack Calculator</title>
  <style>
    body { font-family: system-ui, sans-serif; max-width: 32rem; margin: 2rem auto; padding: 0 1rem; }
    label { display: block; margin-bottom: .5rem; font-weight: 600; }
    input, button { font-size: 1rem; padding: .4rem .6rem; }
    table { border-collapse: collapse; margin-top: 1rem; }
    th, td { border: 1px solid #ccc; padding: .4rem .8rem; text-align: left; }
    .error { color: #b00020; }
  </style>
</head>
<body>
  <h1>Pack Calculator</h1>
  <form id="pack-form">
    <label for="quantity">Quantity</label>
    <input id="quantity" name="quantity" type="number" min="1" value="1" required>
    <button type="submit">Calculate</button>
  </form>
  <div id="result"></div>

  <script>
    const form = document.getElementById('pack-form');
    const result = document.getElementById('result');

    form.addEventListener('submit', async (event) => {
      event.preventDefault();
      const quantity = Number(document.getElementById('quantity').value);
      result.innerHTML = '';

      try {
        const response = await fetch('/pack', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ quantity }),
        });
        const body = await response.json();

        if (!response.ok) {
          result.innerHTML = '<p class="error"></p>';
          result.querySelector('.error').textContent = body.detail || 'Something went wrong.';
          return;
        }
        renderPacks(quantity, body.packs);
      } catch (err) {
        result.innerHTML = '<p class="error"></p>';
        result.querySelector('.error').textContent = 'Could not reach the server.';
      }
    });

    function renderPacks(quantity, packs) {
      const sizes = Object.keys(packs).map(Number).sort((a, b) => b - a);
      const rows = sizes.map(s => '<tr><td>' + s + '</td><td>' + packs[s] + '</td></tr>').join('');
      result.innerHTML =
        '<p>For ' + quantity + ' item(s), ship:</p>' +
        '<table><thead><tr><th>Pack size</th><th>Count</th></tr></thead>' +
        '<tbody>' + rows + '</tbody></table>';
    }
  </script>
</body>
</html>
`

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := io.WriteString(w, indexHTML); err != nil {
		log.Printf("write page: %v", err)
	}
}
