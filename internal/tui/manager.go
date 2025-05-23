package tui

import (
	"fmt"
	"github.com/Nolions/s3Viewer/internal/aws"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"path"
)

var selectedFile aws.FileInfo

func (appCTX *S3App) ManagerLayout() *tview.Flex {
	consoleLayout = appCTX.ConsoleLayout()
	btnLayout := appCTX.ButtonsLayout(consoleLayout)

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

	// Asynchronously load root prefixes and files
	go appCTX.loadPrefixes(console, prefixTreeView, fileListView)

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

func (appCTX *S3App) ButtonsLayout(console *tview.TextView) *tview.Flex {
	inputField := tview.NewInputField().SetLabel("Upload Path: ").SetFieldWidth(55)
	selectBtn := tview.NewButton("Select").SetSelectedFunc(func() {
		appCTX.Pages.ShowPage("filepicker")
		appCTX.Pages.SendToFront("filepicker")
		appCTX.App.SetFocus(filePicker) // 可選
	})
	uploadBtn := tview.NewButton("Upload").SetSelectedFunc(func() {
	})
	downloadBtn := tview.NewButton("Download").SetSelectedFunc(func() {
		if selectedFile.Key != "" && selectedFile.Name != "" {
			appCTX.Pages.ShowPage("dirPicker")
			appCTX.Pages.SendToFront("dirPicker")
			appCTX.App.SetFocus(dirPicker) // 可選
		} else {
			console.SetText("no select file")
		}
	})
	//deleteBtn := tview.NewButton("Delete").SetSelectedFunc(func() {
	//})
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
		//AddItem(deleteBtn, 10, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(exitBtn, 10, 0, false).
		AddItem(tview.NewBox(), 2, 0, false)

	layout.SetBorder(true).SetTitle("Buttons")

	// Focus 切換處理
	focusables := []tview.Primitive{
		inputField, selectBtn, uploadBtn, downloadBtn, exitBtn,
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
	placeholder := tview.NewTreeNode("Loading...").
		SetReference("").
		SetColor(tcell.ColorGray)

	tree := tview.NewTreeView().SetRoot(placeholder).SetCurrentNode(placeholder)
	tree.SetBorder(true).SetTitle("Prefixes")

	return tree
}

func (appCTX *S3App) FileListLayout() *tview.List {
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Files")

	return list
}

func (appCTX *S3App) loadPrefixes(console *tview.TextView, prefixTreeView *tview.TreeView, fileListView *tview.List) {
	res, err := appCTX.S3Client.ListPrefix("")
	if err != nil {
		appCTX.App.QueueUpdateDraw(func() {
			console.SetText(fmt.Sprintf("[red]ListPrefix error:[-] %v", err))
		})
		return
	}

	appCTX.App.QueueUpdateDraw(func() {
		// Build tree root
		root := tview.NewTreeNode(appCTX.AwsConf.Bucket).SetColor(tcell.ColorGreen).SetReference("")
		// Add top-level directories
		for _, dir := range res.Dirs {
			if dir == "" {
				continue
			}
			child := tview.NewTreeNode(dir).SetReference(dir + "/")
			root.AddChild(child)
		}
		prefixTreeView.SetRoot(root).SetCurrentNode(root)

		// Populate file list
		fileListView.Clear()
		for _, f := range res.Files {
			fileListView.AddItem(f.Name, "", 0, nil)
		}
		console.SetText(fmt.Sprintf("Loaded %d dirs, %d files", len(res.Dirs), len(res.Files)))

		// Drill down: when a prefix is selected, list sub dirs and files
		prefixTreeView.SetSelectedFunc(func(node *tview.TreeNode) {
			ref := node.GetReference()
			if ref == nil {
				return
			}
			prefix := ref.(string)
			go appCTX.loadSubDirs(prefix, node, fileListView, console)
		})
	})
}

func (appCTX *S3App) loadSubDirs(
	prefix string,
	node *tview.TreeNode,
	fileListView *tview.List,
	console *tview.TextView,
) {
	res, err := appCTX.S3Client.ListPrefix(prefix)
	appCTX.App.QueueUpdateDraw(func() {
		if err != nil {
			console.SetText(fmt.Sprintf("[red]ListPrefix error:[-] %v", err))
			return
		}
		node.ClearChildren()
		for _, d := range res.Dirs {
			if d == "" {
				continue
			}
			childPrefix := path.Join(prefix, d) + "/"
			child := tview.NewTreeNode(d).
				SetReference(childPrefix)
			node.AddChild(child)
		}
		node.SetExpanded(true)

		fileListView.Clear()
		for _, f := range res.Files {
			fileListView.AddItem(f.Name, "", 0, func() {
				selectedFile = aws.FileInfo{
					Name: f.Name,
					Key:  f.Key,
					Size: f.Size,
					Time: f.Time,
				}
			})
		}
		//console.SetText(fmt.Sprintf("Under '%s': %d dirs, %d files", prefix, len(res.Dirs), len(res.Files)))
	})
}
