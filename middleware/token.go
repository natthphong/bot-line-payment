package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

type JwtTokenBody struct {
	IsLogin     bool   `json:"isLogin"`
	UserId      string `json:"userId"`
	DisplayName string `json:"displayName"`
	ExpiresAt   int64  `json:"expiresAt"`
	ClientId    string `json:"clientId"`
}

type JwtTokenGetUserFunc func(c *fiber.Ctx, logger *zap.Logger) (*JwtTokenBody, error)

func JwtTokenGetUser(db *pgxpool.Pool) JwtTokenGetUserFunc {

	return func(c *fiber.Ctx, logger *zap.Logger) (user *JwtTokenBody, err error) {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, errors.New("Invalid Authorization Header")
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		sql := `
			select tlu.userId,tlu.expires_in,tlu.client_id,tlu.is_login
			from tbl_line_user tlu
			where tlu.id_token = $1;
			`
		rows, err := db.Query(c.Context(), sql, token)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			user = &JwtTokenBody{}
			err = rows.Scan(&user.UserId, &user.ExpiresAt, &user.ClientId, &user.IsLogin)
			if err != nil {
				return nil, err
			}
		}
		if user == nil {
			return nil, errors.New("Token Not Found")
		}
		nowUnix := time.Now().Unix()
		if !user.IsLogin {
			return nil, errors.New("Logout Already")
		}
		if nowUnix > user.ExpiresAt {
			return nil, errors.New("Token Expired")
		}
		return user, err
	}
}
