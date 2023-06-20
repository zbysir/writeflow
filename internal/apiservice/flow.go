package apiservice

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/repo"
	"time"
)

type IdReq struct {
	Id  int64   `json:"id" form:"id"`
	Ids []int64 `json:"ids" form:"ids"`
}

type KeyReq struct {
	Key string `json:"key" form:"key"`
}

type RunFlowReq struct {
	Id           int64                  `json:"id"`
	Params       map[string]interface{} `json:"params"`
	Graph        *model.Graph           `json:"graph"`
	Parallel     int                    `json:"parallel"`
	OutputNodeId string                 `json:"output_node_id"`
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
		//cs = cs.Upgrade()
		ctx.JSON(200, cs)
	})

	router.POST("/flow", func(ctx *gin.Context) {
		var params model.Flow
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		id, err := a.flowRepo.CreateFlow(ctx, &params)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(200, id)
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
		if params.Id != 0 {
			params.Ids = append(params.Ids, params.Id)
		}

		for _, id := range params.Ids {
			err = a.flowRepo.DeleteFlow(ctx, id)
			if err != nil {
				ctx.Error(err)
				return
			}
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
		if params.Id == 0 && params.Graph == nil {
			ctx.Error(fmt.Errorf("id or graph must be set"))
			return
		}
		if params.Graph != nil {
			r, err := a.flowUsecase.RunFlowByDetail(context.Background(), &model.Flow{
				Graph: *params.Graph,
			}, params.Params, params.Parallel)
			if err != nil {
				ctx.Error(err)
				return
			}
			ctx.JSON(200, r)
		} else {
			r, err := a.flowUsecase.RunFlow(context.Background(), params.Id, params.Params, params.Parallel)
			if err != nil {
				ctx.Error(err)
				return
			}
			ctx.JSON(200, r)
		}
	})

	router.POST("/flow/run_sync", func(ctx *gin.Context) {
		var params RunFlowReq
		err := ctx.Bind(&params)
		if err != nil {
			ctx.Error(err)
			return
		}
		if params.Id == 0 && params.Graph == nil {
			ctx.Error(fmt.Errorf("id or graph must be set"))
			return
		}
		start := time.Now()
		if params.Graph != nil {
			r, err := a.flowUsecase.RunFlowByDetailSync(context.Background(), &model.Flow{
				Graph: *params.Graph,
			}, params.Params, params.Parallel)
			if err != nil {
				ctx.Error(err)
				return
			}
			ctx.Header("x-spend", fmt.Sprintf("%v", time.Since(start)))
			ctx.JSON(200, r.Raw())
		} else {
			r, err := a.flowUsecase.RunFlowSync(context.Background(), params.Id, params.Params, params.Parallel, params.OutputNodeId)
			if err != nil {
				ctx.Error(err)
				return
			}
			ctx.Header("x-spend", fmt.Sprintf("%v", time.Since(start)))
			ctx.JSON(200, r.Raw())
		}
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
