# GoMod

A modern, high-performance TUI (Terminal User Interface) tracker music player for macOS, built with Go.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21%2B-00ADD8)

## Features

### Playback & Audio
- **High-Quality Audio** - Windowed sinc interpolation for pristine sound
- **Instant Mute/Solo** - Channel changes take effect immediately via flush+seek
- **Live Stereo Control** - Adjust stereo separation in real-time (0-200%)
- **Hardware-Synced UI** - Pattern view locked precisely to audio output

### File Browser
- **Built-in File Browser** - Launch without arguments to browse for modules
- **Hot-Swap Modules** - Press Tab to switch songs without restarting
- **Shared Audio Context** - No driver reinit between tracks
- **Smart Filtering** - Directories first, hidden files excluded, modules highlighted

### Visualization
- **Real-time Pattern View** - Typewriter-style scrolling tracker display
- **Channel VU Meters** - 3-row vertical bars with smooth gravity physics
- **8 Color Themes** - Default, Amber, Green, Ocean, Peachy, Purple, Pastel, Matrix
- **Active Instrument Tracking** - See which instruments are playing

### Channel Control
- **Channel Muting** - Mute/unmute individual channels (1-9, 0, -, =)
- **Channel Soloing** - Solo channels with Shift+key
- **Visual Feedback** - Muted channels shown dimmed in pattern and VU meters

### Performance
- **Optimized CGo** - Pattern caching eliminates ~2,600 CGo calls per frame
- **Metadata Caching** - One-time fetch of immutable module data
- **60 FPS UI** - Smooth, responsive interface

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

## Usage

```bash
# Play a specific file
gomod path/to/module.xm

# Or launch and browse
gomod
```

### Controls

| Key | Action |
|-----|--------|
| **Space** | Play/Pause |
| **Tab** | Toggle file browser |
| **Q** | Quit |
| **[ ]** | Adjust stereo separation (0-200%) |
| **1-9, 0, -, =** | Mute/unmute channels (1=Ch1, 0=Ch10, -=Ch11, ==Ch12) |
| **Shift + 1-9, 0, -, =** | Solo channel (unmute one, mute all others) |

### File Browser

| Key | Action |
|-----|--------|
| **Up/Down** or **j/k** | Navigate |
| **Enter** | Open directory or play file |
| **Backspace** or **h** | Go up one directory |
| **Tab** or **Esc** | Close browser (if playing) |

### Themes

GoMod includes 8 color themes:
- **default** - Classic tracker aesthetic
- **amber** - Retro CRT amber
- **green** - Classic green terminal
- **ocean** - Calm blue/teal
- **peachy** - Soft coral gradient
- **purple** - Vibrant purple-pink
- **pastel** - Soft multi-color
- **matrix** - Bright neon green

Launch with a theme: `gomod -theme matrix song.mod`

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

### Key Implementation Details
- **Instant Mute**: Flush audio buffer + seek to playback position = no delay
- **Hardware Sync**: `UnplayedBufferSize()` tracks exact audio latency
- **Pattern Cache**: Full patterns stored in Go memory after first CGo fetch
- **Shared Context**: One `oto.Context` reused across module loads

## Contributing

Contributions welcome!

## License

MIT License - See LICENSE file

---

*Optimized for macOS - Built with Go - Powered by libopenmpt*