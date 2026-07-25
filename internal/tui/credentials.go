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

	var buildForm func(selectedType string)
	buildForm = func(selectedType string) {
		form.Clear(true)
		appCTX.AwsConf.Type = selectedType

		typeOptions := []string{"AWS S3", "MinIO"}
		typeIdx := 0
		if selectedType == "MinIO" {
			typeIdx = 1
		}

		form.AddDropDown("Type", typeOptions, typeIdx, func(option string, idx int) {
			if option != appCTX.AwsConf.Type {
				appCTX.AwsConf.AccessKey = ""
				appCTX.AwsConf.SecretKey = ""
				appCTX.AwsConf.Bucket = ""
				appCTX.AwsConf.Acl = false
				if option == "MinIO" {
					appCTX.AwsConf.Endpoint = "http://localhost:9000"
				} else {
					appCTX.AwsConf.Endpoint = ""
				}
				buildForm(option)
			}
		})

		if selectedType == "MinIO" {
			appCTX.AwsConf.UsePathStyle = true
			if appCTX.AwsConf.Region == "" {
				appCTX.AwsConf.Region = "us-east-1"
			}
			if appCTX.AwsConf.Endpoint == "" {
				appCTX.AwsConf.Endpoint = "http://localhost:9000"
			}

			form.AddInputField("Host", appCTX.AwsConf.Endpoint, 35, nil, func(text string) {
				appCTX.AwsConf.Endpoint = text
			}).
				AddPasswordField("Username", appCTX.AwsConf.AccessKey, 35, '*', func(text string) {
					appCTX.AwsConf.AccessKey = text
				}).
				AddPasswordField("Password", appCTX.AwsConf.SecretKey, 35, '*', func(text string) {
					appCTX.AwsConf.SecretKey = text
				}).
				AddInputField("Bucket", appCTX.AwsConf.Bucket, 35, nil, func(text string) {
					appCTX.AwsConf.Bucket = text
				})
		} else {
			appCTX.AwsConf.Endpoint = ""
			appCTX.AwsConf.UsePathStyle = false

			regionIdx := 15
			for i, r := range aws.Regions {
				if r == appCTX.AwsConf.Region {
					regionIdx = i
					break
				}
			}

			form.AddDropDown("Region", aws.Regions, regionIdx, func(text string, idx int) {
				appCTX.AwsConf.Region = aws.Regions[idx]
			}).
				AddPasswordField("AccessKey", appCTX.AwsConf.AccessKey, 35, '*', func(text string) {
					appCTX.AwsConf.AccessKey = text
				}).
				AddPasswordField("SecretKey", appCTX.AwsConf.SecretKey, 35, '*', func(text string) {
					appCTX.AwsConf.SecretKey = text
				}).
				AddInputField("Bucket", appCTX.AwsConf.Bucket, 35, nil, func(text string) {
					appCTX.AwsConf.Bucket = text
				}).
				AddCheckbox("Acl", appCTX.AwsConf.Acl, func(checked bool) {
					appCTX.AwsConf.Acl = checked
				})
		}

		form.AddButton("OK", func() {
			s3c, err := aws.NewS3Client(appCTX.Ctx, *appCTX.AwsConf)
			if err != nil {
				// TODO
				return
			}
			appCTX.S3Client = s3c

			err = s3c.CheckHeadBucket()
			if err != nil {
				// TODO
			} else {
				managerPage = appCTX.ManagerLayout()
				appCTX.Pages.AddPage("manager", managerPage, true, false)
				appCTX.Pages.SwitchToPage(switchTo)
			}
		}).
			AddButton("Reset", func() {
				appCTX.AwsConf.AccessKey = ""
				appCTX.AwsConf.SecretKey = ""
				appCTX.AwsConf.Bucket = ""
				if selectedType == "MinIO" {
					appCTX.AwsConf.Endpoint = "http://localhost:9000"
				}
				buildForm(selectedType)
			}).
			AddButton("Exit", func() {
				exitFun(appCTX.App)
			})
	}

	initialType := appCTX.AwsConf.Type
	if initialType == "" {
		initialType = "AWS S3"
	}
	buildForm(initialType)

	form.SetTitle("Credentials").SetBorder(true)

	return form
}
