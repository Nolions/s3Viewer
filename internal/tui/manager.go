package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (appCTX *S3App) ManagerLayout() *tview.Flex {
	btnLayout := appCTX.ButtonsLayout()
	consoleLayout := appCTX.ConsoleLayout()
	topLayout := appCTX.TopLayout()
	browserLayout := appCTX.BrowserLayout(consoleLayout)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topLayout, 1, 0, false).
		AddItem(browserLayout, 0, 5, true).
		AddItem(btnLayout, 3, 0, false).
		AddItem(consoleLayout, 5, 0, false)

	layout.SetBorder(true)

	// Tab 切換 layout 區塊
	focusables := []tview.Primitive{browserLayout, btnLayout, consoleLayout}
	currentFocus := 0

	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			currentFocus = (currentFocus + 1) % len(focusables)
			appCTX.App.SetFocus(focusables[currentFocus])
			return nil
		}
		return event
	})

	return layout
}

func (appCTX *S3App) TopLayout() *tview.Flex {

	regionLayout := appCTX.LabelLayout("Region", appCTX.AwsConf.Region)
	bucketLayout := appCTX.LabelLayout("Bucket", appCTX.AwsConf.Bucket)

	flex := tview.NewFlex().
		AddItem(regionLayout, 0, 1, true).
		AddItem(bucketLayout, 0, 2, false)

	return flex
}

func (appCTX *S3App) BrowserLayout(console *tview.TextView) *tview.Flex {
	prefixTreeView := appCTX.PrefixTreeLayout()
	fileListView := appCTX.FileListLayout()

	prefixTreeView.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref != nil {
			prefix := ref.(string)
			fileListView.Clear().
				AddItem(fmt.Sprintf("%s_File_1", prefix), "", 0, nil).
				AddItem(fmt.Sprintf("%s_File_2", prefix), "", 0, nil).
				AddItem(fmt.Sprintf("%s_File_3", prefix), "", 0, nil)
			console.SetText(fmt.Sprintf("Selected Prefix: %s", prefix))
		}
	})

	flex := tview.NewFlex().
		AddItem(prefixTreeView, 0, 1, true).
		AddItem(fileListView, 0, 2, false)

	focusables := []tview.Primitive{prefixTreeView, fileListView}
	currentFocus := 0

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			currentFocus = 0
			appCTX.App.SetFocus(focusables[currentFocus])
			return nil
		case tcell.KeyRight:
			currentFocus = 1
			appCTX.App.SetFocus(focusables[currentFocus])
			return nil
		}
		return event
	})

	return flex
}

func (appCTX *S3App) LabelLayout(label, cont string) *tview.TextView {
	return tview.NewTextView().SetText(label + ": " + cont)
}

func (appCTX *S3App) ConsoleLayout() *tview.TextView {
	console := tview.NewTextView().
		SetText("console...").
		SetDynamicColors(true).
		SetScrollable(true)

	console.SetTitle(" Console ").SetBorder(true)

	return console
}

func (appCTX *S3App) ButtonsLayout() *tview.Flex {
	inputField := tview.NewInputField().SetLabel("Upload Path: ").SetFieldWidth(55)
	selectBtn := tview.NewButton("Select").SetSelectedFunc(func() {
		appCTX.Pages.ShowPage("filepicker")
		//app.SetFocus(filePicker) // 可選
	})
	uploadBtn := tview.NewButton("Upload").SetSelectedFunc(func() {
	})
	downloadBtn := tview.NewButton("Download").SetSelectedFunc(func() {
	})
	deleteBtn := tview.NewButton("Delete").SetSelectedFunc(func() {
	})
	exitBtn := tview.NewButton("Exit").SetSelectedFunc(func() {
		appCTX.Pages.SwitchToPage("credentials")
	})

	layout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 2, 0, false).
		AddItem(inputField, 70, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(selectBtn, 10, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(uploadBtn, 10, 0, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(downloadBtn, 12, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(deleteBtn, 10, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(exitBtn, 10, 0, false).
		AddItem(tview.NewBox(), 2, 0, false)

	layout.SetBorder(true).SetTitle("Buttons")

	// Focus 切換處理
	focusables := []tview.Primitive{
		inputField, selectBtn, uploadBtn, downloadBtn, deleteBtn, exitBtn,
	}
	currentFocus := 0

	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			currentFocus = (currentFocus - 1 + len(focusables)) % len(focusables)
			appCTX.App.SetFocus(focusables[currentFocus])
			return nil
		case tcell.KeyRight:
			currentFocus = (currentFocus + 1) % len(focusables)
			appCTX.App.SetFocus(focusables[currentFocus])
			return nil
		}
		return event
	})

	return layout
}

func (appCTX *S3App) PrefixTreeLayout() *tview.TreeView {
	root := tview.NewTreeNode("Prefixes").SetColor(tcell.ColorGreen)
	for i := 1; i <= 5; i++ {
		prefix := fmt.Sprintf("Prefix_%d", i)
		node := tview.NewTreeNode(prefix).SetReference(prefix)
		root.AddChild(node)
	}

	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetBorder(true).SetTitle("Prefixes")

	return tree
}

func (appCTX *S3App) FileListLayout() *tview.List {
	list := tview.NewList().
		AddItem("File_1", "", 0, nil).
		AddItem("File_2", "", 0, nil).
		AddItem("File_3", "", 0, nil)

	list.SetBorder(true).SetTitle("Objects")

	return list
}
