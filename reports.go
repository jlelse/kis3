package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jordan-wright/email"
	"github.com/whiteshtef/clockwork"
	"io/ioutil"
	"net"
	"net/http"
	"net/smtp"
)

type report struct {
	Name     string `json:"name"`
	Time     string `json:"time"`
	Query    string `json:"query"`
	Type     string `json:"type"`
	To       string `json:"to"`
	TGUserId int64  `json:"tgUserId"`
}

func startReports() {
	scheduler := clockwork.NewScheduler()
	for _, r := range appConfig.Reports {
		scheduledReport := r
		scheduler.Schedule().Every().Day().At(scheduledReport.Time).Do(func() {
			executeReport(&scheduledReport)
		})
	}
	scheduler.Run()
}

func executeReport(r *report) {
	fmt.Println("Execute report:", r.Name)
	req, e := http.NewRequest("GET", "http://localhost:"+appConfig.Port+"/stats?"+r.Query, nil)
	if e != nil {
		fmt.Println("Executing report failed:", e)
		return
	}
	req.SetBasicAuth(appConfig.StatsUsername, appConfig.StatsPassword)
	res, e := http.DefaultClient.Do(req)
	if e != nil {
		fmt.Println("Executing report failed:", e)
		return
	}
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		fmt.Println("Executing report failed:", e)
		return
	}
	sendReport(r, body)
}

func sendReport(r *report, content []byte) {
	switch r.Type {
	case "telegram":
		sendTelegram(r, content)
	default:
		sendMail(r, content)
	}
}

func sendMail(r *report, content []byte) {
	if r.To == "" || appConfig.SmtpFrom == "" || appConfig.SmtpUser == "" || appConfig.SmtpHost == "" {
		fmt.Println("No valid report configuration")
		return
	}
	smtpHostNoPort, _, _ := net.SplitHostPort(appConfig.SmtpHost)
	mail := email.NewEmail()
	mail.From = appConfig.SmtpFrom
	mail.To = []string{r.To}
	mail.Subject = "KISSS report: " + r.Name
	mail.Text = content
	e := mail.Send(appConfig.SmtpHost, smtp.PlainAuth("", appConfig.SmtpUser, appConfig.SmtpPassword, smtpHostNoPort))
	if e != nil {
		fmt.Println("Sending report failed:", e)
		return
	} else {
		fmt.Println("Report sent")
	}
}

func sendTelegram(r *report, content []byte) {
	if r.TGUserId == 0 || app.tgBot == nil {
		fmt.Println("No valid report configuration")
		return
	}
	msg := tgbotapi.NewMessage(r.TGUserId, r.Name+"\n\n"+string(content))
	_, e := app.tgBot.Send(msg)
	if e != nil {
		fmt.Println("Sending report failed:", e)
		return
	} else {
		fmt.Println("Report sent")
	}
}
