package helpers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPagingData(ctx *gin.Context) (page, limit int, err error) {
	page, err = strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		err = errors.New("incorrect page format")
		return
	}
	limit, err = strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		err = errors.New("incorrect limit format")
		return
	}
	return
}

func GetOffset(page, limit int) int {
	if page < 1 || limit < 1 {
		return 0
	}
	return (page - 1) * limit
}
