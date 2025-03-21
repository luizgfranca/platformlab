package tooleventdisplay

import (
	"slices"
)

type DisplayType string

const (
	DisplayTypePrompt      DisplayType = "prompt"
	DisplayTypeInformation DisplayType = "information"
	DisplayTypeTextBox     DisplayType = "textbox"
    DisplayTypeSelelction  DisplayType = "selection"
)

type DisplayDefinition struct {
	Type        DisplayType         `json:"type"`
	Prompt      *PromptDisplay      `json:"prompt"`
	Information *InformationDisplay `json:"information"`
	TextBox     *TextBoxDisplay     `json:"textBox"`
	Selection   *SelectionDisplay   `json:"selection"`
}

func (d *DisplayDefinition) IsValid() bool {

	if !slices.Contains([]DisplayType{DisplayTypePrompt, DisplayTypeInformation, DisplayTypeTextBox}, d.Type) {
		return false
	}

	if d.Type == DisplayTypePrompt && (d.Prompt == nil || !d.Prompt.IsValid()) {
		return false
	}

	if d.Type == DisplayTypeInformation && (d.Information == nil || !d.Information.IsValid()) {
		return false
	}

	if d.Type == DisplayTypeTextBox && (d.TextBox == nil || !d.TextBox.IsValid()) {
		return false
	}

    if d.Type == DisplayTypeSelelction && (d.Selection == nil || !d.Selection.IsValid()) {
		return false
	}

	return true
}
