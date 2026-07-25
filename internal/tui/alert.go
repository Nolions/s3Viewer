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

func ConfirmModal(
	title string,
	content string,
	confirmLabel string,
	cancelLabel string,
	onConfirm func(),
	onCancel func(),
) *tview.Modal {
	if confirmLabel == "" {
		confirmLabel = "Confirm"
	}
	if cancelLabel == "" {
		cancelLabel = "Cancel"
	}

	m := tview.NewModal().
		SetText(content).
		AddButtons([]string{confirmLabel, cancelLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				if onConfirm != nil {
					onConfirm()
				}
			} else {
				if onCancel != nil {
					onCancel()
				}
			}
		})

	m.SetTitle(title)

	return m
}

