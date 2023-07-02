package apiservice

import (
	"github.com/gin-gonic/gin"
	"github.com/zbysir/writeflow/internal/model"
)

func (a *ApiService) RegisterSys(router gin.IRoutes) {
	router.GET("/system/setting", func(ctx *gin.Context) {
		var params model.Setting
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}

		set, err := a.sysRepo.GetSetting(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, set)
	})
	router.GET("/system/plugin_status", func(ctx *gin.Context) {
		ctx.JSON(200, a.flowUsecase.PluginStatus)
	})
	router.PUT("/system/setting", func(ctx *gin.Context) {
		var params model.Setting
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		err = a.sysRepo.SaveSetting(ctx, &params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, "ok")
	})
}
