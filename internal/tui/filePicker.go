package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var selectedPath string

// FilePickerModal 把 FilePicker 包成置中 Modal
func FilePickerModal(picker *tview.TreeView, width, height int, closeFunc func(), confirmFunc func()) tview.Primitive {
	confirmBtn := tview.NewButton("Confirm").SetSelectedFunc(func() {
		if selectedPath != "" && confirmFunc != nil {
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

	selectedPath = startDir

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref != nil {
			selectedPath = ref.(string)
		}
	})

	refreshFileTree(tree, rootNode, startDir, opt)

	return tree
}

func refreshFileTree(tree *tview.TreeView, rootNode *tview.TreeNode, dir string, opt FilePickerOption) {
	tree.SetTitle(" File Picker - " + dir)
	rootNode.ClearChildren()
	rootNode.SetReference(dir)

	parent := filepath.Dir(dir)
	upNode := tview.NewTreeNode("[..]").
		SetColor(tcell.ColorYellow).
		SetReference(parent).
		SetSelectable(true).
		SetSelectedFunc(func() {
			refreshFileTree(tree, rootNode, parent, opt)
		})
	rootNode.AddChild(upNode)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		name := entry.Name()
		childPath, _ := filepath.Abs(filepath.Join(dir, name))

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

		if entry.IsDir() {
			childNode.SetColor(tcell.ColorGreen)
			childNode.SetSelectedFunc(func(path string) func() {
				return func() {
					refreshFileTree(tree, rootNode, path, opt)
					tree.SetCurrentNode(childNode)
					tree.SetTitle(" File Picker - " + path)
				}
			}(childPath))
		} else {
			childNode.SetColor(tcell.ColorWhite)
			childNode.SetSelectedFunc(func(path string) func() {
				return func() {
					opt.OnSelect(path)
				}
			}(childPath))
		}

		rootNode.AddChild(childNode)
	}
}
