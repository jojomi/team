# Job-Format

Jobs sind YAML-Dateien unter `jobs/`. Deaktivierte Jobs enden auf `.yml.disabled`.

## Grundstruktur

```yaml
metadata:
  type: shell        # Pflichtfeld
  name: "Job-Name"
  # weitere Felder...

# typ-spezifische Felder
```

---

## `metadata`-Felder (alle Typen)

| Feld          | Typ      | Standard | Bedeutung |
|---------------|----------|----------|-----------|
| `type`        | string   | —        | **Pflicht.** `shell`, `open`, `text`, `feed-download`, `http-download` |
| `name`        | string   | —        | Anzeigename im Runner |
| `description` | string   | —        | Zusätzliche Beschreibung |
| `next`        | string   | —        | Scheduling-Ausdruck, z.B. `lastSuccess + "7d"` |
| `unattended`  | bool     | `false`  | `true` = Job läuft ohne User-Interaktion |
| `enabled`     | bool     | `true`   | `false` = Job wird ignoriert (alternativ: `.disabled`-Suffix) |
| `weight`      | int      | `0`      | Sortiergewichtung: **höherer Wert = früher** in der Liste |
| `sudo`        | bool     | `false`  | Job mit sudo ausführen |
| `possible`    | list     | —        | Bedingungen (Shell-Befehle), die alle Exit 0 liefern müssen |
| `output`      | string   | —        | `debug` für verbose Ausgabe, sonst leer |
| `pre.command` | string   | —        | Befehl vor dem Job |
| `post.command`| string   | —        | Befehl nach dem Job |

### Scheduling (`next`)

```yaml
next: lastSuccess + "7d"   # 7 Tage nach letztem Erfolg
next: lastSuccess + "24h"  # 24 Stunden nach letztem Erfolg
next: lastSuccess + "30m"  # 30 Minuten nach letztem Erfolg
next: lastSuccess          # immer fällig (kein Offset)
```

Einheiten: `m` (Minuten), `h` (Stunden), `d` (Tage)

### Bedingungen (`possible`)

Liste von Shell-Befehlen. Liefert einer Exit ≠ 0, gilt der Job als nicht ausführbar.

```yaml
possible:
  - assert weekday mon-fri              # nur Mo–Fr
  - assert weekday friday saturday sunday  # bestimmte Tage
  - assert ssh-reachable-noninteractive homer  # SSH-Host erreichbar?
  - assert time-after 18:00            # erst ab 18 Uhr
```

### Pre/Post

```yaml
pre:
  command: toolbox vpn omg
post:
  command: toolbox vpn reset
```

---

## Typ: `shell`

Führt einen Shell-Befehl aus.

```yaml
metadata:
  type: shell
  name: "Beispiel-Job"
  unattended: true
  next: lastSuccess + "1d"

execute: sh scripts/mein-script.sh
execute-dry-run: sh scripts/mein-script.sh --dry-run  # optional
wait: false  # optional: true = fragt nach Bestätigung
```

| Feld              | Bedeutung |
|-------------------|-----------|
| `execute`         | **Pflicht.** Auszuführender Befehl |
| `execute-dry-run` | Alternativer Befehl bei `--dry-run` |
| `wait`            | `true` = Runner fragt nach Abschluss: „Task done?" |

**Working directory:** Verzeichnis der Job-Datei (`jobs/`).

**Env-Variablen** (vom Runner gesetzt):
- `TEAM_DRYRUN` – `"true"` / `"false"`
- `TEAM_VERBOSE` – `"true"` / `"false"`

**VPN-Pattern:**
```yaml
execute: toolbox vpn with ökorenta -- ssh oeresoft-root 'df -h /'
execute: toolbox vpn with ökorenta -- sh scripts/mein-script/execute.sh
```

---

## Typ: `open`

Öffnet URLs, Dateien oder Apps.

```yaml
metadata:
  type: open
  name: "Dashboard öffnen"
  next: lastSuccess + "1d"

targets:
  - https://example.com/dashboard
  - /Applications/Cyberduck.app/Contents/MacOS/Cyberduck
  - /Users/me/.Trash/

wait: true  # optional
```

`targets` kann enthalten: HTTP(S)-URLs, Dateipfade, App-Binaries.

---

## Typ: `text`

Zeigt einen Text an (manuelle Erinnerung).

```yaml
metadata:
  type: text
  name: "Erinnerung"
  unattended: false
  next: lastSuccess + "14d"

text: "Hier steht, was zu tun ist."
wait: true  # optional
```

---

## Typ: `feed-download`

Lädt einen RSS/Atom-Feed herunter.

```yaml
metadata:
  type: feed-download
  name: "Feed laden"

remote-url: https://example.com/feed.xml
local-dir: ~/Downloads/feeds/
```

Benötigt `feeddownload` im PATH.

---

## Typ: `http-download`

Lädt Dateien per HTTP herunter.

```yaml
metadata:
  type: http-download
  name: "Datei laden"

files:
  - remote-url: https://example.com/file.zip
    local-dir: ~/Downloads/
    local-filename: file.zip  # optional
```

---

## Scripts-Konvention

Komplexe Shell-Logik gehört in `scripts/<job-name>/execute.sh`:

```
jobs/oeresoft-cert-check.yml
scripts/oeresoft-cert-check/execute.sh
```

Das Job-YAML ruft dann `sh scripts/<name>/execute.sh` auf.

---

## Beispiele

**Wöchentlicher unattended Shell-Job mit VPN:**
```yaml
metadata:
  type: shell
  name: "Server: Zertifikat prüfen"
  unattended: true
  next: lastSuccess + "7d"

execute: toolbox vpn with ökorenta -- sh scripts/cert-check/execute.sh
```

**Manueller Job mit Vorbedingung (Wochentag + SSH):**
```yaml
metadata:
  type: shell
  name: "Backup"
  unattended: false
  next: lastSuccess + "3d"
  possible:
    - assert weekday mon-fri
    - assert ssh-reachable-noninteractive homer

execute: rsync -avz /local/ homer:/remote/
execute-dry-run: rsync -avzn /local/ homer:/remote/
```

**Offene URL mit Pre-Command:**
```yaml
metadata:
  type: open
  name: "Internes Dashboard"
  next: lastSuccess + "12h"
  pre:
    command: toolbox vpn omg
  post:
    command: toolbox vpn reset

targets:
  - https://intern.example.com/dashboard
```
