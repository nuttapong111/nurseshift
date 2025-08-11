package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct{}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// Notification represents notification data structure
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"` // schedule, leave, system, payment, reminder, holiday
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Priority  string    `json:"priority"` // high, medium, low
	IsRead    bool      `json:"isRead"`
	ActionURL *string   `json:"actionUrl"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Mock data
var mockNotifications = []Notification{
	{
		ID:        "1",
		UserID:    "user-1",
		Type:      "schedule",
		Title:     "ตารางเวรใหม่ประจำเดือนมีนาคม",
		Message:   "ตารางเวรประจำเดือนมีนาคม 2024 ได้รับการอนุมัติและเผยแพร่แล้ว กรุณาตรวจสอบตารางเวรของคุณ",
		Priority:  "high",
		IsRead:    false,
		ActionURL: stringPtr("/dashboard/schedule"),
		Timestamp: time.Now().Add(-time.Hour * 2),
		CreatedAt: time.Now().Add(-time.Hour * 2),
		UpdatedAt: time.Now().Add(-time.Hour * 2),
	},
	{
		ID:        "2",
		UserID:    "user-1",
		Type:      "leave",
		Title:     "คำขอลาป่วยได้รับการอนุมัติ",
		Message:   "คำขอลาป่วยวันที่ 15 มีนาคม 2024 ได้รับการอนุมัติจากหัวหน้าแผนกแล้ว",
		Priority:  "medium",
		IsRead:    true,
		ActionURL: stringPtr("/dashboard/employee-leaves"),
		Timestamp: time.Now().Add(-time.Hour * 6),
		CreatedAt: time.Now().Add(-time.Hour * 6),
		UpdatedAt: time.Now().Add(-time.Hour * 6),
	},
	{
		ID:        "3",
		UserID:    "user-1",
		Type:      "system",
		Title:     "การอัปเดตระบบ",
		Message:   "ระบบจะมีการปรับปรุงในวันที่ 20 มีนาคม 2024 เวลา 02:00-04:00 น. อาจมีการหยุดให้บริการชั่วคราว",
		Priority:  "low",
		IsRead:    false,
		ActionURL: nil,
		Timestamp: time.Now().Add(-time.Hour * 12),
		CreatedAt: time.Now().Add(-time.Hour * 12),
		UpdatedAt: time.Now().Add(-time.Hour * 12),
	},
	{
		ID:        "4",
		UserID:    "user-1",
		Type:      "payment",
		Title:     "การชำระเงินแพ็คเกจ",
		Message:   "การชำระเงินแพ็คเกจมาตรฐานได้รับการอนุมัติแล้ว บัญชีของคุณได้รับการต่ออายุเป็น 30 วัน",
		Priority:  "high",
		IsRead:    true,
		ActionURL: stringPtr("/dashboard/packages"),
		Timestamp: time.Now().Add(-time.Hour * 24),
		CreatedAt: time.Now().Add(-time.Hour * 24),
		UpdatedAt: time.Now().Add(-time.Hour * 24),
	},
	{
		ID:        "5",
		UserID:    "user-1",
		Type:      "reminder",
		Title:     "แจ้งเตือนเปลี่ยนเวร",
		Message:   "คุณมีเวรเช้าในวันพรุ่งนี้ (16 มีนาคม 2024) เวลา 07:00-15:00 น. ที่แผนกผู้ป่วยใน",
		Priority:  "medium",
		IsRead:    false,
		ActionURL: stringPtr("/dashboard/schedule"),
		Timestamp: time.Now().Add(-time.Hour * 4),
		CreatedAt: time.Now().Add(-time.Hour * 4),
		UpdatedAt: time.Now().Add(-time.Hour * 4),
	},
	{
		ID:        "6",
		UserID:    "user-1",
		Type:      "holiday",
		Title:     "วันหยุดประจำปี",
		Message:   "วันสงกรานต์ (13-15 เมษายน 2024) ได้รับการอนุมัติเป็นวันหยุดประจำปีแล้ว",
		Priority:  "low",
		IsRead:    true,
		ActionURL: stringPtr("/dashboard/department-settings"),
		Timestamp: time.Now().Add(-time.Hour * 48),
		CreatedAt: time.Now().Add(-time.Hour * 48),
		UpdatedAt: time.Now().Add(-time.Hour * 48),
	},
}

// GetNotifications returns notifications with optional filters
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Query("userId")
	notificationType := c.Query("type")
	isRead := c.Query("isRead")
	priority := c.Query("priority")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filteredNotifications := mockNotifications

	// Filter by user
	if userID != "" {
		var filtered []Notification
		for _, notification := range filteredNotifications {
			if notification.UserID == userID {
				filtered = append(filtered, notification)
			}
		}
		filteredNotifications = filtered
	}

	// Filter by type
	if notificationType != "" {
		var filtered []Notification
		for _, notification := range filteredNotifications {
			if notification.Type == notificationType {
				filtered = append(filtered, notification)
			}
		}
		filteredNotifications = filtered
	}

	// Filter by read status
	if isRead != "" {
		readStatus := isRead == "true"
		var filtered []Notification
		for _, notification := range filteredNotifications {
			if notification.IsRead == readStatus {
				filtered = append(filtered, notification)
			}
		}
		filteredNotifications = filtered
	}

	// Filter by priority
	if priority != "" {
		var filtered []Notification
		for _, notification := range filteredNotifications {
			if notification.Priority == priority {
				filtered = append(filtered, notification)
			}
		}
		filteredNotifications = filtered
	}

	// Pagination
	start := (page - 1) * limit
	end := start + limit
	if start > len(filteredNotifications) {
		start = len(filteredNotifications)
	}
	if end > len(filteredNotifications) {
		end = len(filteredNotifications)
	}

	paginatedNotifications := filteredNotifications[start:end]
	totalPages := (len(filteredNotifications) + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลการแจ้งเตือนสำเร็จ",
		"data": fiber.Map{
			"notifications": paginatedNotifications,
			"total":         len(filteredNotifications),
			"page":          page,
			"limit":         limit,
			"totalPages":    totalPages,
		},
	})
}

// CreateNotification creates a new notification
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	var req struct {
		UserID    string  `json:"userId" validate:"required"`
		Type      string  `json:"type" validate:"required"`
		Title     string  `json:"title" validate:"required"`
		Message   string  `json:"message" validate:"required"`
		Priority  string  `json:"priority"`
		ActionURL *string `json:"actionUrl"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Set default priority if not provided
	if req.Priority == "" {
		req.Priority = "medium"
	}

	// Create new notification
	newNotification := Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Message:   req.Message,
		Priority:  req.Priority,
		IsRead:    false,
		ActionURL: req.ActionURL,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Add to mock data
	mockNotifications = append(mockNotifications, newNotification)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้างการแจ้งเตือนสำเร็จ",
		"data":    newNotification,
	})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")

	// Find and mark notification as read
	for i, notification := range mockNotifications {
		if notification.ID == notificationID {
			mockNotifications[i].IsRead = true
			mockNotifications[i].UpdatedAt = time.Now()

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "ทำเครื่องหมายอ่านแล้วสำเร็จ",
				"data":    mockNotifications[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบการแจ้งเตือน",
	})
}

