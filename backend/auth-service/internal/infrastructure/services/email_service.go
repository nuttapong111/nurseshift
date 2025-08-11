package services

import (
	"fmt"
	"net/smtp"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendPasswordResetEmail(to, resetToken string) error
}

// GmailEmailService implements EmailService using Gmail SMTP
type GmailEmailService struct {
	fromEmail    string
	fromPassword string
	smtpHost     string
	smtpPort     string
}

// NewGmailEmailService creates a new Gmail email service
func NewGmailEmailService(fromEmail, fromPassword string) EmailService {
	return &GmailEmailService{
		fromEmail:    fromEmail,
		fromPassword: fromPassword,
		smtpHost:     "smtp.gmail.com",
		smtpPort:     "587",
	}
}

// SendEmail sends a generic email
func (s *GmailEmailService) SendEmail(to, subject, body string) error {
	// Create message
	message := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, s.fromEmail, subject, body)

	// Authentication
	auth := smtp.PlainAuth("", s.fromEmail, s.fromPassword, s.smtpHost)

	// Send email
	err := smtp.SendMail(
		s.smtpHost+":"+s.smtpPort,
		auth,
		s.fromEmail,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email
func (s *GmailEmailService) SendPasswordResetEmail(to, resetToken string) error {
	subject := "รีเซ็ตรหัสผ่าน - NurseShift"
	
	// Create HTML body for password reset
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>รีเซ็ตรหัสผ่าน</title>
    <style>
        body { font-family: 'Sarabun', Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .button { display: inline-block; background: #667eea; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
        .token { background: #e9ecef; padding: 15px; border-radius: 5px; font-family: monospace; font-size: 18px; text-align: center; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 รีเซ็ตรหัสผ่าน</h1>
            <p>NurseShift Management System</p>
        </div>
        <div class="content">
            <h2>สวัสดีครับ/ค่ะ</h2>
            <p>เราได้รับคำขอให้รีเซ็ตรหัสผ่านสำหรับบัญชีของคุณ</p>
            <p>หากคุณไม่ได้ขอให้รีเซ็ตรหัสผ่าน กรุณาละเว้นอีเมลนี้</p>
            
            <h3>รหัสยืนยัน:</h3>
            <div class="token">%s</div>
            
            <p><strong>คำแนะนำ:</strong></p>
            <ul>
                <li>รหัสยืนยันนี้จะหมดอายุใน 15 นาที</li>
                <li>กรุณาใส่รหัสยืนยันในหน้าเว็บเพื่อรีเซ็ตรหัสผ่าน</li>
                <li>ห้ามแชร์รหัสยืนยันนี้กับผู้อื่น</li>
            </ul>
            
            <p>หากคุณมีคำถามหรือต้องการความช่วยเหลือ กรุณาติดต่อทีมสนับสนุน</p>
            
            <p>ขอบคุณที่ใช้บริการ NurseShift</p>
        </div>
        <div class="footer">
            <p>© 2024 NurseShift. All rights reserved.</p>
            <p>อีเมลนี้ถูกส่งโดยระบบอัตโนมัติ กรุณาอย่าตอบกลับ</p>
        </div>
    </div>
</body>
</html>`, resetToken)

	// Send email
	return s.SendEmail(to, subject, body)
}

// MockEmailService implements EmailService for testing
type MockEmailService struct{}

// NewMockEmailService creates a new mock email service
func NewMockEmailService() EmailService {
	return &MockEmailService{}
}

// SendEmail mocks sending an email
func (s *MockEmailService) SendEmail(to, subject, body string) error {
	fmt.Printf("📧 Mock Email Sent:\nTo: %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}

// SendPasswordResetEmail mocks sending a password reset email
func (s *MockEmailService) SendPasswordResetEmail(to, resetToken string) error {
	fmt.Printf("📧 Mock Password Reset Email Sent:\nTo: %s\nReset Token: %s\n", to, resetToken)
	return nil
}
