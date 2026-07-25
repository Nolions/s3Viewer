package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FilePickerModal 把 FilePicker 包成置中 Modal
func (appCTX *S3App) FilePickerModal(picker *tview.TreeView, width, height int, closeFunc func(), confirmFunc func()) tview.Primitive {
	confirmBtn := tview.NewButton("OK").SetSelectedFunc(func() {
		if appCTX.selectedPath != "" && confirmFunc != nil {
			confirmFunc()
		}
	})
	confirmBtn.SetBorder(true)

	cancelBtn := tview.NewButton("Cancel").SetSelectedFunc(func() {
		if closeFunc != nil {
			closeFunc()
		}
	})
	cancelBtn.SetBorder(true)

	btnRow := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(confirmBtn, 12, 0, true).
		AddItem(tview.NewBox(), 2, 0, false).
		AddItem(cancelBtn, 12, 0, false).
		AddItem(nil, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(picker, width, 0, true).
				AddItem(nil, 0, 1, false),
			height, 0, true).
		AddItem(btnRow, 3, 0, false).
		AddItem(nil, 0, 1, false)

	// 支援 Esc 關閉
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			if closeFunc != nil {
				closeFunc()
			}
			return nil
		}
		return event
	})

	return flex
}

// FilePickerOption 是 file picker 可配置的參數
type FilePickerOption struct {
	StartDir          string            // 起始目錄
	AllowFolderSelect bool              // 是否允許選資料夾
	AllowShowFile     bool              // 是否允許顯示檔案
	ExtensionFilter   []string          // 允許的副檔名（例如 .txt）
	OnSelect          func(path string) // 當使用者選擇一個檔案或資料夾時觸發

}

// FilePickerLayout 回傳可配置選項的 FilePicker
func (appCTX *S3App) FilePickerLayout(opt FilePickerOption) *tview.TreeView {
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle(" File Picker ")

	startDir := opt.StartDir
	if startDir == "" {
		startDir, _ = os.Getwd()
	}
	startDir, _ = filepath.Abs(startDir)

	rootNode := tview.NewTreeNode(startDir).SetReference(startDir).SetExpanded(true)
	tree.SetRoot(rootNode).SetCurrentNode(rootNode)

	appCTX.selectedPath = startDir

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref != nil {
			appCTX.selectedPath = ref.(string)
			if opt.OnSelect != nil {
				opt.OnSelect(ref.(string))
			}
		}
	})

	appCTX.refreshFileTree(tree, rootNode, startDir, opt)

	return tree
}

func (appCTX *S3App) refreshFileTree(tree *tview.TreeView, rootNode *tview.TreeNode, dir string, opt FilePickerOption) {
	absDir, err := filepath.Abs(dir)
	if err == nil {
		dir = absDir
	}

	appCTX.selectedPath = dir
	if opt.OnSelect != nil {
		opt.OnSelect(dir)
	}

	tree.SetTitle(" File Picker - " + dir)
	rootNode.SetText(dir).SetReference(dir)
	rootNode.ClearChildren()

	parent := filepath.Dir(dir)
	if parent != dir {
		targetParent := parent
		upNode := tview.NewTreeNode("[..]").
			SetColor(tcell.ColorYellow).
			SetReference(targetParent).
			SetSelectable(true).
			SetSelectedFunc(func() {
				appCTX.refreshFileTree(tree, rootNode, targetParent, opt)
			})
		rootNode.AddChild(upNode)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		name := entry.Name()
		childPath := filepath.Join(dir, name)

		if !entry.IsDir() && !opt.AllowShowFile {
			continue
		}

		if !entry.IsDir() && len(opt.ExtensionFilter) > 0 {
			allowed := false
			for _, ext := range opt.ExtensionFilter {
				if strings.HasSuffix(strings.ToLower(name), strings.ToLower(ext)) {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}

		childNode := tview.NewTreeNode(name).
			SetReference(childPath).
			SetSelectable(true)

		targetPath := childPath
		if entry.IsDir() {
			childNode.SetColor(tcell.ColorGreen)
			childNode.SetSelectedFunc(func() {
				appCTX.refreshFileTree(tree, rootNode, targetPath, opt)
			})
		} else {
			childNode.SetColor(tcell.ColorWhite)
			childNode.SetSelectedFunc(func() {
				appCTX.selectedPath = targetPath
				if opt.OnSelect != nil {
					opt.OnSelect(targetPath)
				}
			})
		}

		rootNode.AddChild(childNode)
	}

	tree.SetCurrentNode(rootNode)
}
