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

	subject := "Email Verification - Blessing BE"
	body := fmt.Sprintf(`
		Hi %s,
		
		Welcome to Blessing BE! Please verify your email address by clicking the link below:
		
		%s
		
		This link will expire in 24 hours.
		
		If you didn't create an account, please ignore this email.
		
		Best regards,
		The Blessing BE Team
	`, username, verificationURL)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailHelper) SendPasswordResetEmail(toEmail, username, token, baseURL string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	subject := "Password Reset - Blessing BE"
	body := fmt.Sprintf(`
		Hi %s,
		
		You requested a password reset for your Blessing BE account. Click the link below to reset your password:
		
		%s
		
		This link will expire in 1 hour.
		
		If you didn't request this, please ignore this email.
		
		Best regards,
		The Blessing BE Team
	`, username, resetURL)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailHelper) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	addr := e.SMTPHost + ":" + e.SMTPPort
	return smtp.SendMail(addr, auth, e.FromEmail, []string{to}, []byte(msg))
}
