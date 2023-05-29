使用配置描述多个任务的依赖关系。
## 特性

> 目前大部份特性都待开发。

- 可配置多个节点的参数关系
- 支持并行节点
- 所以节点逻辑都是 Golang 代码，支持直接编辑来热更新节点。
  - 可以考虑 JS，不过通常 Golang 更容易管理些，不用再去下载其他包（管理包的功能很难），自带包就能满足需求。
- 应用市场: 支持任何 Golang Github 仓库作为应用，支持热插拔，基于 [yaegi](https://github.com/traefik/yaegi)。
- 节点之间可以通过“插座”或者“连线”方式链接。
- 节点自己不会携带输入，应该使用 输入节点 来定义任何输入。

## 目标场景

由于 AI 的工作流需要快速的更改来发掘更有效的使用方式，所以需要一个更简单的编排工具，其有很多开箱即用的工具，能快速通过连线来构建工作流程。

引用 langchina 的一句话：

> Large language models (LLMs) are emerging as a transformative technology, enabling developers to build applications that they previously could not. But using these LLMs in isolation is often not enough to create a truly powerful app - the real power comes when you can combine them with other sources of computation or knowledge.
>

> 大型语言模型 (LLM) 正在成为一种变革性技术，使开发人员能够构建他们以前无法构建的应用程序。但是单独使用这些 LLM 往往不足以创建一个真正强大的应用程序——当您可以将它们与其他计算或知识来源相结合时，真正的力量就来了。
>


## 参考项目

- UI 可以直接抄 [https://github.com/FlowiseAI/Flowise](https://github.com/FlowiseAI/Flowise)，包括他使用的 reactflow 框架。

## 使用方法

定义命令

```go
f.RegisterCmd("hello", func (ctx context.Context, name string) string {
return "hello: " + name
})

f.RegisterCmd("append", func (ctx context.Context, args []string) string {
return strings.Join(args, " ")
})

```

定义流程

```yaml
version: 1

flow:
  append:
    inputs:
      - _args

  hello:
    inputs:
      - append[0]

  END:
    inputs:
      - hello[0]
```

有两种方法配置任务间的依赖关系
- 通过 inputs，inputs 声明的所有参数都会被自动计算并作为 CMD 的参数传入
- 通过 depends，depends 声明的任务会更先执行，如果你不需要传递参数，则可以使用 depends。

inputs 和 depends 可以同时使用。inputs 会被 depend 更先执行。

默认情况下依赖任务会并行执行，你可以通过配置: parallel 来定义可以并行执行的任务数量，默认为 10.

## 边界

### 类型转换

除了使用 CMDer 接口来定义 CMD，你也可以直接使用 CMDFun 来定义 CMD。它通过反射来调用方法，更适用于简单场景。

使用 CMDFun 时你需要尽量保证 flow 定义的输入类型和函数的参数类型一致，如果不一致 writeflow 会尽量帮你做类型转换：

| 输入类型     | CMD 输入类型 | 是否支持                  |
|----------|----------|-----------------------|
| string   | int      | No                    |
| string   | any      | Yes                   |
| any      | string   | 当 any 是 string 时才支持   |
| struct{} | struct{} | 将通过 json copy 的方式尝试装换 |

你也可以通过 配置：objectConversion 来自定义类型转换器。

### 错误处理
每个 CMD 都可以返回 error, 有任何一个 CMD 产生 error 都会停止整个流程的调度，你可以通过配置 : retry, 来配置重试策略，默认会重试 3 次。

## 计划
- [ ] 可视化 UI
  - [ ] 流程配置 - reactflow
  - [ ] 运行状态
- [ ] 分布式调度，支持重启恢复；这不是最优先级的，因为这个项目的编排能力是我最感兴趣的，我要优先实现它。
