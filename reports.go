package main

import (
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/whiteshtef/clockwork"
	"io/ioutil"
	"net"
	"net/http"
	"net/smtp"
)

func setupReports() {
	scheduler := clockwork.NewScheduler()
	for _, r := range appConfig.Reports {
		scheduledReport := r
		scheduler.Schedule().Every().Day().At(scheduledReport.Time).Do(func() {
			executeReport(&scheduledReport)
		})
	}
	go scheduler.Run()
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
	sendMail(r, body)
}

func sendMail(r *report, content []byte) {
	smtpHostNoPort, _, _ := net.SplitHostPort(r.SmtpHost)
	mail := email.NewEmail()
	mail.From = r.From
	mail.To = []string{r.To}
	mail.Subject = "KISSS report: " + r.Name
	mail.Text = content
	e := mail.Send(r.SmtpHost, smtp.PlainAuth("", r.SmtpUser, r.SmtpPassword, smtpHostNoPort))
	if e != nil {
		fmt.Println("Sending report failed:", e)
		return
	} else {
		fmt.Println("Report sent")
	}
}
