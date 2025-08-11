package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct{}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

// Package represents a package data structure
type Package struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Price        int       `json:"price"`
	Duration     int       `json:"duration"`
	Description  string    `json:"description"`
	Features     []string  `json:"features"`
	IsPopular    bool      `json:"isPopular"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Payment represents a payment data structure
type Payment struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	PackageID    string    `json:"packageId"`
	PackageName  string    `json:"packageName"`
	Amount       int       `json:"amount"`
	Status       string    `json:"status"` // pending, approved, rejected
	Evidence     *string   `json:"evidence"`
	PaymentDate  string    `json:"paymentDate"`
	ApprovedDate *string   `json:"approvedDate"`
	ExtendedDays *int      `json:"extendedDays"`
	RejectReason *string   `json:"rejectReason"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Mock data
var mockPackages = []Package{
	{
		ID:          "1",
		Name:        "แพ็คเกจมาตรฐาน",
		Price:       990,
		Duration:    30,
		Description: "เหมาะสำหรับแผนกขนาดกลาง",
		Features: []string{
			"จัดการหลายแผนก",
			"พนักงานไม่จำกัด",
			"ตารางเวรอัตโนมัติ",
			"การแจ้งเตือนแบบเรียลไทม์",
			"รายงานและสถิติ",
		},
		IsPopular: true,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:          "2",
		Name:        "แพ็คเกจระดับองค์กร",
		Price:       2990,
		Duration:    90,
		Description: "เหมาะสำหรับองค์กรขนาดใหญ่",
		Features: []string{
			"จัดการหลายแผนกไม่จำกัด",
			"พนักงานไม่จำกัด",
			"ตารางเวรอัตโนมัติด้วย AI",
			"การแจ้งเตือนแบบเรียลไทม์",
			"รายงานและสถิติขั้นสูง",
			"การสำรองข้อมูล",
			"การสนับสนุนลูกค้าแบบพิเศษ",
		},
		IsPopular: false,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var mockPayments = []Payment{
	{
		ID:           "1",
		UserID:       "user-1",
		PackageID:    "1",
		PackageName:  "แพ็คเกจมาตรฐาน",
		Amount:       990,
		Status:       "approved",
		Evidence:     stringPtr("payment_evidence_1.jpg"),
		PaymentDate:  "2024-03-01",
		ApprovedDate: stringPtr("2024-03-02"),
		ExtendedDays: intPtr(30),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
	{
		ID:          "2",
		UserID:      "user-1",
		PackageID:   "2",
		PackageName: "แพ็คเกจระดับองค์กร",
		Amount:      2990,
		Status:      "pending",
		Evidence:    stringPtr("payment_evidence_2.jpg"),
		PaymentDate: "2024-03-15",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:           "3",
		UserID:       "user-1",
		PackageID:    "1",
		PackageName:  "แพ็คเกจมาตรฐาน",
		Amount:       990,
		Status:       "rejected",
		Evidence:     stringPtr("payment_evidence_3.jpg"),
		PaymentDate:  "2024-02-15",
		RejectReason: stringPtr("หลักฐานการโอนเงินไม่ชัดเจน"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
}

// GetPayments returns all payments with optional filters
func (h *PaymentHandler) GetPayments(c *fiber.Ctx) error {
	userID := c.Query("userId")
	status := c.Query("status")

	filteredPayments := mockPayments

	// Filter by user
	if userID != "" {
		var filtered []Payment
		for _, payment := range filteredPayments {
			if payment.UserID == userID {
				filtered = append(filtered, payment)
			}
		}
		filteredPayments = filtered
	}

	// Filter by status
	if status != "" {
		var filtered []Payment
		for _, payment := range filteredPayments {
			if payment.Status == status {
				filtered = append(filtered, payment)
			}
		}
		filteredPayments = filtered
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลการชำระเงินสำเร็จ",
		"data":    filteredPayments,
	})
}

// CreatePayment creates a new payment record
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var req struct {
		UserID      string  `json:"userId" validate:"required"`
		PackageID   string  `json:"packageId" validate:"required"`
		PackageName string  `json:"packageName" validate:"required"`
		Amount      int     `json:"amount" validate:"required"`
		Evidence    *string `json:"evidence"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Create new payment
	newPayment := Payment{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		PackageID:   req.PackageID,
		PackageName: req.PackageName,
		Amount:      req.Amount,
		Status:      "pending",
		Evidence:    req.Evidence,
		PaymentDate: time.Now().Format("2006-01-02"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add to mock data
	mockPayments = append(mockPayments, newPayment)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้างรายการชำระเงินสำเร็จ",
		"data":    newPayment,
	})
}

// GetPayment returns specific payment details
func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")

	for _, payment := range mockPayments {
		if payment.ID == paymentID {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "ดึงข้อมูลการชำระเงินสำเร็จ",
				"data":    payment,
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลการชำระเงิน",
	})
}

// UpdatePayment updates payment information
func (h *PaymentHandler) UpdatePayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")

	var req struct {
		Evidence *string `json:"evidence"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Find and update payment
	for i, payment := range mockPayments {
		if payment.ID == paymentID {
			if req.Evidence != nil {
				mockPayments[i].Evidence = req.Evidence
			}
			mockPayments[i].Status = "pending" // Reset to pending when updated
			mockPayments[i].UpdatedAt = time.Now()
			mockPayments[i].PaymentDate = time.Now().Format("2006-01-02")
			mockPayments[i].RejectReason = nil // Clear previous reject reason

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "อัปเดตการชำระเงินสำเร็จ",
				"data":    mockPayments[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลการชำระเงิน",
	})
}

// ApprovePayment approves a payment
func (h *PaymentHandler) ApprovePayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")

	var req struct {
		ExtendedDays int `json:"extendedDays" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Find and approve payment
	for i, payment := range mockPayments {
		if payment.ID == paymentID {
			mockPayments[i].Status = "approved"
			mockPayments[i].ApprovedDate = stringPtr(time.Now().Format("2006-01-02"))
			mockPayments[i].ExtendedDays = &req.ExtendedDays
			mockPayments[i].UpdatedAt = time.Now()

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "อนุมัติการชำระเงินสำเร็จ",
				"data":    mockPayments[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลการชำระเงิน",
	})
}

// RejectPayment rejects a payment
func (h *PaymentHandler) RejectPayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Find and reject payment
	for i, payment := range mockPayments {
		if payment.ID == paymentID {
			mockPayments[i].Status = "rejected"
			mockPayments[i].RejectReason = &req.Reason
			mockPayments[i].UpdatedAt = time.Now()

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "ปฏิเสธการชำระเงินสำเร็จ",
				"data":    mockPayments[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลการชำระเงิน",
	})
}

// GetPaymentStats returns payment statistics
func (h *PaymentHandler) GetPaymentStats(c *fiber.Ctx) error {
	totalPayments := len(mockPayments)
	approved := 0
	pending := 0
	rejected := 0
	totalRevenue := 0

	for _, payment := range mockPayments {
		switch payment.Status {
		case "approved":
			approved++
			totalRevenue += payment.Amount
		case "pending":
			pending++
		case "rejected":
			rejected++
		}
	}

	stats := fiber.Map{
		"totalPayments": totalPayments,
		"approved":      approved,
		"pending":       pending,
		"rejected":      rejected,
		"totalRevenue":  totalRevenue,
		"packages":      mockPackages,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถิติการชำระเงินสำเร็จ",
		"data":    stats,
	})
}

// Health returns service health status
func (h *PaymentHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "payment-service",
		"timestamp": time.Now(),
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}