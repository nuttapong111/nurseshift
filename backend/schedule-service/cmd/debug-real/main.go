package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	cfgpkg "nurseshift/schedule-service/internal/infrastructure/config"
	dbpkg "nurseshift/schedule-service/internal/infrastructure/database"
	"nurseshift/schedule-service/internal/optimizer"
)

func main() {
	var departmentID, departmentName, month string
	flag.StringVar(&departmentID, "departmentId", "", "department UUID")
	flag.StringVar(&departmentName, "departmentName", "", "department name (optional)")
	flag.StringVar(&month, "month", "", "YYYY-MM")
	flag.Parse()
	if month == "" {
		log.Fatalf("usage: go run . -departmentId <uuid> -month YYYY-MM  (or provide -departmentName to resolve id)")
	}

	os.Setenv("SCHEDULE_DEBUG", "1")

	cfg, err := cfgpkg.Load()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := dbpkg.NewConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	repo := dbpkg.NewScheduleRepository(conn)

	// Resolve department id by name if not provided
	if departmentID == "" {
		rows, err := conn.DB.QueryContext(context.Background(), "SELECT id, name FROM nurse_shift.departments")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		type row struct{ id, name string }
		var list []row
		for rows.Next() {
			var r row
			if err := rows.Scan(&r.id, &r.name); err != nil {
				log.Fatal(err)
			}
			list = append(list, r)
		}
		if len(list) == 0 {
			log.Fatal("no departments found")
		}
		if departmentName != "" {
			for _, r := range list {
				if strings.Contains(strings.ToLower(r.name), strings.ToLower(departmentName)) {
					departmentID = r.id
					break
				}
			}
		}
		if departmentID == "" {
			departmentID = list[0].id
		}
		fmt.Println("Resolved departmentId:", departmentID)
	}

	shifts, err := repo.ListShifts(context.Background(), departmentID)
	if err != nil {
		log.Fatal(err)
	}
	staff, err := repo.ListDepartmentStaff(context.Background(), departmentID)
	if err != nil {
		log.Fatal(err)
	}
	working, err := repo.ListWorkingDays(context.Background(), departmentID)
	if err != nil {
		log.Fatal(err)
	}
	holidays, err := repo.ListHolidaysForMonth(context.Background(), departmentID, month)
	if err != nil {
		log.Fatal(err)
	}
	leaves, err := repo.ListLeavesForMonth(context.Background(), departmentID, month)
	if err != nil {
		log.Fatal(err)
	}
	maxDiff := 1
	if v, err := repo.GetPriorityValue(context.Background(), departmentID, "จำนวนเวรเท่าเทียมในแต่ละประเภท"); err == nil && v.Valid {
		if v.Int64 >= 0 && v.Int64 <= 5 {
			maxDiff = int(v.Int64)
		}
	}

	out, err := optimizer.SolveMonth(optimizer.Input{
		DepartmentID:   departmentID,
		Month:          month,
		Shifts:         shifts,
		Staff:          staff,
		WorkingDays:    working,
		Holidays:       holidays,
		Leaves:         leaves,
		MaxDiffAllowed: maxDiff,
	})
	if err != nil {
		log.Fatal(err)
	}

	// สรุปผล
	name := map[string]string{}
	for _, s := range staff {
		name[s.ID] = s.Name
	}
	cnt := map[string]int{}
	for _, a := range out {
		cnt[a.StaffID]++
	}
	fmt.Println("=== RESULT (REAL DATA) ===")
	for _, s := range staff {
		fmt.Printf("%-20s -> %d\n", name[s.ID], cnt[s.ID])
	}
}
