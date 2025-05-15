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
	form.AddDropDown("Region", aws.Regions, 1, func(text string, idx int) { appCTX.AwsConf.Region = aws.Regions[idx] }).
		AddInputField("AccessKey", appCTX.AwsConf.AccessKey, 35, nil, func(text string) { appCTX.AwsConf.AccessKey = text }).
		AddInputField("SecretKey", appCTX.AwsConf.SecretKey, 35, nil, func(text string) { appCTX.AwsConf.SecretKey = text }).
		AddInputField("Bucket", appCTX.AwsConf.Bucket, 35, nil, func(text string) { appCTX.AwsConf.Bucket = text }).
		AddCheckbox("Acl", appCTX.AwsConf.Acl, func(checked bool) { appCTX.AwsConf.Acl = checked }).
		AddButton("Save", func() {
			s3c, err := aws.NewS3Client(appCTX.Ctx, *appCTX.AwsConf)
			if err != nil {
				// TODO
			}
			appCTX.S3Client = s3c

			err = s3c.CheckHeadBucket()
			if err != nil {
				// TODO
			} else {
				newManager := appCTX.ManagerLayout()
				appCTX.Pages.AddPage("manager", newManager, true, false)
				appCTX.Pages.SwitchToPage(switchTo)
			}

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
