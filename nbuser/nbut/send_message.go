package nbut

import (
	"crypto/tls"
	"errors"
	"net/smtp"

	"github.com/NoneBorder/tasker"
	"github.com/astaxie/beego"
)

const (
	SendMessageTypeEmail = "email"
)

type SendMessage struct {
	Type  string // message type
	To    string // message receiver
	Title string // message title
	Body  string // message body
}

func init() {
	tasker.RegisterTask(new(SendMessage))
}

func (self *SendMessage) New() tasker.MsgQ {
	return &SendMessage{}
}

func (self *SendMessage) Topic() string {
	return "send_message"
}

func (self *SendMessage) TaskSpec() string {
	return "*/2 * * * * *"
}

func (self *SendMessage) Concurency() int {
	return 3
}

func (self *SendMessage) Exec(workerID uint64) (err error) {
	switch self.Type {
	case SendMessageTypeEmail:
		return self.SendEmail()

	default:
		return errors.New("unkown message type " + self.Type)
	}
	return nil
}

func (self *SendMessage) SendEmail() error {
	auth := smtp.PlainAuth("",
		beego.AppConfig.String("notifySMTPUsername"),
		beego.AppConfig.String("notifySMTPPassword"),
		beego.AppConfig.String("notifySMTPHostname"),
	)

	msg := []byte("From:" + beego.AppConfig.String("notifySMTPUsername") + "\r\n" +
		"To:" + self.To + "\r\n" +
		"Subject:" + self.Title + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		self.Body + "\r\n",
	)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         beego.AppConfig.String("notifySMTPHostname"),
	}

	// TLS connect
	conn, err := tls.Dial("tcp",
		beego.AppConfig.String("notifySMTPHostname")+":"+beego.AppConfig.String("notifySMTPPort"),
		tlsconfig,
	)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, beego.AppConfig.String("notifySMTPHostname"))
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail(beego.AppConfig.String("notifySMTPUsername")); err != nil {
		return err
	}
	if err := c.Rcpt(self.To); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = wc.Write(msg); err != nil {
		return err
	}
	if err = wc.Close(); err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	return c.Quit()
}
