package repbak

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

// EmailNotifier sends emails based on repliaction failure
type EmailNotifier struct {
	config *Config
}

// NewEmailNotifier creates a EmailNotifier using the config
func NewEmailNotifier(config *Config) *EmailNotifier {
	return &EmailNotifier{
		config: config,
	}
}

// Notify sends a failure notification
func (n *EmailNotifier) Notify(stat Stat) error {
	message := gomail.NewMessage()
	message.SetHeader("From", n.config.Email.From)
	message.SetHeader("To", n.config.Email.To...)
	message.SetHeader("Subject", n.config.Email.Subject)
	message.Attach(n.config.LogPath)

	return n.send(message)
}

func (n *EmailNotifier) NotifyHistory(statMap map[string][]Stat) error {
	if n.config.Retention < 0 || n.config.Email == nil {
		return nil
	}

	if len(statMap) == 0 {
		return errors.New("Email Notifier: no stats found when trying to send stats email")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", n.config.Email.From)
	message.SetHeader("To", n.config.Email.To...)
	message.SetHeader("Subject", n.config.Email.HistorySubject)

	var emailTmpl *template.Template
	var err error
	if n.config.Email.HistoryTemplate != "" {
		emailTmpl, err = template.ParseFiles(n.config.Email.HistoryTemplate)
		if err != nil {
			return fmt.Errorf("Email Notifier: failed to parse custom email template %s: %w", n.config.Email.HistoryTemplate, err)
		}
	} else {
		tmpl := template.New("history")
		emailTmpl, err = tmpl.Parse(emailTemplate)
		if err != nil {
			return fmt.Errorf("Email Notifier: failed to parse email template: %w", err)
		}
	}

	var tpl bytes.Buffer
	if err := emailTmpl.Execute(&tpl, statMap); err != nil {
		return fmt.Errorf("Email Notifier: failed to execute email template: %w", err)
	}

	message.SetBody("text/html", tpl.String())

	return n.send(message)
}

func (n *EmailNotifier) send(message *gomail.Message) error {
	dialer := gomail.NewDialer(
		n.config.Email.Host,
		n.config.Email.Port,
		n.config.Email.User,
		n.config.Email.Pass,
	)

	if n.config.Email.StartTLS {
		dialer.TLSConfig = &tls.Config{
			ServerName:         n.config.Email.Host,
			InsecureSkipVerify: n.config.Email.InsecureSkipVerify,
		}
	}
	dialer.SSL = n.config.Email.SSL

	err := dialer.DialAndSend(message)
	if err != nil {
		return fmt.Errorf("Email Notifer: failed to send email: %w", err)
	}
	return nil
}

var emailTemplate = `<html>
<head>

<style type="text/css">
.tg {
  border-collapse:separate;
  border-spacing:0;
  border-radius:10px;
  width: 100%;
  font-family:Roboto,"Helvetica Neue",sans-serif;
}

.tg td {
  color:#444;
  font-size:14px;
  overflow:hidden;
  padding:3px 10px 3px 0px;
  word-break:normal;
  border-bottom: 1px solid;
  border-bottom-color: #BDBDBD;
  height: 40px;
}

.tg th {
  background-color:#424242;
  color:#FFFFFF;
  font-family:Roboto,"Helvetica Neue",sans-serif;
  font-size:12px;
  font-weight:bold;
  overflow:hidden;
  padding:3px 10px 3px 0px;
  word-break:normal;
  border:0;
  height: 60px;
}

.tg .tg-title {
  text-align:center;
  font-size: 14px;
}

.tg .tg-data {
  text-align:left;
}

.tg .tg-header {
  text-align:left;
  font-weight: bold;
  color: #9E9E9E;
}

.success {
  color: #2E7D32 !important;
  font-weight: bold;
}

.failure {
  color: #D32F2F !important;
  font-weight: bold;
}
</style>

</head>
<body>

{{ range $name, $stats := . }}
<table class="tg">
        <thead>
                <tr>
                        <th class="tg-title" colspan="6">{{$name}}</th>
                </tr>
        </thead>
        <tbody>
                <tr>
                        <td class="tg-header">Status</td>
                        <td class="tg-header">Start</td>
                        <td class="tg-header">End</td>
                        <td class="tg-header">Duration</td>
                </tr>
                {{ range $stats}}
                <tr>
                        {{if .Success}}
                          <td class="success">Success</td>
                        {{else}}
                          <td class="failure">Failed</td>
                        {{end}}
                        <td class="tg-data">{{.Start}}</td>
                        <td class="tg-data">{{.End}}</td>
                        <td class="tg-data">{{.Duration}}</td>
                </tr>
                {{ end}}
        </tbody>
</table>
<br>
{{ end }}

</body>
</html>`
