package mailer

import "testing"

func TestMail_SendSMTPMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@here.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	err := mailer.SendSMTPMessage(msg)
	if err != nil {
		t.Errorf("Error sending message: %s", err)
	}
}

func TestMail_SendUsingChan(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@here.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.Jobs <- msg
	res := <-mailer.Results
	if res.Error != nil {
		t.Errorf("Error sending message: %s", res.Error)
	}

	msg.To = "not_an_email"
	mailer.Jobs <- msg
	res = <-mailer.Results
	if res.Error == nil {
		t.Errorf("Expected error sending message")
	}
}

func TestMail_SendUsingAPI(t *testing.T) {
	msg := Message{
		To:          "you@here.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unknown"
	mailer.APIKey = "abcd123"
	mailer.APIUrl = "http://www.fake.com"

	err := mailer.SendUsingAPI(msg, "unknown")
	if err == nil {
		t.Errorf("Expected error sending message")
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_buildHTMLMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@here.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.buildHTMLMessage(msg)
	if err != nil {
		t.Errorf("Error building message: %s", err)
	}
}

func TestMail_buildPlainMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@here.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.buildPlainTextMessage(msg)
	if err != nil {
		t.Errorf("Error building message: %s", err)
	}
}

func TestMail_Send(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@here.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	err := mailer.Send(msg)
	if err != nil {
		t.Errorf("Error sending message: %s", err)
	}

	mailer.API = "unknown"
	mailer.APIKey = "abcd123"
	mailer.APIUrl = "http://www.fake.com"

	err = mailer.Send(msg)
	if err == nil {
		t.Errorf("Expected error sending message")
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_ChooseAPI(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@here.com",
		Subject:     "Test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unknown"
	mailer.APIKey = "abcd123"
	mailer.APIUrl = "http://www.fake.com"

	err := mailer.ChooseAPI(msg)
	if err == nil {
		t.Errorf("Expected error choosing API")
	}

	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}
