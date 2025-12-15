package ui

import (
	"github.com/slimewell/GoMod/internal/player"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hajimehoshi/oto/v2"
)

type AppState int

const (
	StateBrowsing AppState = iota
	StatePlaying
)

// AppModel is the main application container
type AppModel struct {
	state        AppState
	playerModel  *PlayerModel
	browserModel *FileBrowserModel

	// Shared audio context
	audioContext *oto.Context

	// Global config to persist across module loads
	stereoSep int
	themeName string

	width  int
	height int
}

// NewModel creates the main application model
func NewModel(filename string, stereoSep int, themeName string) (AppModel, error) {
	// Initialize audio context once
	ac, err := player.NewAudioContext()
	if err != nil {
		return AppModel{}, err
	}

	// Initialize with default dimensions
	w, h := 80, 24

	var state AppState
	var pm *PlayerModel

	if filename != "" {
		state = StatePlaying
		pm = NewPlayerModel(ac, filename, stereoSep, themeName, w, h)
	} else {
		state = StateBrowsing
		// Player is nil initially
	}

	return AppModel{
		state:        state,
		playerModel:  pm,
		browserModel: NewFileBrowserModel(w, h),
		audioContext: ac,
		stereoSep:    stereoSep,
		themeName:    themeName,
		width:        w,
		height:       h,
	}, nil
}

// Init initializes the application
func (m AppModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.playerModel != nil {
		cmds = append(cmds, m.playerModel.Init())
	}
	// Browser doesn't need async Init currently
	return tea.Batch(cmds...)
}

// Update handles global messages and routes others to sub-models
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global Key Handling
		switch msg.String() {
		case "q", "ctrl+c":
			// If browsing and we have a player, just close browser (esc behavior)
			// But 'q' is usually quit. Let's make 'q' quit always for now.
			if m.playerModel != nil {
				m.playerModel.Close()
			}
			return m, tea.Quit
		
		case "tab":
			// Toggle browser if we have a player
			if m.playerModel != nil {
				if m.state == StateBrowsing {
					m.state = StatePlaying
				} else {
					m.state = StateBrowsing
					// Update browser size just in case
					m.browserModel.width = m.width
					m.browserModel.height = m.height
				}
				return m, nil
			}
		
		case "esc":
			if m.state == StateBrowsing && m.playerModel != nil {
				m.state = StatePlaying
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.browserModel.width = msg.Width
		m.browserModel.height = msg.Height
		
		// Propagate to player if it exists
		if m.playerModel != nil {
			newPlayer, cmd := m.playerModel.Update(msg)
			m.playerModel = newPlayer.(*PlayerModel)
			cmds = append(cmds, cmd)
		}
		// Also propagate to browser (it uses width/height for layout)
		newBrowser, cmd := m.browserModel.Update(msg)
		m.browserModel = newBrowser.(*FileBrowserModel)
		cmds = append(cmds, cmd)
		
		return m, tea.Batch(cmds...)
	}

	// Route based on state
	if m.state == StateBrowsing {
		newBrowser, cmd := m.browserModel.Update(msg)
		m.browserModel = newBrowser.(*FileBrowserModel)
		cmds = append(cmds, cmd)

        // Fix Frozen UI: Update Player background (Ticks only)
        if m.playerModel != nil {
             switch msg.(type) {
             case tickMsg:
                 newPlayer, pCmd := m.playerModel.Update(msg)
                 m.playerModel = newPlayer.(*PlayerModel)
                 cmds = append(cmds, pCmd)
             }
        }

		// Check if file was selected
		if m.browserModel.Selected != "" {
			filename := m.browserModel.Selected
			m.browserModel.Selected = "" // Reset

			// If we had a player, close it
			if m.playerModel != nil {
				// Save state
				m.stereoSep = m.playerModel.stereoSep
				
				// Close old player
				m.playerModel.Close()
				m.playerModel = nil
			}

			// Create new player with current config and SHARED AUDIO CONTEXT
			m.playerModel = NewPlayerModel(m.audioContext, filename, m.stereoSep, m.themeName, m.width, m.height)
			m.state = StatePlaying
			
			cmds = append(cmds, m.playerModel.Init())
		}
		
		return m, tea.Batch(cmds...)

	} else {
		// StatePlaying
		if m.playerModel != nil {
			newPlayer, cmd := m.playerModel.Update(msg)
			m.playerModel = newPlayer.(*PlayerModel)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the application
func (m AppModel) View() string {
	if m.state == StateBrowsing {
		browserView := m.browserModel.View()

		// Overlay on top of player if player exists
		if m.playerModel != nil {
			// Use lipgloss.Place to center the browser
			// Note: This temporarily hides the player view behind the browser overlay
			return lipgloss.Place(
				m.width, m.height,
				lipgloss.Center, lipgloss.Center,
				browserView,
			)
		}
		
		// No player, just render browser centered
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			browserView,
		)
	}

	// StatePlaying
	if m.playerModel != nil {
		return m.playerModel.View()
	}

	return "No module loaded."
}

// errMsg is shared within the ui package
type errMsg struct {
	err error
}
