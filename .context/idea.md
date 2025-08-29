# ASCII Name Transliteration API — Go Spec & Starter (with Hugo demo)

## Goals

1. Accept user-provided names in any script (Unicode).
2. Normalize + transliterate to ASCII using best-effort rules.
3. Parse into `{ title, given/first, middle[], family }` with optional `gender` inference (best-effort + confidence).
4. Produce an **ASCII-only** representation suitable for legacy systems (boarding passes, labels, etc.).
5. Provide a dead-simple public demo via a static **Hugo** site that can directly call the API while the API enforces **client-origin allowlists** and **rate limits** per client.

---

## Example

**Input** → `"Doctor Nguyễn Văn Minh"`

**Output** →

```json
{
  "name": {
    "family": "NGUYEN",
    "first": "Minh",
    "middle": ["Van"],
    "full_ascii": "DR NGUYEN MINH VAN",
    "full_unicode": "Doctor Nguyễn Văn Minh",
    "order": "western"
  },
  "title": "Dr.",
  "gender": { "value": "M", "confidence": 0.65, "reason": "Vietnamese middle marker 'Văn' often male" },
  "script": "Latin+Vietnamese",
  "notes": ["Diacritics removed via NFKD/ASCII mapping", "Title normalized to 'Dr.'"],
  "version": "v1"
}
```

> Note: We keep **`Văn`** as a middle token but do **not** equate it to an honorific. For airline-style strict uppercase ASCII we expose `full_ascii` (ICAO-like approximations).

---

## API Surface

### Base URL

```
https://api.nameservice.example.com
```

### Endpoint

```
GET /v1/transliterate/{client_id}
```

### Query Params

- `q` (required): the input name string.
- `locale` (optional): user hint (e.g., `vi`, `zh`, `ja`, `ko`, `de`, `es`).
- `order` (optional): desired output order for `full_ascii` (`western` or `eastern`). Default `western`.
- `upper_family` (optional bool): uppercase family name in ASCII (default `true`).
- `guess_gender` (optional bool): enable heuristic gender inference (default `true`).

### Headers (from static site)

- `X-Client-Id`: same as `{client_id}` in path. Used for logging/defense-in-depth.
- `Origin`: enforced by backend against the client’s **allowlisted origins**.
- (Optional) `Referer`: secondary check; useful when some static hosts omit `Origin`.

### Response `200 application/json`

```json
{
  "name": {
    "family": "...",
    "first": "...",
    "middle": ["..."],
    "suffix": "Jr.",
    "prefix": null,
    "order": "western|eastern",
    "full_ascii": "...",
    "full_unicode": "..."
  },
  "title": "Mr.|Ms.|Mrs.|Mx.|Dr.|Prof.|...",
  "gender": { "value": "M|F|X|U", "confidence": 0.0, "reason": "..." },
  "script": "DetectedScriptLabel",
  "notes": ["..."],
  "version": "v1",
  "timing_ms": 3
}
```

### Errors

- `400` invalid input / empty `q`.
- `401` unknown client id.
- `403` origin/referrer not allowed for client.
- `429` rate limited.
- `5xx` internal.

---

## Data Model & Rules

### 1) Title Detection / Normalization

- Map variants to canonical set: `Dr., Mr., Mrs., Ms., Mx., Prof., Rev., Hon., Sir, Lady`.
- Recognize full words: `Doctor → Dr.`, `Professor → Prof.`
- Handle language variants where safe (e.g., `Herr → Mr.`, `Frau → Ms./Mrs.` ambiguous → default `Ms.` with note).

### 2) Tokenization & Order

- Split on whitespace and common punctuation.
- Script detection (simple): count Unicode blocks per token; label overall script (e.g., `Han`, `Hiragana`, `Katakana`, `Latin`, `Cyrillic`).
- Order heuristics:
  - If locale in `{zh, ja, ko}` and tokens look like **FamilyGiven** → output `eastern` order but expose a `western` `full_ascii` variant if requested.
  - Default to Western: `title | given | middle* | family`.

