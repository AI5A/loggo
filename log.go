package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
)

type QSO struct {
	gorm.Model
	Callsign  string
	Frequency int
	Mode      string
	Sent      int
	Received  int
	Comment   string
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("log.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&QSO{})
	return db
}

func inputForm(app *tview.Application) *tview.Form {
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

	frequencyField := tview.NewInputField().
		SetPlaceholder("Frequency KHz").
		SetFieldWidth(20)

	modeField := tview.NewDropDown().
		SetOptions([]string{"SSB", "CW", "FT8", "FT4", "RTTY", "PSK31", "FM", "AM"}, nil).
		SetCurrentOption(0).
		SetFieldWidth(10)

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

	commentField := tview.NewInputField().
		SetPlaceholder("Comment").
		SetFieldWidth(100)

	form := tview.NewForm().
		AddTextView("", "QSO: ", 4, 2, true, false).
		AddFormItem(callsignField).
		AddFormItem(frequencyField).
		AddFormItem(modeField).
		AddFormItem(sentField).
		AddFormItem(receivedField).
		AddFormItem(commentField).
		SetHorizontal(true)

		/*
			SetFinishedFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					// Show modal with the entered data.
					modal := tview.NewModal()
					text := fmt.Sprintf("Callsign: %s Frequency: %s Mode: %s Sent: %s Received: %s Comment: %s",
						callsignField.GetText(),
						frequencyField.GetText(),
						modeField.GetCurrentOption(),
						sentField.GetText(),
						receivedField.GetText(),
						commentField.GetText(),
					)
					modal.SetText(text).Show()
					app.SetRoot(modal, true)
				}
			})*/

	return form
}

// Return a dictionary of tag data from the comment field text given.
func commentToTags(comment string) map[string]string {
	tags := make(map[string]string)
	re := regexp.MustCompile(`([\w\d-]+)\s*:\s*(?:"([^"]+)"|([^"\s]+))`)
	matches := re.FindAllStringSubmatch(comment, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		if value == "" {
			value = match[3]
		}
		tags[key] = value
	}
	return tags
}

// Color the keys and values of tags in the comment field.
func colorTags(comment string) string {
	re := regexp.MustCompile(`([\w\d-]+)\s*:\s*(?:"([^"]+)"|([^"\s]+))`)
	matches := re.FindAllStringSubmatch(comment, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		if value == "" {
			value = match[3]
		}
		comment = strings.Replace(comment, key, "[red]"+key+"[white]", 1)
		comment = strings.Replace(comment, value, "[green]"+value+"[white]", 1)
	}
	return comment
}

func renderHz(hz int) string {
	if hz < 1000 {
		return fmt.Sprintf("%d Hz", hz)
	} else if hz < 1000000 {
		return fmt.Sprintf("%.3f KHz", float64(hz)/1000)
	} else {
		return fmt.Sprintf("%.3f MHz", float64(hz)/1000000)
	}
}

func addQSOToTable(table *tview.Table, qso QSO) {
	row := table.GetRowCount()
	table.SetCellSimple(row, 0, fmt.Sprintf("%d", qso.ID))
	table.SetCellSimple(row, 1, qso.Callsign)
	table.SetCellSimple(row, 2, renderHz(qso.Frequency))
	table.SetCellSimple(row, 3, qso.Mode)
	table.SetCellSimple(row, 4, fmt.Sprintf("%d", qso.Sent))
	table.SetCellSimple(row, 5, fmt.Sprintf("%d", qso.Received))
	table.SetCellSimple(row, 6, colorTags(qso.Comment))
}

func renderLogTable(db *gorm.DB) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 0).
		SetCell(0, 0, tview.NewTableCell("[::bu]QSO[::]").SetSelectable(false)).
		SetCell(0, 1, tview.NewTableCell("[::bu]Callsign[::]").SetSelectable(false)).
		SetCell(0, 2, tview.NewTableCell("[::bu]Frequency[::]").SetSelectable(false)).
		SetCell(0, 3, tview.NewTableCell("[::bu]Mode[::]").SetSelectable(false)).
		SetCell(0, 4, tview.NewTableCell("[::bu]Sent[::]").SetSelectable(false)).
		SetCell(0, 5, tview.NewTableCell("[::bu]Rcv'd[::]").SetSelectable(false)).
		SetCell(0, 6, tview.NewTableCell("[::bu]Comment[::]").SetExpansion(2).SetSelectable(false))

	// Get all the QSOs from the database.
	var qsos []QSO
	db.Find(&qsos)

	// Add each QSO to the table.
	for _, qso := range qsos {
		addQSOToTable(table, qso)
	}

	return table
}

