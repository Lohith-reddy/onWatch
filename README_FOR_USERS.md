# onWatch - Quick Use Guide

This is the shortest path to run onWatch and start tracking Codex + Claude usage.

## 1) Requirements

- macOS/Linux
- Go installed (`go version`)

## 2) Get the code

```bash
git clone https://github.com/Lohith-reddy/onWatch.git
cd onWatch
```

## 3) Configure environment

```bash
cp .env.example .env
```

Edit `.env` with:

- `ONWATCH_ADMIN_USER`
- `ONWATCH_ADMIN_PASS`
- At least one provider token (`CODEX_TOKEN` or `ANTHROPIC_TOKEN`) for initial startup

Optional multi-account tracking:

```env
ONWATCH_MULTI_ACCOUNTS=[{"name":"codex-work","provider":"codex","auth_file":"~/.codex/auth.json"},{"name":"claude-personal","provider":"anthropic","credentials_file":"~/.claude/.credentials.json"}]
```

## 4) Start app

```bash
chmod +x ./launch_onwatch.sh "./OnWatch Launcher.command"
./launch_onwatch.sh
```

Then open:

- http://127.0.0.1:9211

Login with your admin username/password from `.env`.

## 5) Add accounts with browser OAuth

In Settings -> Providers -> OAuth Account Login:

- Choose provider (`Codex` or `Claude`)
- Click **Start OAuth Login**
- onWatch will:
  - open your default browser (optional)
  - show the raw OAuth URL/code so you can use any browser profile manually

## 6) macOS app icon launcher (optional)

```bash
chmod +x ./scripts/create_onwatch_app.sh
./scripts/create_onwatch_app.sh
```

This creates `OnWatch Launcher.app` which you can double-click.

## 7) Stop app

If running in the foreground, use `Ctrl+C`.

---

## Open Source

- License: GPL-3.0 (`LICENSE`)
- Repo: https://github.com/Lohith-reddy/onWatch
- Contributions via pull requests are welcome.
