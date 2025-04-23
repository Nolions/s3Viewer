package tui

import (
	"github.com/Nolions/s3Viewer/internal/aws"
	"github.com/rivo/tview"
)

func (appCtx *S3App) CredentialsLayout() *tview.Flex {
	credentialsForm := appCtx.CredentialsForm("manager", func(app *tview.Application) {
		app.Stop()
	})
	return wrapCentered(credentialsForm)
}

func (appCtx *S3App) CredentialsForm(switchTo string, exitFun func(app *tview.Application)) *tview.Form {
	form := tview.NewForm()
	form.AddDropDown("Region", aws.Regions, appCtx.AwsConf.Region, nil).
		AddInputField("AccessKey", appCtx.AwsConf.AccessKey, 35, nil, nil).
		AddInputField("SecretKey", appCtx.AwsConf.SecretKey, 35, nil, nil).
		AddInputField("Bucket", appCtx.AwsConf.Bucket, 35, nil, nil).
		AddCheckbox("Acl", appCtx.AwsConf.Acl, nil).
		AddButton("Save", func() {
			appCtx.Pages.SwitchToPage(switchTo)
		}).
		AddButton("Reset", func() {
			appCtx.AwsConf.AccessKey = ""
			appCtx.AwsConf.SecretKey = ""
			appCtx.AwsConf.Bucket = ""
			form.GetFormItem(1).(*tview.InputField).SetText("")
			form.GetFormItem(2).(*tview.InputField).SetText("")
			form.GetFormItem(3).(*tview.InputField).SetText("")
		}).
		AddButton("Exit", func() {
			exitFun(appCtx.App)
		})

	form.SetTitle("Credentials").SetBorder(true)

	return form
}
