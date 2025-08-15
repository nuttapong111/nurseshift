package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"nurseshift/schedule-service/internal/infrastructure/database"
	"nurseshift/schedule-service/internal/optimizer"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ScheduleHandler handles schedule-related HTTP requests
type ScheduleHandler struct{ repo *database.ScheduleRepository }

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler(repo *database.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{repo: repo}
}

// GetSchedules returns schedules for authenticated user's departments
func (h *ScheduleHandler) GetSchedules(c *fiber.Ctx) error {
	_ = c.Locals("userID").(string)
	departmentId := c.Query("departmentId")
	month := c.Query("month")

	// 1) พยายามแบบ staff-based ก่อน
	itemsStaff, err := h.repo.ListWithStaff(c.Context(), departmentId, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	out := make([]fiber.Map, 0, len(itemsStaff))
	for _, r := range itemsStaff {
		var notesPtr *string
		if r.Notes.Valid {
			v := r.Notes.String
			notesPtr = &v
		}
		role := "nurse"
		if strings.Contains(strings.ToLower(r.StaffRole), "assist") || strings.Contains(r.StaffRole, "ผู้ช่วย") {
			role = "assistant"
		}
		out = append(out, fiber.Map{
			"id":             r.ID,
			"departmentId":   r.DepartmentID,
			"staffId":        r.StaffID,
			"shiftId":        r.ShiftID,
			"scheduleDate":   r.ScheduleDate,
			"status":         r.Status,
			"notes":          notesPtr,
			"departmentRole": role,
			"userName":       r.StaffName,
		})
	}
	// 2) ถ้ายังว่าง ให้ fallback ไปแบบ user-based (รองรับข้อมูลเก่า)
	if len(out) == 0 {
		itemsUser, err2 := h.repo.ListWithRole(c.Context(), departmentId, month)
		if err2 != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err2.Error()})
		}
		for _, r := range itemsUser {
			var notesPtr *string
			if r.Notes.Valid {
				v := r.Notes.String
				notesPtr = &v
			}
			role := "nurse"
			if strings.Contains(strings.ToLower(r.DepartmentRole), "assist") || strings.Contains(r.DepartmentRole, "ผู้ช่วย") {
				role = "assistant"
			}
			out = append(out, fiber.Map{
				"id":             r.ID,
				"departmentId":   r.DepartmentID,
				"userId":         r.UserID,
				"shiftId":        r.ShiftID,
				"scheduleDate":   r.ScheduleDate,
				"status":         r.Status,
				"notes":          notesPtr,
				"departmentRole": role,
				"userName":       strings.TrimSpace(r.UserFirstName + " " + r.UserLastName),
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "ดึงข้อมูลตารางเวรสำเร็จ", "data": out})
}

// CalendarMeta returns working/holiday flags for given month
func (h *ScheduleHandler) CalendarMeta(c *fiber.Ctx) error {
	// ไม่ต้องบังคับมี userID เพราะเส้นนี้ไม่ได้ติด middleware เสมอไป
	if v := c.Locals("userID"); v != nil {
		if _, ok := v.(string); !ok {
			// ignore silently; this endpoint does not require auth
		}
	}
	departmentId := c.Query("departmentId")
	month := c.Query("month")
	if departmentId == "" || month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ต้องระบุ departmentId และ month"})
	}

	working, err := h.repo.ListWorkingDays(c.Context(), departmentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	holidays, err := h.repo.ListHolidaysForMonth(c.Context(), departmentId, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	t, err := time.Parse("2006-01", month)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "รูปแบบเดือนไม่ถูกต้อง"})
	}
	year, m, _ := t.Date()
	first := time.Date(year, m, 1, 0, 0, 0, 0, time.UTC)
	next := first.AddDate(0, 1, 0)
	days := int(next.Sub(first).Hours() / 24)
	isHoliday := func(d time.Time) bool {
		ds := d.Format("2006-01-02")
		for _, h := range holidays {
			if ds >= h.Start && ds <= h.End {
				return true
			}
		}
		return false
	}
	data := []fiber.Map{}
	for day := 1; day <= days; day++ {
		d := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
		wd := int(d.Weekday())
		w := true
		if v, ok := working[wd]; ok {
			w = v
		}
		data = append(data, fiber.Map{
			"date":      d.Format("2006-01-02"),
			"isWorking": w,
			"isHoliday": isHoliday(d),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "ปฏิทินการทำงาน", "data": data})
}

