package apiservice

import (
	"github.com/gin-gonic/gin"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/repo"
)

func (a *ApiService) RegisterDocument(router gin.IRoutes) {
	router.POST("/document/document", func(ctx *gin.Context) {
		var params model.Document
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}

		id, err := a.documentRepo.SaveDocument(ctx, params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, id)
	})

	router.GET("/document/document_list", func(ctx *gin.Context) {
		var params repo.GetArticleListParams
		err := ctx.ShouldBind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ds, total, err := a.documentRepo.GetDocumentList(ctx, params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"list":  ds,
			"total": total,
		})
	})
}
