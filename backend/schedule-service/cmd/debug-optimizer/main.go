package main

import (
	"fmt"
	"os"
	"time"

	db "nurseshift/schedule-service/internal/infrastructure/database"
	"nurseshift/schedule-service/internal/optimizer"
)

func main() {
	os.Setenv("SCHEDULE_DEBUG", "1")

	month := time.Now().Format("2006-01")
	// Mock shifts: 2 shifts/day, eachต้องการพยาบาล 1 และผู้ช่วย 1
	shifts := []db.ShiftRecord{
		{ID: "S1", DepartmentID: "D1", Name: "เช้า", StartTime: "08:00", EndTime: "18:00", RequiredNurse: 1, RequiredAsst: 1},
		{ID: "S2", DepartmentID: "D1", Name: "เย็น", StartTime: "10:00", EndTime: "20:00", RequiredNurse: 1, RequiredAsst: 1},
	}
	// Mock staff: 6 พยาบาล + 6 ผู้ช่วย (มีชื่อซ้ำเพื่อทดสอบ)
	nurses := []db.DepartmentStaff{
		{ID: "N1", DepartmentID: "D1", Name: "นักรพงศ์ ศิลปวุฒิ", Position: "พยาบาล"},
		{ID: "N2", DepartmentID: "D1", Name: "A A", Position: "พยาบาล"},
		{ID: "N3", DepartmentID: "D1", Name: "B B", Position: "พยาบาล"},
		{ID: "N4", DepartmentID: "D1", Name: "C C", Position: "พยาบาล"},
		{ID: "N5", DepartmentID: "D1", Name: "D D", Position: "พยาบาล"},
		{ID: "N6", DepartmentID: "D1", Name: "E E", Position: "พยาบาล"},
	}
	assts := []db.DepartmentStaff{
		{ID: "A1", DepartmentID: "D1", Name: "A A", Position: "ผู้ช่วยพยาบาล"},
		{ID: "A2", DepartmentID: "D1", Name: "B B", Position: "ผู้ช่วยพยาบาล"},
		{ID: "A3", DepartmentID: "D1", Name: "S S", Position: "ผู้ช่วยพยาบาล"},
		{ID: "A4", DepartmentID: "D1", Name: "T T", Position: "ผู้ช่วยพยาบาล"},
		{ID: "A5", DepartmentID: "D1", Name: "test1 test1", Position: "ผู้ช่วยพยาบาล"},
		{ID: "A6", DepartmentID: "D1", Name: "test test", Position: "ผู้ช่วยพยาบาล"},
	}
	staff := append(nurses, assts...)

	// Working days: จันทร์-ศุกร์
	working := map[int]bool{0: false, 1: true, 2: true, 3: true, 4: true, 5: true, 6: false}

	out, err := optimizer.SolveMonth(optimizer.Input{
		DepartmentID:   "D1",
		Month:          month,
		Shifts:         shifts,
		Staff:          staff,
		WorkingDays:    working,
		Holidays:       nil,
		Leaves:         nil,
		MaxDiffAllowed: 1,
	})
	if err != nil {
		panic(err)
	}
	// สรุปจำนวนเวรต่อคน
	cnt := map[string]int{}
	name := map[string]string{}
	for _, s := range staff {
		name[s.ID] = s.Name
	}
	for _, a := range out {
		cnt[a.StaffID]++
	}
	fmt.Println("=== RESULT ===")
	for _, s := range staff {
		fmt.Printf("%-15s -> %d\n", name[s.ID], cnt[s.ID])
	}
}
