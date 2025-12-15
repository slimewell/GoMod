package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/slimewell/GoMod/internal/player"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajimehoshi/oto/v2"
)

// tickMsg is sent periodically to update the UI
type tickMsg time.Time

// PlayerModel handles the music playback view
type PlayerModel struct {
	audioContext      *oto.Context
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
	lastVolumes       []float64
	currentTime       float64
	ctx               context.Context
	cancel            context.CancelFunc
	ready             bool
	err               error
}

// NewPlayerModel creates a new player model
func NewPlayerModel(audioContext *oto.Context, filename string, stereoSep int, themeName string, width, height int) *PlayerModel {
	ctx, cancel := context.WithCancel(context.Background())
	return &PlayerModel{
		audioContext:      audioContext,
		filename:          filename,
		stereoSep:         stereoSep,
		palette:           GetPalette(themeName),
		activeInstruments: make(map[int]int),
		visibleRows:       21,
		width:             width,
		height:            height,
		ctx:               ctx,
		cancel:            cancel,
	}
}

// Init initializes the player
func (m *PlayerModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadModule,
		m.tickCmd,
	)
}

func (m *PlayerModel) loadModule() tea.Msg {
	mod, err := player.LoadModule(m.filename)
	if err != nil {
		return errMsg{err}
	}

	_ = mod.SetStereoSeparation(m.stereoSep)
	_ = mod.SetInterpolationFilter(8)

	m.instruments = mod.GetInstrumentList()

	// Use shared audio context
	p, err := player.NewPlayer(m.audioContext, mod)
	if err != nil {
		return errMsg{err}
	}

	if err := p.Play(m.ctx); err != nil {
		return errMsg{err}
	}

	m.module = mod
	m.player = p
	m.ready = true

	// Recalculate layout now that we have instruments
	m.recalculateVisibleRows()

	return nil
}

func (m *PlayerModel) tickCmd() tea.Msg {
	time.Sleep(17 * time.Millisecond)
	return tickMsg(time.Now())
}

