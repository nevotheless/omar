# omar — Immutable OCI-Image-Based omarchy OS

omar baut **täglich** ein bootbares OCI-Image aus Arch Linux + [omarchy](https://omarchy.org/)
und stellt es via `ghcr.io` bereit. Das CLI `omar` verwaltet Conversion, Updates und Rollback.

## Architektur

```
┌─────────────────────────────────────────────────┐
│  ghcr.io/basecamp/omar:rolling-2026-05-14       │
│  ┌─────────────────────────────────────────────┐ │
│  │ Arch Linux + Kernel + systemd-boot          │ │
│  │ Hyprland + Waybar + Mako + … (omarchy-DE)  │ │
│  │ NetworkManager, Firmware, Flatpak           │ │
│  │ /usr/bin/omar (unser CLI)                   │ │
│  └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘
         │
         ▼ tägliches pull via bootc
┌─────────────────────────────────────────────────┐
│  Laufendes System                                │
│  /usr → read-only (aus OCI-Image)               │
│  /etc → mutable (Configs)                       │
│  /var → mutable (Daten)                         │
│  /home → mutable (User)                         │
│  Flatpak + Distrobox für User-Pakete            │
└─────────────────────────────────────────────────┘
```

## Repository-Struktur

```
omar/
├── cmd/omar/           # CLI Entrypoint (Go + Cobra)
├── internal/           # Packages
│   ├── bootc/          # bootc-Wrapper (switch, upgrade, rollback, install)
│   ├── convert/        # mutable → immutable migration
│   ├── pkg/            # Paket-Backends (Flatpak, Distrobox)
│   ├── image/          # Image-Status, Registry-Info
│   └── update/         # Update-Logik
├── images/             # Build-Pipeline (mkosi)
│   ├── mkosi.conf      # Hauptkonfiguration
│   ├── packages/       # Paketlisten (base, hyprland, omarchy)
│   ├── scripts/        # Postinstall-, Konfigurations-Skripte
│   └── mkosi.extra/    # Dateien die 1:1 ins Image wandern
└── .github/workflows/  # CI/CD (täglicher Build, PR-Checks, Releases)
```

## CLI

| Befehl | Beschreibung |
|--------|-------------|
| `omar install [--from <image>]` | Läuft **auf einem bestehenden omarchy** → wandelt in immutable um |
| `omar install --fresh <disk>` | Fresh-Install auf leerer Platte |
| `omar update` | Prüft auf neues Image, staged es, fragt nach Reboot |
| `omar rollback` | Bootet zum vorherigen Deployment |
| `omar status` | Zeigt aktuelles Image, Deployments, Update-Status |
| `omar pkg add <pkg>` | Installiert smart: Flatpak (GUI) oder Distrobox (CLI) |
| `omar pkg list` | Listet installierte Pakete mit Herkunft |
| `omar version` | Zeigt CLI + Image-Version |

## Build-Pipeline

Täglicher Cron (06:00 UTC) in `.github/workflows/build-rolling.yml`:

1. `mkosi -d arch --format=oci build`
2. Image pushen als `rolling-YYYY-MM-DD` + `rolling`
3. QEMU-Boot-Test (optional)

## Entwicklung

```bash
# CLI bauen
go build ./cmd/omar

# CLI testen
go test ./...

# Image lokal bauen (braucht mkosi + root)
sudo mkosi -d arch --format=oci -C images build
```