func clearForm(app *tview.Application, inputForm *tview.Form) {
	for i := 0; i < inputForm.GetFormItemCount(); i++ {
		switch item := inputForm.GetFormItem(i).(type) {
		case *tview.InputField:
			item.SetText("")
		case *tview.DropDown:
			item.SetCurrentOption(0)
		}
	}

	// Take me back to the callsign field.
	inputForm.SetFocus(1)
	app.SetFocus(inputForm.GetFormItem(1))
}

func globalInputHandler(inputForm *tview.Form, app *tview.Application, pages *tview.Pages, db *gorm.DB, table *tview.Table) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		// Handle alt-c to clear the form and focus the callsign field.
		if event.Modifiers()&tcell.ModAlt > 0 && event.Rune() == 'c' {
			clearForm(app, inputForm)
			// Eat the event, otherwise we'll get a "c" in the current field.
			return nil
		}

		// Handle Enter to submit the form.
		if event.Key() == tcell.KeyEnter {
			// If we're in the dropdown, select the current option instead of submitting.
			if inputForm.HasFocus() && inputForm.GetFormItem(3).HasFocus() {
				return event
			}

			// Show modal with the entered data.
			modal := tview.NewModal().
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					pages.RemovePage("modal")
				})

			callsign := inputForm.GetFormItem(1).(*tview.InputField).GetText()
			frequency, err := strconv.ParseFloat(inputForm.GetFormItem(2).(*tview.InputField).GetText(), 64)
			_, mode := inputForm.GetFormItem(3).(*tview.DropDown).GetCurrentOption()
			sent, err := strconv.Atoi(inputForm.GetFormItem(4).(*tview.InputField).GetText())
			received, err := strconv.Atoi(inputForm.GetFormItem(5).(*tview.InputField).GetText())
			comment := inputForm.GetFormItem(6).(*tview.InputField).GetText()

			if err != nil {
				modal.SetText(fmt.Sprintf("Error: %s", err))
				pages.AddPage("modal", modal, true, true)
				pages.SwitchToPage("modal").ShowPage("main")
				return nil
			} else {
				qso := QSO{
					Callsign:  callsign,
					Frequency: int(frequency * 1000),
					Mode:      mode,
					Sent:      sent,
					Received:  received,
					Comment:   comment,
				}

				db = db.Create(&qso)

				if db.Error != nil {
					modal.SetText(fmt.Sprintf("Error saving QSO: %s", db.Error))
					pages.AddPage("modal", modal, true, true)
					pages.SwitchToPage("modal").ShowPage("main")
					return nil
				} else {
					addQSOToTable(table, qso)
					clearForm(app, inputForm)
				}
			}

			return nil
		}

		return event
	}
}

func main() {
	db := initDB()

	app := tview.NewApplication()
	pages := tview.NewPages()

	table := renderLogTable(db)
	form := inputForm(app)
	form.
		SetFocus(1).
		SetInputCapture(globalInputHandler(form, app, pages, db, table))
	grid := tview.NewGrid().
		SetRows(0, 3).
		SetColumns(0).
		SetBorders(false).
		AddItem(table, 0, 0, 1, 1, 0, 0, false).
		AddItem(form, 1, 0, 1, 1, 0, 0, true)

	pages.AddPage("main", grid, true, true)
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
