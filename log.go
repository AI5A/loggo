package main

import (
	"flag"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

type QSO struct {
	gorm.Model
	Callsign  string
	Frequency int
	Mode      string
	Sent      int
	Received  int
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LogType interface {
	// Generate the input form for the log type.
	// The contact table is passed in so that the form can filter/edit it.
	inputForm(contactTable *tview.Table) *tview.Form

	// Clear the input form.
	clearForm(inputForm *tview.Form)

	// Commit the QSO to the database.
	commitQSO(inputForm *tview.Form, page *tview.Pages, table *tview.Table)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("log.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&QSO{})
	return db
}

// Return a dictionary of tag data from the comment field text given.
func commentToTags(comment string) map[string]string {
	tags := make(map[string]string)
	re := regexp.MustCompile(`([\w\d-]+):(?:"([^"]+)"|([^"\s]+))`)
	matches := re.FindAllStringSubmatch(comment, -1)
	for _, match := range matches {
		// convert key to lowercase
		key := strings.ToLower(match[1])
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
	re := regexp.MustCompile(`([\w\d-]+):(?:"([^"]+)"|([^"\s]+))`)
	matches := re.FindAllStringSubmatch(comment, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		if value == "" {
			value = match[3]
		}
		comment = strings.Replace(comment, key, "[red]"+key+"[-]", 1)
		comment = strings.Replace(comment, value, "[green]"+value+"[-]", 1)
	}
	return comment
}

func addQSOToTable(table *tview.Table, qso QSO) {
	row := table.GetRowCount()
	table.SetCellSimple(row, 0, fmt.Sprintf("%d", qso.ID))
	table.SetCellSimple(row, 1, qso.CreatedAt.Format("2006-01-02 15:04:05 MST"))
	table.SetCellSimple(row, 2, fmt.Sprintf("[::b]%s[-]", qso.Callsign))
	table.SetCellSimple(row, 3, renderHz(qso.Frequency))
	table.SetCellSimple(row, 4, qso.Mode)
	table.SetCellSimple(row, 5, fmt.Sprintf("%d", qso.Sent))
	table.SetCellSimple(row, 6, fmt.Sprintf("%d", qso.Received))
	table.SetCellSimple(row, 7, colorTags(qso.Comment))
}

func renderLogTable(table *tview.Table, qsos []QSO) *tview.Table {
	table.SetBorders(true).
		SetFixed(1, 0).
		SetCell(0, 0, tview.NewTableCell("[::bu]QSO[::]").SetSelectable(false)).
		SetCell(0, 1, tview.NewTableCell("[::bu]Time/Date[::]").SetSelectable(false)).
		SetCell(0, 2, tview.NewTableCell("[::bu]Callsign[::]").SetSelectable(false)).
		SetCell(0, 3, tview.NewTableCell("[::bu]Frequency[::]").SetSelectable(false)).
		SetCell(0, 4, tview.NewTableCell("[::bu]Mode[::]").SetSelectable(false)).
		SetCell(0, 5, tview.NewTableCell("[::bu]Sent[::]").SetSelectable(false)).
		SetCell(0, 6, tview.NewTableCell("[::bu]Rcv'd[::]").SetSelectable(false)).
		SetCell(0, 7, tview.NewTableCell("[::bu]Comment[::]").SetExpansion(2).SetSelectable(false))

	// Add each QSO to the table.
	for _, qso := range qsos {
		addQSOToTable(table, qso)
	}

	table.ScrollToEnd()

	return table
}

func globalInputHandler(
	log LogType,
	inputForm *tview.Form,
	app *tview.Application,
	pages *tview.Pages,
	db *gorm.DB,
	table *tview.Table) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		// Handle alt-c or escape to clear the form and focus the callsign field.
		if event.Key() == tcell.KeyEscape || (event.Modifiers()&tcell.ModAlt > 0 && event.Rune() == 'c') {
			log.clearForm(inputForm)
			// Eat the event, otherwise we'll get a "c" in the current field.
			return nil
		}

		// Handle Enter to submit the form.
		if event.Key() == tcell.KeyEnter {
			// If we're in the dropdown, select the current option instead of submitting.
			if inputForm.HasFocus() && inputForm.GetFormItem(3).HasFocus() {
				return event
			}

			log.commitQSO(inputForm, pages, table)
			return nil
		}

		return event
	}
}

func main() {
	db := initDB()

	wantsExport := flag.Bool("export-adif", false, "Export the log in ADIF format.")

	// Allow users to optionally limit which contact IDs are exported.
	exportFrom := flag.Int("export-from", 0, "Export contacts from this ID onwards.")
	exportTo := flag.Int("export-to", 0, "Export contacts up to this ID.")

	flag.Parse()

	if *wantsExport {
		fmt.Println("ADIF export from loggo by AI5A")
		fmt.Println("<ADIF_VER:5>3.1.4")
		fmt.Println("<programid:5>loggo")
		fmt.Println("<EOH>")
		var qsos []QSO
		if *exportFrom > 0 && *exportTo > 0 {
			db.Where("id >= ? AND id <= ?", *exportFrom, *exportTo).Find(&qsos)
		} else if *exportFrom > 0 {
			db.Where("id >= ?", *exportFrom).Find(&qsos)
		} else if *exportTo > 0 {
			db.Where("id <= ?", *exportTo).Find(&qsos)
		} else {
			db.Find(&qsos)
		}
		for _, qso := range qsos {
			fmt.Println(qso.ToADIF())
		}
		return
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	var qsos []QSO
	db.Find(&qsos)

	table := tview.NewTable()
	renderLogTable(table, qsos)
	log := GeneralQSOLog{app, db}
	form := log.inputForm(table)
	form.
		SetFocus(1).
		SetInputCapture(globalInputHandler(&log, form, app, pages, db, table))
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
