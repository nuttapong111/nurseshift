package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PackageHandler handles package-related HTTP requests
type PackageHandler struct{}

// NewPackageHandler creates a new package handler
func NewPackageHandler() *PackageHandler {
	return &PackageHandler{}
}

// GetPackages returns available packages
func (h *PackageHandler) GetPackages(c *fiber.Ctx) error {
	// Mock packages data matching frontend
	packages := []fiber.Map{
		{
			"id":          1,
			"name":        "แพ็คเกจทดลองใช้",
			"type":        "trial",
			"price":       0.00,
			"duration":    90,
			"description": "ทดลองใช้ฟรี 90 วัน",
			"features": []string{
				"จัดการแผนกพื้นฐาน",
				"พนักงานสูงสุด 5 คน",
				"ตารางเวรแบบง่าย",
			},
			"maxUsers":       5,
			"maxDepartments": 2,
			"isPopular":      false,
			"isActive":       true,
		},
		{
			"id":          2,
			"name":        "แพ็คเกจมาตรฐาน",
			"type":        "standard",
			"price":       990.00,
			"duration":    30,
			"description": "เหมาะสำหรับแผนกขนาดกลาง",
			"features": []string{
				"จัดการหลายแผนก",
				"พนักงานไม่จำกัด",
				"ตารางเวรอัตโนมัติ",
				"การแจ้งเตือนแบบเรียลไทม์",
				"รายงานและสถิติ",
			},
			"maxUsers":       25,
			"maxDepartments": 5,
			"isPopular":      true,
			"isActive":       true,
		},
		{
			"id":          3,
			"name":        "แพ็คเกจระดับองค์กร",
			"type":        "enterprise",
			"price":       2990.00,
			"duration":    90,
			"description": "เหมาะสำหรับองค์กรขนาดใหญ่",
			"features": []string{
				"จัดการหลายแผนกไม่จำกัด",
				"พนักงานไม่จำกัด",
				"ตารางเวรอัตโนมัติด้วย AI",
				"การแจ้งเตือนแบบเรียลไทม์",
				"รายงานและสถิติขั้นสูง",
				"การสำรองข้อมูล",
				"การสนับสนุนลูกค้าแบบพิเศษ",
			},
			"maxUsers":       1000,
			"maxDepartments": 50,
			"isPopular":      false,
			"isActive":       true,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลแพ็คเกจสำเร็จ",
		"data":    packages,
	})
}

// GetPackage returns specific package details
func (h *PackageHandler) GetPackage(c *fiber.Ctx) error {
	packageID := c.Params("id")

	// Mock package detail
	pkg := fiber.Map{
		"id":          packageID,
		"name":        "แพ็คเกจมาตรฐาน",
		"type":        "standard",
		"price":       990.00,
		"duration":    30,
		"description": "เหมาะสำหรับแผนกขนาดกลาง",
		"features": []string{
			"จัดการหลายแผนก",
			"พนักงานไม่จำกัด",
			"ตารางเวรอัตโนมัติ",
			"การแจ้งเตือนแบบเรียลไทม์",
			"รายงานและสถิติ",
		},
		"maxUsers":       25,
		"maxDepartments": 5,
		"isPopular":      true,
		"isActive":       true,
		"createdAt":      time.Now().Add(-90 * 24 * time.Hour),
		"updatedAt":      time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลแพ็คเกจสำเร็จ",
		"data":    pkg,
	})
}

// GetCurrentUserPackage returns user's current package
func (h *PackageHandler) GetCurrentUserPackage(c *fiber.Ctx) error {
	// Mock current package data
	currentPackage := fiber.Map{
		"id":         2,
		"name":       "แพ็คเกจมาตรฐาน",
		"type":       "standard",
		"price":      990.00,
		"duration":   30,
		"isActive":   true,
		"startDate":  "2025-08-01",
		"expireDate": "2025-08-31",
		"daysLeft":   22,
		"autoRenew":  false,
		"features": []string{
			"จัดการหลายแผนก",
			"พนักงานไม่จำกัด",
			"ตารางเวรอัตโนมัติ",
			"การแจ้งเตือนแบบเรียลไทม์",
			"รายงานและสถิติ",
		},
		"usage": fiber.Map{
			"users":          18,
			"maxUsers":       25,
			"departments":    3,
			"maxDepartments": 5,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลแพ็คเกจปัจจุบันสำเร็จ",
		"data":    currentPackage,
	})
}

// CreatePackageOrder creates a new package order
func (h *PackageHandler) CreatePackageOrder(c *fiber.Ctx) error {
	var req struct {
		PackageID int `json:"packageId" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Mock order response
	order := fiber.Map{
		"id":         uuid.New().String(),
		"packageId":  req.PackageID,
		"status":     "pending_payment",
		"totalPrice": 990.00,
		"createdAt":  time.Now(),
		"bankInfo": fiber.Map{
			"bankName":      "ธนาคารกสิกรไทย",
			"accountName":   "บริษัท เนิร์สชิฟท์ จำกัด",
			"accountNumber": "123-4-56789-0",
		},
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้างคำสั่งซื้อสำเร็จ",
		"data":    order,
	})
}

// GetPackageStats returns package statistics
func (h *PackageHandler) GetPackageStats(c *fiber.Ctx) error {
	stats := fiber.Map{
		"totalPackages":  3,
		"activePackages": 3,
		"popularPackage": "แพ็คเกจมาตรฐาน",
		"packageUsage": []fiber.Map{
			{
				"packageName": "แพ็คเกจทดลองใช้",
				"userCount":   15,
				"percentage":  25.0,
			},
			{
				"packageName": "แพ็คเกจมาตรฐาน",
				"userCount":   35,
				"percentage":  58.3,
			},
			{
				"packageName": "แพ็คเกจระดับองค์กร",
				"userCount":   10,
				"percentage":  16.7,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถิติแพ็คเกจสำเร็จ",
		"data":    stats,
	})
}

// UpdatePackageSettings updates package settings
func (h *PackageHandler) UpdatePackageSettings(c *fiber.Ctx) error {
	var req fiber.Map
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลไม่ถูกต้อง",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตการตั้งค่าแพ็คเกจสำเร็จ",
		"data":    req,
	})
}

// Health returns service health status
func (h *PackageHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "package-service",
		"timestamp": time.Now(),
	})
}


