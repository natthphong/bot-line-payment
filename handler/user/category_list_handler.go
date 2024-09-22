package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natthphong/bot-line-payment/api"
)

type CategoryListRequest struct {
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	BranchCode string `json:"branchCode"`
}

type CategoryObject struct {
	CategoryNameTh      string `json:"categoryNameTh"`
	CategoryNameEng     string `json:"categoryNameEng"`
	CategoryDescription string `json:"categoryDescription"`
	BranchCode          string `json:"branchCode"`
}

type CategoryListResponse struct {
	TotalCount int `json:"totalCount"`
	Message    struct {
		Categories []CategoryObject `json:"categories"`
	} `json:"message"`
}

func CategoryListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req CategoryListRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid request format")
		}

		resp := CategoryListResponse{
			TotalCount: 5,
			Message: struct {
				Categories []CategoryObject `json:"categories"`
			}{
				Categories: []CategoryObject{
					{
						CategoryNameTh:      "หมวดหมู่ 1",
						CategoryNameEng:     "Category 1",
						CategoryDescription: "Description 1",
						BranchCode:          "B001",
					},
					{
						CategoryNameTh:      "หมวดหมู่ 2",
						CategoryNameEng:     "Category 2",
						CategoryDescription: "Description 2",
						BranchCode:          "B002",
					},
				},
			},
		}

		return api.Ok(c, &resp)
	}
}
