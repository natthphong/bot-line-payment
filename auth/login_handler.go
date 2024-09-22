package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/bot-line-payment/api"
	"github.com/natthphong/bot-line-payment/internal/logz"
	"github.com/natthphong/bot-line-payment/model/request"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type LineErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
type LineResponse struct {
	Iss     string   `json:"iss"`
	Sub     string   `json:"sub"`
	Aud     string   `json:"aud"`
	Exp     int64    `json:"exp"`
	Iat     int64    `json:"iat"`
	Amr     []string `json:"amr"`
	Name    string   `json:"name"`
	Picture string   `json:"picture"`
}

func LoginHandler(
	clientId string,
	insertAndMergeUserLineLoginFunc InsertAndMergeUserLineLoginFunc,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := logz.NewLogger()
		var req request.LoginRequest
		err := c.BodyParser(&req)
		if err != nil {
			logger.Error(err.Error())
			return api.BadRequest(c, "BAD_REQUEST")
		}
		//cal api
		LineRes, err := verifyIDToken(req.IDToken, clientId)
		if err != nil {
			logger.Error(err.Error())
			return api.BadRequest(c, err.Error())
		}

		var lineResponse LineResponse
		err = json.Unmarshal(LineRes, &lineResponse)
		if err != nil {
			return api.InternalError(c, err.Error())

		}

		err = insertAndMergeUserLineLoginFunc(c.Context(), logger, lineResponse)
		if err != nil {
			return api.InternalError(c, "DB Error")

		}
		fmt.Println(lineResponse.Exp > time.Now().Unix())
		return api.Ok(c, lineResponse)
	}
}

func verifyIDToken(idToken, clientID string) ([]byte, error) {
	apiURL := "https://api.line.me/oauth2/v2.1/verify"
	data := url.Values{}
	data.Set("id_token", idToken)
	data.Set("client_id", clientID)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return body, fmt.Errorf("%s", body)
	}

	return body, nil
}

type InsertAndMergeUserLineLoginFunc func(ctx context.Context, logger *zap.Logger, lineResponse LineResponse) error

func InsertAndMergeUserLineLogin(db *pgxpool.Pool) InsertAndMergeUserLineLoginFunc {
	return func(ctx context.Context, logger *zap.Logger, lineResponse LineResponse) error {

		query := `
			INSERT INTO tbl_line_user (
				userId, client_id, display_name, id_token, picture, expires_in, create_date, create_by, is_login
			) VALUES (
				$1, $2, $3, $4, $5, $6, now(), 'SYSTEM', 'Y'
			) ON CONFLICT (userId)
			DO UPDATE SET
				id_token = EXCLUDED.id_token,
				picture = EXCLUDED.picture,
				expires_in = EXCLUDED.expires_in,
				update_date = now(),
				update_by = 'SYSTEM',
				is_login = 'Y'
		`

		_, err := db.Exec(ctx, query,
			lineResponse.Sub,
			lineResponse.Aud,
			lineResponse.Name,
			lineResponse.Sub,
			lineResponse.Picture,
			lineResponse.Exp,
		)

		if err != nil {
			logger.Error("Failed to insert or update user data", zap.Error(err))
			return err
		}

		logger.Debug("User data inserted or updated successfully", zap.String("userId", lineResponse.Sub))
		return nil
	}
}