### 3) Vietnamese-Specific Heuristics

- Preserve tokens like `Văn` (often male marker) and `Thị` (often female marker) as **middle elements**.
- Gender hints (soft): if middle contains `Văn` → `M` (confidence ~0.6). If contains `Thị` → `F` (confidence ~0.7). Document as heuristic.

### 4) CJK & Other Scripts

- Use best-effort transliteration (e.g., Unidecode-style) to Latin; do **not** auto-invert without hints.
- For Chinese: if locale `zh` and Han-only tokens length=2–3, assume `family-first` unless Romanized input detected.

### 5) Diacritics & ASCII Projection

- Use NFKD → drop combining marks → ASCII.
- Apply language-specific expansions:
  - German: `Ä→AE`, `Ö→OE`, `Ü→UE`, `ß→SS`.
  - Scandinavian: `Å→AA`, `Æ→AE`, `Ø→OE`.
  - Vietnamese diacritics removed but base letters kept.
  - Spanish: `Ñ→N` (note loss), accents removed.
- Uppercase family if `upper_family=true`.

### 6) Suffixes / Prefixes

- Recognize `Jr., Sr., II, III, IV` → `suffix`.
- Recognize nobiliary particles (`de`, `del`, `von`, `van`, `da`, `di`, `bin`, `binti`) and **keep case** for ASCII but keep them as part of `family` where culturally correct. Special case: Dutch `van` may be middle or family particle; keep with family if at penultimate/initial position.

### 7) Safety / Ambiguity

- Never infer title → from `Thi`, `Van`, `-san`, etc.
- Emit `notes[]` and `confidence` when applying heuristics.

---

## Go Implementation (Starter)

### Module Layout

```
/ namesvc
  /cmd/api
    main.go
  /internal/http
    router.go
    handlers.go
    cors.go
    ratelimit.go
    clientauth.go
  /internal/names
    normalize.go
    parse.go
    translit.go
    language.go
    rules.go
    detect.go
  /internal/store
    clients.go      (in-memory or Redis)
  /testdata
    cases.json
  go.mod
```

### Dependencies (Go)

- `golang.org/x/text/unicode/norm` (NFKD)
- `golang.org/x/text/transform`
- `github.com/mozillazg/go-unidecode` (fallback transliteration)
- `github.com/julienschmidt/httprouter` or `net/http` + chi
- Optional: Redis for rate limiting `github.com/redis/go-redis/v9`

#### `go.mod`

```go
module example.com/namesvc

go 1.22

require (
    golang.org/x/text v0.17.0
    github.com/mozillazg/go-unidecode v0.2.0
    github.com/go-chi/chi/v5 v5.1.0
    github.com/redis/go-redis/v9 v9.5.1 // optional
)
```

### HTTP Entrypoint (`cmd/api/main.go`)

```go
package main

import (
    "log"
    "net/http"
    "os"

    "example.com/namesvc/internal/http"
)

func main() {
    addr := ":8080"
    if v := os.Getenv("ADDR"); v != "" { addr = v }

    srv := &http.Server{ Addr: addr, Handler: http.NewRouter() }
    log.Printf("namesvc listening on %s", addr)
    log.Fatal(srv.ListenAndServe())
}
```

### Router & Middleware (`internal/http/router.go`)

```go
package http

import (
    chi "github.com/go-chi/chi/v5"
    "net/http"
)

func NewRouter() http.Handler {
    r := chi.NewRouter()
    r.Use(Recoverer, RequestID, Logging)
    r.Route("/v1", func(v chi.Router) {
        v.With(ClientGate, CORS, RateLimit).Get("/transliterate/{client}", Transliterate)
    })
    return r
}
```

### Client Gate & CORS (`internal/http/clientauth.go`, `cors.go`)

