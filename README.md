# GoMod 

A modern, high-performance TUI (Terminal User Interface) tracker music player for macOS, built with Go.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25%2B-00ADD8)

##  Features

### Playback & Audio
-  **High-Quality Audio** - Windowed sinc interpolation for pristine sound
-  **Live Stereo Control** - Adjust stereo separation in real-time (0-200%)
-  **Playback Controls** - Pause/resume with spacebar
-  **Accurate Timing** - Direct position tracking from libopenmpt

### Visualization
-  **Real-time Pattern View** - Typewriter-style scrolling tracker display
-  **Channel VU Meters** - 3-row vertical bars with smooth gravity physics
-  **9 Color Themes** - Peachy, Purple, Pastel, Matrix, Cyberpunk, and more
-  **Active Instrument Tracking** - See which instruments are playing

### Channel Control
-  **Channel Muting** - Mute/unmute individual channels (1-9, 0, -, =)
-  **Channel Soloing** - Solo channels with Shift+key
-  **Visual Feedback** - Muted channels shown dimmed in pattern and VU meters

### Performance
-  **Optimized CGo** - Pattern caching eliminates ~2,600 CGo calls per frame
-  **Metadata Caching** - One-time fetch of immutable module data
-  **60 FPS UI** - Smooth, responsive interface

## Supported Formats

- `.mod` (ProTracker)
- `.xm` (FastTracker II)
- `.it` (Impulse Tracker)
- `.s3m` (ScreamTracker 3)
- And 20+ more via libopenmpt!

## Installation

### Prerequisites

Install `libopenmpt`:

```bash
brew install libopenmpt
```

### Build from Source

```bash
git clone https://github.com/slimewell/GoMod.git
cd GoMod
go build -o gomod ./cmd/modtui
```

### Install to PATH

```bash
cp gomod ~/.local/bin/
# Or
sudo cp gomod /usr/local/bin/
```

##  Usage

```bash
go build -o gomod ./cmd/modtui
./gomod path/to/your/module.xm
```

Or run directly:
```bash
go run ./cmd/modtui path/to/your/module.xm
```

When in PATH:
```bash
gomod <module-file>
```

### Controls

| Key | Action |
|-----|--------|
| **Space** | Play/Pause |
| **Q** | Quit |
| **T** | Cycle themes (9 available) |
| **[ ]** | Adjust stereo separation (0-200%) |
| **1-9, 0, -, =** | Mute/unmute channels (hex-style: 1=Ch0, 0=Ch9, -=Ch10, ==Ch11) |
| **Shift + 1-9, 0, -, =** | Solo channel (unmute one, mute all others) |

### Themes

GoMod includes 9 carefully crafted color themes:
- **Default** - Classic tracker aesthetic
- **Gruvbox** - Warm, retro palette
- **Nord** - Cool, minimal blues
- **Dracula** - Purple-tinted dark mode
- **Monokai** - Vibrant syntax-inspired
- **Solarized** - Precision-balanced contrast
- **Cyberpunk** - Neon-soaked future
- **Matrix** - Green-on-black terminal
- **Sunset** - Warm orange gradients

Press **T** to cycle through them in real-time.

## Configuration

GoMod saves preferences to `~/.gomod.json`:
- Theme choice
- Stereo separation
- Last played file

## Architecture

### Tech Stack
- **[libopenmpt](https://lib.openmpt.org/)** - Tracker file decoding (Extended API for channel control)
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Styling
- **[Oto](https://github.com/hajimehoshi/oto)** - Audio output

### Performance Optimizations
- **Pattern Caching** - Entire patterns cached in Go memory on first visit
- **Metadata Caching** - Title, artist, duration fetched once
- **Lazy Interface Fetching** - Extended API interfaces retrieved on-demand
- **VU Smoothing** - Exponential decay (0.92) with instant attack

## Contributing

Contributions welcome!

## 📄 License

MIT License - See LICENSE file

---

**Made with ❤️ for tracker music fans**

*Optimized for macOS • Built with Go • Powered by libopenmpt*
