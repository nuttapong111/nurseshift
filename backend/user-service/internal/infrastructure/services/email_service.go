package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

// EmailService interface for sending emails
type EmailService interface {
	SendVerificationEmail(to, token string) error
}

// SMTPEmailService implements EmailService using SMTP
type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService() *SMTPEmailService {
	return &SMTPEmailService{
		host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		port:     getEnv("SMTP_PORT", "587"),
		username: getEnv("SMTP_USERNAME", ""),
		password: getEnv("SMTP_PASSWORD", ""),
		from:     getEnv("SMTP_FROM", ""),
	}
}

// SendVerificationEmail sends a verification email
func (s *SMTPEmailService) SendVerificationEmail(to, token string) error {
	if s.username == "" || s.password == "" {
		// Fallback to console output if SMTP not configured
		fmt.Printf("⚠️  SMTP not configured - Email will not be sent\n")
		fmt.Printf("📧 Verification email would be sent to: %s\n", to)
		fmt.Printf("🔑 Token: %s\n", token)
		fmt.Printf("🌐 Verification URL: http://localhost:3000/verify-email?token=%s\n", token)
		fmt.Printf("💡 To enable real email sending, configure SMTP settings in config.env\n")
		return nil
	}

	fmt.Printf("📧 Sending verification email to: %s\n", to)
	fmt.Printf("🔧 SMTP Configuration: %s:%s (from: %s)\n", s.host, s.port, s.from)

	subject := "ยืนยันอีเมล - NurseShift"
	body := fmt.Sprintf(`
<html>
<head>
    <meta charset="UTF-8">
    <title>ยืนยันอีเมล</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <div style="text-align: center; margin-bottom: 30px;">
            <h1 style="color: #2563eb;">NurseShift</h1>
            <h2 style="color: #1f2937;">ยืนยันอีเมลของคุณ</h2>
        </div>
        
        <div style="background-color: #f8fafc; padding: 30px; border-radius: 10px; margin-bottom: 20px;">
            <p>สวัสดีครับ/ค่ะ</p>
            <p>คุณได้ขอให้ยืนยันอีเมลสำหรับบัญชี NurseShift ของคุณ</p>
            
            <div style="background-color: #dbeafe; padding: 20px; border-radius: 8px; margin: 20px 0; text-align: center;">
                <h3 style="color: #1e40af; margin: 0;">รหัสยืนยัน</h3>
                <div style="font-size: 24px; font-weight: bold; color: #1e40af; letter-spacing: 2px; margin: 10px 0;">
                    %s
                </div>
                <p style="margin: 0; color: #64748b;">รหัสนี้มีอายุการใช้งาน 24 ชั่วโมง</p>
            </div>
            
            <p>หรือคลิกลิงก์ด้านล่างเพื่อยืนยันอีเมล:</p>
            <div style="text-align: center; margin: 20px 0;">
                <a href="http://localhost:3000/email-verification?token=%s" 
                   style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
                    ยืนยันอีเมล
                </a>
            </div>
        </div>
        
        <div style="text-align: center; color: #64748b; font-size: 14px;">
            <p>หากคุณไม่ได้ขอให้ยืนยันอีเมล กรุณาละเว้นอีเมลนี้</p>
            <p>© 2025 NurseShift. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, token, token)

	err := s.sendEmail(to, subject, body)
	if err != nil {
		fmt.Printf("❌ Failed to send email: %v\n", err)
		return err
	}

	fmt.Printf("✅ Verification email sent successfully to: %s\n", to)
	return nil
}

// sendEmail sends an email using SMTP
func (s *SMTPEmailService) sendEmail(to, subject, body string) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Set up the message
	message := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, s.from, subject, body)

	// Send the email
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// Use TLS for port 587
	if s.port == "587" {
		return s.sendEmailWithTLS(to, addr, auth, message)
	}

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
}

// sendEmailWithTLS sends email with TLS encryption
func (s *SMTPEmailService) sendEmailWithTLS(to, addr string, auth smtp.Auth, message string) error {
	// Connect to SMTP server
	host := s.host
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Start TLS
	if err = conn.StartTLS(&tls.Config{ServerName: host}); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err = conn.Mail(s.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = conn.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send data
	w, err := conn.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return conn.Quit()
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
