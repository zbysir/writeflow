package apiservice

import (
	"github.com/gin-gonic/gin"
	"github.com/zbysir/writeflow/internal/repo"
)

func (a *ApiService) RegisterRepo(router gin.IRoutes) {
	// 获取所有的 repo
	router.GET("/flow", func(ctx *gin.Context) {
		var params repo.GetFlowListParams
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		cs, total, err := a.flowRepo.GetComponentList(ctx, params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"total": total,
			"list":  cs,
		})
	})

}
