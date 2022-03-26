package pkg

type MailSender struct {
	from        string
	to          string
	password    string
	sizeLimitMb uint16
}

func (*MailSender) sendMail(subject string, filePathPrefix string) {
	println("Email sent")
}
