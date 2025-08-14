package optimizer

import (
	"fmt"
	"os"
	"sort"
	"time"

	"nurseshift/schedule-service/internal/infrastructure/database"

	"github.com/google/uuid"
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
	debug := os.Getenv("SCHEDULE_DEBUG") == "1"
	dlog := func(format string, a ...any) {
		if debug {
			fmt.Printf("[optimizer] "+format+"\n", a...)
		}
	}
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
	staffName := map[string]string{}
	for _, s := range in.Staff {
		role := s.Position
		staffName[s.ID] = s.Name
		if role == "assistant" || role == "ผู้ช่วยพยาบาล" || role == "ผู้ช่วย" {
			assistantIDs = append(assistantIDs, s.ID)
		} else {
			nurseIDs = append(nurseIDs, s.ID)
		}
	}
	dlog("inputs: nurses=%d assistants=%d shifts=%d month=%s", len(nurseIDs), len(assistantIDs), len(in.Shifts), in.Month)

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
	// allow multiple non-overlapping shifts/day with max contiguous-hour limit
	assignedIntervals := map[string]map[string][][2]int{} // staffID -> date -> list of [start,end] minutes

	// helpers for time
	parseHM := func(hm string) (int, bool) {
		var hh, mm int
		if _, err := fmt.Sscanf(hm, "%d:%d", &hh, &mm); err != nil {
			return 0, false
		}
		return hh*60 + mm, true
	}
	shiftInterval := func(sh database.ShiftRecord) (int, int, bool) {
		s, ok1 := parseHM(sh.StartTime)
		e, ok2 := parseHM(sh.EndTime)
		if !ok1 || !ok2 {
			return 0, 0, false
		}
		if e <= s {
			e += 24 * 60
		}
		return s, e, true
	}
	overlaps := func(aS, aE, bS, bE int) bool { return aS < bE && bS < aE }
	mergeAndMaxContiguous := func(ivals [][2]int) int {
		if len(ivals) == 0 {
			return 0
		}
		sort.Slice(ivals, func(i, j int) bool { return ivals[i][0] < ivals[j][0] })
		curS, curE := ivals[0][0], ivals[0][1]
		maxDur := curE - curS
		for i := 1; i < len(ivals); i++ {
			s, e := ivals[i][0], ivals[i][1]
			if s <= curE {
				if e > curE {
					curE = e
				}
			} else {
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
	maxContiguousMinutes := 16 * 60

	// helper
	isEligible := func(staffID, date string, d time.Time, sh database.ShiftRecord) bool {
		if leave[staffID][date] {
			return false
		}
		if prev, ok := lastDay[staffID]; ok {
			pd := time.Date(year, m, prev, 0, 0, 0, 0, time.UTC)
			if pd.AddDate(0, 0, 1).Equal(d) {
				return false
			}
		}
		s, e, ok := shiftInterval(sh)
		if !ok {
			return false
		}
		if assignedIntervals[staffID] == nil {
			assignedIntervals[staffID] = map[string][][2]int{}
		}
		cur := assignedIntervals[staffID][date]
		for _, iv := range cur {
			if overlaps(iv[0], iv[1], s, e) {
				return false
			}
		}
		merged := append(append([][2]int{}, cur...), [2]int{s, e})
		if mergeAndMaxContiguous(merged) > maxContiguousMinutes {
			return false
		}
		return true
	}
	reason := func(staffID, date string, d time.Time, sh database.ShiftRecord) string {
		if leave[staffID][date] {
			return "leave"
		}
		if prev, ok := lastDay[staffID]; ok {
			pd := time.Date(year, m, prev, 0, 0, 0, 0, time.UTC)
			if pd.AddDate(0, 0, 1).Equal(d) {
				return "consecutive-day"
			}
		}
		s, e, ok := shiftInterval(sh)
		if !ok {
			return "invalid-shift"
		}
		for _, iv := range assignedIntervals[staffID][date] {
			if overlaps(iv[0], iv[1], s, e) {
				return "overlap"
			}
		}
		merged := append(append([][2]int{}, assignedIntervals[staffID][date]...), [2]int{s, e})
		if mergeAndMaxContiguous(merged) > maxContiguousMinutes {
			return "exceed-contiguous-hours"
		}
		return "unknown"
	}
	cost := func(role, staffID string) int {
		t := roleTarget[role]
		diff := count[staffID] - t
		if diff <= 0 { // under target preferred
			return (diff) * 5 // negative or zero; greedy picks minimum
		}
		return diff * diff * 10 // square penalty
	}

	// Seed pass: assure at least 1 shift for everyone if capacity allows
	type needNA struct{ n, a int }
	capacity := map[string]map[string]*needNA{}
	for day := 1; day <= days; day++ {
		d := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
		if w, ok := in.WorkingDays[int(d.Weekday())]; ok && !w {
			continue
		}
		if isHoliday(d) {
			continue
		}
		ds := d.Format("2006-01-02")
		capacity[ds] = map[string]*needNA{}
		for _, sh := range in.Shifts {
			capacity[ds][sh.ID] = &needNA{n: sh.RequiredNurse, a: sh.RequiredAsst}
		}
	}

	seedOnce := func(ids []string, role string) {
		for _, id := range ids {
			if count[id] > 0 {
				continue
			}
			assigned := false
			dlog("seed: try %s(%s)", staffName[id], role)
			for day := 1; day <= days && !assigned; day++ {
				d := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
				if w, ok := in.WorkingDays[int(d.Weekday())]; ok && !w {
					continue
				}
				if isHoliday(d) {
					continue
				}
				dateStr := d.Format("2006-01-02")
				for _, sh := range in.Shifts {
					rem := 0
					if role == "assistant" {
						rem = capacity[dateStr][sh.ID].a
					} else {
						rem = capacity[dateStr][sh.ID].n
					}
					if rem <= 0 {
						continue
					}
					if !isEligible(id, dateStr, d, sh) {
						dlog("seed: block %s %s shift=%s reason=%s", staffName[id], dateStr, sh.Name, reason(id, dateStr, d, sh))
						continue
					}
					assignments = append(assignments, database.Assignment{ID: RandID(), DepartmentID: in.DepartmentID, StaffID: id, ShiftID: sh.ID, ScheduleDate: dateStr, Status: "assigned"})
					count[id]++
					lastDay[id] = day
					if assignedIntervals[id] == nil {
						assignedIntervals[id] = map[string][][2]int{}
					}
					if s, e, ok := shiftInterval(sh); ok {
						assignedIntervals[id][dateStr] = append(assignedIntervals[id][dateStr], [2]int{s, e})
					}
					if role == "assistant" {
						capacity[dateStr][sh.ID].a--
					} else {
						capacity[dateStr][sh.ID].n--
					}
					assigned = true
					dlog("seed: assigned %s on %s shift=%s", staffName[id], dateStr, sh.Name)
					break
				}
			}
			if !assigned {
				dlog("seed: fail %s no slot or constraints", staffName[id])
			}
		}
	}
	seedOnce(nurseIDs, "nurse")
	seedOnce(assistantIDs, "assistant")

	for day := 1; day <= days; day++ {
		d := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
		if w, ok := in.WorkingDays[int(d.Weekday())]; ok && !w {
			continue
		}
		if isHoliday(d) {
			continue
		}
		dateStr := d.Format("2006-01-02")

		for _, sh := range in.Shifts {
			// Nurses
			need := 0
			if capacity[dateStr][sh.ID] != nil {
				need = capacity[dateStr][sh.ID].n
			}
			for need > 0 {
				best := ""
				bestCost := 1 << 30
				for _, id := range nurseIDs {
					if !isEligible(id, dateStr, d, sh) {
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
				if assignedIntervals[best] == nil {
					assignedIntervals[best] = map[string][][2]int{}
				}
				if s, e, ok := shiftInterval(sh); ok {
					assignedIntervals[best][dateStr] = append(assignedIntervals[best][dateStr], [2]int{s, e})
				}
				if capacity[dateStr][sh.ID] != nil {
					capacity[dateStr][sh.ID].n--
				}
				need--
			}
			// Assistants
			needA := 0
			if capacity[dateStr][sh.ID] != nil {
				needA = capacity[dateStr][sh.ID].a
			}
			for needA > 0 {
				best := ""
				bestCost := 1 << 30
				for _, id := range assistantIDs {
					if !isEligible(id, dateStr, d, sh) {
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
				if assignedIntervals[best] == nil {
					assignedIntervals[best] = map[string][][2]int{}
				}
				if s, e, ok := shiftInterval(sh); ok {
					assignedIntervals[best][dateStr] = append(assignedIntervals[best][dateStr], [2]int{s, e})
				}
				if capacity[dateStr][sh.ID] != nil {
					capacity[dateStr][sh.ID].a--
				}
				needA--
			}
		}
	}

	// Build helper maps for rebalancing
	shiftByID := map[string]database.ShiftRecord{}
	for _, sh := range in.Shifts {
		shiftByID[sh.ID] = sh
	}
	staffRole := map[string]string{}
	for _, s := range in.Staff {
		r := s.Position
		if r == "assistant" || r == "ผู้ช่วยพยาบาล" || r == "ผู้ช่วย" {
			staffRole[s.ID] = "assistant"
		} else {
			staffRole[s.ID] = "nurse"
		}
	}

	// Ensure everyone gets at least one shift if possible by swapping
	ensureMinimumOne := func(roleIDs []string) {
		for _, lowID := range roleIDs {
			if count[lowID] > 0 {
				continue
			}
			// try to steal from someone with count > 1
			for i := range assignments {
				a := assignments[i]
				if staffRole[a.StaffID] != staffRole[lowID] {
					continue
				}
				if count[a.StaffID] <= 1 {
					continue
				}
				// check eligibility
				d, _ := time.Parse("2006-01-02", a.ScheduleDate)
				sh := shiftByID[a.ShiftID]
				if !isEligible(lowID, a.ScheduleDate, d, sh) {
					continue
				}
				// update intervals maps
				if s, e, ok := shiftInterval(sh); ok {
					// remove from holder
					old := assignedIntervals[a.StaffID][a.ScheduleDate]
					nxt := make([][2]int, 0, len(old))
					for _, iv := range old {
						if !(iv[0] == s && iv[1] == e) {
							nxt = append(nxt, iv)
						}
					}
					assignedIntervals[a.StaffID][a.ScheduleDate] = nxt
					// add to lowID
					assignedIntervals[lowID][a.ScheduleDate] = append(assignedIntervals[lowID][a.ScheduleDate], [2]int{s, e})
				}
				assignments[i].StaffID = lowID
				count[lowID]++
				count[a.StaffID]--
				break
			}
		}
	}
	ensureMinimumOne(nurseIDs)
	ensureMinimumOne(assistantIDs)
	if debug {
		zeros := []string{}
		for _, id := range nurseIDs {
			if count[id] == 0 {
				zeros = append(zeros, staffName[id])
			}
		}
		dlog("post-minimum nurses zero=%v", zeros)
		zeros = []string{}
		for _, id := range assistantIDs {
			if count[id] == 0 {
				zeros = append(zeros, staffName[id])
			}
		}
		dlog("post-minimum assistants zero=%v", zeros)
	}

	// Gentle re-balance to reduce spread towards role targets
	rebalance := func(role string, ids []string) {
		// iterate limited times to avoid long loops
		for iter := 0; iter < 400; iter++ {
			// find high and low
			highID, lowID := "", ""
			highCnt := -1
			lowCnt := 1<<31 - 1
			for _, id := range ids {
				if count[id] > highCnt {
					highCnt, highID = count[id], id
				}
				if count[id] < lowCnt {
					lowCnt, lowID = count[id], id
				}
			}
			if highCnt-lowCnt <= in.MaxDiffAllowed {
				break
			}
			// try to move one assignment from high to low
			moved := false
			for i := range assignments {
				if assignments[i].StaffID != highID {
					continue
				}
				a := assignments[i]
				d, _ := time.Parse("2006-01-02", a.ScheduleDate)
				sh := shiftByID[a.ShiftID]
				if !isEligible(lowID, a.ScheduleDate, d, sh) {
					continue
				}
				// swap
				if s, e, ok := shiftInterval(sh); ok {
					// remove from high
					old := assignedIntervals[highID][a.ScheduleDate]
					nxt := make([][2]int, 0, len(old))
					for _, iv := range old {
						if !(iv[0] == s && iv[1] == e) {
							nxt = append(nxt, iv)
						}
					}
					assignedIntervals[highID][a.ScheduleDate] = nxt
					// add to low
					assignedIntervals[lowID][a.ScheduleDate] = append(assignedIntervals[lowID][a.ScheduleDate], [2]int{s, e})
				}
				assignments[i].StaffID = lowID
				count[lowID]++
				count[highID]--
				moved = true
				break
			}
			if !moved {
				break
			}
		}
	}
	if in.MaxDiffAllowed <= 0 {
		in.MaxDiffAllowed = 1
	}
	rebalance("nurse", nurseIDs)
	rebalance("assistant", assistantIDs)

	// Stronger distribution: try to donate from high to any lower-count eligible candidate on the same day/shift
	distribute := func(role string, ids []string) {
		for iter := 0; iter < 600; iter++ {
			// find current high and low counts
			highID := ""
			highCnt := -1
			lowCnt := 1<<31 - 1
			for _, id := range ids {
				if count[id] > highCnt {
					highCnt, highID = count[id], id
				}
				if count[id] < lowCnt {
					lowCnt = count[id]
				}
			}
			if highCnt-lowCnt <= in.MaxDiffAllowed {
				break
			}
			moved := false
			for i := range assignments {
				a := assignments[i]
				if a.StaffID != highID {
					continue
				}
				// candidate list = everyone in same role with strictly less count than highCnt
				cands := []string{}
				for _, id := range ids {
					if count[id] < highCnt {
						cands = append(cands, id)
					}
				}
				if len(cands) == 0 {
					continue
				}
				// try assign to the lowest-count eligible candidate
				d, _ := time.Parse("2006-01-02", a.ScheduleDate)
				sh := shiftByID[a.ShiftID]
				bestID := ""
				bestCnt := 1<<31 - 1
				for _, cid := range cands {
					if !isEligible(cid, a.ScheduleDate, d, sh) {
						continue
					}
					if count[cid] < bestCnt {
						bestCnt, bestID = count[cid], cid
					}
				}
				if bestID == "" {
					continue
				}
				// perform reassignment and update interval maps
				if s, e, ok := shiftInterval(sh); ok {
					old := assignedIntervals[highID][a.ScheduleDate]
					nxt := make([][2]int, 0, len(old))
					for _, iv := range old {
						if !(iv[0] == s && iv[1] == e) {
							nxt = append(nxt, iv)
						}
					}
					assignedIntervals[highID][a.ScheduleDate] = nxt
					assignedIntervals[bestID][a.ScheduleDate] = append(assignedIntervals[bestID][a.ScheduleDate], [2]int{s, e})
				}
				assignments[i].StaffID = bestID
				count[bestID]++
				count[highID]--
				moved = true
				break
			}
			if !moved {
				break
			}
		}
	}
	distribute("nurse", nurseIDs)
	distribute("assistant", assistantIDs)
	if debug {
		zeros := []string{}
		for _, id := range nurseIDs {
			if count[id] == 0 {
				zeros = append(zeros, staffName[id])
			}
		}
		dlog("summary nurses zero=%v", zeros)
		zeros = []string{}
		for _, id := range assistantIDs {
			if count[id] == 0 {
				zeros = append(zeros, staffName[id])
			}
		}
		dlog("summary assistants zero=%v", zeros)
	}
	return assignments, nil
}

// RandID returns a pseudo-random ID using time.Now() to avoid extra deps; acceptable for internal bulk inserts.
func RandID() string { return uuid.New().String() }
