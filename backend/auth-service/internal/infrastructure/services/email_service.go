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

	// Create reset link
	resetLink := fmt.Sprintf("http://localhost:3000/auth/reset-password?token=%s", resetToken)

	// Create HTML body for password reset
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>รีเซ็ตรหัสผ่าน</title>
    <style>
        body {
            font-family: 'Sarabun', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px;
            text-align: center;
            border-radius: 10px 10px 0 0;
        }
        .content {
            background-color: #f8f9fa;
            padding: 30px;
            border-radius: 0 0 10px 10px;
        }
        .button {
            background-color: #007bff;
            color: white;
            padding: 15px 30px;
            text-decoration: none;
            border-radius: 5px;
            display: inline-block;
            margin: 20px 0;
            font-weight: bold;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #666;
            font-size: 14px;
        }
        .token {
            background-color: #e3f2fd;
            padding: 15px;
            border-radius: 5px;
            font-family: monospace;
            font-size: 18px;
            text-align: center;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>รีเซ็ตรหัสผ่าน</h1>
            <p>NurseShift Management System</p>
        </div>
        
        <div class="content">
            <h2>สวัสดีครับ/ค่ะ</h2>
            <p>เราได้รับคำขอให้รีเซ็ตรหัสผ่านสำหรับบัญชีของคุณในระบบจัดตารางเวร</p>
            
            <p>หากคุณเป็นผู้ทำการขอรีเซ็ตรหัสผ่าน กรุณาคลิกลิงก์ด้านล่างเพื่อสร้างรหัสผ่านใหม่:</p>
            
            <div style="text-align: center;">
                <a href="%s" class="button">รีเซ็ตรหัสผ่าน</a>
            </div>
            
            <p><strong>ข้อมูลสำคัญ:</strong></p>
            <ul>
                <li>ลิงก์นี้จะหมดอายุภายใน 15 นาที</li>
                <li>ลิงก์สามารถใช้ได้เพียง 1 ครั้งเท่านั้น</li>
                <li>หากไม่ได้เป็นผู้ขอรีเซ็ต กรุณาเพิกเฉยต่ออีเมลนี้</li>
            </ul>
            
            <p><strong>ต้องการความช่วยเหลือ?</strong></p>
            <p>หากมีปัญหาหรือข้อสงสัย กรุณาติดต่อทีมสนับสนุน:</p>
            <ul>
                <li>อีเมล: support@nurseshift.com</li>
                <li>โทรศัพท์: 02-xxx-xxxx</li>
                <li>เวลาทำการ: จันทร์-ศุกร์ 8:00-17:00 น.</li>
            </ul>
            
            <p>ขอบคุณที่ใช้บริการ NurseShift</p>
            <p><em>อีเมลนี้ถูกส่งโดยอัตโนมัติจากระบบจัดตารางเวร กรุณาอย่าตอบกลับอีเมลนี้</em></p>
        </div>
        
        <div class="footer">
            <p>© 2024 NurseShift. สงวนลิขสิทธิ์.</p>
        </div>
    </div>
</body>
</html>`, resetLink)

	// Send email with HTML content type
	return s.SendHTMLEmail(to, subject, body)
}

// SendHTMLEmail sends an HTML email
func (s *GmailEmailService) SendHTMLEmail(to, subject, body string) error {
	// Create message with HTML content type
	message := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
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
		return fmt.Errorf("failed to send HTML email: %w", err)
	}

	return nil
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
	resetLink := fmt.Sprintf("http://localhost:3000/auth/reset-password?token=%s", resetToken)
	fmt.Printf("📧 Mock Password Reset Email Sent:\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Reset Link: %s\n", resetLink)
	fmt.Printf("Email Template: HTML Email with clickable button\n")
	return nil
}
