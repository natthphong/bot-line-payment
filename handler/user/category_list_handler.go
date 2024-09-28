package user

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/bot-line-payment/api"
	"github.com/natthphong/bot-line-payment/internal/logz"
	"github.com/natthphong/bot-line-payment/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CategoryListRequest struct {
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	BranchCode string `json:"branchCode"`
}

func (c *CategoryListRequest) Validate() error {
	if c.Page <= 0 {
		return errors.New("page must be greater than zero")
	}
	if c.Size <= 0 {
		return errors.New("size must be greater than zero")
	}
	return nil
}

type CategoryObject struct {
	CategoryNameTh      string `json:"categoryNameTh"`
	CategoryNameEng     string `json:"categoryNameEng"`
	CategoryDescription string `json:"categoryDescription"`
	CategoryCode        string `json:"categoryCode"`
}

type CategoryListResponse struct {
	TotalCount int `json:"totalCount"`
	Message    struct {
		Categories []CategoryObject `json:"categories"`
	} `json:"message"`
}

func CategoryListHandler(
	tokenFunc middleware.JwtTokenGetUserFunc,
	findAllCategoryFunc FindAllCategoryFunc,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := logz.NewLogger()
		var req CategoryListRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid request format")
		}
		if err := req.Validate(); err != nil {
			return api.BadRequest(c, err.Error())
		}
		_, err := tokenFunc(c, logger)
		if err != nil {
			return api.JwtError(c, err.Error())
		}

		list, total, err := findAllCategoryFunc(c.Context(), logger, req)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		resp := CategoryListResponse{
			TotalCount: total,
			Message: struct {
				Categories []CategoryObject `json:"categories"`
			}{
				Categories: list,
			},
		}

		return api.Ok(c, &resp)
	}
}

type FindAllCategoryFunc func(ctx context.Context, logger *zap.Logger, req CategoryListRequest) (resp []CategoryObject, total int, err error)

func FindAllCategory(db *pgxpool.Pool) FindAllCategoryFunc {
	return func(ctx context.Context, logger *zap.Logger, req CategoryListRequest) (resp []CategoryObject, total int, err error) {
		countSql := `
				select count(*) from tbl_category_branch tcb	where tcb.branch_code = $1;
				`
		err = db.QueryRow(ctx, countSql, req.BranchCode).Scan(&total)
		if err != nil {
			return resp, 0, err
		}
		sql := `
			select tc.category_code, tc.category_name_th,tc.category_name_eng,tc.category_description
			from tbl_category_branch tcb
			inner join tbl_category tc on tcb.category_code = tc.category_code
			where tcb.branch_code = $1
			order by tc.category_id asc
				offset $2 limit $3
				
				`
		rows, err := db.Query(ctx, sql, req.BranchCode, (req.Page-1)*req.Size, req.Size)
		if err != nil {
			logger.Error(err.Error())
			return resp, total, err
		}
		for rows.Next() {
			var category CategoryObject
			err = rows.Scan(&category.CategoryCode, &category.CategoryNameTh, &category.CategoryNameEng, &category.CategoryDescription)
			if err != nil {
				return resp, total, err
			}
			resp = append(resp, category)
		}

		return resp, total, err
	}
}
