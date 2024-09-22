package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natthphong/bot-line-payment/api"
)

type BranchListRequest struct {
	Page        int    `json:"page"`
	Size        int    `json:"size"`
	CompanyCode string `json:"companyCode"`
	Internal    bool   `json:"internal"`
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

func BranchListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req BranchListRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid request format")
		}

		resp := BranchListResponse{
			TotalCount: 10,
			Message: struct {
				Branches []BranchObject `json:"branches"`
			}{
				Branches: []BranchObject{
					{
						BranchName:        "Branch 1",
						BranchDescription: "Description 1",
						BranchCode:        "B001",
						InActive:          "N",
					},
					{
						BranchName:        "Branch 2",
						BranchDescription: "Description 2",
						BranchCode:        "B002",
						InActive:          "Y",
					},
				},
			},
		}

		return api.Ok(c, &resp)
	}
}
