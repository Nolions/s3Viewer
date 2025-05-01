package tui

import (
	"fmt"
	"github.com/Nolions/s3Viewer/config"
	"github.com/rivo/tview"
)

type S3App struct {
	App     *tview.Application
	Pages   *tview.Pages
	AwsConf *config.AWSConfig
}

func NewS3App(conf *config.AWSConfig) *S3App {
	app := tview.NewApplication()
	app.EnableMouse(true)

	pages := tview.NewPages()

	return &S3App{
		App:     app,
		Pages:   pages,
		AwsConf: conf,
	}
}

func (appCTX *S3App) BuildUI() {
	credentialsPage := appCTX.CredentialsLayout() // credentials 頁面
	managerPage := appCTX.ManagerLayout()         // manager 頁面

	filePicker := FilePickerLayout(FilePickerOption{
		StartDir:          ".",
		AllowFolderSelect: false,
		ExtensionFilter:   []string{},
		OnSelect: func(path string) {
			fmt.Println("你選擇了：", path)
			appCTX.Pages.HidePage("filepicker")
		},
	})
	filePicker.SetBorder(true).SetTitle("Select a file")

	modal := FilePickerModal(filePicker, 60, 15, func() {
		appCTX.Pages.HidePage("filepicker")
	})

	appCTX.Pages.AddPage("credentials", credentialsPage, true, true)
	appCTX.Pages.AddPage("manager", managerPage, true, false)
	appCTX.Pages.AddPage("filepicker", modal, true, false)

	focusMap := map[string]tview.Primitive{
		"credentials": credentialsPage.GetItem(1).(tview.Primitive),
		"manager":     managerPage.GetItem(1).(*tview.Flex),
	}

	setFocusOnPage(appCTX.App, "credentials", focusMap)

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
