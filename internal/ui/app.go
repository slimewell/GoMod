package ui

import (
	"context"
	"fmt"
	"github.com/slimewell/GoMod/internal/player"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tickMsg is sent periodically to update the UI
type tickMsg time.Time

// Model is the main Bubble Tea model
type Model struct {
	module            *player.Module
	player            *player.Player
	filename          string
	stereoSep         int
	width             int
	height            int
	visibleRows       int
	palette           ColorPalette
	instruments       []player.Instrument
	activeInstruments map[int]int
	patternData       player.PatternSnapshot
	lastVolumes       []float64 // For smooth VU decay
	currentTime       float64
	ctx               context.Context
	cancel            context.CancelFunc
	ready             bool
	err               error
}

// NewModel creates a new TUI model
func NewModel(filename string, stereoSep int, themeName string) *Model {
	ctx, cancel := context.WithCancel(context.Background())
	return &Model{
		filename:          filename,
		stereoSep:         stereoSep,
		palette:           GetPalette(themeName),
		activeInstruments: make(map[int]int),
		visibleRows:       21, // Default
		ctx:               ctx,
		cancel:            cancel,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadModule,
		m.tickCmd,
	)
}

func (m *Model) loadModule() tea.Msg {
	// Load module
	mod, err := player.LoadModule(m.filename)
	if err != nil {
		return errMsg{err}
	}

	// Set stereo separation (non-fatal if it fails)
	_ = mod.SetStereoSeparation(m.stereoSep)
	// Ignore error - stereo separation might not be supported on all modules

	// Set high-quality interpolation (8 = windowed sinc, best quality)
	_ = mod.SetInterpolationFilter(8)

	// Load instrument list
	m.instruments = mod.GetInstrumentList()

	// Create player
	p, err := player.NewPlayer(mod)
	if err != nil {
		return errMsg{err}
	}

	// Start playback
	if err := p.Play(m.ctx); err != nil {
		return errMsg{err}
	}

	m.module = mod
	m.player = p
	m.ready = true

	return nil
}

func (m *Model) tickCmd() tea.Msg {
	time.Sleep(17 * time.Millisecond) // ~60 FPS update rate
	return tickMsg(time.Now())
}

type errMsg struct {
	err error
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.cancel()
			if m.player != nil {
				m.player.Close()
			}
			if m.module != nil {
				m.module.Close()
			}
			return m, tea.Quit

		case " ":
			if m.player != nil {
				m.player.TogglePause()
			}
			return m, nil

		case "[":
			// Decrease stereo separation
			if m.player != nil && m.module != nil {
				m.stereoSep -= 10
				if m.stereoSep < 0 {
					m.stereoSep = 0
				}
				_ = m.module.SetStereoSeparation(m.stereoSep)
			}
			return m, nil

		case "]":
			// Increase stereo separation
			if m.player != nil && m.module != nil {
				m.stereoSep += 10
				if m.stereoSep > 200 {
					m.stereoSep = 200
				}
				_ = m.module.SetStereoSeparation(m.stereoSep)
			}
			return m, nil

		// Channel Muting (1-9, 0, -, =)
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=":
			if m.module != nil {
				ch := -1
				switch msg.String() {
				case "1":
					ch = 0
				case "2":
					ch = 1
				case "3":
					ch = 2
				case "4":
					ch = 3
				case "5":
					ch = 4
				case "6":
					ch = 5
				case "7":
					ch = 6
				case "8":
					ch = 7
				case "9":
					ch = 8
				case "0":
					ch = 9
				case "-":
					ch = 10
				case "=":
					ch = 11
				}
				if ch != -1 {
					m.module.ToggleChannelMute(ch)
				}
			}
			return m, nil

		// Channel Soloing (Shift + Mute keys)
		case "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+":
			if m.module != nil {
				ch := -1
				switch msg.String() {
				case "!":
					ch = 0
				case "@":
					ch = 1
				case "#":
					ch = 2
				case "$":
					ch = 3
				case "%":
					ch = 4
				case "^":
					ch = 5
				case "&":
					ch = 6
				case "*":
					ch = 7
				case "(":
					ch = 8
				case ")":
					ch = 9
				case "_":
					ch = 10
				case "+":
					ch = 11
				}
				if ch != -1 {
					m.module.SoloChannel(ch)
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Calculate visible rows: height - overhead
		// Overhead breakdown:
		// - Header: 2 (title + info)
		// - Spacing: 1
		// - Active instruments: 2-3 (can wrap)
		// - Spacing: 1
		// - VU meters: 3
		// - Pattern header/sep: 2
		// - Spacing: 1
		// - Controls: 1
		// Total: ~16 lines
		available := m.height - 16
		if available < 5 {
			available = 5
		}
		if available%2 == 0 {
			available--
		}
		m.visibleRows = available
		return m, nil

	case tickMsg:
		if m.ready {
			// Only update estimated time if playing
			if m.player != nil && m.player.IsPlaying() {
				// Update current time (accurate)
				m.currentTime = m.module.GetPositionSeconds()
			}

			// Update active instruments
			currentRow := m.module.GetCurrentRow()
			newActives := TrackActiveInstruments(m.module, currentRow)

			// Merge with existing (for fade effect later)
			for id, brightness := range newActives {
				m.activeInstruments[id] = brightness
			}

			// Fade existing active instruments
			// Slower fade for more realistic sustain visualization
			for id := range m.activeInstruments {
				if _, isNew := newActives[id]; !isNew {
					// Slower fade: -10 per tick (was -20)
					// This means ~500ms of visible sustain after note off
					m.activeInstruments[id] -= 10
					if m.activeInstruments[id] <= 0 {
						delete(m.activeInstruments, id)
					}
				}
			}

			// Update pattern data snapshot (safely locked)
			m.patternData = m.module.GetPatternSnapshot(m.visibleRows)

			// Apply gravity/smoothing to VU meters
			newVolumes := m.patternData.ChannelVolumes

			// Initialize lastVolumes if needed
			if len(m.lastVolumes) != len(newVolumes) {
				m.lastVolumes = make([]float64, len(newVolumes))
			}

			// Smooth decay: volumes fall gradually instead of instantly
			decayRate := 0.92 // Even slower decay (gravity)

			for i := range newVolumes {
				if newVolumes[i] >= m.lastVolumes[i] {
					// Rising: instant update for responsiveness (snappy)
					// We use >= to prevent oscillation when values are constant (e.g. paused)
					m.lastVolumes[i] = newVolumes[i]
				} else {
					// Falling: gradual decay for smooth gravity effect
					m.lastVolumes[i] = m.lastVolumes[i] * decayRate
					// Floor to prevent tiny values
					if m.lastVolumes[i] < 0.001 {
						m.lastVolumes[i] = 0
					}
				}
			}

			// Use smoothed volumes for display
			m.patternData.ChannelVolumes = m.lastVolumes
		}
		return m, m.tickCmd

	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if !m.ready {
		return "Loading module...\n"
	}

	metadata := m.module.GetMetadata()

	// Render header
	header := RenderHeader(metadata, m.filename, m.currentTime, m.stereoSep, m.palette)

	// Render compact active instruments (horizontal)
	activeInstruments := RenderInstrumentsCompact(m.instruments, m.activeInstruments, m.palette)

	// Get channel mute states
	mutedChannels := make([]bool, m.patternData.NumChannels)
	for i := 0; i < m.patternData.NumChannels; i++ {
		mutedChannels[i] = m.module.IsChannelMuted(i)
	}

	// Render Channel VU Meters with mute indicators
	vuMeters := RenderVUMeters(m.patternData.ChannelVolumes, mutedChannels, m.width, m.palette)

	// Render pattern view with mute indicators
	pattern := RenderPattern(m.patternData, mutedChannels, m.palette)

	// Render controls
	controls := lipgloss.NewStyle().
		Foreground(m.palette.Controls).
		Render("[q] quit  [space] pause  [[ ]] stereo  [1-9,0,-,=] mute  [Shift+] solo")

	// Build the view
	var sections []string
	sections = append(sections, header)
	sections = append(sections, "", activeInstruments)
	sections = append(sections, "", vuMeters) // Add VU meters
	sections = append(sections, pattern, "", controls)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
