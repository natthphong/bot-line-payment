package user

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/bot-line-payment/api"
	"github.com/natthphong/bot-line-payment/internal/logz"
	"github.com/natthphong/bot-line-payment/middleware"
	"go.uber.org/zap"
)

type BranchListRequest struct {
	Page        int    `json:"page"`
	Size        int    `json:"size"`
	CompanyCode string `json:"companyCode"`
	Internal    string `json:"internal"`
}

func (c *BranchListRequest) Validate() error {
	if c.Page <= 0 {
		return errors.New("page must be greater than zero")
	}
	if c.Size <= 0 {
		return errors.New("size must be greater than zero")
	}
	return nil
}

type BranchObject struct {
	BranchName        string `json:"branchName"`
	BranchDescription string `json:"branchDescription"`
	BranchCode        string `json:"branchCode"`
	InActive          string `json:"inActive"`
}

type BranchListResponse struct {
	TotalCount int `json:"totalCount"`
	Message    struct {
		Branches []BranchObject `json:"branches"`
	} `json:"message"`
}

func BranchListHandler(
	tokenFunc middleware.JwtTokenGetUserFunc,
	getListBranchesFunc GetListBranchesFunc,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := logz.NewLogger()
		var req BranchListRequest
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

		list, total, err := getListBranchesFunc(c.Context(), logger, req)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		resp := BranchListResponse{
			TotalCount: total,
			Message: struct {
				Branches []BranchObject `json:"branches"`
			}{
				Branches: list,
			},
		}

		return api.Ok(c, &resp)
	}
}

type GetListBranchesFunc func(ctx context.Context, logger *zap.Logger, req BranchListRequest) (resp []BranchObject, total int, err error)

func GetListBranches(db *pgxpool.Pool) GetListBranchesFunc {
	return func(ctx context.Context, logger *zap.Logger, req BranchListRequest) (resp []BranchObject, total int, err error) {
		countSql := `SELECT COUNT(*) FROM tbl_company_branch tcb where ($1='' or tcb.company_code=$1) and ($2='' or tcb.internal=$2)`
		err = db.QueryRow(ctx, countSql, req.CompanyCode, req.Internal).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		sql := `
			select tcb.branch_code, tcb.branch_name, tcb.branch_description, tcb.in_active 
			FROM tbl_company_branch tcb where ($1='' or tcb.company_code=$1) and ($2='' or tcb.internal=$2)
			offset $3 limit $4
			`
		rows, err := db.Query(ctx, sql, req.CompanyCode, req.Internal, (req.Page-1)*req.Size, req.Size)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		for rows.Next() {
			var temp BranchObject
			err = rows.Scan(&temp.BranchCode, &temp.BranchName, &temp.BranchDescription, &temp.InActive)
			if err != nil {
				return nil, 0, err
			}
			resp = append(resp, temp)
		}
		return resp, total, err
	}
}
