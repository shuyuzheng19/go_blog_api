package handler

import (
	"blog/internal/service"
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"bufio"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Define database table enum
type DataBaseTable string

const (
	BLOG     DataBaseTable = "BLOG"
	TAG      DataBaseTable = "TAG"
	BlogTag  DataBaseTable = "BLOG_TAG"
	FILE     DataBaseTable = "FILE"
	FileMd5  DataBaseTable = "FILE_MD5"
	Category DataBaseTable = "CATEGORY"
	Topic    DataBaseTable = "TOPICS"
	Role     DataBaseTable = "ROLE"
	USER     DataBaseTable = "USER"
)

// DataBaseController struct
type DataBaseController struct {
	service service.DataBaseService
	mu      sync.Mutex
}

// Create a new DataBaseController
func NewDataBaseController() *DataBaseController {
	return &DataBaseController{
		service: service.NewDataBaseService(),
	}
}

// Execute SQL handler
func (d *DataBaseController) ExecSQL(ctx fiber.Ctx) error {
	// Validate API Key
	apiKey := ctx.Query("apiKey")
	if apiKey != configs.CONFIG.DataBaseKey {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "Invalid API key")
	}

	// Read request body
	body := ctx.Body()
	if len(body) == 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "Failed to bind parameters")
	}

	// Asynchronously execute SQL
	go func(sql string) {
		d.mu.Lock()
		defer d.mu.Unlock()
		d.service.ExecDataBaseSQL(sql)
		logger.Info("Executing SQL", zap.String("sql", sql))
	}(string(body))

	return ResultSuccessToResponse(nil, ctx)
}

// Handle insert SQL with streaming
func (d *DataBaseController) handleInsertSQL(
	w *bufio.Writer,
	sqlFunc func(page int) []string,
	size int,
) error {
	for page := 1; ; page++ {
		sqls := sqlFunc(page)
		if len(sqls) == 0 {
			break
		}

		for _, sql := range sqls {
			if _, err := w.WriteString("data: " + sql + "\n\n"); err != nil {
				return err
			}
		}

		w.Flush()

		if page*50 > size {
			break
		}
	}

	return nil
}

// Get table insert SQL
func (d *DataBaseController) GetTableInsertSQL(ctx fiber.Ctx) error {
	// Validate API Key
	apiKey := ctx.Query("apiKey")
	if apiKey != configs.CONFIG.DataBaseKey {
		return ResultErrorToResponse(common.Unauthorized, ctx, "Invalid API key")
	}

	// Parse parameters
	tableType := DataBaseTable(ctx.Query("TYPE"))
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		size = 50
	}

	// Use a map to replace switch
	sqlFuncMap := map[DataBaseTable]func(page int) []string{
		BLOG:     d.service.GetBlogInsertSQL,
		TAG:      d.service.GetTagInsertSQL,
		BlogTag:  d.service.GetBlogTagInsertSQL,
		FILE:     d.service.GetFileInsertSQL,
		FileMd5:  d.service.GetFileMd5InsertSQL,
		Category: d.service.GetCategoryInsertSQL,
		Topic:    d.service.GetTopicInsertSQL,
		Role:     d.service.GetRoleInsertSQL,
		USER:     d.service.GetUserInsertSQL,
	}

	// Find the corresponding SQL generation function
	sqlFunc, exists := sqlFuncMap[tableType]
	if !exists {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "Processing failed")
	}

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")

	ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := d.handleInsertSQL(w, sqlFunc, size); err != nil {
			logger.Error("Error streaming SQL", zap.Error(err))
		}
	}))

	return nil
}
