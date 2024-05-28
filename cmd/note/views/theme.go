package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ThemeRosepine function  
func ThemeRosepine() *huh.Theme {
	var (
		base = lipgloss.AdaptiveColor{Light: "#fafafa", Dark: "#0f111A"}
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
	t := copyTheme(*huh.ThemeBase())
	t.FieldSeparator = lipgloss.NewStyle().SetString("\n")
	f := &t.Focused
	f.Base.BorderForeground(highlightMed).BorderStyle(lipgloss.OuterHalfBlockBorder())
	f.Title.Foreground(iris).Bold(true)
	f.NoteTitle.Foreground(rose).Bold(true)
	f.Description.Foreground(subtle)
	f.ErrorIndicator.Foreground(love)
	f.ErrorMessage = f.ErrorMessage.SetString(" ").Foreground(love)
	f.Option.Foreground(text)
	f.SelectSelector = f.SelectSelector.Foreground(pine).SetString("▍ ")
	f.Option.Foreground(text)
	f.MultiSelectSelector = f.MultiSelectSelector.Foreground(pine).SetString("▍ ")
	f.SelectedPrefix = f.SelectedPrefix.Foreground(foam).SetString(" ")
	f.SelectedOption.Foreground(foam)
	f.UnselectedPrefix = f.UnselectedPrefix.Foreground(highlightMed).SetString(" ")
	f.UnselectedOption.Foreground(text)
	f.FocusedButton.Foreground(text).Background(pine)
	f.BlurredButton.Foreground(text).Background(highlightLow)
	f.Next = f.FocusedButton

	f.TextInput.Text = f.TextInput.Text.BorderBottom(true).
		BorderForeground(highlightLow).
		Background(base)
	f.TextInput.Cursor.Foreground(highlightHigh)
	f.TextInput.Placeholder.Foreground(highlightLow)
	f.TextInput.Prompt.Foreground(highlightMed)

	t.Blurred = copyFieldStyle(*f)
	t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder()).
		Background(lipgloss.NoColor{})
	t.Blurred.MultiSelectSelector = t.Blurred.MultiSelectSelector.SetString("  ")
	t.Help.Ellipsis.Foreground(overlay)
	t.Help.ShortKey.Foreground(highlightHigh)
	t.Help.ShortDesc.Foreground(highlightMed)
	t.Help.ShortSeparator.Foreground(overlay)
	t.Help.FullKey.Foreground(highlightHigh)
	t.Help.FullDesc.Foreground(highlightMed)
	t.Help.FullSeparator.Foreground(overlay)

	return &t
}

func copyTheme(t huh.Theme) huh.Theme {
	return huh.Theme{
		Form:           t.Form,
		Group:          t.Group,
		FieldSeparator: t.FieldSeparator,
		Blurred:        copyFieldStyle(t.Blurred),
		Focused:        copyFieldStyle(t.Focused),
		Help: help.Styles{
			Ellipsis:       t.Help.Ellipsis,
			ShortKey:       t.Help.ShortKey,
			ShortDesc:      t.Help.ShortDesc,
			ShortSeparator: t.Help.ShortSeparator,
			FullKey:        t.Help.FullKey,
			FullDesc:       t.Help.FullDesc,
			FullSeparator:  t.Help.FullSeparator,
		},
	}
}

func copyFieldStyle(f huh.FieldStyles) huh.FieldStyles {
	return huh.FieldStyles{
		Base:                f.Base,
		Title:               f.Title,
		Description:         f.Description,
		ErrorIndicator:      f.ErrorIndicator,
		ErrorMessage:        f.ErrorMessage,
		SelectSelector:      f.SelectSelector,
		Option:              f.Option,
		MultiSelectSelector: f.MultiSelectSelector,
		SelectedOption:      f.SelectedOption,
		SelectedPrefix:      f.SelectedPrefix,
		UnselectedOption:    f.UnselectedOption,
		UnselectedPrefix:    f.UnselectedPrefix,
		FocusedButton:       f.FocusedButton,
		BlurredButton:       f.BlurredButton,
		TextInput:           copyTextInputStyles(f.TextInput),
		Card:                f.Card,
		NoteTitle:           f.NoteTitle,
		Next:                f.Next,
	}
}

func copyTextInputStyles(t huh.TextInputStyles) huh.TextInputStyles {
	return huh.TextInputStyles{
		Cursor:      t.Cursor,
		Placeholder: t.Placeholder,
		Prompt:      t.Prompt,
		Text:        t.Text,
	}
}
