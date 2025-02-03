package services

import (
	"context"
	"errors"
	"poc/initializer"
	"poc/model"
	"poc/utils"
	"time"

	"gorm.io/gorm"
)

// CreateUser creates a new user in the database.
func CreateUser(ctx context.Context, email, password, firstName, lastName string, isPayer, isPayee bool) (*model.User, error) {
	db := initializer.GetDB()

	// Check if the email already exists
	var existingUser model.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := utils.HashPasswordWrapper(password)
	if err != nil {
		return nil, err
	}

	// Create new user model
	user := &model.User{
		UserID:       utils.GenerateUniqueID(),
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert user into the database using GORM
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	// Create Payer or Payee record if necessary
	if isPayer {
		payer := &model.Payer{
			PayerID: user.UserID,
			Name:    user.FirstName + " " + user.LastName,
			Email:   user.Email,
			Balance: 0.0, // Initial balance
			Status:  "active",
		}
		if err := db.Create(payer).Error; err != nil {
			return nil, err
		}
	}

	if isPayee {
		payee := &model.Payee{
			PayeeID: user.UserID,
			Name:    user.FirstName + " " + user.LastName,
			Email:   user.Email,
			Balance: 0.0, // Initial balance
			Status:  "active",
		}
		if err := db.Create(payee).Error; err != nil {
			return nil, err
		}
	}

	// Return user and nil error
	return user, nil
}

// LoginUser checks the credentials and returns the user.
func LoginUser(ctx context.Context, email, password string) (*model.User, error) {
	db := initializer.GetDB()
	var user model.User

	// Fetch user by email
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Return the user if credentials are valid
	return &user, nil
}

// LogoutUser handles user logout by invalidating the session (or token).
// Note: This method doesn't have much logic for now, as session invalidation is handled in the controller.
func LogoutUser(ctx context.Context) error {
	// If you're using sessions or JWT, you'd invalidate here
	// For example, clear a JWT from the context or database session, if applicable.
	return nil
}

// UpdateUser updates the details of a user in the database.
func UpdateUser(ctx context.Context, userID, email, firstName, lastName string, isPayee bool, isPayer bool) error {
	db := initializer.GetDB()
	// Fetch the user by ID
	var user model.User
	if err := db.First(&user, "UserID = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Update the user fields
	user.Email = email
	user.FirstName = firstName
	user.LastName = lastName
	user.UpdatedAt = time.Now()

	// Save the changes
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// UpdatePayer updates the balance of a payer in the database.
func UpdatePayer(ctx context.Context, payerID string, balance float64) error {
	db := initializer.GetDB()
	// Fetch the payer by ID
	var payer model.Payer
	if err := db.First(&payer, "PayerID = ?", payerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payer not found")
		}
		return err
	}

	// Update the payer balance
	payer.Balance = balance
	payer.UpdatedAt = time.Now()

	// Save the changes
	if err := db.Save(&payer).Error; err != nil {
		return err
	}

	return nil
}

// UpdatePayee updates the balance of a payee in the database.
func UpdatePayee(ctx context.Context, payeeID string, balance float64) error {
	db := initializer.GetDB()
	// Fetch the payee by ID
	var payee model.Payee
	if err := db.First(&payee, "PayeeID = ?", payeeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payee not found")
		}
		return err
	}

	// Update the payee balance
	payee.Balance = balance
	payee.UpdatedAt = time.Now()

	// Save the changes
	if err := db.Save(&payee).Error; err != nil {
		return err
	}

	return nil
}
