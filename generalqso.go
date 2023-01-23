package main

import (
	"fmt"
	"github.com/rivo/tview"
	"gorm.io/gorm"
	"strconv"
)

type GeneralQSOLog struct {
	app *tview.Application
	db  *gorm.DB
}

func (log *GeneralQSOLog) inputForm(contactTable *tview.Table) *tview.Form {
	callsignField := callsignField(log.db, contactTable)
	frequencyField := frequencyField()
	modeField := modeField()
	sentField := sentField(modeField)
	receivedField := receivedField(modeField)
	commentField := commentField()

	return tview.NewForm().
		AddTextView("", "QSO: ", 4, 2, true, false).
		AddFormItem(callsignField).
		AddFormItem(frequencyField).
		AddFormItem(modeField).
		AddFormItem(sentField).
		AddFormItem(receivedField).
		AddFormItem(commentField).
		SetHorizontal(true)
}

func (log *GeneralQSOLog) clearForm(inputForm *tview.Form) {
	for i := 0; i < inputForm.GetFormItemCount(); i++ {
		if i == 2 || i == 3 {
			// Skip frequency and mode
			continue
		}
		switch item := inputForm.GetFormItem(i).(type) {
		case *tview.InputField:
			item.SetText("")
		case *tview.DropDown:
			item.SetCurrentOption(0)
		}
	}

	// Take me back to the callsign field.
	inputForm.SetFocus(1)
	log.app.SetFocus(inputForm.GetFormItem(1))
}

func (log *GeneralQSOLog) commitQSO(inputForm *tview.Form, pages *tview.Pages, table *tview.Table) {
	callsign := inputForm.GetFormItem(1).(*tview.InputField).GetText()
	frequency, err := strconv.ParseFloat(inputForm.GetFormItem(2).(*tview.InputField).GetText(), 64)
	_, mode := inputForm.GetFormItem(3).(*tview.DropDown).GetCurrentOption()
	sent, err := strconv.Atoi(inputForm.GetFormItem(4).(*tview.InputField).GetText())
	received, err := strconv.Atoi(inputForm.GetFormItem(5).(*tview.InputField).GetText())
	comment := inputForm.GetFormItem(6).(*tview.InputField).GetText()

	if err != nil {
		modal := tview.NewModal().
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("modal")
			}).
			SetText(fmt.Sprintf("Error: %s", err))
		pages.AddPage("modal", modal, true, true)
		pages.SwitchToPage("modal").ShowPage("main")
		return
	}

	qso := QSO{
		Callsign:  callsign,
		Frequency: int(frequency * 1000),
		Mode:      mode,
		Sent:      sent,
		Received:  received,
		Comment:   comment,
	}

	if log.db.Create(&qso).Error != nil {
		// TODO: Generalize modal stuff
		modal := tview.NewModal().
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("modal")
			}).
			SetText(fmt.Sprintf("Error: %s", err))
		pages.AddPage("modal", modal, true, true)
		pages.SwitchToPage("modal").ShowPage("main")
		return
	}

	addQSOToTable(table, qso)
	log.clearForm(inputForm)
}
