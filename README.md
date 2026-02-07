# 🎵 Terminal Metronome

A beautiful, feature-rich metronome for your terminal with an animated visual interface. Built with Go, Bubble Tea, and love for musicians and developers alike.

## ✨ Features

- 🎨 **Beautiful TUI** - Gorgeous terminal interface with smooth animations
- 🌊 **Wave Animation** - Visual wave effect on each beat
- 🔊 **High-Quality Sound** - Crisp, clean click generated programmatically (1200 Hz sine wave with exponential decay)
- ⏯️  **Play/Pause** - Start and stop whenever you need
- ⚡ **Adjustable BPM** - Real-time tempo changes from 20 to 300 BPM
- 🎯 **Downbeat Highlighting** - Beat 1 visually distinct with double border
- 🎵 **Musical Symbols** - Downbeat (♩) vs regular beats (♪)
- ⌨️  **Keyboard Controls** - Everything at your fingertips
- 🎨 **Color-Coded** - Active beats highlighted with beautiful styling

## 📦 Installation

### Homebrew (macOS/Linux)

```bash
brew install luca-filipponi/tap/metronome
```

### Go Install

```bash
go install github.com/luca-filipponi/metronome@latest
```

### From Source

```bash
git clone https://github.com/luca-filipponi/metronome.git
cd metronome
go mod download
go build -o metronome metronome.go
./metronome
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/luca-filipponi/metronome/releases).

## 🚀 Usage

### Basic Usage

```bash
# Start with default 120 BPM
metronome

# Start with custom BPM
metronome 140

# Valid BPM range: 20-300
metronome 80
```

### Controls

| Key | Action |
|-----|--------|
| `↑` or `+` | Increase BPM by 5 |
| `↓` or `-` | Decrease BPM by 5 |
| `Space` | Play/Pause |
| `q` or `Ctrl+C` | Quit |

## 🎬 Demo

```
🎵 Terminal Metronome

╭────────────╮  ╭────────────╮  ╭────────────╮  ╭────────────╮
║            ║  ║            ║  ║            ║  ║            ║
║     ♩      ║  ║     ♪      ║  ║     ♪      ║  ║     ♪      ║
║            ║  ║            ║  ║            ║  ║            ║
║     1      ║  ║     2      ║  ║     3      ║  ║     4      ║
╰────────────╯  ╰────────────╯  ╰────────────╯  ╰────────────╯

BPM: 120  •  Status: Playing

↑/+ increase BPM • ↓/- decrease BPM • space play/pause • q quit
```