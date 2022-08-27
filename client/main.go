package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AuthRes struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type TimeRes struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

type UserAgentRes struct {
	UserAgent string `json:"user-agent"`
	Message   string `json:"message"`
}

type ReverseRes struct {
	Text    string `json:"text"`
	Message string `json:"message"`
}

func main() {
	app := app.New()
	window := app.NewWindow("Client Program")
	//	app.SetIcon(resourceIconPng)
	window.Resize(fyne.NewSize(300, 0))
	window.SetMaster()

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Enter your id ...")
	pwEntry := widget.NewEntry()
	pwEntry.SetPlaceHolder("Enter your password ...")

	loginForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idEntry},
			{Text: "PW", Widget: pwEntry},
		},
	}

	tokenLabel := widget.NewLabel("")

	sepContent := widget.NewSeparator()
	sepContent.Hide()

	bodyArea := widget.NewMultiLineEntry()
	bodyContent := container.NewMax(
		bodyArea,
	)
	bodyContent.Hide()

	reqContent := container.NewGridWithColumns(
		3,
		widget.NewButton("GET /time", func() {
			resp, _ := http.Get("http://localhost:3000/time?token=" + tokenLabel.Text)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			var res TimeRes
			json.Unmarshal(body, &res)

			if resp.StatusCode == http.StatusOK {
				bodyArea.SetText(res.Time)
			} else {
				bodyArea.SetText(res.Message)
			}
		}),

		widget.NewButton("GET /user-agent", func() {
			resp, _ := http.Get("http://localhost:3000/user-agent?token=" + tokenLabel.Text)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			var res UserAgentRes
			json.Unmarshal(body, &res)

			if resp.StatusCode == http.StatusOK {
				bodyArea.SetText(res.UserAgent)
			} else {
				bodyArea.SetText(res.Message)
			}
		}),

		widget.NewButton("POST /reverse", func() {
			params := url.Values{}
			params.Add("text", bodyArea.Text)
			resp, _ := http.PostForm("http://localhost:3000/reverse?token="+tokenLabel.Text, params)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			var res ReverseRes
			json.Unmarshal(body, &res)

			if resp.StatusCode == http.StatusOK {
				bodyArea.SetText(res.Text)
			} else {
				bodyArea.SetText(res.Message)
			}
		}),
	)
	reqContent.Hide()

	authContent := container.NewGridWithColumns(
		3,
		tokenLabel,
		widget.NewLabel(""),
		widget.NewButton("Auth", func() {
			params := url.Values{}
			params.Add("id", idEntry.Text)
			params.Add("pw", pwEntry.Text)
			resp, _ := http.PostForm("http://localhost:3000/auth", params)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			var res AuthRes
			json.Unmarshal(body, &res)

			if resp.StatusCode == http.StatusOK {
				tokenLabel.SetText(res.Token)

				sepContent.Show()
				reqContent.Show()
				bodyContent.Show()
			} else {
				tokenLabel.SetText(res.Message)
			}
		}),
	)

	window.SetContent(container.NewVBox(
		loginForm,
		authContent,
		sepContent,
		reqContent,
		bodyContent,
	))

	window.ShowAndRun()
}
