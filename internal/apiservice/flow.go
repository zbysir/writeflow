package apiservice

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/repo"
)

type IdReq struct {
	Id int64 `json:"id" form:"id"`
}

type KeyReq struct {
	Key string `json:"key" form:"key"`
}

type RunFlowReq struct {
	Id     int64                  `json:"id"`
	Params map[string]interface{} `json:"params"`
}

func (a *ApiService) RegisterFlow(router gin.IRoutes) {
	// 获取所有的 repo
	router.GET("/flow", func(ctx *gin.Context) {
		var params repo.GetFlowListParams
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		cs, total, err := a.flowRepo.GetFlowList(ctx, params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"total": total,
			"list":  cs,
		})
	})

	router.GET("/flow_one", func(ctx *gin.Context) {
		var params IdReq
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		cs, exist, err := a.flowRepo.GetFlowById(ctx, params.Id)
		if err != nil {
			ctx.Error(err)
			return
		}
		if !exist {
			ctx.JSON(404, "not found")
			return
		}

		ctx.JSON(200, cs)
	})

	router.POST("/flow", func(ctx *gin.Context) {
		var params model.Flow
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		err = a.flowRepo.CreateFlow(ctx, &params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, "ok")
	})
	router.PUT("/flow", func(ctx *gin.Context) {
		var params model.Flow
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		err = a.flowRepo.UpdateFlow(ctx, &params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, "ok")
	})
	router.DELETE("/flow", func(ctx *gin.Context) {
		var params IdReq
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		err = a.flowRepo.DeleteFlow(ctx, params.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, "ok")
	})

	router.POST("/flow/run", func(ctx *gin.Context) {
		var params RunFlowReq
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		r, err := a.flowUsecase.RunFlow(context.Background(), params.Id, params.Params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, r)
	})

	type GetComponentsParams struct {
	}

	// component
	router.GET("/component", func(ctx *gin.Context) {
		var params GetComponentsParams
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		cs, err := a.flowUsecase.GetComponents(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, cs)
	})

	router.GET("/component_one", func(ctx *gin.Context) {
		var params KeyReq
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		c, exist, err := a.flowUsecase.GetComponentByKey(ctx, params.Key)
		if err != nil {
			ctx.Error(err)
			return
		}
		if !exist {
			ctx.JSON(404, fmt.Sprintf("not found component by key: %s", params.Key))
			return
		}

		ctx.JSON(200, c)
	})

	router.POST("/component", func(ctx *gin.Context) {
		//var params model.Component
		//err := ctx.Bind(&params)
		//if err != nil {
		//	ctx.Error(err)
		//	return
		//}
		//err = a.flowRepo.CreateComponent(ctx, &params)
		//if err != nil {
		//	ctx.Error(err)
		//	return
		//}

		ctx.JSON(200, "ok")
	})
	router.DELETE("/component", func(ctx *gin.Context) {
		//var params KeyReq
		//err := ctx.Bind(&params)
		//if err != nil {
		//	ctx.Error(err)
		//	return
		//}
		//err = a.flowRepo.DeleteComponent(ctx, params.Key)
		//if err != nil {
		//	ctx.Error(err)
		//	return
		//}

		ctx.JSON(200, "ok")
	})
}
