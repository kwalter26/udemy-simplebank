package mail

import (
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGmailSender_SendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("../", false)
	require.NoError(t, err)

	//sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	sender2 := NewSendGridSender("Simple Bank", "no-reply@fusionkoding.com", config.SendGridApiKey)

	subject := "test subject"
	content := `
		<html>
			<body>
				<h1 style="color:red;">test body</h1>
				<p>test paragraph</p>
			</body>
		</html>
	`
	to := []string{"kwalter@fusionkoding.com"}
	attachFiles := []string{"../coverage.out"}
	err = sender2.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
