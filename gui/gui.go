package gui

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"school_ruijie/core"
)

type Config struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ServicesPasswd string `json:"servicesPasswd"`
}

func MainWindow() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RuijieLogin")
	myWindow.Resize(fyne.NewSize(300, 200)) // 设置窗口大小

	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	servicesPasswdEntry := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Services Password", Widget: servicesPasswdEntry},
		},
		OnSubmit: func() {
			config := Config{
				Username:       usernameEntry.Text,
				Password:       passwordEntry.Text,
				ServicesPasswd: servicesPasswdEntry.Text,
			}
			writeConfigToFile(config, "config.json")
			core.ExecLoginRuijie()
			myApp.Quit()
		},
	}

	content := container.NewVBox(
		form,
		widget.NewButton("Quit", func() {
			myApp.Quit()
		}),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
	return
}

func writeConfigToFile(config Config, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("Config has been written to", filename)
}
