package tui

import (
	"context"
	"fmt"
	"github.com/Nolions/s3Viewer/config"
	"github.com/Nolions/s3Viewer/internal/aws"
	"github.com/rivo/tview"
)

type S3App struct {
	Ctx      context.Context
	App      *tview.Application
	Pages    *tview.Pages
	AwsConf  *config.AWSConfig
	S3Client *aws.S3Client
}

var (
	dirPicker  *tview.TreeView
	filePicker *tview.TreeView
	//filePickerModal tview.Primitive
	credentialsPage *tview.Flex
	managerPage     *tview.Flex
	consoleLayout   *tview.TextView
)

func NewS3App(ctx context.Context, conf *config.AWSConfig) *S3App {
	app := tview.NewApplication()
	app.EnableMouse(true)

	pages := tview.NewPages()

	return &S3App{
		Ctx:     ctx,
		App:     app,
		Pages:   pages,
		AwsConf: conf,
	}
}

func (appCTX *S3App) BuildUI() {
	credentialsPage = appCTX.CredentialsLayout() // credentials 頁面

	dirPicker = FilePickerLayout(FilePickerOption{
		StartDir:          ".",
		AllowFolderSelect: false,
		AllowShowFile:     false,
		ExtensionFilter:   []string{},
		OnSelect: func(path string) {
			selectedPath = path
			//appCTX.Pages.HidePage("filepicker")
		},
	})
	dirPicker.SetBorder(true).SetTitle("Select a dir")

	dirPickerModal := FilePickerModal(dirPicker, 60, 15,
		func() {
			appCTX.Pages.HidePage("dirPicker")
			selectedPath = ""
		},
		func() {
			consoleLayout.SetText(fmt.Sprintf("你按下 Confirm，選擇了：%s", selectedPath))
			downloadPath := selectedPath + "\\" + selectedFile.Name
			err := appCTX.S3Client.DownloadFile(selectedFile.Key, downloadPath)
			consoleLayout.SetText(fmt.Sprintf("你按下 Confirm，選擇了：%s", selectedPath))
			if err != nil {
				consoleLayout.SetText(fmt.Sprintf("download file: %s to %s fail, error:%s", selectedFile.Name, downloadPath, err.Error()))
			} else {
				consoleLayout.SetText(fmt.Sprintf("download file: %s to %s, success, size:%d", selectedFile.Name, downloadPath, selectedFile.Size))
			}
			appCTX.Pages.HidePage("dirPicker")
			selectedPath = ""
		},
	)

	filePicker = FilePickerLayout(FilePickerOption{
		StartDir:          ".",
		AllowFolderSelect: false,
		AllowShowFile:     true,
		ExtensionFilter:   []string{},
		OnSelect: func(path string) {
			selectedPath = path
			//appCTX.Pages.HidePage("filepicker")
		},
	})
	filePicker.SetBorder(true).SetTitle("Select a file")

	filePickerModal := FilePickerModal(filePicker, 60, 15,
		func() {
			appCTX.Pages.HidePage("filepicker")
		},
		func() {
			consoleLayout.SetText(fmt.Sprintf("你按下 Confirm，選擇了：%s", selectedPath))
			appCTX.Pages.HidePage("filepicker")
		},
	)

	appCTX.Pages.AddPage("credentials", credentialsPage, true, true)
	appCTX.Pages.AddPage("filepicker", filePickerModal, true, false)
	appCTX.Pages.AddPage("dirPicker", dirPickerModal, true, false)

	if err := appCTX.App.SetRoot(appCTX.Pages, true).Run(); err != nil {
		panic(err)
	}
}

// 泛用 Focus 切換器
func setFocusOnPage(app *tview.Application, pageName string, focusMap map[string]tview.Primitive) {
	if view, ok := focusMap[pageName]; ok && view != nil {
		app.SetFocus(view)
	}
}

// 通用：將元件置中包住（使用 Flex）
func wrapCentered(content tview.Primitive) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(content, 50, 1, true).
				AddItem(nil, 0, 1, false),
			0, 2, true).
		AddItem(nil, 0, 1, false)
}
