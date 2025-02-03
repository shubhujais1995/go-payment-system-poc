package controller

import (
	"net/http"
	"poc/services"
	"poc/utils"

	"github.com/kataras/iris/v12"
)

// Signup handles user registration (signup).
func Signup(ctx iris.Context) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		IsPayer   bool   `json:"is_payer"`
		IsPayee   bool   `json:"is_payee"`
	}

	// Decode the incoming request
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Call the UserService to create the user
	user, err := services.CreateUser(ctx.Request().Context(), req.Email, req.Password, req.FirstName, req.LastName, req.IsPayer, req.IsPayee)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Respond with the created user data
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(user)
}

// Login handles user login.
func Login(ctx iris.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the incoming request
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Call the UserService to verify the credentials
	user, err := services.LoginUser(ctx.Request().Context(), req.Email, req.Password)
	if err != nil {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Generate a JWT or session token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Respond with user data and the generated token
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]interface{}{
		"user_id":    user.UserID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"token":      token,
	})
}

// Logout handles user logout.
func Logout(ctx iris.Context) {
	// Invalidate session or token here

	// Respond with success message
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{
		"message": "Logged out successfully",
	})
}

// UpdateUser updates user details.
func UpdateUser(ctx iris.Context) {
	userID := ctx.Values().GetString("UserID")
	var req struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		IsPayee   bool   `json:"isPayee"`
		IsPayer   bool   `json:"isPayer"`
	}

	// Decode the incoming request
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Call the service to update the user
	err := services.UpdateUser(ctx, userID, req.Email, req.FirstName, req.LastName, req.IsPayee, req.IsPayer)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{"message": "User updated successfully"})
}

// UpdatePayer updates payer-specific details like balance.
func UpdatePayer(ctx iris.Context) {
	var req struct {
		PayerID string  `json:"payer_id"`
		Amount  float64 `json:"amount"`
	}

	// Decode the incoming request
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Call the service to update the payer
	err := services.UpdatePayer(ctx, req.PayerID, req.Amount)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{"message": "Payer updated successfully"})
}

// UpdatePayee updates payee-specific details like balance.
func UpdatePayee(ctx iris.Context) {
	var req struct {
		PayeeID string  `json:"payee_id"`
		Amount  float64 `json:"amount"`
	}

	// Decode the incoming request
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Call the service to update the payee
	err := services.UpdatePayee(ctx, req.PayeeID, req.Amount)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{"message": "Payee updated successfully"})
}
