package httphandler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)



func AuthMiddleWare(next echo.HandlerFunc) echo.HandlerFunc{
	return func(c echo.Context) error {
        // Get token from header
        authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
        splitAuth := strings.Split(authHeader, " ")
        if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
            return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
        }
        
        accessToken := splitAuth[1]
        log.Println("Token:", accessToken)

        // Decode token
        var claim model.CustomClaims
        token, err := jwt.ParseWithClaims(accessToken, &claim, func(token *jwt.Token) (interface{}, error) {
            // Verify signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }
            return []byte(config.JWTSigningKey()), nil
        })

        if err != nil {
            // Check if the error is because the token is expired
            if strings.Contains(err.Error(), "token is expired") {
                return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
            }
            log.Printf("Token Error: %v", err)
            return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
        }

        // Additional validation
        if !token.Valid {
            return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
        }

        // Menggunakan helper function untuk cek expired
        if helper.IsTokenExpired(claim.ExpiresAt) {
            return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
        }

        // Set claims to context
        ctx := context.WithValue(c.Request().Context(), model.BearerAuthKey, claim)
        req := c.Request().WithContext(ctx)
        c.SetRequest(req)
        
        return next(c)
    }
}