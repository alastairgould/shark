# Pack Calculator

A small HTTP service that works out the fewest whole packs needed to fulfil an
order, given a set of pack sizes. Comes with a small web UI and a JSON API.

**Live demo:** <https://shark-51ub.onrender.com>
(free hosting, so the first request after a while can take 30-60s to wake up).

## Requirements

- [Go](https://go.dev/dl/) 1.24 or newer (`go version` to check)

## Quick start

Easiest way to run it is the bundled script, which sets the env vars, builds the
binary and starts the server:

```bash
./run.sh
```

Then open <http://localhost:8080> in your browser.

If the script isn't executable yet:

```bash
chmod +x run.sh
./run.sh
```

Or build and run it yourself:

```bash
go build -o shark .
./shark
```

## Configuration

Config comes from environment variables. They're all optional and fall back to
defaults.

| Variable       | Default                    | Description                                          |
| -------------- | -------------------------- | ---------------------------------------------------- |
| `PORT`         | `8080`                     | Port the HTTP server listens on.                     |
| `PACK_SIZES`   | `250,500,1000,2000,5000`   | Comma-separated list of available pack sizes.        |
| `MAX_QUANTITY` | `1000000`                  | Largest quantity a single request may ask for.       |

Edit the defaults at the top of `run.sh`, or set them inline:

```bash
PORT=9090 PACK_SIZES=10,25,100 MAX_QUANTITY=5000 ./run.sh
```

Invalid values (non-numeric, zero, or negative) get ignored and fall back to the
default, with a note logged to stderr.

## Using the API

Send a `POST` to `/pack` with a JSON body containing the quantity:

```bash
curl -s http://localhost:8080/pack \
  -H 'Content-Type: application/json' \
  -d '{"quantity": 501}'
```

Response:

```json
{"packs":{"250":1,"500":1}}
```

The `packs` object maps each pack size to how many of that pack to ship.

> **Why `POST` and not `GET`?** It's a pure calculation, so `GET` would be fine
> too. I went with `POST` mainly to leave room to store results later and hand
> them back by id, and it stops responses being cached along the way.

Errors come back as [problem+json](https://www.rfc-editor.org/rfc/rfc9457) with a
matching HTTP status, for example a quantity below `1` or above `MAX_QUANTITY`:

```json
{"type":"about:blank","title":"Bad Request","status":400,"detail":"quantity must be at least 1"}
```

## Running the tests

```bash
go test ./...
```

## Design notes

**Correctness.** Pack sizes are configurable so they might not be neat multiples
of each other, which would have let me get away with something simpler. Greedy
approaches (just grab the biggest pack that fits) can over-ship on odd sizes, so
I went with dynamic programming instead. It always follows the rules: fewest
items first, then fewest packs. There are tests covering it.

**Flat layout.** Everything sits in one `main` package at the root, with files
split by what they do (`main`, `api`, `ui`, `packer`). For something this small
the extra `cmd/` or `internal/` folders don't really earn their keep yet. If the
packing logic grew, `internal/pack/` would be the next step.

**Graceful shutdown** is left out on purpose. The service is stateless and
requests are quick, so there's nothing much to drain on exit, it just stops. If
in-flight work ever mattered you'd add `http.Server.Shutdown` behind a signal
context.
