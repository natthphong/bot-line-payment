package user

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/bot-line-payment/api"
	"github.com/natthphong/bot-line-payment/internal/logz"
	"github.com/natthphong/bot-line-payment/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ProductListRequest struct {
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	CategoryCode string `json:"categoryCode"`
	BranchCode   string `json:"branchCode"`
}

type ProductObject struct {
	ProductId          string  `json:"productId"`
	ProductNameTh      string  `json:"productNameTh"`
	ProductNameEng     string  `json:"productNameEng"`
	ProductDescription string  `json:"productDescription"`
	ProductCode        string  `json:"productCode"`
	Amount             string  `json:"amount"`
	InActive           string  `json:"inActive"`
	ProductType        string  `json:"productType"`
	ProductQuantity    int     `json:"productQuantity"`
	ProductImage       *string `json:"productImage"`
	CategoryCode       string  `json:"categoryCode"`
}

func (p *ProductListRequest) Validate() error {
	if p.Page <= 0 {
		return errors.New("page must be greater than zero")
	}
	if p.Size <= 0 {
		return errors.New("size must be greater than zero")
	}
	if p.BranchCode == "" {
		return errors.New("branchCode cannot be empty")
	}
	return nil
}

type ProductListResponse struct {
	TotalCount int `json:"totalCount"`
	Message    struct {
		Products []ProductObject `json:"products"`
	} `json:"message"`
}

func ProductListHandler(
	tokenFunc middleware.JwtTokenGetUserFunc,
	findAllProductFunc FindAllProductFunc,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := logz.NewLogger()
		var req ProductListRequest
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

		list, total, err := findAllProductFunc(c.Context(), logger, req)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		resp := ProductListResponse{
			TotalCount: total,
			Message: struct {
				Products []ProductObject `json:"products"`
			}{
				Products: list,
			},
		}

		return api.Ok(c, &resp)
	}
}

type FindAllProductFunc func(ctx context.Context, logger *zap.Logger, req ProductListRequest) (resp []ProductObject, total int, err error)

func FindAllProduct(db *pgxpool.Pool) FindAllProductFunc {
	return func(ctx context.Context, logger *zap.Logger, req ProductListRequest) (resp []ProductObject, total int, err error) {
		var queryArgs []interface{}
		i := 1 // Start SQL placeholder counter
		queryArgs = append(queryArgs, req.BranchCode)

		countSql := `
			SELECT count(*)
			FROM tbl_product tp
			LEFT OUTER JOIN tbl_product_branch tpb ON tpb.product_code = tp.product_code`

		if req.CategoryCode != "" {
			countSql += ` LEFT OUTER JOIN tbl_product_category_branch tpcb ON tpb.product_code = tpcb.product_code AND tpb.branch_code = tpcb.branch_code`
			queryArgs = append(queryArgs, req.CategoryCode)
			i++
		}

		countSql += ` WHERE tpb.branch_code = $1`

		if req.CategoryCode != "" {
			countSql += ` AND tpcb.category_code = $` + fmt.Sprint(i)
		}

		err = db.QueryRow(ctx, countSql, queryArgs...).Scan(&total)
		if err != nil {
			return resp, 0, errors.Wrap(err, "failed to count products")
		}

		// Reset the queryArgs for the second query
		queryArgs = []interface{}{req.BranchCode}
		i = 1 // Reset the counter

		sql := `
			SELECT tp.product_code, tp.product_id, tp.amount, tp.product_type,
			       tp.product_quantity, tp.product_image, tp.product_name_th,
			       tp.product_name_eng, tp.product_description, tpb.in_active`

		if req.CategoryCode != "" {
			sql += `, tpcb.category_code`
		}

		sql += `
			FROM tbl_product tp
			LEFT OUTER JOIN tbl_product_branch tpb ON tpb.product_code = tp.product_code`

		if req.CategoryCode != "" {
			sql += ` LEFT OUTER JOIN tbl_product_category_branch tpcb ON tpb.product_code = tpcb.product_code AND tpb.branch_code = tpcb.branch_code`
			queryArgs = append(queryArgs, req.CategoryCode)
			i++
		}

		sql += ` WHERE tpb.branch_code = $1`

		if req.CategoryCode != "" {
			sql += ` AND tpcb.category_code = $` + fmt.Sprint(i)
		}

		// Add pagination and limit
		sql += ` ORDER BY tp.product_id ASC OFFSET $` + fmt.Sprint(i+1) + ` LIMIT $` + fmt.Sprint(i+2)
		queryArgs = append(queryArgs, (req.Page-1)*req.Size, req.Size)

		fmt.Println(sql)
		rows, err := db.Query(ctx, sql, queryArgs...)
		if err != nil {
			logger.Error("failed to query products", zap.Error(err))
			return resp, total, errors.Wrap(err, "failed to retrieve products")
		}

		for rows.Next() {
			var product ProductObject
			if req.CategoryCode != "" {
				err = rows.Scan(
					&product.ProductCode, &product.ProductId, &product.Amount,
					&product.ProductType, &product.ProductQuantity, &product.ProductImage,
					&product.ProductNameTh, &product.ProductNameEng, &product.ProductDescription,
					&product.InActive,
					&product.CategoryCode,
				)
			} else {
				err = rows.Scan(
					&product.ProductCode, &product.ProductId, &product.Amount,
					&product.ProductType, &product.ProductQuantity, &product.ProductImage,
					&product.ProductNameTh, &product.ProductNameEng, &product.ProductDescription,
					&product.InActive,
				)
			}

			if err != nil {
				logger.Error("failed to scan product", zap.Error(err))
				return resp, total, err
			}
			resp = append(resp, product)
		}

		return resp, total, nil
	}
}
