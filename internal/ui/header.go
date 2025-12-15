package ui

import (
	"fmt"
	"github.com/slimewell/GoMod/internal/player"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// RenderHeader creates the metadata header display
func RenderHeader(metadata player.Metadata, filename string, currentTime float64, stereoSep int, palette ColorPalette) string {
	// Create styles with palette
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Title)
	infoStyle := lipgloss.NewStyle().Foreground(palette.InfoLabel)
	valueStyle := lipgloss.NewStyle().Foreground(palette.InfoValue)

	// Format duration
	duration := formatTime(metadata.Duration)
	current := formatTime(currentTime)

	// Build header lines
	title := metadata.Title
	if title == "" {
		title = filename
	}

	// Dynamic info line construction
	var infoParts []string

	// Artist (only if present)
	if metadata.Artist != "" {
		infoParts = append(infoParts, infoStyle.Render("Artist:"), valueStyle.Render(metadata.Artist))
	}

	// Format/Tracker type
	infoParts = append(infoParts, infoStyle.Render("Format:"), valueStyle.Render(metadata.Type))

	// Time / Duration
	infoParts = append(infoParts,
		infoStyle.Render("Time:"),
		valueStyle.Render(fmt.Sprintf("%s / %s", current, duration)),
	)

	// Channels
	infoParts = append(infoParts,
		infoStyle.Render("Ch:"),
		valueStyle.Render(fmt.Sprintf("%d", metadata.Channels)),
	)

	// Stereo
	infoParts = append(infoParts,
		infoStyle.Render("Stereo:"),
		valueStyle.Render(fmt.Sprintf("%d%%", stereoSep)),
	)

	// Join with spacing
	infoLine := ""
	for i, part := range infoParts {
		if i > 0 && i%2 == 0 {
			infoLine += "  " // Spacer between pairs
		} else if i > 0 {
			infoLine += " " // Spacer between label and value
		}
		infoLine += part
	}

	lines := []string{
		headerStyle.Render(fmt.Sprintf("♪ %s", title)),
		infoLine,
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func formatTime(seconds float64) string {
	d := time.Duration(seconds * float64(time.Second))
	minutes := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

func orDefault(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