```go
// clientauth.go
package http

import (
    "context"
    "net/http"
)

type Client struct { ID string; AllowedOrigins []string; RPS int }

type ClientStore interface { Lookup(id string) (*Client, bool) }

var store ClientStore = NewInMemoryClientStore() // swap with Redis/DB

func ClientGate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "client")
        if id == "" { http.Error(w, "missing client", http.StatusUnauthorized); return }
        c, ok := store.Lookup(id)
        if !ok { http.Error(w, "unknown client", http.StatusUnauthorized); return }
        ctx := context.WithValue(r.Context(), clientKey{}, c)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

```go
// cors.go
package http

import (
    "net/http"
    "slices"
)

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c := r.Context().Value(clientKey{}).(*Client)
        origin := r.Header.Get("Origin")
        referer := r.Header.Get("Referer")
        allowed := origin != "" && slices.Contains(c.AllowedOrigins, origin)
        if !allowed && referer != "" {
            // Basic referer host match
            for _, ao := range c.AllowedOrigins {
                if strings.HasPrefix(referer, ao) { allowed = true; break }
            }
        }
        if !allowed { http.Error(w, "origin not allowed", http.StatusForbidden); return }
        if origin != "" { w.Header().Set("Access-Control-Allow-Origin", origin) }
        w.Header().Set("Vary", "Origin")
        next.ServeHTTP(w, r)
    })
}
```

### Rate Limit (`internal/http/ratelimit.go`)

```go
package http

import (
    "net/http"
    "sync"
    "time"
)

type bucket struct { tokens float64; last time.Time; mu sync.Mutex }
var buckets sync.Map // key=clientID

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c := r.Context().Value(clientKey{}).(*Client)
        bAny, _ := buckets.LoadOrStore(c.ID, &bucket{ tokens: float64(c.RPS), last: time.Now() })
        b := bAny.(*bucket)
        b.mu.Lock()
        defer b.mu.Unlock()
        now := time.Now()
        elapsed := now.Sub(b.last).Seconds()
        b.tokens = min(float64(c.RPS), b.tokens + elapsed*float64(c.RPS))
        b.last = now
        if b.tokens < 1.0 { http.Error(w, "rate limited", http.StatusTooManyRequests); return }
        b.tokens -= 1.0
        next.ServeHTTP(w, r)
    })
}
```

### Handler (`internal/http/handlers.go`)

```go
package http

import (
    "encoding/json"
    "net/http"
    "time"

    "example.com/namesvc/internal/names"
)

type resp struct {
    Name    names.Structured `json:"name"`
    Title   string           `json:"title"`
    Gender  names.Gender     `json:"gender"`
    Script  string           `json:"script"`
    Notes   []string         `json:"notes"`
    Version string           `json:"version"`
    Timing  int64            `json:"timing_ms"`
}

