package mailer

import (
	"bytes"
	"embed"
	ht "html/template"
	tt "text/template"
	"time"

	"github.com/wneessen/go-mail"
	// Import the html/template and text/template packages. Because these share the same
	// package name ("template") we need disambiguate them and alias them to ht and tt
	// respectively.
)

// Below we declare a new variable with the type embed.FS (embedded file system) to hold
// our email templates. This has a comment directive in the format `//go:embed <path>`
// IMMEDIATELY ABOVE it, which indicates to Go that we want to store the contents of the
// ./templates directory in the templateFS embedded file system variable.
// ↓↓↓
//
//go:embed "templates"
var templateFS embed.FS

// Define a Mailer struct which contains a mail.Client instance (used to connect to a
// SMTP server) and the sender information for your emails (the name and address you
// want the email to be from, such as "Alice Smith <alice@example.com>").
type Mailer struct {
	client *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (*Mailer, error) {
	// Initialize a new mail.Dialer instance with the given SMTP server settings. We
	// also configure this to use a 5-second timeout whenever we send an email. I've
	// split the NewClient arguments over multiple lines for readability, but you can
	// make this a single line if you prefer.
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}
	mailer := &Mailer{
		client: client,
		sender: sender,
	}
	return mailer, nil
}

func (m *Mailer) Send(recipient string, templateFile string, data any) error {
	// Use the ParseFS() method from text/template to parse the required template file
	// from the embedded file system.
	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	// Follow the same pattern to execute the "plainBody" template and store the result
	// in the plainBody variable.
	plainBody := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	// Use the ParseFS() method from html/template this time to parse the required template
	// file from the embedded file system.
	htmlTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	// Use the mail.NewMsg() function to initialize a new mail.Msg instance.
	// Then we use the To(), From() and Subject() methods to set the email recipient,
	// sender and subject headers, the SetBodyString() method to set the plain-text body,
	// and the AddAlternativeString() method to set the HTML body.
	msg := mail.NewMsg()
	err = msg.To(recipient)
	if err != nil {
		return err
	}
	err = msg.From(m.sender)
	if err != nil {
		return err
	}
	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())
	// Call the DialAndSend() method on the dialer, passing in the message to send. This
	// opens a connection to the SMTP server, sends the message, then closes the
	// connection.
	// return m.client.DialAndSend(msg)

	// Try sending the email up to three times before aborting and returning the final
	// error. We sleep for 500 milliseconds between each attempt.
	for i := 1; i <= 3; i++ {
		err = m.client.DialAndSend(msg)
		if err == nil {
			return nil
		}
		// If it didn't work, sleep for a short time and retry.
		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return err
}
