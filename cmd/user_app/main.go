package main

import (
	"buptcs/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var cli client.Client = client.New()

// 注意：目前只加入登录逻辑并未加入注册逻辑
func main() {
	//create app
	myApp := app.New()

	//创建界面（登录界面+功能界面）
	loginWindow := myApp.NewWindow("loginWindow")
	funcWindow := myApp.NewWindow("funcWindow")

	//login计划布局：文本框（表示欢迎以及登录错误信息）、（文本框）账号（输入框）、（文本框）密码（输入框）、登录按钮，退出按钮。

	//欢迎/显示错误的文本框
	infoLabel := widget.NewLabel("Welcome to the charging app")

	//账号密码输入文本框：
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()

	//创造登录与输出按钮
	loginButton := widget.NewButton("Login", func() {
		//调用函数login（）（未实现）
		if login(username.Text, password.Text) {
			//打开功能界面,关闭登陆界面
			funcWindow.Show()
			loginWindow.Close()
		} else {
			//显示登录错误文本
			infoLabel.SetText("Login failed")
		}
	})

	// register button
	regButton := widget.NewButton("Register", func() {
		if register(username.Text, password.Text) {
			infoLabel.SetText("Registered successfully")
		} else {
			//显示登录错误文本
			infoLabel.SetText("Register failed")
		}
	})

	//创建退出按钮
	CancalButton := widget.NewButton("exit", func() {
		myApp.Quit()
	})

	//创建登录界面布局
	loginForm := container.NewVBox(
		infoLabel,
		widget.NewLabel("Username:"),
		username,
		widget.NewLabel("Password:"),
		password,
		container.NewHBox(
			loginButton,
			regButton,
			CancalButton,
		),
	)

	//func界面计划布局：申请充电按钮，修改充电请求按钮，结束充电按钮，查看充电信息按钮。
	var StartButton, ModificateButton, EndButton *widget.Button

	StartButton = widget.NewButton("Start Charging", func() {
		StartCharging()
	})

	ModificateButton = widget.NewButton("Modification Request", func() {
		ModificationRequest()
	})

	EndButton = widget.NewButton("End Charging", func() {
		EndCharging()
	})

	funcForm := container.NewVBox(
		StartButton,
		ModificateButton,
		EndButton,
		CancalButton,
	)

	//设置界面大小
	loginWindow.Resize(fyne.NewSize(600, 300))
	funcWindow.Resize(fyne.NewSize(600, 300))

	loginWindow.SetContent(loginForm)
	funcWindow.SetContent(funcForm)
	loginWindow.ShowAndRun()
}

func login(username string, passwd string) bool {
	// TODO: 实现登录逻辑
	return cli.RequestLogin(username, passwd)	
}

func register(username string, passwd string) bool {
	return cli.RequestRegister(username, passwd)	
}

func StartCharging() {
	// TODO: 实现开始充电逻辑
}

func ModificationRequest() {
	// TODO: 实现修改逻辑
}

func EndCharging() {
	//TODO: 实现结束充电逻辑
}