func Transliterate(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    q := r.URL.Query().Get("q")
    if q == "" { http.Error(w, "missing q", http.StatusBadRequest); return }
    opts := names.OptionsFromQuery(r.URL.Query())
    out := names.Process(q, opts)
    res := resp{
        Name: out.Name, Title: out.Title, Gender: out.Gender,
        Script: out.Script, Notes: out.Notes, Version: "v1",
        Timing: time.Since(start).Milliseconds(),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}
```

### Core Name Logic (`/internal/names`)

#### Types (`rules.go`)

```go
package names

type Options struct {
    Locale      string
    Order       string // western|eastern
    UpperFamily bool
    GuessGender bool
}

type Gender struct { Value string; Confidence float32; Reason string }

type Structured struct {
    Family      string   `json:"family"`
    First       string   `json:"first"`
    Middle      []string `json:"middle"`
    Prefix      *string  `json:"prefix,omitempty"`
    Suffix      *string  `json:"suffix,omitempty"`
    Order       string   `json:"order"`
    FullASCII   string   `json:"full_ascii"`
    FullUnicode string   `json:"full_unicode"`
}

type Output struct {
    Name   Structured
    Title  string
    Gender Gender
    Script string
    Notes  []string
}
```

#### Options Parser (`language.go`)

```go
func OptionsFromQuery(q url.Values) Options {
    return Options{
        Locale: q.Get("locale"),
        Order:  coalesce(q.Get("order"), "western"),
        UpperFamily: q.Get("upper_family") != "false",
        GuessGender: q.Get("guess_gender") != "false",
    }
}
```

#### Pipeline (`process.go`)

```go
func Process(input string, opt Options) Output {
    tokens, title, notes := detectTitleAndTokenize(input)
    script := detectScript(tokens)
    parts := splitParts(tokens, opt)
    if opt.GuessGender { inferGender(&parts, &notes) }
    asciiParts := toASCII(parts, opt)
    fullASCII := joinASCII(asciiParts, title, opt)
    name := Structured{
        Family: asciiParts.family, First: asciiParts.first,
        Middle: asciiParts.middle, Order: opt.Order,
        FullASCII: fullASCII, FullUnicode: input,
    }
    return Output{ Name: name, Title: title, Gender: parts.gender, Script: script, Notes: notes }
}
```

#### Transliteration (`translit.go`)

```go
func stripDiacritics(s string) string {
    // NFKD → remove combining marks
    t := transform.Chain(norm.NFKD, transform.RemoveFunc(isMn), norm.NFC)
    out, _, _ := transform.String(t, s)
    return out
}

func languageSpecificASCII(r rune) string {
    switch r {
    case 'Ä': return "AE"
    case 'Ö': return "OE"
    case 'Ü': return "UE"
    case 'ß': return "SS"
    case 'Æ': return "AE"
    case 'Ø': return "OE"
    case 'Å': return "AA"
    }
    return ""
}

func toASCIIWord(w string) string {
    // First try explicit mappings
    var b strings.Builder
    for _, r := range w {
        if rep := languageSpecificASCII(r); rep != "" { b.WriteString(rep); continue }
        if r <= 127 { b.WriteRune(r); continue }
        // fallback: unidecode
        b.WriteString(unidecode.Unidecode(string(r)))
    }
    return b.String()
}
```

#### Heuristics (`parse.go`)

```go
var titles = map[string]string{
    "dr": "Dr.", "doctor": "Dr.",
    "mr": "Mr.", "mrs": "Mrs.", "ms": "Ms.", "mx": "Mx.",
    "prof": "Prof.", "professor": "Prof.",
}

func detectTitleAndTokenize(s string) (tokens []string, title string, notes []string) {
    raw := strings.FieldsFunc(s, func(r rune) bool { return unicode.IsSpace(r) || strings.ContainsRune(",.;:\t\n\r()[]{}", r) })
    for _, t := range raw {
        lt := strings.ToLower(strings.Trim(t, "."))
        if v, ok := titles[lt]; ok && title == "" {
            title = v
            continue
        }
        tokens = append(tokens, t)
    }
    return
}

func splitParts(tokens []string, opt Options) (p struct{ first, family string; middle []string; gender Gender }, notes []string) {
    if len(tokens) == 0 { return }
    // Vietnamese hint
    if opt.Locale == "vi" || looksVietnamese(tokens) {
        // Assume Western order for romanized input: Given Middle* Family
        p.first = tokens[0]
        if len(tokens) > 2 { p.middle = tokens[1 : len(tokens)-1] }
        p.family = tokens[len(tokens)-1]
        return
    }
    // Default: Given Middle* Family
    p.first = tokens[0]
    if len(tokens) > 2 { p.middle = tokens[1 : len(tokens)-1] }
    p.family = tokens[len(tokens)-1]
    return
}

func inferGender(p *struct{ first, family string; middle []string; gender Gender }, notes *[]string) {
    for _, m := range p.middle {
        lm := strings.ToLower(stripDiacritics(m))
        if lm == "van" || lm == "van." || lm == "van\u" { // crude
            p.gender = Gender{ Value: "M", Confidence: 0.6, Reason: "Vietnamese marker 'Văn'" }
            *notes = append(*notes, "Gender heuristic: 'Văn' → male (not guaranteed)")
            return
        }
        if lm == "thi" {
            p.gender = Gender{ Value: "F", Confidence: 0.7, Reason: "Vietnamese marker 'Thị'" }
            *notes = append(*notes, "Gender heuristic: 'Thị' → female (not guaranteed)")
            return
        }
    }
    p.gender = Gender{ Value: "U", Confidence: 0.0 }
}

func toASCII(p struct{ first, family string; middle []string; gender Gender }, opt Options) (out struct{ first, family string; middle []string }) {
    out.first = toASCIIWord(p.first)
    out.family = toASCIIWord(p.family)
    for _, m := range p.middle { out.middle = append(out.middle, toASCIIWord(m)) }
    if opt.UpperFamily { out.family = strings.ToUpper(out.family) }
    return
}

func joinASCII(p struct{ first, family string; middle []string }, title string, opt Options) string {
    var segs []string
    if title != "" { segs = append(segs, strings.ToUpper(strings.TrimSuffix(title, "."))) }
    if opt.Order == "western" {
        segs = append(segs, p.family, p.first)
        segs = append(segs, p.middle...)
    } else {
        // eastern order variant
        segs = append(segs, p.first)
        segs = append(segs, p.middle...)
        segs = append(segs, p.family)
    }
    return strings.Join(segs, " ")
}
```

> The Vietnamese detection is intentionally simple here; in production, add language trigram models or Unicode block ratios for better detection.

### In-Memory Client Store (`/internal/store/clients.go`)

```go
type InMemoryClientStore struct { m map[string]*Client }

func NewInMemoryClientStore() *InMemoryClientStore {
    return &InMemoryClientStore{ m: map[string]*Client{
        "demo-public": { ID: "demo-public", AllowedOrigins: []string{"https://nameservice.example.com"}, RPS: 5 },
    }}
}

func (s *InMemoryClientStore) Lookup(id string) (*Client, bool) { c, ok := s.m[id]; return c, ok }
```

### Unit Tests (outline)

- Title normalization table tests.
- Vietnamese cases: `Nguyễn Văn Minh`, `Nguyễn Thị Mai`.
- German umlauts: `Müller`, `Jürgen Groß`.
- CJK hanzi/kanji/hanja → ASCII via unidecode.
- Mixed punctuation & whitespace.

---

## Hugo Demo Site

### Goals

- Single page with an input box; shows JSON from API call.
- Public **client_id** baked into the page (e.g., `demo-public`).
- Deployed to Netlify/Vercel/S3; API restricts requests to the site’s **Origin**.

### Hugo Structure

```
/site
  /config.toml
  /content
    _index.md (landing)
  /layouts
    index.html (single-page app)
  /assets/js/app.js
```

#### `/assets/js/app.js`

```js
const API_BASE = "https://api.nameservice.example.com"
const CLIENT_ID = "demo-public" // public identifier, limited by origin + rate limit

async function run() {
  const input = document.getElementById("name").value.trim()
  if (!input) return
  const url = `${API_BASE}/v1/transliterate/${CLIENT_ID}?q=${encodeURIComponent(input)}`
  const res = await fetch(url, { headers: { "X-Client-Id": CLIENT_ID } })
  const data = await res.json()
  document.getElementById("out").textContent = JSON.stringify(data, null, 2)
}

document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("go").addEventListener("click", run)
})
```

#### Minimal Page (`/layouts/index.html`)

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Name Transliteration — Demo</title>
    <style>
      body {
        font-family: system-ui, Segoe UI, Roboto, Helvetica, Arial, sans-serif;
        max-width: 720px;
        margin: 2rem auto;
        padding: 0 1rem;
      }
      pre {
        background: #f6f6f6;
        padding: 1rem;
        border-radius: 12px;
        overflow: auto;
      }
    </style>
  </head>
  <body>
    <h1>ASCII Name Transliteration</h1>
    <p>Try: <code>Doctor Nguyễn Văn Minh</code>, <code>李小龍</code>, <code>Müller</code></p>
    <input id="name" placeholder="Enter a name" style="width:100%;padding:.75rem;border-radius:10px;border:1px solid #ddd" />
    <button id="go" style="margin-top:1rem;padding:.6rem 1rem;border-radius:10px">Transliterate</button>
    <h2>Result</h2>
    <pre id="out">{}</pre>
    <script src="/js/app.js" type="module"></script>
  </body>
</html>
```

