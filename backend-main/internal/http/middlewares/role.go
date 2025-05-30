package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/libs/jwt"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func isJWTExpire(tokenExp int64) bool {
	return time.Now().Unix() >= tokenExp
}

func validateAuthHeader(authToken []string) (string, error) {
	if len(authToken) == 0 || authToken[0] == "" {
		return "", errors.New("missing authorization token")
	}

	tokenData := strings.Split(authToken[0], " ")

	if len(tokenData) != 2 {
		return "", errors.New("wrong auth token format")
	}

	return tokenData[1], nil
}

// role - the minimal role level which will have access to resource
func RoleBasedAccess(secret string, requiredJob entities.UserJob, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		op := slog.String("op", "middlewares.RoleBasedAccess")

		token, err := validateAuthHeader(r.Header["Authorization"])
		if err != nil {
			slog.Info("validate token error", op, slog.Any("auth header", r.Header["Authorization"]),
				slog.String("error", err.Error()))

			if err = utils.RespondWith400(w, err.Error()); err != nil {
				slog.Error("failed respond with 400: validateAuthHeader", op, slog.String("error", err.Error()))
			}
			return
		}

		claims, err := extractClaims(token, secret)
		if err != nil {
			slog.Error("failed extract token claims: token is invalid", op, slog.String("token", token), slog.String("error", err.Error()))
			if err := utils.RespondWith400(w, err.Error()); err != nil {
				slog.Error("failed respond with 400: parse claims", op, slog.String("error", err.Error()))
			}
			return
		}

		jwtData, err := jwt.GetDataFromClaims(claims)
		if err != nil {
			slog.Error("failed parse jwt data", slog.String("error", err.Error()))

			if err = utils.RespondWith400(w, "wrong jwt token data"); err != nil {
				slog.Error("failed respond with 400: parse claims", op, slog.String("error", err.Error()))
			}
			return
		}

		if isJWTExpire(jwtData.Exp) {
			if err := utils.RespondWith401(w, "token is expired"); err != nil {
				slog.Error("failed respond with 400: parse claims", op, slog.String("error", err.Error()))
			}
			return
		}

		if !hasUserPermission(jwtData.JobPosition, requiredJob) {
			slog.Info("user has no permission to data", slog.String("path", r.URL.Path),
				slog.String("user_job", jwtData.JobPosition))

			if err := utils.RespondWith400(w, "user have no access to this resource"); err != nil {
				slog.Error("failed respond with 400: parse claims", op, slog.String("error", err.Error()))
			}
			return
		}

		// user has access to data
		slog.Info("handle request with auth",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		)

		ctx := context.WithValue(context.Background(), "user_id", jwtData.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
