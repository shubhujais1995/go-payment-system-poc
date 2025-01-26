package controller

import (
	"net/http"
	"poc/services"
	"poc/utils"

	"github.com/kataras/iris/v12"
)

// UserController handles HTTP requests for user operations.
type UserController struct {
	UserService *services.UserService
}

// Signup handles user registration (signup).
func (uc *UserController) Signup(ctx iris.Context) {
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
	user, err := uc.UserService.CreateUser(ctx.Request().Context(), req.Email, req.Password, req.FirstName, req.LastName, req.IsPayer, req.IsPayee)
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
func (uc *UserController) Login(ctx iris.Context) {
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
	user, err := uc.UserService.LoginUser(ctx.Request().Context(), req.Email, req.Password)
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
		"user":  user,
		"token": token,
	})
}

// Logout handles user logout.
func (uc *UserController) Logout(ctx iris.Context) {
	// Invalidate session or token here

	// Respond with success message
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{
		"message": "Logged out successfully",
	})
}

// UpdateUser updates user details.
func (uc *UserController) UpdateUser(ctx iris.Context) {
	// userID := ctx.Params().Get("id")
	userID := ctx.Values().GetString("UserID")
	var req struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		isPayee   bool   `json:"isPayee"`
		isPayer   bool   `json:"isPayer"`
	}

	// Decode the incoming request
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Call the service to update the user
	err := uc.UserService.UpdateUser(ctx, userID, req.Email, req.FirstName, req.LastName, req.isPayee, req.isPayer)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{"message": "User updated successfully"})
}

// UpdatePayer updates payer-specific details like balance.
func (uc *UserController) UpdatePayer(ctx iris.Context) {
	var req struct {
		PayerID string  `json:"payer_id"`
		Amount  float64 `json:"amount"`
	}

	// fmt.Println("req ", req.PayerID, req.Amount, req)
	// Decode the incoming request
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Call the service to update the payer
	err := uc.UserService.UpdatePayer(ctx, req.PayerID, req.Amount)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{"message": "Payer updated successfully"})
}

// UpdatePayee updates payee-specific details like balance.
func (uc *UserController) UpdatePayee(ctx iris.Context) {
	// payeeID := ctx.Params().Get("id")
	// var req struct {
	// 	Balance float64 `json:"balance"`
	// }

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
	err := uc.UserService.UpdatePayee(ctx, req.PayeeID, req.Amount)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(map[string]string{"message": "Payee updated successfully"})
}
