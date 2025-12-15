# GoMod 🎵

A modern, high-performance TUI (Terminal User Interface) tracker music player for macOS, built with Go.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25%2B-00ADD8)

## ✨ Features

### Playback & Audio
- 🎵 **High-Quality Audio** - Windowed sinc interpolation for pristine sound
- 🎛️ **Live Stereo Control** - Adjust stereo separation in real-time (0-200%)
- ⏯️ **Playback Controls** - Pause/resume with spacebar
- 🎯 **Accurate Timing** - Direct position tracking from libopenmpt

### Visualization
- 📊 **Real-time Pattern View** - Typewriter-style scrolling tracker display
- 📈 **Channel VU Meters** - 3-row vertical bars with smooth gravity physics
- 🎨 **9 Color Themes** - Peachy, Purple, Pastel, Matrix, Cyberpunk, and more
- 🎹 **Active Instrument Tracking** - See which instruments are playing

### Channel Control
- 🔇 **Channel Muting** - Mute/unmute individual channels (1-9, 0, -, =)
- 🎚️ **Channel Soloing** - Solo channels with Shift+key
- 👁️ **Visual Feedback** - Muted channels shown dimmed in pattern and VU meters

### Performance
- ⚡ **Optimized CGo** - Pattern caching eliminates ~2,600 CGo calls per frame
- 💾 **Metadata Caching** - One-time fetch of immutable module data
- 🚀 **60 FPS UI** - Smooth, responsive interface

## 🎼 Supported Formats

- `.mod` (ProTracker)
- `.xm` (FastTracker II)
- `.it` (Impulse Tracker)
- `.s3m` (ScreamTracker 3)
- And 20+ more via libopenmpt!

## 📦 Installation

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

## 🚀 Usage

```bash
# Play a module
gomod song.xm

# With custom stereo separation (0-200, default 100)
gomod -separation 150 song.mod

# Choose a theme
gomod -theme peachy song.it
gomod -t cyberpunk song.s3m
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `Space` | Pause/Resume |
| `[` / `]` | Decrease/Increase Stereo Separation |
| `1-9,0,-,=` | Mute/Unmute Channels 1-12 |
| `Shift+1-9,0,-,=` | Solo Channels 1-12 |
| `q` / `Ctrl+C` | Quit |

### Available Themes

**Gradient Palettes:**
- `peachy` - Soft peachy-pink gradient
- `purple` - Vibrant purple-to-pink gradient
- `pastel` - Soft multi-color pastels
- `matrix` - Bright neon matrix green

**Classic Palettes:**
- `default` - Clean cyan/yellow/purple
- `cyberpunk` - Hot pink and neon cyan
- `amber` - Retro monochrome amber CRT
- `green` - Classic green screen terminal
- `ocean` - Calm blues and teals

## ⚙️ Configuration

GoMod saves preferences to `~/.modtui.json`:
- Theme choice
- Stereo separation
- Last played file

## 🏗️ Architecture

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

## 🗺️ Roadmap

- [ ] Oscilloscope/Waveform Visualizer
- [ ] Seeking (arrow keys)
- [ ] Playlist Support
- [ ] Volume Ramping (click reduction)
- [ ] Hard Mute Mode (instant silence)
- [ ] Export to WAV

## 🤝 Contributing

Contributions welcome! This project follows standard Go conventions.

## 📄 License

MIT License - See LICENSE file

---

**Made with ❤️ for tracker music fans**

*Optimized for macOS • Built with Go • Powered by libopenmpt*