> Configure Hugo to pipe `/assets/js/app.js` → `/js/app.js` via the asset pipeline.

### CORS/Origin Setup

- In the API, set client `demo-public.AllowedOrigins = ["https://nameservice.example.com"]`.
- If you publish under a different domain, update the allowlist.

---

## Deployment

### Dockerfile (API)

```dockerfile
FROM golang:1.22 as build
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/namesvc ./cmd/api

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/namesvc /namesvc
ENV ADDR=":8080"
EXPOSE 8080
ENTRYPOINT ["/namesvc"]
```

### Dockerfile (Hugo site)

```dockerfile
FROM klakegg/hugo:0.126.1-ext as build
WORKDIR /site
COPY site/ .
RUN hugo --minify

FROM caddy:2
COPY --from=build /site/public /srv
```

### docker-compose.yml

```yaml
services:
  api:
    build: .
    environment:
      - ADDR=:8080
    ports: ["8080:8080"]
  web:
    build: ./site
    ports: ["8081:80"]
```

---

## Operational Notes

- **Client IDs are public**; treat them as **identifiers**, not secrets. Restrict by **Origin/Referer** and apply **rate limits**. For paid tiers, issue distinct client IDs per site and increase RPS.
- For **private server-side integrations**, also support a separate `Authorization: Bearer <secret>` flow, but that’s outside the static-site scenario.
- **Logging**: log `{client_id, origin, ip_prefix, q_length}` not the full name (PII) unless you have consent/legal basis.
- **Privacy**: if you store examples for improvement, hash or truncate inputs.