// CreateSchedule creates a new schedule
func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		DepartmentID string `json:"departmentId" validate:"required"`
		Date         string `json:"date" validate:"required"`
		Shifts       []struct {
			Name               string   `json:"name"`
			StartTime          string   `json:"startTime"`
			EndTime            string   `json:"endTime"`
			RequiredNurses     int      `json:"requiredNurses"`
			RequiredAssistants int      `json:"requiredAssistants"`
			Nurses             []string `json:"nurses"`
			Assistants         []string `json:"assistants"`
		} `json:"shifts"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	id := uuid.New().String()
	rec := &database.ScheduleRecord{ID: id, DepartmentID: req.DepartmentID, UserID: userID, ShiftID: "", ScheduleDate: req.Date, Status: "assigned"}
	if err := h.repo.Create(c.Context(), rec); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "สร้างตารางเวรสำเร็จ", "data": rec})
}

// GetScheduleStats returns schedule statistics for user's departments
func (h *ScheduleHandler) GetScheduleStats(c *fiber.Ctx) error {
	_ = c.Locals("userID").(string)
	departmentId := c.Query("departmentId")
	month := c.Query("month")

	withRole, err := h.repo.ListWithRole(c.Context(), departmentId, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// Aggregate basic stats
	total := len(withRole)
	nurses := 0
	assistants := 0
	for _, r := range withRole {
		if r.DepartmentRole == "nurse" {
			nurses++
		} else if r.DepartmentRole == "assistant" {
			assistants++
		}
	}
	stats := fiber.Map{
		"totalSchedules":     total,
		"thisMonthSchedules": total,
		"totalNurses":        nurses,
		"totalAssistants":    assistants,
		"totalShifts":        total, // simplistic count
		"totalDepartments":   1,
		"departmentStats":    []fiber.Map{},
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "ดึงสถิติตารางเวรสำเร็จ", "data": stats})
}

// GetSchedule returns specific schedule details
func (h *ScheduleHandler) GetSchedule(c *fiber.Ctx) error {
	_ = c.Locals("userID").(string)
	_ = uuid.New().String()
	// For now, return 404 to keep minimal; can be expanded later
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "ไม่พบข้อมูลตารางเวรที่ระบุ"})
}

// UpdateSchedule updates schedule information
func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	scheduleID := c.Params("id")

	var req fiber.Map
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	_ = userID
	var statusPtr *string
	var notesPtr *string
	var shiftPtr *string
	if v, ok := req["status"].(string); ok {
		statusPtr = &v
	}
	if v, ok := req["notes"].(string); ok {
		notesPtr = &v
	}
	if v, ok := req["shiftId"].(string); ok {
		shiftPtr = &v
	}
	if err := h.repo.Update(c.Context(), scheduleID, statusPtr, notesPtr, shiftPtr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "อัปเดตตารางเวรสำเร็จ", "data": fiber.Map{"id": scheduleID, "updatedAt": time.Now()}})
}

// DeleteSchedule deletes a schedule
func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	scheduleID := c.Params("id")

	_ = userID
	if err := h.repo.Delete(c.Context(), scheduleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "ลบตารางเวรสำเร็จ"})
}

// Health returns service health status
func (h *ScheduleHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "schedule-service",
		"timestamp": time.Now(),
	})
}

// ListShifts returns shift definitions for a department
func (h *ScheduleHandler) ListShifts(c *fiber.Ctx) error {
	_ = c.Locals("userID").(string)
	departmentId := c.Query("departmentId")
	items, err := h.repo.ListShifts(c.Context(), departmentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	out := make([]fiber.Map, 0, len(items))
	for _, s := range items {
		out = append(out, fiber.Map{
			"id":            s.ID,
			"departmentId":  s.DepartmentID,
			"name":          s.Name,
			"type":          s.Type,
			"startTime":     s.StartTime,
			"endTime":       s.EndTime,
			"requiredNurse": s.RequiredNurse,
			"requiredAsst":  s.RequiredAsst,
			"color":         s.Color,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "ดึงข้อมูลกะสำเร็จ", "data": out})
}

// AutoGenerate creates schedules using simple backend logic
func (h *ScheduleHandler) AutoGenerate(c *fiber.Ctx) error {
	var req struct {
		DepartmentID string `json:"departmentId"`
		Month        string `json:"month"` // YYYY-MM
	}
	if err := c.BodyParser(&req); err != nil || req.DepartmentID == "" || req.Month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง ต้องระบุ departmentId และ month"})
	}

	shifts, err := h.repo.ListShifts(c.Context(), req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	// switch to department_staff (candidates)
	if err := h.repo.EnsureStaffSchedulingSchema(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	staffList, err := h.repo.ListDepartmentStaff(c.Context(), req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	var nurses []string
	var assistants []string
	for _, s := range staffList {
		role := strings.ToLower(s.Position)
		if role == "assistant" || strings.Contains(s.Position, "ผู้ช่วย") {
			assistants = append(assistants, s.ID)
		} else {
			nurses = append(nurses, s.ID)
		}
	}
	if len(nurses) == 0 && len(assistants) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ไม่มีพนักงานในแผนกนี้"})
	}

	log.Printf("=== AUTO GENERATE START: dept=%s, month=%s ===", req.DepartmentID, req.Month)

	// calculate days in month
	t, err := time.Parse("2006-01", req.Month)
	if err != nil {
		log.Printf("=== ERROR: Parse month failed: %v ===", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "รูปแบบเดือนไม่ถูกต้อง"})
	}
	year, month, _ := t.Date()
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	nextMonth := first.AddDate(0, 1, 0)
	days := int(nextMonth.Sub(first).Hours() / 24)

	log.Printf("=== PARSED: year=%d, month=%d, days=%d ===", year, month, days)

	// If there are existing schedules for this month and department, clear them before re-generate
	if err := h.repo.DeleteByDepartmentAndMonth(c.Context(), req.DepartmentID, req.Month); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// Use Enhanced Dynamic Priority Algorithm instead of old algorithm
	return h.runEnhancedAlgorithm(c, req, nurses, assistants, shifts, staffList, year, month, days)
}

// runEnhancedAlgorithm implements the Enhanced Dynamic Priority Algorithm
func (h *ScheduleHandler) runEnhancedAlgorithm(c *fiber.Ctx, req struct {
	DepartmentID string `json:"departmentId"`
	Month        string `json:"month"`
}, nurses, assistants []string, shifts []database.ShiftRecord, staffList []database.DepartmentStaff, year int, month time.Month, days int) error {

	log.Printf("=== START ENHANCED ALGORITHM: %s-%s ===", req.DepartmentID, req.Month)

	// Enhanced Dynamic Priority Algorithm with Progressive Relaxation
	var items []database.Assignment
	assignmentCount := map[string]int{}
	lastAssignedDate := map[string]time.Time{}

	// Track assigned time intervals per staff per date to allow multiple non-overlapping shifts
	// intervals are in minutes from 00:00 of the day; end may exceed 1440 for overnight shifts
	assignedIntervals := map[string]map[string][][2]int{}

	// helpers for time/intervals
	parseHM := func(hm string) (int, bool) {
		// hm format HH:MM
		if len(hm) < 4 {
			return 0, false
		}
		var hh, mm int
		if _, err := fmt.Sscanf(hm, "%d:%d", &hh, &mm); err != nil {
			return 0, false
		}
		return hh*60 + mm, true
	}
	shiftInterval := func(sh database.ShiftRecord) (int, int, bool) {
		start, ok1 := parseHM(sh.StartTime)
		end, ok2 := parseHM(sh.EndTime)
		if !ok1 || !ok2 {
			return 0, 0, false
		}
		if end <= start { // overnight or 00:00 next day
			end += 24 * 60
		}
		return start, end, true
	}
	overlaps := func(aStart, aEnd, bStart, bEnd int) bool {
		return aStart < bEnd && bStart < aEnd
	}
	mergeAndMaxContiguous := func(ivals [][2]int) int {
		if len(ivals) == 0 {
			return 0
		}
		sort.Slice(ivals, func(i, j int) bool { return ivals[i][0] < ivals[j][0] })
		curS, curE := ivals[0][0], ivals[0][1]
		maxDur := curE - curS
		for i := 1; i < len(ivals); i++ {
			s, e := ivals[i][0], ivals[i][1]
			if s <= curE { // overlap or touch → contiguous
				if e > curE {
					curE = e
				}
			} else { // gap → new block
				if curE-curS > maxDur {
					maxDur = curE - curS
				}
				curS, curE = s, e
			}
		}
		if curE-curS > maxDur {
			maxDur = curE - curS
		}
		return maxDur
	}
	// read max contiguous-hours policy (default 16h)
	maxContiguousHours := 16
	if v, err := h.repo.GetPriorityValue(c.Context(), req.DepartmentID, "ชั่วโมงติดต่อกันสูงสุด"); err == nil && v.Valid {
		if v.Int64 > 0 && v.Int64 <= 24 {
			maxContiguousHours = int(v.Int64)
		}
	}
	maxContiguousMinutes := maxContiguousHours * 60

	log.Printf("=== SETUP COMPLETE: maxContiguousHours=%d ===", maxContiguousHours)

	canAssignShift := func(staffID string, d time.Time, sh database.ShiftRecord) bool {
		start, end, ok := shiftInterval(sh)
		if !ok {
			return false
		}
		dateStr := d.Format("2006-01-02")
		if assignedIntervals[staffID] == nil {
			assignedIntervals[staffID] = map[string][][2]int{}
		}
		existing := assignedIntervals[staffID][dateStr]
		for _, iv := range existing {
			if overlaps(iv[0], iv[1], start, end) {
				return false // overlap not allowed
			}
		}
		// check contiguous hours after adding this shift
		merged := append(append([][2]int{}, existing...), [2]int{start, end})
		if mergeAndMaxContiguous(merged) > maxContiguousMinutes {
			return false
		}
		return true
	}

	log.Printf("=== START BACKEND ALGORITHM: %s-%s ===", req.DepartmentID, req.Month)

	// pull working days & holidays & leaves first
	workingDays, _ := h.repo.ListWorkingDays(c.Context(), req.DepartmentID)
	holidays, _ := h.repo.ListHolidaysForMonth(c.Context(), req.DepartmentID, req.Month)
	leaves, _ := h.repo.ListLeavesForMonth(c.Context(), req.DepartmentID, req.Month)

	log.Printf("=== DATA LOADED: %d working days, %d holidays, %d leaves ===",
		len(workingDays), len(holidays), len(leaves))

	log.Printf("=== WORKING DAYS MAP: %+v ===", workingDays)

	isHoliday := func(d time.Time) bool {
		ds := d.Format("2006-01-02")
		for _, h := range holidays {
			if ds >= h.Start && ds <= h.End {
				return true
			}
		}
		return false
	}
	isOnLeave := func(staffID string, d time.Time) bool {
		ds := d.Format("2006-01-02")
		for _, lv := range leaves {
			if lv.StaffID == staffID && ds >= lv.Start && ds <= lv.End {
				return true
			}
		}
		return false
	}

	// Get priorities from database to create dynamic algorithm
	type Priority struct {
		Name   string
		Order  int
		Value  interface{}
		Active bool
	}

	priorities := []Priority{}
	// This should be replaced with actual API call to priority service
	// For now, hardcode based on the current setup
	priorities = append(priorities, Priority{"วันที่ขอหยุด", 1, nil, true})
	priorities = append(priorities, Priority{"จำนวนเวรเท่ากันในแต่ละประเภท", 2, 1, true})
	priorities = append(priorities, Priority{"จำนวนชั่วโมงการทำงานทั้งหมด", 3, 20, true})
	priorities = append(priorities, Priority{"จำนวนเวรติดต่อกัน", 4, 5, true})

	log.Printf("=== DYNAMIC PRIORITIES LOADED: %d priorities ===", len(priorities))

	// Dynamic pick function that respects priority order
	pickWithPriorities := func(cands []string, date time.Time, sh database.ShiftRecord, relaxLevel int) (string, bool) {
		best := ""
		bestScore := int(^uint(0) >> 1)

		log.Printf("=== PICK: candidates=%d, relaxLevel=%d ===", len(cands), relaxLevel)

		for _, uid := range cands {
			score := 0
			canAssign := true
			skipReasons := []string{}

			// Apply priority constraints based on relaxLevel
			switch {
			case relaxLevel <= 1: // Strictest - respect all priorities
				// Priority 1: วันที่ขอหยุด (NEVER relax this)
				if isOnLeave(uid, date) {
					skipReasons = append(skipReasons, "on-leave")
					canAssign = false
					break
				}

				// Priority 2: จำนวนเวรเท่ากันในแต่ละประเภท
				score += assignmentCount[uid] * 100

				// Priority 3: จำนวนชั่วโมงการทำงานทั้งหมด (time overlap)
				if !canAssignShift(uid, date, sh) {
					skipReasons = append(skipReasons, "time-overlap")
					canAssign = false
					break
				}

				// Priority 4: จำนวนเวรติดต่อกัน
				if d, exists := lastAssignedDate[uid]; exists {
					if d.AddDate(0, 0, 1).Equal(date) {
						skipReasons = append(skipReasons, "consecutive-day")
						score += 1000 // High penalty
					}
				}

			case relaxLevel <= 3: // Relax consecutive days
				// Priority 1: วันที่ขอหยุด (NEVER relax)
				if isOnLeave(uid, date) {
					skipReasons = append(skipReasons, "on-leave")
					canAssign = false
					break
				}

				// Priority 2: จำนวนเวรเท่ากันในแต่ละประเภท
				score += assignmentCount[uid] * 100

				// Priority 3: จำนวนชั่วโมงการทำงานทั้งหมด (time overlap)
				if !canAssignShift(uid, date, sh) {
					skipReasons = append(skipReasons, "time-overlap")
					canAssign = false
					break
				}

				// Priority 4: จำนวนเวรติดต่อกัน (RELAXED - just add penalty)
				if d, exists := lastAssignedDate[uid]; exists {
					if d.AddDate(0, 0, 1).Equal(date) {
						score += 500 // Lower penalty
					}
				}

			case relaxLevel <= 6: // Relax time overlap
				// Priority 1: วันที่ขอหยุด (NEVER relax)
				if isOnLeave(uid, date) {
					skipReasons = append(skipReasons, "on-leave")
					canAssign = false
					break
				}

				// Priority 2: จำนวนเวรเท่ากันในแต่ละประเภท
				score += assignmentCount[uid] * 50

				// Priority 3: จำนวนชั่วโมงการทำงานทั้งหมด (RELAXED - just penalty)
				if !canAssignShift(uid, date, sh) {
					score += 200 // Add penalty but don't block
				}

				// Priority 4: จำนวนเวรติดต่อกัน (RELAXED)
				if d, exists := lastAssignedDate[uid]; exists {
					if d.AddDate(0, 0, 1).Equal(date) {
						score += 100 // Small penalty
					}
				}

			default: // relaxLevel > 6 - Emergency mode, only respect leaves
				// Priority 1: วันที่ขอหยุด (NEVER relax)
				if isOnLeave(uid, date) {
					skipReasons = append(skipReasons, "on-leave")
					canAssign = false
					break
				}

				// All other priorities ignored - just try to balance assignments
				score += assignmentCount[uid] * 10
			}

			if canAssign && score < bestScore {
				bestScore = score
				best = uid
				log.Printf("=== PICK: New best candidate %s (score=%d) ===", uid, score)
			} else if !canAssign {
				log.Printf("=== PICK: Skipped %s (reasons: %v) ===", uid, skipReasons)
			}
		}

		if best != "" {
			log.Printf("=== PICK: Selected %s (score=%d, relaxLevel=%d) ===", best, bestScore, relaxLevel)
		} else {
			log.Printf("=== PICK: No candidates available at relaxLevel %d ===", relaxLevel)
		}

		return best, best != ""
	}

	// (duplicate code removed - already defined above)

	log.Printf("=== STARTING MAIN ASSIGNMENT LOOP: %d days ===", days)

	for day := 1; day <= days; day++ {
		d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dateStr := d.Format("2006-01-02")

		log.Printf("=== DAY %d (%s): checking conditions ===", day, dateStr)

		// skip non-working day or holidays
		weekday := int(d.Weekday())
		w, ok := workingDays[weekday]
		log.Printf("=== DAY %d: weekday=%d, working=%t, exists=%t ===", day, weekday, w, ok)

		// Skip if: not in map OR explicitly marked as not working
		if !ok || !w {
			log.Printf("=== DAY %d: SKIPPED (not working day) ===", day)
			continue
		}
		if isHoliday(d) {
			log.Printf("=== DAY %d: SKIPPED (holiday) ===", day)
			continue
		}

		log.Printf("=== DAY %d: PROCESSING %d shifts ===", day, len(shifts))

		for _, sh := range shifts {
			log.Printf("=== PROCESSING %s %s (need %d nurses, %d assistants) ===",
				dateStr, sh.Name, sh.RequiredNurse, sh.RequiredAsst)

			// nurses - try with progressive relaxation of priorities
			nurseAssigned := 0
			for relaxLevel := 1; relaxLevel <= 10 && nurseAssigned < sh.RequiredNurse && len(nurses) > 0; relaxLevel++ {
				log.Printf("=== NURSES: RelaxLevel %d, assigned %d/%d ===",
					relaxLevel, nurseAssigned, sh.RequiredNurse)

				attemptCount := 0
				maxAttempts := len(nurses) + 5 // Prevent infinite loops

				for nurseAssigned < sh.RequiredNurse && attemptCount < maxAttempts {
					attemptCount++
					sid, ok := pickWithPriorities(nurses, d, sh, relaxLevel)
					if !ok {
						log.Printf("=== NURSES: No candidates at relax level %d ===", relaxLevel)
						break // Try next relax level
					}

					log.Printf("=== NURSES: Assigned %s at relax level %d ===", sid, relaxLevel)
					items = append(items, database.Assignment{ID: uuid.New().String(), DepartmentID: req.DepartmentID, StaffID: sid, ShiftID: sh.ID, ScheduleDate: dateStr, Status: "assigned"})
					assignmentCount[sid]++
					lastAssignedDate[sid] = d
					// record interval for this day
					if assignedIntervals[sid] == nil {
						assignedIntervals[sid] = map[string][][2]int{}
					}
					if s, e, ok := shiftInterval(sh); ok {
						assignedIntervals[sid][dateStr] = append(assignedIntervals[sid][dateStr], [2]int{s, e})
					}
					nurseAssigned++
				}
			}

			if nurseAssigned < sh.RequiredNurse {
				log.Printf("=== WARNING: %s %s still needs %d nurses ===",
					dateStr, sh.Name, sh.RequiredNurse-nurseAssigned)
			}

			// assistants - try with progressive relaxation of priorities
			assistantAssigned := 0
			for relaxLevel := 1; relaxLevel <= 10 && assistantAssigned < sh.RequiredAsst && len(assistants) > 0; relaxLevel++ {
				log.Printf("=== ASSISTANTS: RelaxLevel %d, assigned %d/%d ===",
					relaxLevel, assistantAssigned, sh.RequiredAsst)

				attemptCount := 0
				maxAttempts := len(assistants) + 5 // Prevent infinite loops

				for assistantAssigned < sh.RequiredAsst && attemptCount < maxAttempts {
					attemptCount++
					sid, ok := pickWithPriorities(assistants, d, sh, relaxLevel)
					if !ok {
						log.Printf("=== ASSISTANTS: No candidates at relax level %d ===", relaxLevel)
						break // Try next relax level
					}

					log.Printf("=== ASSISTANTS: Assigned %s at relax level %d ===", sid, relaxLevel)
					items = append(items, database.Assignment{ID: uuid.New().String(), DepartmentID: req.DepartmentID, StaffID: sid, ShiftID: sh.ID, ScheduleDate: dateStr, Status: "assigned"})
					assignmentCount[sid]++
					lastAssignedDate[sid] = d
					if assignedIntervals[sid] == nil {
						assignedIntervals[sid] = map[string][][2]int{}
					}
					if s, e, ok := shiftInterval(sh); ok {
						assignedIntervals[sid][dateStr] = append(assignedIntervals[sid][dateStr], [2]int{s, e})
					}
					assistantAssigned++
				}
			}

			if assistantAssigned < sh.RequiredAsst {
				log.Printf("=== WARNING: %s %s still needs %d assistants ===",
					dateStr, sh.Name, sh.RequiredAsst-assistantAssigned)
			}
		}
	}

	// Post-balance pass: ลดความต่างจำนวนเวรต่อคน (ตามตำแหน่ง) ให้ใกล้กันมากที่สุดภายใต้กฏ
	// อ่านค่าจาก scheduling_priorities ถ้ามี (priority ชื่อ: "จำนวนเวรเท่าเทียมในแต่ละประเภท")
	maxDiffAllowed := 1
	if v, err := h.repo.GetPriorityValue(c.Context(), req.DepartmentID, "จำนวนเวรเท่าเทียมในแต่ละประเภท"); err == nil && v.Valid {
		if v.Int64 >= 0 && v.Int64 <= 5 {
			maxDiffAllowed = int(v.Int64)
		}
	}

	// สร้าง map ช่วยเหลือ
	staffRole := map[string]string{}
	for _, s := range staffList {
		role := strings.ToLower(s.Position)
		if role == "assistant" || strings.Contains(s.Position, "ผู้ช่วย") {
			staffRole[s.ID] = "assistant"
		} else {
			staffRole[s.ID] = "nurse"
		}
	}
	assignedDates := map[string]map[string]bool{} // staffID -> set(date)
	for _, a := range items {
		if assignedDates[a.StaffID] == nil {
			assignedDates[a.StaffID] = map[string]bool{}
		}
		assignedDates[a.StaffID][a.ScheduleDate] = true
	}
	parseDate := func(s string) time.Time {
		d, _ := time.Parse("2006-01-02", s)
		return d
	}
	canAssignOn := func(staffID string, dateStr string) bool {
		// ห้ามชนวันเดียวกัน และห้ามติดกันวันก่อน/วันถัดไป
		if assignedDates[staffID] != nil && assignedDates[staffID][dateStr] {
			return false
		}
		d := parseDate(dateStr)
		prev := d.AddDate(0, 0, -1).Format("2006-01-02")
		next := d.AddDate(0, 0, 1).Format("2006-01-02")
		if (assignedDates[staffID] != nil && assignedDates[staffID][prev]) || (assignedDates[staffID] != nil && assignedDates[staffID][next]) {
			return false
		}
		return true
	}
	// index ของ items ต่อวัน เพื่อหาความขัดแย้ง
	dayToIndices := map[string][]int{}
	for idx, a := range items {
		dayToIndices[a.ScheduleDate] = append(dayToIndices[a.ScheduleDate], idx)
	}
	// ฟังก์ชันคำนวณค่ามากสุด/น้อยสุด และคนที่เกี่ยวข้อง ตามตำแหน่ง
	findExtremes := func(role string) (maxID string, maxCnt int, minID string, minCnt int) {
		maxCnt = -1
		minCnt = int(^uint(0) >> 1)
		for id, r := range staffRole {
			if r != role {
				continue
			}
			c := assignmentCount[id]
			if c > maxCnt {
				maxCnt = c
				maxID = id
			}
			if c < minCnt {
				minCnt = c
				minID = id
			}
		}
		return
	}
	// พยายามสลับจากคนที่ได้เวรเยอะ -> คนน้อย ภายใต้กฏไม่ซ้ำวัน/ไม่ติดวัน
	limitIterations := 500
	rebalanceRole := func(role string) {
		for k := 0; k < limitIterations; k++ {
			highID, highCnt, lowID, lowCnt := findExtremes(role)
			if highCnt-lowCnt <= maxDiffAllowed {
				break
			}
			// มองหางานในวันใดวันหนึ่งที่ highID ถูกมอบหมาย และ lowID ไม่ถูกมอบหมาย
			found := false
			for dateStr, idxList := range dayToIndices {
				// ข้ามถ้า lowID ลางานวันนั้น
				if isOnLeave(lowID, parseDate(dateStr)) {
					continue
				}
				if !canAssignOn(lowID, dateStr) {
					continue
				}
				// ตรวจว่ามี assignment ของ highID ในวันนี้ และไม่มีของ lowID
				hasHigh := false
				hasLow := false
				for _, idx := range idxList {
					if items[idx].StaffID == highID {
						hasHigh = true
					}
					if items[idx].StaffID == lowID {
						hasLow = true
					}
				}
				if !hasHigh || hasLow {
					continue
				}
				// หา assignment ตัวแรกของ highID ในวันนั้น แล้วย้ายให้ lowID หากตำแหน่งตรงกัน
				for _, idx := range idxList {
					if items[idx].StaffID != highID {
						continue
					}
					// เช็คตำแหน่งตรงกัน
					if staffRole[highID] != role || staffRole[lowID] != role {
						continue
					}
					// ป้องกันย้ายทับวันลาของ lowID
					if isOnLeave(lowID, parseDate(dateStr)) {
						continue
					}
					// ย้าย
					items[idx].StaffID = lowID
					// ปรับชุดข้อมูลช่วย
					assignmentCount[highID]--
					assignmentCount[lowID]++
					if assignedDates[lowID] == nil {
						assignedDates[lowID] = map[string]bool{}
					}
					assignedDates[lowID][dateStr] = true
					if assignedDates[highID] != nil {
						delete(assignedDates[highID], dateStr)
					}
					found = true
					break
				}
				if found {
					break
				}
			}
			if !found {
				break
			}
		}
	}
	// Rebalance สำหรับ nurse และ assistant แยกกัน
	rebalanceRole("nurse")
	rebalanceRole("assistant")

	// Enhanced Algorithm complete - save results directly
	log.Printf("=== ENHANCED ALGORITHM COMPLETE: Total items=%d ===", len(items))

	// Skip old algorithm - comment out
	/*
		tryFill := func(role string, cands []string) {
			for dateStr, needMap := range demand {
				for shId, need := range needMap {
					if assignedNA[dateStr] == nil {
						assignedNA[dateStr] = map[string]struct {
							n int
							a int
						}{}
					}
					got := assignedNA[dateStr][shId]
					deficit := 0
					if role == "assistant" {
						deficit = need.a - got.a
					} else {
						deficit = need.n - got.n
					}
					for deficit > 0 {
						best := ""
						bestCnt := int(^uint(0) >> 1)
						d := parseDate(dateStr)
						for _, id := range cands {
							if staffRole[id] != role {
								continue
							}
							if isOnLeave(id, d) {
								continue
							}
							if !canAssignOn(id, dateStr) {
								continue
							}
							if assignmentCount[id] < bestCnt {
								bestCnt = assignmentCount[id]
								best = id
							}
						}
						if best == "" {
							break
						}
						items = append(items, database.Assignment{ID: uuid.New().String(), DepartmentID: req.DepartmentID, StaffID: best, ShiftID: shId, ScheduleDate: dateStr, Status: "assigned"})
						assignmentCount[best]++
						if assignedDates[best] == nil {
							assignedDates[best] = map[string]bool{}
						}
						assignedDates[best][dateStr] = true
						if assignedNA[dateStr] == nil {
							assignedNA[dateStr] = map[string]struct {
								n int
								a int
							}{}
						}
						got = assignedNA[dateStr][shId]
						if role == "assistant" {
							got.a++
						} else {
							got.n++
						}
						assignedNA[dateStr][shId] = got
						deficit--
					}
				}
			}
		}
		tryFill("nurse", nurses)
		tryFill("assistant", assistants)
	*/ // End of old algorithm comment
	if err := h.repo.BulkInsertAssignmentsStaff(c.Context(), items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "สร้างตารางเวรอัตโนมัติสำเร็จ", "data": fiber.Map{"inserted": len(items)}})
}

// AIGenerate delegates schedule generation to Gemini Flash
func (h *ScheduleHandler) AIGenerate(c *fiber.Ctx) error {
	var req struct {
		DepartmentID string `json:"departmentId"`
		Month        string `json:"month"`
	}
	if err := c.BodyParser(&req); err != nil || req.DepartmentID == "" || req.Month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง ต้องระบุ departmentId และ month"})
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ไม่พบ GEMINI_API_KEY ใน environment"})
	}

	shifts, err := h.repo.ListShifts(c.Context(), req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	users, err := h.repo.ListDepartmentUsers(c.Context(), req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// Build prompt with strict JSON instruction
	prompt := strings.Builder{}
	prompt.WriteString("You are a scheduling assistant for hospital nurse shifts.\n")
	prompt.WriteString("Return ONLY valid JSON with this schema: {\"assignments\":[{\"userId\":string,\"shiftId\":string,\"date\":\"YYYY-MM-DD\"}]}\n")
	prompt.WriteString("Users (id, role):\n")
	for _, u := range users {
		prompt.WriteString(u.UserID + "," + u.DepartmentRole + "\n")
	}
	prompt.WriteString("Shifts (id,name,type,start,end,needNurse,needAssistant):\n")
	for _, s := range shifts {
		prompt.WriteString(s.ID + "," + s.Name + "," + s.Type + "," + s.StartTime + "," + s.EndTime + "," + fmtInt(s.RequiredNurse) + "," + fmtInt(s.RequiredAsst) + "\n")
	}
	prompt.WriteString("Target month: " + req.Month + "\n")
	prompt.WriteString("Constraints: balance total hours and contiguous days, respect staff role requirements, fill all required positions per shift per day.\n")

	payload := map[string]any{
		"contents": []map[string]any{{
			"parts": []map[string]string{{"text": prompt.String()}},
		}},
	}
	body, _ := json.Marshal(payload)
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey
	httpClient := &http.Client{Timeout: 30 * time.Second}
	reqHttp, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	reqHttp.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(reqHttp)
	if err != nil || resp == nil || resp.StatusCode >= 300 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "เรียก Gemini ไม่สำเร็จ"})
	}
	defer resp.Body.Close()
	var aiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "แปลงผลลัพธ์ AI ไม่สำเร็จ"})
	}
	var textOut string
	if len(aiResp.Candidates) > 0 && len(aiResp.Candidates[0].Content.Parts) > 0 {
		textOut = aiResp.Candidates[0].Content.Parts[0].Text
	}
	// attempt to extract JSON
	jsonStr := extractJSON(textOut)
	var parsed struct {
		Assignments []struct {
			UserID  string `json:"userId"`
			ShiftID string `json:"shiftId"`
			Date    string `json:"date"`
		} `json:"assignments"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "รูปแบบผลลัพธ์ AI ไม่เป็น JSON ที่กำหนด"})
	}
	var items []database.Assignment
	for _, a := range parsed.Assignments {
		items = append(items, database.Assignment{ID: uuid.New().String(), DepartmentID: req.DepartmentID, UserID: a.UserID, ShiftID: a.ShiftID, ScheduleDate: a.Date, Status: "assigned"})
	}
	if err := h.repo.BulkInsertAssignments(c.Context(), items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "สร้างตารางเวรด้วย AI สำเร็จ", "data": fiber.Map{"inserted": len(items)}})
}

