package server

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

func sendEmail(logger *zap.Logger, name string, address string, plainTextContent string, htmlContent string) error {
	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail(name, address)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")) //@todo this is the only os.Getenv()... should we use?

	response, err := client.Send(message)
	if err != nil {
		logger.Error("failed to send", zap.Error(err))
		return err
	}
	logger.Debug("Email Success",
		zap.Int("Status Code", response.StatusCode),
		zap.String("Body", response.Body),
	)

	return nil
}
