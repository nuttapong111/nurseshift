package optimizer

import (
	"time"

	"nurseshift/schedule-service/internal/infrastructure/database"
)

type Input struct {
	DepartmentID   string
	Month          string // YYYY-MM
	Shifts         []database.ShiftRecord
	Staff          []database.DepartmentStaff
	WorkingDays    map[int]bool          // 0=Sun..6=Sat
	Holidays       []database.Holiday    // Start/End = YYYY-MM-DD
	Leaves         []database.LeaveRange // StaffID, Start/End
	MaxDiffAllowed int
}

// SolveMonth builds assignments using fairness-weighted greedy with hard constraints (no same-day, no consecutive-day, no leave/holiday/non-working).
func SolveMonth(in Input) ([]database.Assignment, error) {
	// Parse month metadata
	t, err := time.Parse("2006-01", in.Month)
	if err != nil {
		return nil, err
	}
	year, m, _ := t.Date()
	first := time.Date(year, m, 1, 0, 0, 0, 0, time.UTC)
	next := first.AddDate(0, 1, 0)
	days := int(next.Sub(first).Hours() / 24)

	// Build holiday and leave maps
	isHoliday := func(d time.Time) bool {
		ds := d.Format("2006-01-02")
		for _, h := range in.Holidays {
			if ds >= h.Start && ds <= h.End {
				return true
			}
		}
		return false
	}
	leave := map[string]map[string]bool{}
	for _, lv := range in.Leaves {
		if leave[lv.StaffID] == nil {
			leave[lv.StaffID] = map[string]bool{}
		}
		// mark all days in range
		start, _ := time.Parse("2006-01-02", lv.Start)
		end, _ := time.Parse("2006-01-02", lv.End)
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			leave[lv.StaffID][d.Format("2006-01-02")] = true
		}
	}

	// Split roles
	nurseIDs := []string{}
	assistantIDs := []string{}
	for _, s := range in.Staff {
		role := s.Position
		if role == "assistant" || role == "ผู้ช่วยพยาบาล" || role == "ผู้ช่วย" {
			assistantIDs = append(assistantIDs, s.ID)
		} else {
			nurseIDs = append(nurseIDs, s.ID)
		}
	}

	// Targets per role
	countRoleSlots := func(role string) int {
		total := 0
		for day := 1; day <= days; day++ {
			d := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
			if w, ok := in.WorkingDays[int(d.Weekday())]; ok && !w {
				continue
			}
			if isHoliday(d) {
				continue
			}
			for _, sh := range in.Shifts {
				if role == "assistant" {
					total += sh.RequiredAsst
				} else {
					total += sh.RequiredNurse
				}
			}
		}
		return total
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	ceilDiv := func(a, b int) int {
		if b <= 0 {
			return a
		}
		return (a + b - 1) / b
	}
	roleTarget := map[string]int{
		"nurse":     ceilDiv(countRoleSlots("nurse"), max(len(nurseIDs), 1)),
		"assistant": ceilDiv(countRoleSlots("assistant"), max(len(assistantIDs), 1)),
	}

	// State trackers
	assignments := []database.Assignment{}
	count := map[string]int{}
	lastDay := map[string]int{}
	assignedOnDate := map[string]map[string]bool{} // staffID->date->true

	// helper
	isEligible := func(staffID, date string, d time.Time) bool {
		if leave[staffID][date] {
			return false
		}
		if assignedOnDate[staffID] != nil && assignedOnDate[staffID][date] {
			return false
		}
		if prev, ok := lastDay[staffID]; ok {
			// no consecutive day
			pd := time.Date(year, m, prev, 0, 0, 0, 0, time.UTC)
			if pd.AddDate(0, 0, 1).Equal(d) {
				return false
			}
		}
		return true
	}
	cost := func(role, staffID string) int {
		t := roleTarget[role]
		diff := count[staffID] - t
		if diff <= 0 { // under target preferred
			return (diff) * 5 // negative or zero; greedy picks minimum
		}
		return diff * diff * 10 // square penalty
	}

	for day := 1; day <= days; day++ {
		d := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
		if w, ok := in.WorkingDays[int(d.Weekday())]; ok && !w {
			continue
		}
		if isHoliday(d) {
			continue
		}
		dateStr := d.Format("2006-01-02")

		// For each shift do role-wise greedy with fairness cost
		for _, sh := range in.Shifts {
			// Nurses
			need := sh.RequiredNurse
			for need > 0 {
				best := ""
				bestCost := 1 << 30
				for _, id := range nurseIDs {
					if !isEligible(id, dateStr, d) {
						continue
					}
					c := cost("nurse", id)
					if c < bestCost {
						bestCost = c
						best = id
					}
				}
				if best == "" {
					break
				}
				assignments = append(assignments, database.Assignment{ID: RandID(), DepartmentID: in.DepartmentID, StaffID: best, ShiftID: sh.ID, ScheduleDate: dateStr, Status: "assigned"})
				count[best]++
				lastDay[best] = day
				if assignedOnDate[best] == nil {
					assignedOnDate[best] = map[string]bool{}
				}
				assignedOnDate[best][dateStr] = true
				need--
			}
			// Assistants
			needA := sh.RequiredAsst
			for needA > 0 {
				best := ""
				bestCost := 1 << 30
				for _, id := range assistantIDs {
					if !isEligible(id, dateStr, d) {
						continue
					}
					c := cost("assistant", id)
					if c < bestCost {
						bestCost = c
						best = id
					}
				}
				if best == "" {
					break
				}
				assignments = append(assignments, database.Assignment{ID: RandID(), DepartmentID: in.DepartmentID, StaffID: best, ShiftID: sh.ID, ScheduleDate: dateStr, Status: "assigned"})
				count[best]++
				lastDay[best] = day
				if assignedOnDate[best] == nil {
					assignedOnDate[best] = map[string]bool{}
				}
				assignedOnDate[best][dateStr] = true
				needA--
			}
		}
	}
	return assignments, nil
}

// RandID returns a pseudo-random ID using time.Now() to avoid extra deps; acceptable for internal bulk inserts.
func RandID() string {
	return time.Now().UTC().Format("20060102150405.000000000")
}

