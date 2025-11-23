# ssh-portfolio

A small terminal-based developer portfolio that runs **entirely over SSH**.  
It’s heavily inspired by [terminal.shop](https://terminal.shop) – instead of a store, this is a personal profile you can `ssh` into.

---

## What this is

When someone runs:

```bash
ssh your-user@your-host -p <port on which custom ssh server is running>
````

they don’t get a normal shell.

Instead, they land in a full-screen **TUI (text user interface)** with:

* An intro animation (typewriter-style name reveal)
* A fixed “card” in the middle of the screen
* **Tabs** for:

  * Overview
  * Experience
  * Projects
  * Contact
* Keyboard navigation (`h/l` or arrows for tabs, `j/k` for paging inside lists)
* Clickable links in supporting terminals (GitHub, LinkedIn, etc.)

It’s essentially a small landing page for backend / DevOps folks, except it lives in the terminal and speaks SSH.

---

## Tech stack

Built on the Charm & Go ecosystem:

* **[gliderlabs/ssh](https://github.com/gliderlabs/ssh)**
  Low-level SSH server used under the hood by Wish.

* **[charmbracelet/wish](https://github.com/charmbracelet/wish)**
  Makes it easy to plug Bubble Tea apps into SSH sessions (handles PTY, sessions, host keys, etc.).

* **[charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)**
  The TUI framework that powers the application state, update loop and key handling.

* **[charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)**

  * `paginator` for paging through Experience / Projects.

* **[charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)**
  Styling, colors and layout (centering the card, tab bar, etc.).

* Standard Go YAML + SSH libraries for config and connections.

---

## How it works over SSH (high level)

1. **SSH server startup**

   * The Go binary starts an SSH server using Wish on a configurable address (e.g. `:23234`).
   * A host key (e.g. `ssh_host_ed25519`) is generated on first run and reused afterwards.

2. **Client connects**

   * When a user runs `ssh host -p 23234`, Wish:

     * Performs the SSH handshake.
     * Allocates a **PTY** (so we get key events and full-screen rendering).
     * Starts a Bubble Tea `Program` bound to that SSH session.

3. **Bubble Tea app**

   * The TUI reads your portfolio data from `internal/config/data.yaml`.
   * It renders:

     * the intro animation (typewriter name),
     * then the main card with tabs.
   * All interaction happens via keyboard:

     * `h` / `l` or `←` / `→` – switch tabs
     * `1–4` – jump directly to a tab
     * `j` / `k` – move between experiences/projects
     * `q` / `ctrl+c` – quit

4. **No shell**

   * The user never gets a system shell; the SSH session is **bound only to this app**.
   * When the Bubble Tea program exits, Wish closes the SSH session.

---

## Running locally

```bash
# clone the repo
git clone git@github.com:shbhom/ssh-portfolio.git
cd ssh-portfolio

# build & run
go build -o ssh-portfolio ./cmd/ssh-portfolio
./ssh-portfolio
```

By default it might listen on something like `:23234`.
Then from another terminal:

```bash
ssh localhost -p 23234
```

You should see the intro animation and the portfolio tabs.

> ⚠️ **Host keys** (`ssh_host_ed25519`, `ssh_host_ed25519.pub`) are generated locally and **must not be committed** to version control.

---

## Making your own portfolio

The application is designed so that the **content lives in YAML**, not in the Go code.

To create your own:

1. Add a file at:

```text
internal/config/data.yaml
```

2. Fill it with your own data following the structure below.

The app will read `data.yaml` at startup and render the tabs from it.

---

## `data.yaml` structure (example)

Here’s a sample `data.yaml` that you can use as a reference.
The content is **intentionally random and not similar** to the author’s actual profile.

```yaml
overview:
  name: ""
  headline: "Backend Engineer | Building reliable systems for messy real-world problems"
  summary: |
    Backend engineer with a focus on data-heavy systems, internal tools and
    developer experience. Enjoys debugging production incidents, cleaning up
    flaky infrastructure, and turning ad-hoc scripts into maintainable services.

  bullets:
    - "Languages: Go, TypeScript, Python, SQL"
    - "Frameworks: Fiber, FastAPI, NestJS"
    - "Technologies: Docker, Kubernetes, gRPC, Redis, Kafka, GitHub Actions"
    - "Datastores & Vector Stores: Postgres, MongoDB, ClickHouse, Redis, Qdrant"

socials:
  github:   "https://github.com/<>"
  linkedin: "https://www.linkedin.com/in/<>"
  website:  ""

contact:
  email: ""
  phone: ""

experience:
  - company:  "Northwind Analytics"
    role:     "Backend Engineer"
    period:   "2023 — Present"
    location: "Remote"
    bullets:
      - "Designed and maintained ETL services that ingest >50M events per day."
      - "Replaced a legacy cron-based data sync with a streaming pipeline using Kafka."
      - "Improved API latency by 40% through query tuning and caching."
    stack: "Go, Postgres, Kafka, Redis, Docker, Kubernetes"

  - company:  "Lumos Delivery"
    role:     "Software Engineer"
    period:   "2021 — 2023"
    location: "Bengaluru, India"
    bullets:
      - "Built internal tooling to replay production requests in staging environments."
      - "Implemented a feature flag system used by multiple teams for rollouts."
      - "Collaborated closely with SREs during on-call rotations and incident reviews."
    stack: "TypeScript, Node.js, FastAPI, MongoDB, Grafana, Prometheus"

projects:
  - name: "Logbook"
    bullets:
      - "CLI + web dashboard for aggregating logs from multiple microservices."
      - "Supports searchable archives with per-service retention policies."
      - "Used by a small team to debug staging and production incidents."
    stack: "Go, SQLite, S3-compatible storage, TUI for the CLI"
    links:
      code: ""
      demo: ""

  - name: "Slow Query Inspector"
    bullets:
      - "Tool that parses Postgres slow-query logs and suggests indexes."
      - "Produces HTML reports sorted by potential impact and execution time."
    stack: "Python, Flask, Postgres"
    links:
      code: ""

  - name: "FocusTimer TUI"
    bullets:
      - "Terminal-based Pomodoro timer with session statistics."
      - "Stores state locally without any external services."
    stack: "Go, Bubble Tea, Lipgloss"
    links: {}  # no public links for this one

```

### Field overview

* `overview`

  * `name` – your display name
  * `headline` – single-line tagline under the name
  * `summary` – a short paragraph describing you
  * `bullets` – a few highlight lines (skills, focus areas, etc.)

* `socials`

  * `github`, `linkedin`, `website` – used to build clickable links in Overview / Contact

* `contact`

  * `email` – primary way to reach you
  * `phone` – optional (shown in Contact tab if present)

* `experience[]`

  * `company`, `role`, `period`, `location`
  * `bullets[]` – up to ~3 bullet points per role work best
  * `stack` – technologies used there

* `projects[]`

  * `name`
  * `bullets[]` – short description of what the project does / why it exists
  * `stack` – tech used
  * `links.code` – link to repository (optional)
  * `links.demo` – link to live demo or docs (optional)

You can change the wording and data freely as long as the structure stays the same.

---

## Customization ideas

* Change colors / styles in Lipgloss to match your personal branding.
* Adjust keybindings if you prefer Vim-style only (`hjkl`) or arrows.
* Add more tabs (e.g. “Writing”, “Talks”) driven from extra sections in `data.yaml`.
* Hook the Contact tab into a small HTTP service that sends you emails via Resend/SES.

---

## Credits

* Inspired by **terminal.shop**, which showed how fun SSH-based UIs can be.
* Built with tons of help from the [**Charmbracelet**](https://charm.land/) ecosystem and open source tools.
