package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var selectedPath string

// FilePickerModal 把 FilePicker 包成置中 Modal
func FilePickerModal(picker *tview.TreeView, width, height int, closeFunc func(), confirmFunc func()) tview.Primitive {
	confirmBtn := tview.NewButton("Confirm").SetSelectedFunc(func() {
		if confirmFunc != nil {
			confirmFunc()
		}
	})

	btnRow := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(confirmBtn, 10, 1, true).
		AddItem(nil, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(picker, width, 0, true).
				AddItem(nil, 0, 1, false),
			height, 0, true).
		AddItem(btnRow, 1, 0, false).
		AddItem(nil, 0, 1, false)

	// 支援 Esc 關閉
	picker.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			closeFunc()
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
func FilePickerLayout(opt FilePickerOption) *tview.TreeView {
	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle(" File Picker ")

	startDir := opt.StartDir
	if startDir == "" {
		startDir, _ = os.Getwd()
	}

	// 記錄目前瀏覽的位置
	//var currentPath = startDir

	rootNode := tview.NewTreeNode(startDir).SetReference(startDir).SetExpanded(true)
	tree.SetRoot(rootNode).SetCurrentNode(rootNode)

	var addChildren func(node *tview.TreeNode, fullPath string)
	addChildren = func(node *tview.TreeNode, fullPath string) {
		node.ClearChildren()

		// 上層節點
		//parent := filepath.Clean(filepath.Join(fullPath, ".."))
		//upNode := tview.NewTreeNode("[..]").
		//	SetColor(tcell.ColorYellow).
		//	SetReference(parent).
		//	SetSelectable(true)
		//
		//upNode.SetSelectedFunc(func() {
		//	addChildren(node, parent)
		//	node.SetReference(parent)
		//	node.SetExpanded(true)
		//	tree.SetCurrentNode(node)
		//})
		//node.AddChild(upNode)

		// 讀取資料夾內容
		entries, err := ioutil.ReadDir(fullPath)
		if err != nil {
			return
		}

		// 排序：按名稱排序
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name() < entries[j].Name()
		})

		for _, entry := range entries {
			name := entry.Name()
			//childPath := filepath.Join(fullPath, name)
			childPath, _ := filepath.Abs(filepath.Join(fullPath, name))

			if !opt.AllowShowFile && !entry.IsDir() {
				continue
			}

			// 過濾副檔名
			if len(opt.ExtensionFilter) > 0 && !entry.IsDir() {
				matched := false
				for _, ext := range opt.ExtensionFilter {
					if strings.HasSuffix(strings.ToLower(name), strings.ToLower(ext)) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			childNode := tview.NewTreeNode(name).SetReference(childPath)
			p := childPath

			if entry.IsDir() {
				childNode.SetColor(tcell.ColorGreen)
				if opt.AllowFolderSelect {
					childNode.SetSelectedFunc(func() {
						opt.OnSelect(p)
					})
				} else {
					childNode.SetSelectedFunc(func() {
						addChildren(childNode, p)
						childNode.SetExpanded(true)
						tree.SetCurrentNode(childNode)
						opt.OnSelect(p)
					})
				}
			} else {
				childNode.SetColor(tcell.ColorWhite)
				childNode.SetSelectedFunc(func() {
					opt.OnSelect(p)
				})
			}

			node.AddChild(childNode)
		}
	}

	addChildren(rootNode, startDir)

	return tree
}