---

## Rust Option (sketch)

- Use `axum` or `actix-web` for HTTP, `whatlang` or `lingua` for detection, `deunicode` for transliteration, and `unicode-normalization` for NFKD.
- Rate limit with `tower` + `governor`.

---

## Test Cases

1. `Doctor Nguyễn Văn Minh` → title `Dr.`, family `NGUYEN`, first `Minh`, middle `["Van"]`, gender `M?` (~0.6).
2. `Prof. Jürgen Groß` → title `Prof.`, family `GROSS`, first `Jurgen`, middle `[]`.
3. `李小龍` (zh) → likely `Li Xiaolong` (unidecode), family `LI`, given `Xiaolong`.
4. `Tanaka-san Yoko` → do **not** map `-san` to title; output note.
5. `Maria del Carmen Núñez` → family `NUNEZ` (note: lost tilde), keep `del` particle with family in ASCII: `DEL CARMEN NUNEZ` if input implies that structure.

---

## Roadmap Enhancements

- Add ICAO Doc 9303 tables for more faithful passport-style transliteration.
- Per-locale toggles (Hàn-Việt vs Pinyin; Hepburn for Japanese).
- Better CJK name-order detection using frequency lists.
- Confidence scoring per field (not just gender).
- Admin API to manage clients, origins, and rate limits.

---

## Ready-to-Build Checklist

- [ ] Implement modules as sketched.
- [ ] Add unit tests for the provided cases.
- [ ] Wire CORS + rate limiting with in-memory store; later swap to Redis.
- [ ] Stand up Hugo site with baked `client_id` and test CORS allowlist.
- [ ] Containerize and deploy behind your preferred TLS proxy (Caddy, Nginx, Cloudflare Tunnel).

---
