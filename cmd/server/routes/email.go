package routes

import (
	"github.com/rs/zerolog/log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendEmail(name string, address string, plainTextContent string, htmlContent string) error {
	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail(name, address)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")) //@todo this is the only os.Getenv()... should we use?

	response, err := client.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("failed to send")
		return err
	}
	log.Debug().Int("Status Code", response.StatusCode).Str("Body", response.Body).Msg("email success")

	return nil
}