// MarkAllAsRead marks all notifications as read for a user
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"userId" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	updatedCount := 0

	// Mark all user's notifications as read
	for i, notification := range mockNotifications {
		if notification.UserID == req.UserID && !notification.IsRead {
			mockNotifications[i].IsRead = true
			mockNotifications[i].UpdatedAt = time.Now()
			updatedCount++
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ทำเครื่องหมายอ่านแล้วทั้งหมดสำเร็จ",
		"data": fiber.Map{
			"updatedCount": updatedCount,
		},
	})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	notificationID := c.Params("id")

	// Find and remove notification
	for i, notification := range mockNotifications {
		if notification.ID == notificationID {
			mockNotifications = append(mockNotifications[:i], mockNotifications[i+1:]...)

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "ลบการแจ้งเตือนสำเร็จ",
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบการแจ้งเตือน",
	})
}

// GetNotificationStats returns notification statistics
func (h *NotificationHandler) GetNotificationStats(c *fiber.Ctx) error {
	userID := c.Query("userId")

	filteredNotifications := mockNotifications
	if userID != "" {
		var filtered []Notification
		for _, notification := range filteredNotifications {
			if notification.UserID == userID {
				filtered = append(filtered, notification)
			}
		}
		filteredNotifications = filtered
	}

	totalNotifications := len(filteredNotifications)
	unreadCount := 0
	readCount := 0
	highPriorityCount := 0

	typeStats := make(map[string]int)

	for _, notification := range filteredNotifications {
		if notification.IsRead {
			readCount++
		} else {
			unreadCount++
		}

		if notification.Priority == "high" {
			highPriorityCount++
		}

		typeStats[notification.Type]++
	}

	stats := fiber.Map{
		"total":           totalNotifications,
		"unread":          unreadCount,
		"read":            readCount,
		"highPriority":    highPriorityCount,
		"typeBreakdown":   typeStats,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถิติการแจ้งเตือนสำเร็จ",
		"data":    stats,
	})
}

// Health returns service health status
func (h *NotificationHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "notification-service",
		"timestamp": time.Now(),
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}