// Update handles player-specific messages
func (m *PlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Note: Global keys like q/ctrl+c are handled by AppModel, 
		// but we handle player controls here.
		
		case " ":
			if m.player != nil {
				m.player.TogglePause()
			}
			return m, nil

		case "[":
			if m.player != nil && m.module != nil {
				m.stereoSep -= 10
				if m.stereoSep < 0 {
					m.stereoSep = 0
				}
				_ = m.module.SetStereoSeparation(m.stereoSep)
			}
			return m, nil

		case "]":
			if m.player != nil && m.module != nil {
				m.stereoSep += 10
				if m.stereoSep > 200 {
					m.stereoSep = 200
				}
				_ = m.module.SetStereoSeparation(m.stereoSep)
			}
			return m, nil

		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=":
			if m.module != nil {
				ch := -1
				switch msg.String() {
				case "1": ch = 0
				case "2": ch = 1
				case "3": ch = 2
				case "4": ch = 3
				case "5": ch = 4
				case "6": ch = 5
				case "7": ch = 6
				case "8": ch = 7
				case "9": ch = 8
				case "0": ch = 9
				case "-": ch = 10
				case "=": ch = 11
				}
				if ch != -1 && m.player != nil {
					m.player.InstantMute(ch)
				}
			}
			return m, nil

		case "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+":
			if m.module != nil {
				ch := -1
				switch msg.String() {
				case "!": ch = 0
				case "@": ch = 1
				case "#": ch = 2
				case "$": ch = 3
				case "%": ch = 4
				case "^": ch = 5
				case "&": ch = 6
				case "*": ch = 7
				case "(": ch = 8
				case ")": ch = 9
				case "_": ch = 10
				case "+": ch = 11
				}
				if ch != -1 && m.player != nil {
					m.player.InstantSolo(ch)
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalculateVisibleRows()
		return m, nil

	case tickMsg:
		if m.ready {
			var currentRow, currentPattern int
			var currentVolumes []float64

			if m.player != nil && m.player.IsPlaying() {
				m.currentTime = m.player.GetSyncedTime()
				currentPattern, currentRow, currentVolumes = m.player.GetSyncedState()
			} else {
				currentRow = m.module.GetCurrentRow()
				currentPattern = m.module.GetCurrentPattern()
			}

			newActives := m.module.GetRowInstruments(currentRow)
			for id, brightness := range newActives {
				m.activeInstruments[id] = brightness
			}

			for id := range m.activeInstruments {
				if _, isNew := newActives[id]; !isNew {
					m.activeInstruments[id] -= 10
					if m.activeInstruments[id] <= 0 {
						delete(m.activeInstruments, id)
					}
				}
			}

			m.patternData = m.module.GetPatternView(currentPattern, currentRow, m.module.GetNumChannels(), m.visibleRows, currentVolumes)

			newVolumes := m.patternData.ChannelVolumes
			if len(m.lastVolumes) != len(newVolumes) {
				m.lastVolumes = make([]float64, len(newVolumes))
			}

			decayRate := 0.92
			for i := range newVolumes {
				if newVolumes[i] >= m.lastVolumes[i] {
					m.lastVolumes[i] = newVolumes[i]
				} else {
					m.lastVolumes[i] = m.lastVolumes[i] * decayRate
					if m.lastVolumes[i] < 0.001 {
						m.lastVolumes[i] = 0
					}
				}
			}
			m.patternData.ChannelVolumes = m.lastVolumes
		}
		return m, m.tickCmd

	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	return m, nil
}

func (m *PlayerModel) recalculateVisibleRows() {
	// Calculate exact overhead to maximize pattern view
	
	// 1. Instrument Panel Height
	// RenderInstrumentsCompact uses max 24 items, 8 per line.
	// Header "♪ Instruments:" is always 1 line.
	// Then chips: 1 to 3 lines.
	// If 0 instruments, it shows " (none)" on same line or next? 
	// Code says: return header + " (none)" if len=0. So 1 line total.
	// Else: Header \n Lines...
	
	instLines := 1 // Header
	if len(m.instruments) > 0 {
		count := len(m.instruments)
		if count > 24 {
			count = 24
		}
		// Ceiling division for lines
		lines := (count + 7) / 8
		instLines += lines
	}
	
	// 2. Total Overhead Calculation
	// Header: 2
	// Spacing: 1
	// Instruments: instLines
	// Spacing: 1
	// VU Meters: 3
	// Pattern Header (RenderPattern adds 2 lines): 2
	// Spacing: 1
	// Controls: 1
	
	overhead := 2 + 1 + instLines + 1 + 3 + 2 + 1 + 1
	
available := m.height - overhead
	if available < 5 {
		available = 5
	}
	// Ensure odd number for centering
	if available%2 == 0 {
		available--
	}
	m.visibleRows = available
}

// View renders the player UI
func (m *PlayerModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if !m.ready {
		return "Loading module...\n"
	}

	metadata := m.module.GetMetadata()
	header := RenderHeader(metadata, m.filename, m.currentTime, m.stereoSep, m.palette)
	activeInstruments := RenderInstrumentsCompact(m.instruments, m.activeInstruments, m.palette)

	mutedChannels := make([]bool, m.patternData.NumChannels)
	for i := 0; i < m.patternData.NumChannels; i++ {
		mutedChannels[i] = m.module.IsChannelMuted(i)
	}

	
	vuMeters := RenderVUMeters(m.patternData.ChannelVolumes, mutedChannels, m.width, m.palette)
	pattern := RenderPattern(m.patternData, mutedChannels, m.palette)
	controls := lipgloss.NewStyle().
		Foreground(m.palette.Controls).
		Render("[q] quit  [space] pause  [[ ]] stereo  [1-9,0,-,=] mute  [Shift+] solo")

	var sections []string
	sections = append(sections, header)
	sections = append(sections, "", activeInstruments)
	sections = append(sections, "", vuMeters)
	sections = append(sections, pattern, "", controls)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Close cleans up resources
func (m *PlayerModel) Close() {
	m.cancel()
	if m.player != nil {
		m.player.Close()
	}
	if m.module != nil {
		m.module.Close()
	}
}