package utils

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)


func CreatePagination(page, limit int32, totalCount int64) *entities.Pagination{
	if limit == -1 {
		return nil
	}
	
	totalPages, nextPage, prevPage := calculatePaginationValues(int32(totalCount), limit, page)

	return &entities.Pagination{
		Total:       totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}

func calculatePaginationValues(totalCount, limit, page int32) (totalPages int32, nextPage, prevPage *int32) {
	if limit <= 0 {
		limit = 1
	}

	totalPages = (totalCount + limit - 1) / limit

	if page < totalPages {
		next := page + 1
		nextPage = &next
	}

	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	if page == totalPages {
		nextPage = nil
	}

	if page == 1 {
		prevPage = nil
	}

	return totalPages, nextPage, prevPage
}
