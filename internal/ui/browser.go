package ui

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// parentDirEntry implements os.DirEntry for the ".." entry
type parentDirEntry struct{}

func (e parentDirEntry) Name() string               { return ".." }
func (e parentDirEntry) IsDir() bool                { return true }
func (e parentDirEntry) Type() os.FileMode          { return os.ModeDir }
func (e parentDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

// FileBrowserModel handles file selection
type FileBrowserModel struct {
	CurrentPath string
	Files       []os.DirEntry
	Cursor      int
	Selected    string // The selected file path (set when user presses enter on a file)

	width  int
	height int
	err    error
	
	palette ColorPalette
}

// NewFileBrowserModel creates a new file browser starting at current directory
func NewFileBrowserModel(width, height int) *FileBrowserModel {
	cwd, _ := os.Getwd()
	m := &FileBrowserModel{
		CurrentPath: cwd,
		width:       width,
		height:      height,
		palette:     DefaultPalette(),
	}
	m.refreshFiles()
	return m
}

// Init satisfies tea.Model
func (m *FileBrowserModel) Init() tea.Cmd {
	return nil
}

func (m *FileBrowserModel) refreshFiles() {
	entries, err := os.ReadDir(m.CurrentPath)
	if err != nil {
		m.err = err
		return
	}

	// Filter and Sort
	var filtered []os.DirEntry
	for _, e := range entries {
		// Hide .DS_Store and hidden files (starting with .)
		if e.Name() == ".DS_Store" || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		filtered = append(filtered, e)
	}

	// Sort: Directories first, then files
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].IsDir() && !filtered[j].IsDir() {
			return true
		}
		if !filtered[i].IsDir() && filtered[j].IsDir() {
			return false
		}
		return strings.ToLower(filtered[i].Name()) < strings.ToLower(filtered[j].Name())
	})

	// Prepend ".." if not at root
	// A simple check: if parent dir is different from current
	parent := filepath.Dir(m.CurrentPath)
	if parent != m.CurrentPath {
		filtered = append([]os.DirEntry{parentDirEntry{}}, filtered...)
	}

	m.Files = filtered
	m.Cursor = 0
}

// Update handles browser navigation
func (m *FileBrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Files)-1 {
				m.Cursor++
			}
		case "enter":
			if len(m.Files) == 0 {
				return m, nil
			}
			selected := m.Files[m.Cursor]
			
			// Handle ".."
			if selected.Name() == ".." {
				parent := filepath.Dir(m.CurrentPath)
				if parent != m.CurrentPath {
					m.CurrentPath = parent
					m.refreshFiles()
				}
				return m, nil
			}

			fullPath := filepath.Join(m.CurrentPath, selected.Name())

			if selected.IsDir() {
				m.CurrentPath = fullPath
				m.refreshFiles()
			} else {
				// Select file!
				m.Selected = fullPath
			}
		case "backspace", "left", "h":
			parent := filepath.Dir(m.CurrentPath)
			if parent != m.CurrentPath {
				m.CurrentPath = parent
				m.refreshFiles()
			}
		}
	}
	return m, nil
}

// View renders the file list with a border
func (m *FileBrowserModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error accessing directory: %v", m.err)
	}

	// Styles using DefaultPalette colors
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.palette.Border).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Foreground(m.palette.Title).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(m.palette.Border).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(m.palette.CurrentRowBg). // Dark slate (background-ish) -> Maybe invert?
		Background(m.palette.CurrentRow).   // Light text
		Bold(true)

	// Since we can't easily invert efficiently without defining new colors, let's use:
	// Selected: CurrentRow (Cream) on CurrentRowBg (Dark Slate)
	selectedStyle = lipgloss.NewStyle().
		Foreground(m.palette.CurrentRow).
		Background(m.palette.CurrentRowBg).
		Bold(true)

	dirStyle := lipgloss.NewStyle().
		Foreground(m.palette.InfoLabel). // Slate Gray
		Bold(true)

	modStyle := lipgloss.NewStyle().
		Foreground(m.palette.Title). // Cream (Highlight)
		Bold(false)

	fileStyle := lipgloss.NewStyle().
		Foreground(m.palette.RowNumber) // Medium Slate

	// Calculate view window
	listHeight := m.height - 10 // Increased overhead for help text
	if listHeight < 5 {
		listHeight = 5
	}

	start := m.Cursor - (listHeight / 2)
	if start < 0 {
		start = 0
	}
	end := start + listHeight
	if end > len(m.Files) {
		end = len(m.Files)
	}

	var content strings.Builder
	
	// Header Path (Truncate if too long?)
	pathStr := m.CurrentPath
	if len(pathStr) > m.width-10 {
		// Basic truncation from left
		pathStr = "..." + pathStr[len(pathStr)-(m.width-10):]
	}
	content.WriteString(fmt.Sprintf("%s\n\n", pathStr))

	// List
	for i := start; i < end; i++ {
		entry := m.Files[i]
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}

		cursor := "  "
		if i == m.Cursor {
			cursor = "> "
		}

		line := cursor + name
		// Pad line to full width for nicer selection bar
		lineWidth := m.width - 10 // Approx inner width
		if len(line) < lineWidth {
			line += strings.Repeat(" ", lineWidth-len(line))
		} else {
			line = line[:lineWidth]
		}
		
		var styledLine string
		if i == m.Cursor {
			styledLine = selectedStyle.Render(line)
		} else if entry.IsDir() {
			styledLine = dirStyle.Render(line)
		} else {
			// Check extension for highlight
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".mod" || ext == ".xm" || ext == ".it" || ext == ".s3m" {
				styledLine = modStyle.Render(line)
			} else {
				styledLine = fileStyle.Render(line)
			}
		}

		content.WriteString(styledLine + "\n")
	}
	
	// Footer / Help
	helpText := lipgloss.NewStyle().
		Foreground(m.palette.InfoLabel).
		Italic(true).
		Render("\n[↑/↓] Move  [Enter] Select  [Back] Up Dir  [Tab] Close")

	// Assemble
	return boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			titleStyle.Render(" Select Module "),
			content.String(),
			helpText,
		),
	)
}
