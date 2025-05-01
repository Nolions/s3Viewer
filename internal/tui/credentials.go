package tui

import (
	"github.com/Nolions/s3Viewer/internal/aws"
	"github.com/rivo/tview"
)

func (appCTX *S3App) CredentialsLayout() *tview.Flex {
	credentialsForm := appCTX.CredentialsForm("manager", func(app *tview.Application) {
		app.Stop()
	})
	return wrapCentered(credentialsForm)
}

func (appCTX *S3App) CredentialsForm(switchTo string, exitFun func(app *tview.Application)) *tview.Form {
	form := tview.NewForm()
	form.AddDropDown("Region", aws.Regions, 1, nil).
		AddInputField("AccessKey", appCTX.AwsConf.AccessKey, 35, nil, nil).
		AddInputField("SecretKey", appCTX.AwsConf.SecretKey, 35, nil, nil).
		AddInputField("Bucket", appCTX.AwsConf.Bucket, 35, nil, nil).
		AddCheckbox("Acl", appCTX.AwsConf.Acl, nil).
		AddButton("Save", func() {
			appCTX.Pages.SwitchToPage(switchTo)
		}).
		AddButton("Reset", func() {
			appCTX.AwsConf.AccessKey = ""
			appCTX.AwsConf.SecretKey = ""
			appCTX.AwsConf.Bucket = ""
			form.GetFormItem(1).(*tview.InputField).SetText("")
			form.GetFormItem(2).(*tview.InputField).SetText("")
			form.GetFormItem(3).(*tview.InputField).SetText("")
		}).
		AddButton("Exit", func() {
			exitFun(appCTX.App)
		})

	form.SetTitle("Credentials").SetBorder(true)

	return form
}
