package ui

import (
	"fmt"
	"github.com/slimewell/GoMod/internal/player"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderInstrumentsCompact renders a compact grid of all instruments, highlighting active ones
func RenderInstrumentsCompact(instruments []player.Instrument, activeInstruments map[int]int, palette ColorPalette) string {
	if len(instruments) == 0 {
		return lipgloss.NewStyle().
			Foreground(palette.InfoLabel).
			Bold(true).
			Render("♪ Instruments: (none)")
	}

	// Styles
	labelStyle := lipgloss.NewStyle().
		Foreground(palette.InfoLabel).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(palette.RowNumber).
		Padding(0, 1)

	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Note).
		Background(palette.CurrentRowBg).
		Padding(0, 1)

	fadingStyle := lipgloss.NewStyle().
		Foreground(palette.InfoValue).
		Padding(0, 1)

	// Build instrument chips
	var lines []string
	var currentLine []string
	maxShow := 24     // Show first 24 instruments (easier to fit in 2-3 lines)
	chipsPerLine := 8 // ~8 instruments per line for readability

	for i, inst := range instruments {
		if i >= maxShow {
			break
		}

		name := inst.Name
		if len(name) > 12 {
			name = name[:9] + "..."
		}

		display := fmt.Sprintf("%02X:%s", inst.ID, name)

		// Check if active
		brightness, isActive := activeInstruments[inst.ID]

		var chip string
		if isActive && brightness > 60 {
			// Fully active (just triggered or sustaining)
			chip = activeStyle.Render("★" + display)
		} else if isActive && brightness > 0 {
			// Fading out
			chip = fadingStyle.Render("○" + display)
		} else {
			// Inactive - empty but dot for no shift
			chip = normalStyle.Render("·" + display)
		}

		currentLine = append(currentLine, chip)

		// Wrap to new line after chipsPerLine instruments
		if len(currentLine) >= chipsPerLine {
			lines = append(lines, "  "+strings.Join(currentLine, " "))
			currentLine = []string{}
		}
	}

	// Add remaining chips if any
	if len(currentLine) > 0 {
		lines = append(lines, "  "+strings.Join(currentLine, " "))
	}

	// Build final output with header
	header := labelStyle.Render("♪ Instruments:")

	if len(lines) == 0 {
		return header + " (none)"
	}

	// Join header and lines
	return header + "\n" + strings.Join(lines, "\n")
}

// RenderInstruments renders the full instrument/sample list (legacy, can be used for detailed view)
func RenderInstruments(instruments []player.Instrument, activeInstruments map[int]int, palette ColorPalette) string {
	if len(instruments) == 0 {
		return ""
	}

	// Create styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.InfoLabel)

	normalStyle := lipgloss.NewStyle().
		Foreground(palette.RowNumber)

	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.Note).
		Background(palette.CurrentRowBg)

	var lines []string
	lines = append(lines, headerStyle.Render(fmt.Sprintf("Instruments (%d)", len(instruments))))
	lines = append(lines, strings.Repeat("─", 40))

	// Show first 20 instruments with scrolling display
	maxShow := 20
	if len(instruments) > maxShow {
		lines = append(lines, normalStyle.Render("(Showing first 20)"))
	}

	for i, inst := range instruments {
		if i >= maxShow {
			break
		}

		// Check if this instrument is currently active
		brightness, isActive := activeInstruments[inst.ID]

		// Format: "01 Sample Name" or "01 Sample Name ★"
		idStr := fmt.Sprintf("%02X", inst.ID)
		nameStr := inst.Name
		if len(nameStr) > 30 {
			nameStr = nameStr[:27] + "..."
		}

		line := fmt.Sprintf("%s  %s", idStr, nameStr)

		// Apply style based on activity
		if isActive && brightness > 0 {
			// Active with varying brightness (0-100)
			// For now, just show as active/inactive
			line = activeStyle.Render(line + " ✦")
		} else {
			line = normalStyle.Render(line)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// TrackActiveInstruments scans pattern data to find which instruments are currently playing
func TrackActiveInstruments(mod *player.Module, currentRow int) map[int]int {
	if mod == nil {
		return nil
	}

	activeMap := make(map[int]int)
	currentPattern := mod.GetCurrentPattern()
	channels := mod.GetMetadata().Channels

	// Check current row for instrument triggers
	for ch := 0; ch < channels; ch++ {
		// Get instrument from current row
		instrument := mod.GetPatternRowChannelCommand(currentPattern, currentRow, ch, 1) // CommandInstrument = 1

		if instrument > 0 {
			// Set brightness to 100 (full) when triggered
			activeMap[instrument] = 100
		}
	}

	return activeMap
}
