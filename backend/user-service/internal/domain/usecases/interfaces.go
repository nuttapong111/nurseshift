package usecases

// PasswordService defines the interface for password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) bool
	GenerateRandomPassword(length int) string
}

