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
	inputForm() *tview.Form

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

func renderLogTable(db *gorm.DB) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 0).
		SetCell(0, 0, tview.NewTableCell("[::bu]QSO[::]").SetSelectable(false)).
		SetCell(0, 1, tview.NewTableCell("[::bu]Time/Date[::]").SetSelectable(false)).
		SetCell(0, 2, tview.NewTableCell("[::bu]Callsign[::]").SetSelectable(false)).
		SetCell(0, 3, tview.NewTableCell("[::bu]Frequency[::]").SetSelectable(false)).
		SetCell(0, 4, tview.NewTableCell("[::bu]Mode[::]").SetSelectable(false)).
		SetCell(0, 5, tview.NewTableCell("[::bu]Sent[::]").SetSelectable(false)).
		SetCell(0, 6, tview.NewTableCell("[::bu]Rcv'd[::]").SetSelectable(false)).
		SetCell(0, 7, tview.NewTableCell("[::bu]Comment[::]").SetExpansion(2).SetSelectable(false))

	// Get all the QSOs from the database.
	var qsos []QSO
	db.Find(&qsos)

	// Add each QSO to the table.
	for _, qso := range qsos {
		addQSOToTable(table, qso)
	}

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
		// Handle alt-c to clear the form and focus the callsign field.
		if event.Modifiers()&tcell.ModAlt > 0 && event.Rune() == 'c' {
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
	flag.Parse()

	if *wantsExport {
		fmt.Println("ADIF export from loggo by AI5A")
		fmt.Println("<ADIF_VER:5>3.1.4")
		fmt.Println("<programid:5>loggo")
		fmt.Println("<EOH>")
		var qsos []QSO
		db.Find(&qsos)
		for _, qso := range qsos {
			fmt.Println(qso.ToADIF())
		}
		return
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	table := renderLogTable(db)
	log := GeneralQSOLog{app, db}
	form := log.inputForm()
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
