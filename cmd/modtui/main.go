package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slimewell/GoMod/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load saved config
	cfg, err := ui.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
		defaultCfg := ui.DefaultConfig()
		cfg = &defaultCfg // Fallback to default
	}

	// Parse command-line flags
	stereoSep := flag.Int("separation", cfg.StereoSep, "Stereo separation percentage (0-100)")
	theme := flag.String("theme", cfg.Theme, "Color theme")
	flag.StringVar(theme, "t", cfg.Theme, "Color theme (shorthand)")
	flag.Parse()

	filename := ""
	if flag.NArg() > 0 {
		filename = flag.Arg(0)

		// Validate file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: File not found: %s\n", filename)
			os.Exit(1)
		}
	}

	// Validate stereo separation range
	if *stereoSep < 0 || *stereoSep > 100 {
		fmt.Fprintf(os.Stderr, "Error: Stereo separation must be between 0 and 100\n")
		os.Exit(1)
	}

	// Save config for next time (only updates startup args, not dynamic file loads yet)
	cfg.Theme = *theme
	cfg.StereoSep = *stereoSep
	if filename != "" {
		cfg.LastUsed = filename
	}
	if err := ui.SaveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save config: %v\n", err)
	}

	// Create and run the TUI
	// NewModel now handles empty filename by opening the browser
	model, err := ui.NewModel(filename, *stereoSep, *theme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing audio: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
