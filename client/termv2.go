package client

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TermV2 struct {
	App            *tview.Application
	RequestTable   *tview.Table
	TunnelStatus   *tview.TextView
	RequestChannel chan []string
}

const L = tview.AlignLeft
const C = tview.AlignCenter
const R = tview.AlignRight

func NewTermV2(agent *Agent) *TermV2 {
	term := &TermV2{
		App:            tview.NewApplication(),
		RequestChannel: make(chan []string),
	}

	term.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			agent.Quit <- 1
			term.App.Stop()
			return nil
		case tcell.KeyCtrlC:
			agent.Quit <- 1
			term.App.Stop()
			return nil
		}
		return event
	})

	return term
}

func (t *TermV2) Draw() {
	var pages *tview.Pages = tview.NewPages()
	var mainLayout *tview.Flex = tview.NewFlex().SetDirection(tview.FlexRow)

	t.TunnelStatus = t.StatusView()
	mainLayout.AddItem(t.TunnelStatus, 3, 0, false)

	t.RequestTable = t.RequestView()
	t.AddRequestHeader()
	mainLayout.AddItem(t.RequestTable, 0, 10, true)

	var footer *tview.TextView = tview.NewTextView()
	footer.SetBorder(true)
	footer.SetText("ESC/CTRL+C=Exit").SetTextAlign(L).SetTextColor(tcell.ColorWhite)
	mainLayout.AddItem(footer, 3, 0, false)

	pages.AddPage("main", mainLayout, true, true)
	t.App.SetFocus(t.RequestTable)

	go func() {
		for {
			select {
			case row := <-t.RequestChannel:
				t.InsertRequestRow(row)
			}
		}
	}()

	if err := t.App.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

func (t *TermV2) StatusView() *tview.TextView {
	var tunnelInfo *tview.TextView = tview.NewTextView()
	tunnelInfo.SetBorder(true)
	tunnelInfo.SetText("Status: Offline  Address: None").
		SetTextAlign(C).SetTextColor(tcell.ColorRed)

	return tunnelInfo
}

func (t *TermV2) SetTunnelStatus(status string, address string) {
	text := fmt.Sprintf("Status: %s Address: %s", status, address)
	t.App.QueueUpdateDraw(func() {
		color := tcell.ColorWhite
		if status == "Offline" {
			color = tcell.ColorRed
		}

		t.TunnelStatus.SetText(text).SetTextAlign(C).SetTextColor(color)
	})
}

func (t *TermV2) RequestView() *tview.Table {
	var table *tview.Table = tview.NewTable().SetFixed(1, 6).
		SetSelectable(true, false)

	table.SetBorder(true).SetTitle("  Requests ").SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(tcell.ColorDarkOrange)

	return table
}

func (t *TermV2) LogView() {
	var logView *tview.List = tview.NewList()
	logView.SetBorder(true).SetBorderColor(tcell.ColorOrange).
		SetTitle(" Logs ").SetBorderPadding(0, 0, 1, 1)

}

func (t *TermV2) AddRequestHeader() {
	headers := []string{"Method", "Code", "URL"}
	t.insertRequestRow([][]string{headers}, tcell.ColorBlue, true)
}

func (t *TermV2) insertRequestRow(rows [][]string, color tcell.Color, selectable bool) {
	lastRow := t.RequestTable.GetRowCount()
	expansions := []int{1, 1, 3}
	alignment := []int{L, L, L}

	for row, line := range rows {
		for col, text := range line {
			cell := tview.NewTableCell(text).
				SetAlign(alignment[col]).
				SetExpansion(expansions[col]).
				SetTextColor(color).
				SetSelectable(selectable)
			t.RequestTable.SetCell(row+lastRow, col, cell)
		}
	}
	t.App.Sync()
}

func (t *TermV2) InsertRequestRow(row []string) {
	t.App.QueueUpdateDraw(func() {
		statuCode, err := strconv.Atoi(row[1])
		color := tcell.ColorGreen
		if err == nil {
			switch code := statuCode; {
			case code >= 500:
				color = tcell.ColorRed
			case code >= 400 && code <= 500:
				color = tcell.ColorYellow
			}
		}
		t.insertRequestRow([][]string{row}, color, true)
	})
}
