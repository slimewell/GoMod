package ui

import (
	"fmt"
	"github.com/slimewell/GoMod/internal/player"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Pattern command types
const (
	CommandNote        = 0
	CommandInstrument  = 1
	CommandVolume      = 2
	CommandEffect      = 3
	CommandEffectParam = 4
)

// RenderPattern renders the pattern view using a pre-fetched snapshot
// Muted channels are shown dimmed
func RenderPattern(snapshot player.PatternSnapshot, mutedChannels []bool, palette ColorPalette) string {
	if len(snapshot.Rows) == 0 {
		return "No pattern data available"
	}

	currentPattern := snapshot.CurrentPattern
	currentRow := snapshot.CurrentRow
	channels := snapshot.NumChannels

	// Base styles
	rowStyle := lipgloss.NewStyle().Foreground(palette.RowNumber)
	noteStyle := lipgloss.NewStyle().Foreground(palette.Note)
	instrumentStyle := lipgloss.NewStyle().Foreground(palette.Instrument)
	volumeStyle := lipgloss.NewStyle().Foreground(palette.Volume)
	effectStyle := lipgloss.NewStyle().Foreground(palette.Effect)
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Dimmed gray

	// Highlight styles (pre-calculated with background)
	hlRowStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(palette.CurrentRow).
		Background(palette.CurrentRowBg)

	hlNoteStyle := noteStyle.Copy().Background(palette.CurrentRowBg)
	hlInstStyle := instrumentStyle.Copy().Background(palette.CurrentRowBg)
	hlVolStyle := volumeStyle.Copy().Background(palette.CurrentRowBg)
	hlEffStyle := effectStyle.Copy().Background(palette.CurrentRowBg)
	hlSepStyle := lipgloss.NewStyle().Foreground(palette.RowNumber).Background(palette.CurrentRowBg)

	var lines []string

	// Header
	header := fmt.Sprintf("%4s │", "Row")
	for ch := 0; ch < channels; ch++ {
		header += fmt.Sprintf(" Ch%-11d │", ch+1)
	}
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(header))
	lines = append(lines, strings.Repeat("─", len(header)))

	// Render rows
	for _, rowStr := range snapshot.Rows {
		row := rowStr.RowNumber
		isCurrent := (row == currentRow)

		// Select styles for this row
		var (
			rStyle, nStyle, iStyle, vStyle, eStyle, sepStyle lipgloss.Style
		)

		if isCurrent {
			rStyle = hlRowStyle
			nStyle = hlNoteStyle
			iStyle = hlInstStyle
			vStyle = hlVolStyle
			eStyle = hlEffStyle
			sepStyle = hlSepStyle
		} else {
			rStyle = rowStyle
			nStyle = noteStyle
			iStyle = instrumentStyle
			vStyle = volumeStyle
			eStyle = effectStyle
			sepStyle = lipgloss.NewStyle().Foreground(palette.RowNumber) // Default separator
		}

		// 1. Row Number
		var rowNumStr string
		if len(rowStr.Channels) == 0 {
			rowNumStr = "    " // Empty if out of bounds
		} else {
			rowNumStr = fmt.Sprintf("%04X", row)
		}
		// Render the row number cell
		renderedRow := rStyle.Render(rowNumStr)

		// 2. Channel Data
		var renderedChannels []string

		if len(rowStr.Channels) == 0 {
			// Out of bounds ROW -> Blank channels
			// But if it's the current row (unlikely for out of bounds, but possible logic), we might want highlight?
			// Actually, typewriter scrolling means currentRow is always valid, so empty rows are never highlighted.
			// Just render spaces.

			// PatternCell width calculation:
			// Note (3) + space + Inst (2) + space + Vol (2) + space + Eff (3) = 13 chars
			emptyContent := strings.Repeat(" ", 13)

			// We need to replicate the separator structure
			for ch := 0; ch < channels; ch++ {
				renderedChannels = append(renderedChannels, sepStyle.Render(emptyContent))
			}
		} else {
			// Normal Valid Row
			for ch := 0; ch < channels; ch++ {
				if ch >= len(rowStr.Channels) {
					break
				}
				cell := rowStr.Channels[ch]

				// Check if this channel is muted
				isMuted := ch < len(mutedChannels) && mutedChannels[ch]

				// Select styles based on mute status
				var chNStyle, chIStyle, chVStyle, chEStyle lipgloss.Style
				if isMuted {
					// Use dimmed style for all components
					chNStyle = mutedStyle
					chIStyle = mutedStyle
					chVStyle = mutedStyle
					chEStyle = mutedStyle
				} else {
					// Use normal or highlighted styles
					chNStyle = nStyle
					chIStyle = iStyle
					chVStyle = vStyle
					chEStyle = eStyle
				}

				nStr := formatNote(cell.Note)
				iStr := formatInstrument(cell.Instrument)
				vStr := formatVolume(cell.Volume)
				eStr := formatEffect(cell.Effect)

				// Render each component with its style (which might have BG)
				// We join them with spaces, which also need the BG if highlighted
				space := sepStyle.Render(" ")

				part := fmt.Sprintf("%s%s%s%s%s%s%s",
					chNStyle.Render(nStr), space,
					chIStyle.Render(iStr), space,
					chVStyle.Render(vStr), space,
					chEStyle.Render(eStr),
				)

				renderedChannels = append(renderedChannels, part)
			}
		}

		// Join everything with separators
		// Note: The separators " | " need the background if highlighted
		sep := sepStyle.Render(" │ ")

		content := strings.Join(renderedChannels, sep)

		// Final line assembly
		// "ROW | CH1 | CH2 |"
		line := fmt.Sprintf("%s%s%s%s", renderedRow, sep, content, sep)

		// Add pattern info if current
		if isCurrent {
			line += fmt.Sprintf(" Pat: %02X", currentPattern)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

var noteNames = []string{"C-", "C#", "D-", "D#", "E-", "F-", "F#", "G-", "G#", "A-", "A#", "B-"}

func formatNote(note int) string {
	if note == 0 || note == 255 {
		return "..."
	}
	if note == 254 {
		return "===" // Note off
	}

	octave := (note - 1) / 12
	noteIdx := (note - 1) % 12

	return fmt.Sprintf("%s%d", noteNames[noteIdx], octave)
}

func formatInstrument(inst int) string {
	if inst == 0 {
		return ".."
	}
	return fmt.Sprintf("%02X", inst)
}

func formatVolume(vol int) string {
	if vol == 0 {
		return ".."
	}
	return fmt.Sprintf("%02X", vol)
}

func formatEffect(eff int) string {
	if eff == 0 {
		return "..."
	}
	return fmt.Sprintf("%03X", eff)
}
