package tui

import (
	"github.com/Nolions/s3Viewer/internal/aws"
	"github.com/rivo/tview"
)

var ResetCredentialsForm func()

func (appCTX *S3App) CredentialsLayout() *tview.Flex {
	credentialsForm := appCTX.CredentialsForm("manager", func(app *tview.Application) {
		app.Stop()
	})
	return wrapCentered(credentialsForm)
}

func (appCTX *S3App) CredentialsForm(switchTo string, exitFun func(app *tview.Application)) tview.Primitive {
	pages := tview.NewPages()

	awsForm := tview.NewForm()
	minioForm := tview.NewForm()

	var showAccessKey, showSecretKey bool
	var showUsername, showPassword bool

	handleOK := func() {
		s3c, err := aws.NewS3Client(appCTX.Ctx, *appCTX.AwsConf)
		if err != nil {
			return
		}
		appCTX.S3Client = s3c

		err = s3c.CheckHeadBucket()
		if err != nil {
			// CheckHeadBucket error
		} else {
			managerPage = appCTX.ManagerLayout()
			appCTX.Pages.AddPage("manager", managerPage, true, false)
			appCTX.Pages.SwitchToPage(switchTo)
		}
	}

	// Build AWS S3 Form
	var buildAWSForm func()
	buildAWSForm = func() {
		awsForm.Clear(true)
		appCTX.AwsConf.Type = "AWS S3"
		appCTX.AwsConf.Endpoint = ""
		appCTX.AwsConf.UsePathStyle = false

		typeOptions := []string{"AWS S3", "MinIO"}

		awsForm.AddDropDown("Type", typeOptions, 0, func(option string, idx int) {
			if option == "MinIO" {
				appCTX.AwsConf.AccessKey = ""
				appCTX.AwsConf.SecretKey = ""
				appCTX.AwsConf.Bucket = ""
				appCTX.AwsConf.Endpoint = "http://localhost:9000"
				pages.SwitchToPage("minio")
				appCTX.App.SetFocus(minioForm)
			}
		})

		regionIdx := 15
		for i, r := range aws.Regions {
			if r == appCTX.AwsConf.Region {
				regionIdx = i
				break
			}
		}

		awsForm.AddDropDown("Region", aws.Regions, regionIdx, func(text string, idx int) {
			appCTX.AwsConf.Region = aws.Regions[idx]
		})

		if showAccessKey {
			awsForm.AddInputField("AccessKey", appCTX.AwsConf.AccessKey, 35, nil, func(text string) {
				appCTX.AwsConf.AccessKey = text
			})
		} else {
			awsForm.AddPasswordField("AccessKey", appCTX.AwsConf.AccessKey, 35, '*', func(text string) {
				appCTX.AwsConf.AccessKey = text
			})
		}
		awsForm.AddCheckbox("Show AccessKey", showAccessKey, func(checked bool) {
			showAccessKey = checked
			buildAWSForm()
		})

		if showSecretKey {
			awsForm.AddInputField("SecretKey", appCTX.AwsConf.SecretKey, 35, nil, func(text string) {
				appCTX.AwsConf.SecretKey = text
			})
		} else {
			awsForm.AddPasswordField("SecretKey", appCTX.AwsConf.SecretKey, 35, '*', func(text string) {
				appCTX.AwsConf.SecretKey = text
			})
		}
		awsForm.AddCheckbox("Show SecretKey", showSecretKey, func(checked bool) {
			showSecretKey = checked
			buildAWSForm()
		})

		awsForm.AddInputField("Bucket", appCTX.AwsConf.Bucket, 35, nil, func(text string) {
			appCTX.AwsConf.Bucket = text
		}).
			AddCheckbox("Acl", appCTX.AwsConf.Acl, func(checked bool) {
				appCTX.AwsConf.Acl = checked
			})

		awsForm.AddButton("OK", handleOK).
			AddButton("Reset", func() {
				appCTX.AwsConf.AccessKey = ""
				appCTX.AwsConf.SecretKey = ""
				appCTX.AwsConf.Bucket = ""
				showAccessKey = false
				showSecretKey = false
				buildAWSForm()
			}).
			AddButton("Close", func() {
				exitFun(appCTX.App)
			})

		awsForm.SetTitle("Credentials (AWS S3)").SetBorder(true)
	}

	// Build MinIO Form
	var buildMinIOForm func()
	buildMinIOForm = func() {
		minioForm.Clear(true)
		appCTX.AwsConf.Type = "MinIO"
		appCTX.AwsConf.UsePathStyle = true
		if appCTX.AwsConf.Region == "" {
			appCTX.AwsConf.Region = "us-east-1"
		}
		if appCTX.AwsConf.Endpoint == "" {
			appCTX.AwsConf.Endpoint = "http://localhost:9000"
		}

		typeOptions := []string{"AWS S3", "MinIO"}

		minioForm.AddDropDown("Type", typeOptions, 1, func(option string, idx int) {
			if option == "AWS S3" {
				appCTX.AwsConf.AccessKey = ""
				appCTX.AwsConf.SecretKey = ""
				appCTX.AwsConf.Bucket = ""
				appCTX.AwsConf.Endpoint = ""
				pages.SwitchToPage("aws")
				appCTX.App.SetFocus(awsForm)
			}
		})

		minioForm.AddInputField("Host", appCTX.AwsConf.Endpoint, 35, nil, func(text string) {
			appCTX.AwsConf.Endpoint = text
		})

		if showUsername {
			minioForm.AddInputField("Username", appCTX.AwsConf.AccessKey, 35, nil, func(text string) {
				appCTX.AwsConf.AccessKey = text
			})
		} else {
			minioForm.AddPasswordField("Username", appCTX.AwsConf.AccessKey, 35, '*', func(text string) {
				appCTX.AwsConf.AccessKey = text
			})
		}
		minioForm.AddCheckbox("Show Username", showUsername, func(checked bool) {
			showUsername = checked
			buildMinIOForm()
		})

		if showPassword {
			minioForm.AddInputField("Password", appCTX.AwsConf.SecretKey, 35, nil, func(text string) {
				appCTX.AwsConf.SecretKey = text
			})
		} else {
			minioForm.AddPasswordField("Password", appCTX.AwsConf.SecretKey, 35, '*', func(text string) {
				appCTX.AwsConf.SecretKey = text
			})
		}
		minioForm.AddCheckbox("Show Password", showPassword, func(checked bool) {
			showPassword = checked
			buildMinIOForm()
		})

		minioForm.AddInputField("Bucket", appCTX.AwsConf.Bucket, 35, nil, func(text string) {
			appCTX.AwsConf.Bucket = text
		})

		minioForm.AddButton("OK", handleOK).
			AddButton("Reset", func() {
				appCTX.AwsConf.AccessKey = ""
				appCTX.AwsConf.SecretKey = ""
				appCTX.AwsConf.Bucket = ""
				appCTX.AwsConf.Endpoint = "http://localhost:9000"
				showUsername = false
				showPassword = false
				buildMinIOForm()
			}).
			AddButton("Close", func() {
				exitFun(appCTX.App)
			})

		minioForm.SetTitle("Credentials (MinIO)").SetBorder(true)
	}

	buildAWSForm()
	buildMinIOForm()

	ResetCredentialsForm = func() {
		appCTX.AwsConf.AccessKey = ""
		appCTX.AwsConf.SecretKey = ""
		appCTX.AwsConf.Bucket = ""
		showAccessKey = false
		showSecretKey = false
		showUsername = false
		showPassword = false
		buildAWSForm()
		buildMinIOForm()
	}

	pages.AddPage("aws", awsForm, true, true)
	pages.AddPage("minio", minioForm, true, false)

	return pages
}
