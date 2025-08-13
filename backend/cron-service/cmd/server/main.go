package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	db     *sql.DB
	schema string
}

func NewCronService(db *sql.DB, schema string) *CronService {
	return &CronService{
		db:     db,
		schema: schema,
	}
}

// UpdateUserDaysRemaining อัปเดตจำนวนวันใช้งานคงเหลือของผู้ใช้ทุกคน
func (cs *CronService) UpdateUserDaysRemaining() {
	log.Println("🔄 Starting daily update of user days remaining...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// อัปเดตจำนวนวันคงเหลือสำหรับผู้ใช้ที่มี subscription_expires_at
	query := fmt.Sprintf(`
		UPDATE %s.users 
		SET 
			days_remaining = GREATEST(0, 
				EXTRACT(EPOCH FROM (subscription_expires_at - NOW())) / 86400
			),
			updated_at = NOW()
		WHERE 
			subscription_expires_at IS NOT NULL 
			AND subscription_expires_at > NOW()
			AND status = 'active'
	`, cs.schema)

	result, err := cs.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("❌ Error updating days remaining: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Updated days remaining for %d users", rowsAffected)

	// อัปเดตสถานะผู้ใช้ที่หมดอายุแล้ว
	expireQuery := fmt.Sprintf(`
		UPDATE %s.users 
		SET 
			status = 'suspended',
			days_remaining = 0,
			updated_at = NOW()
		WHERE 
			subscription_expires_at IS NOT NULL 
			AND subscription_expires_at <= NOW()
			AND status = 'active'
	`, cs.schema)

	expireResult, err := cs.db.ExecContext(ctx, expireQuery)
	if err != nil {
		log.Printf("❌ Error updating expired users: %v", err)
		return
	}

	expiredRows, _ := expireResult.RowsAffected()
	if expiredRows > 0 {
		log.Printf("⚠️  Suspended %d expired users", expiredRows)
	}

	log.Println("✅ Daily update completed successfully")
}

// LogCurrentStatus แสดงสถานะปัจจุบันของผู้ใช้
func (cs *CronService) LogCurrentStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_users,
			COUNT(CASE WHEN days_remaining > 0 THEN 1 END) as active_users,
			COUNT(CASE WHEN days_remaining = 0 THEN 1 END) as expired_users,
			AVG(days_remaining) as avg_days_remaining
		FROM %s.users 
		WHERE status = 'active'
	`, cs.schema)

	var totalUsers, activeUsers, expiredUsers int
	var avgDaysRemaining sql.NullFloat64

	err := cs.db.QueryRowContext(ctx, query).Scan(&totalUsers, &activeUsers, &expiredUsers, &avgDaysRemaining)
	if err != nil {
		log.Printf("❌ Error getting user status: %v", err)
		return
	}

	log.Printf("📊 User Status Report:")
	log.Printf("   Total Users: %d", totalUsers)
	log.Printf("   Active Users: %d", activeUsers)
	log.Printf("   Expired Users: %d", expiredUsers)
	if avgDaysRemaining.Valid {
		log.Printf("   Average Days Remaining: %.1f", avgDaysRemaining.Float64)
	}
}

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "nurseshift")
	dbSchema := getEnv("DB_SCHEMA", "nurse_shift")

	// Connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database successfully")

	// Create cron service
	cronService := NewCronService(db, dbSchema)

	// Create cron scheduler
	c := cron.New(cron.WithLocation(time.UTC))

	// Schedule daily update at midnight UTC (7 AM Thailand time)
	_, err = c.AddFunc("0 0 * * *", cronService.UpdateUserDaysRemaining)
	if err != nil {
		log.Fatalf("❌ Failed to schedule daily update: %v", err)
	}

	// Schedule status logging every 6 hours
	_, err = c.AddFunc("0 */6 * * *", cronService.LogCurrentStatus)
	if err != nil {
		log.Fatalf("❌ Failed to schedule status logging: %v", err)
	}

	// Run initial update and status check
	log.Println("🚀 Running initial update...")
	cronService.UpdateUserDaysRemaining()
	cronService.LogCurrentStatus()

	// Start cron scheduler
	c.Start()
	log.Println("✅ Cron service started successfully")
	log.Println("📅 Scheduled daily update at midnight UTC (7 AM Thailand time)")
	log.Println("📅 Scheduled status logging every 6 hours")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Shutting down cron service...")
	c.Stop()
	log.Println("✅ Cron service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
