package tui

import (
	"github.com/rivo/tview"
)

func AlertModel(
	title string,
	content string,
	pages *tview.Pages,
	pageName string,
	switchFun func(pages *tview.Pages, pageName string),
) *tview.Modal {
	m := tview.NewModal().
		SetText(content).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(i int, v string) {
			switch i {
			case 0:
				switchFun(pages, pageName)
			}
		})

	m.SetTitle(title)

	return m
}
