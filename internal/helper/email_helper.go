package helper

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
)

type EmailHelper struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

func NewEmailHelper(host, port, username, password, fromEmail string) *EmailHelper {
	return &EmailHelper{
		SMTPHost:     host,
		SMTPPort:     port,
		SMTPUsername: username,
		SMTPPassword: password,
		FromEmail:    fromEmail,
	}
}

func (e *EmailHelper) GenerateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (e *EmailHelper) SendVerificationEmail(toEmail, username, token, baseURL string) error {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)

	subject := "Email Verification - Wastetrack"
	body := fmt.Sprintf(`
		Hi %s,
		
		Welcome to Wastetrack! Please verify your email address by clicking the link below:
		
		%s
		
		This link will expire in 24 hours.
		
		If you didn't create an account, please ignore this email.
		
		Best regards,
		Wastetrack Team
	`, username, verificationURL)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailHelper) SendPasswordResetEmail(toEmail, username, token, baseURL string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	subject := "Password Reset - Wastetrack"
	body := fmt.Sprintf(`
		Hi %s,
		
		You requested a password reset for your Wastetrack account. Click the link below to reset your password:
		
		%s
		
		This link will expire in 1 hour.
		
		If you didn't request this, please ignore this email.
		
		Best regards,
		Wastetrack Team
	`, username, resetURL)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailHelper) SendEmailChangeConfirmation(toEmail, username, token, baseURL string) error {
	confirmationURL := fmt.Sprintf("%s/confirm-email-change?token=%s", baseURL, token)

	subject := "Email Change Confirmation - Wastetrack"
	body := fmt.Sprintf(`
		Hi %s,
		
		You requested to change your email address for your Wastetrack account. 
		
		To confirm this change and update your email to: %s
		
		Please click the link below:
		
		%s
		
		This link will expire in 1 hour.
		
		If you didn't request this email change, please ignore this email and your account will remain unchanged.
		
		Important: After confirming, you'll need to verify this new email address to regain full access to your account.
		
		Best regards,
		Wastetrack Team
	`, username, toEmail, confirmationURL)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailHelper) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	addr := e.SMTPHost + ":" + e.SMTPPort
	return smtp.SendMail(addr, auth, e.FromEmail, []string{to}, []byte(msg))
}
