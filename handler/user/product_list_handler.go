package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natthphong/bot-line-payment/api"
)

type ProductListRequest struct {
	Page        int    `json:"page"`
	Size        int    `json:"size"`
	CompanyCode string `json:"companyCode"` // Optional, default = "ALL"
	BranchCode  string `json:"branchCode"`
}

type ProductObject struct {
	ProductId          string `json:"productId"`
	ProductNameTh      string `json:"productNameTh"`
	ProductNameEng     string `json:"productNameEng"`
	ProductDescription string `json:"productDescription"`
	ProductCode        string `json:"productCode"`
	Amount             string `json:"amount"`
	InActive           bool   `json:"inActive"`
	ProductType        string `json:"productType"`
	ProductQuantity    int    `json:"productQuantity"`
	ProductImage       string `json:"productImage"`
}

type ProductListResponse struct {
	TotalCount int `json:"totalCount"`
	Message    struct {
		Products []ProductObject `json:"products"`
	} `json:"message"`
}

func ProductListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req ProductListRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid request format")
		}

		resp := ProductListResponse{
			TotalCount: 20,
			Message: struct {
				Products []ProductObject `json:"products"`
			}{
				Products: []ProductObject{
					{
						ProductId:          "P001",
						ProductNameTh:      "ผลิตภัณฑ์ 1",
						ProductNameEng:     "Product 1",
						ProductDescription: "Description 1",
						ProductCode:        "PC001",
						Amount:             "1000",
						InActive:           false,
						ProductType:        "Type 1",
						ProductQuantity:    100,
						ProductImage:       "https://example.com/product1.png",
					},
					{
						ProductId:          "P002",
						ProductNameTh:      "ผลิตภัณฑ์ 2",
						ProductNameEng:     "Product 2",
						ProductDescription: "Description 2",
						ProductCode:        "PC002",
						Amount:             "2000",
						InActive:           true,
						ProductType:        "Type 2",
						ProductQuantity:    200,
						ProductImage:       "https://example.com/product2.png",
					},
				},
			},
		}

		return api.Ok(c, &resp)
	}
}