// OptimizeGenerate creates schedules using internal Go optimizer (fairness-weighted greedy)
func (h *ScheduleHandler) OptimizeGenerate(c *fiber.Ctx) error {
	var req struct {
		DepartmentID string `json:"departmentId"`
		Month        string `json:"month"`
	}
	if err := c.BodyParser(&req); err != nil || req.DepartmentID == "" || req.Month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง ต้องระบุ departmentId และ month"})
	}

	shifts, err := h.repo.ListShifts(c.Context(), req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.repo.EnsureStaffSchedulingSchema(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	staffList, err := h.repo.ListDepartmentStaff(c.Context(), req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	working, _ := h.repo.ListWorkingDays(c.Context(), req.DepartmentID)
	holidays, _ := h.repo.ListHolidaysForMonth(c.Context(), req.DepartmentID, req.Month)
	leaves, _ := h.repo.ListLeavesForMonth(c.Context(), req.DepartmentID, req.Month)
	maxDiffAllowed := 1
	if v, err := h.repo.GetPriorityValue(c.Context(), req.DepartmentID, "จำนวนเวรเท่าเทียมในแต่ละประเภท"); err == nil && v.Valid {
		if v.Int64 >= 0 && v.Int64 <= 5 {
			maxDiffAllowed = int(v.Int64)
		}
	}

	out, err := optimizer.SolveMonth(optimizer.Input{
		DepartmentID:   req.DepartmentID,
		Month:          req.Month,
		Shifts:         shifts,
		Staff:          staffList,
		WorkingDays:    working,
		Holidays:       holidays,
		Leaves:         leaves,
		MaxDiffAllowed: maxDiffAllowed,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	if err := h.repo.DeleteByDepartmentAndMonth(c.Context(), req.DepartmentID, req.Month); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	if err := h.repo.BulkInsertAssignmentsStaff(c.Context(), out); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "สร้างตารางเวรด้วย Optimizer (Go) สำเร็จ", "data": fiber.Map{"inserted": len(out)}})
}

func fmtInt(v int) string { return fmt.Sprintf("%d", v) }

func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	// crude extraction between first '{' and last '}'
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
