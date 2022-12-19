package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

func callsignField() *tview.InputField {
	callsignField := tview.NewInputField()
	callsignField.
		SetPlaceholder("Callsign").
		SetFieldWidth(20).
		SetChangedFunc(func(text string) {
			// This is kind of a hack. We check if the field is already
			// uppercase to prevent infinite recursion. There might be a
			// better way to do this.
			if strings.ToUpper(text) != text {
				callsignField.SetText(strings.ToUpper(text))
			}
		})
	return callsignField
}

func frequencyField() *tview.InputField {
	frequencyField := tview.NewInputField().
		SetPlaceholder("Frequency KHz").
		SetFieldWidth(20)
	return frequencyField
}

func modeField() *tview.DropDown {
	modeField := tview.NewDropDown().
		SetOptions([]string{"SSB", "CW", "FT8", "FT4", "RTTY", "PSK31", "FM", "AM"}, nil).
		SetCurrentOption(0).
		SetFieldWidth(10)
	return modeField
}

func sentField(modeField *tview.DropDown) *tview.InputField {
	sentField := tview.NewInputField()
	sentField.
		SetPlaceholder("Sent RST").
		SetFieldWidth(10).
		SetDoneFunc(func(key tcell.Key) {
			if sentField.GetText() == "" {
				_, mode := modeField.GetCurrentOption()
				if mode == "CW" {
					sentField.SetText("599")
				} else {
					sentField.SetText("59")
				}
			}
		}).
		SetFocusFunc(func() {
			_, mode := modeField.GetCurrentOption()
			if mode == "CW" {
				sentField.SetPlaceholder("599")
			} else {
				sentField.SetPlaceholder("59")
			}
		})
	return sentField
}

func receivedField(modeField *tview.DropDown) *tview.InputField {
	receivedField := tview.NewInputField()
	receivedField.
		SetPlaceholder("Rcv'd RST").
		SetFieldWidth(10).
		SetDoneFunc(func(key tcell.Key) {
			if receivedField.GetText() == "" {
				_, mode := modeField.GetCurrentOption()
				if mode == "CW" {
					receivedField.SetText("599")
				} else {
					receivedField.SetText("59")
				}
			}
		}).
		SetFocusFunc(func() {
			_, mode := modeField.GetCurrentOption()
			if mode == "CW" {
				receivedField.SetPlaceholder("599")
			} else {
				receivedField.SetPlaceholder("59")
			}
		})
	return receivedField
}

func commentField() *tview.InputField {
	commentField := tview.NewInputField().
		SetPlaceholder("Comment").
		SetFieldWidth(100)
	return commentField
}
