# Claret: Cross-Platform Build & Installation Guide

This guide explains how to compile and distribute Claret for Windows, macOS, and Linux.

## 1. System Requirements for Users
One of the best features of Claret is that **end-users do not need to install Go, Node.js, or GTK** (except on Linux).

| Platform | Requirement | Notes |
| :--- | :--- | :--- |
| **Windows** | WebView2 Runtime | Pre-installed on Windows 10/11. If missing, Claret will prompt to download it. |
| **macOS** | None | Uses the built-in Safari engine (WKWebView). Works on Intel and Apple Silicon. |
| **Linux** | WebKit2GTK | Users may need: `sudo apt install libwebkit2gtk-4.0-3` (Ubuntu/Debian) or equivalent. |

---

## 2. Compiling for Windows (from Linux)
To build a Windows executable (`.exe`) from your Linux machine, you must install the MinGW cross-compiler.

### Step A: Install MinGW
- **OpenSUSE**: `sudo zypper install mingw64-cross-gcc`
- **Ubuntu/Debian**: `sudo apt install gcc-mingw-w64`
- **Fedora**: `sudo dnf install mingw64-gcc`

### Step B: Build the Binary
Run the following command in the project root:
```bash
wails build -platform windows/amd64 -skipbindings
```
- **Result**: `build/bin/claret.exe`
- **Distribution**: Simply send the `.exe` file to your Windows users.

---

## 3. Compiling for macOS (Requires a Mac)
Due to Apple's proprietary frameworks (Cocoa) and code-signing requirements, you must build the macOS version on a machine running macOS.

### Step A: Setup
1. Install **Xcode Command Line Tools**: `xcode-select --install`
2. Install **Go** and **Node.js** via Homebrew: `brew install go node`
3. Install **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Step B: Build the App Bundle
```bash
wails build -platform darwin/universal -skipbindings
```
- **Note**: `-platform darwin/universal` creates a single app that runs natively on both Intel and M1/M2/M3 Macs.
- **Result**: `build/bin/Claret.app`
- **Distribution**: Right-click `Claret.app`, select "Compress", and share the `.zip` file.

---

## 4. Compiling for Linux
To build for your own system or other Linux users:

```bash
# Ensure dev headers are installed (OpenSUSE example)
sudo zypper install webkit2gtk3-devel gtk3-devel

# Build
go build -tags production,webkit2gtk41 -o build/bin/claret .
```

---

## 5. Automated Multi-Platform Builds (GitHub Actions)
If you don't have access to a Mac, the professional way to build for all platforms is to use GitHub Actions. Create a file at `.github/workflows/build.yml` with this content:

```yaml
name: Wails Build
on: [push]
jobs:
  build:
    strategy:
      matrix:
        os: [windows-latest, macos-latest, ubuntu-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: wailsapp/wails-build-action@v2
        with:
          build-name: claret
          platform: ${{ matrix.os == 'macos-latest' && 'darwin/universal' || (matrix.os == 'windows-latest' && 'windows/amd64' || 'linux/amd64') }}
```
*This workflow will automatically create binaries for all three OSs every time you push code to GitHub.*

---

## 6. Troubleshooting
- **Wails Binding Error**: If you see `encoding/json without types`, always ensure you use the `-skipbindings` flag and build using the `go build` command or the specific `wails build` commands listed above.
- **macOS Permissions**: Since the app isn't "Notarized" by Apple yet, users may need to **Right-Click > Open** the app the first time to bypass the "Unidentified Developer" warning.
