package views

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ThemeRosepine function  
func ThemeRosepine() *huh.Theme {
	var (
		// base = lipgloss.AdaptiveColor{Light: "#fafafa", Dark: "#0f111A"}
		// surface = lipgloss.AdaptiveColor{Light: "#fffaf3", Dark: "#1f1d2e"}
		overlay = lipgloss.AdaptiveColor{Light: "#f2e9e1", Dark: "#26233a"}
		// muted   = lipgloss.AdaptiveColor{Light: "#9893a5", Dark: "#6e6a86"}
		subtle = lipgloss.AdaptiveColor{Light: "#797593", Dark: "#908caa"}
		text   = lipgloss.AdaptiveColor{Light: "#575279", Dark: "#e0def4"}
		love   = lipgloss.AdaptiveColor{Light: "#b4637a", Dark: "#eb6f92"}
		// gold    = lipgloss.AdaptiveColor{Light: "#ea9d34", Dark: "#f6c177"}
		rose          = lipgloss.AdaptiveColor{Light: "#d7827e", Dark: "#ebbcba"}
		pine          = lipgloss.AdaptiveColor{Light: "#286983", Dark: "#31748f"}
		foam          = lipgloss.AdaptiveColor{Light: "#56949f", Dark: "#9ccfd8"}
		iris          = lipgloss.AdaptiveColor{Light: "#907aa9", Dark: "#c4a7e7"}
		highlightLow  = lipgloss.AdaptiveColor{Light: "#f4ede8", Dark: "#21202e"}
		highlightMed  = lipgloss.AdaptiveColor{Light: "#dfdad9", Dark: "#403d52"}
		highlightHigh = lipgloss.AdaptiveColor{Light: "#cecacd", Dark: "#524f67"}
	)
	t := *huh.ThemeBase()
	t.FieldSeparator = lipgloss.NewStyle().SetString("\n")
	f := &t.Focused
	f.Base = f.Base.BorderForeground(highlightMed).
		BorderStyle(lipgloss.OuterHalfBlockBorder()).
		Background(lipgloss.NoColor{})
	f.Title = f.Title.Foreground(iris).Bold(true)
	f.NoteTitle = f.NoteTitle.Foreground(rose).Bold(true)
	f.Description = f.Description.Foreground(subtle).Background(lipgloss.NoColor{})
	f.ErrorIndicator = f.ErrorIndicator.Foreground(love)
	f.ErrorMessage = f.ErrorMessage.SetString(" ").Foreground(love)
	f.Option = f.Option.Foreground(text)
	f.SelectSelector = f.SelectSelector.Foreground(pine).SetString("▍ ")
	f.Option = f.Option.Foreground(text)
	f.MultiSelectSelector = f.MultiSelectSelector.Foreground(pine).SetString("▍ ")
	f.SelectedPrefix = f.SelectedPrefix.Foreground(foam).SetString(" ")
	f.SelectedOption = f.SelectedOption.Foreground(foam)
	f.UnselectedPrefix = f.UnselectedPrefix.Foreground(highlightMed).SetString(" ")
	f.UnselectedOption = f.UnselectedOption.Foreground(text)
	f.FocusedButton = f.FocusedButton.Foreground(text).Background(pine)
	f.BlurredButton = f.BlurredButton.Foreground(text).Background(highlightLow)
	f.Next = f.FocusedButton

	f.TextInput.Text = f.TextInput.Text.BorderBottom(true).
		BorderForeground(highlightLow).
		Background(lipgloss.NoColor{})
	f.TextInput.Cursor = f.TextInput.Cursor.Foreground(highlightHigh)
	f.TextInput.Placeholder = f.TextInput.Placeholder.Foreground(highlightLow).
		Background(lipgloss.NoColor{})
	f.TextInput.Prompt = f.TextInput.Prompt.Foreground(highlightMed)

	t.Blurred = *f
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder()).
		Background(lipgloss.NoColor{})
	t.Blurred.MultiSelectSelector = t.Blurred.MultiSelectSelector.SetString("  ")
	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(overlay)
	t.Help.ShortKey = t.Help.ShortKey.Foreground(highlightHigh)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(highlightMed)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(overlay)
	t.Help.FullKey = t.Help.FullKey.Foreground(highlightHigh)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(highlightMed)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(overlay)

	return &t
}
