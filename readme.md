使用配置描述多个任务的依赖关系。

## 相似产品
- [Concepts](https://conductor.netflix.com/devguide/concepts/index.html)
- [Flowise](https://github.com/FlowiseAI/Flowise)
- [PostMan Flow](https://learning.postman.com/docs/postman-flows/gs/flows-overview/)
## 特性

> 目前大部份特性都待开发。

- 节点之间可以通过“连线”方式链接来描述输入输出关系。
- 默认尽可能的并行节点。
- 所以节点逻辑都是 Golang 代码，支持直接编辑来热更新节点。
  - 可以考虑 JS，不过通常 Golang 更容易管理些，不用再去下载其他包（很难管理包），自带包就能满足需求。
- 应用市场: 支持任何 Golang Github 仓库作为应用，支持热插拔，基于 [yaegi](https://github.com/traefik/yaegi)。
  - 支持按照模块 或者 CMD 来安装扩展

### 暂时不会做的事
- 管理 Task 状态，实现优雅重启、重试等。你可以理解目前 writeflow 是内存型工作流。

## 目标场景

### LangChain UI
由于 AI 的工作流需要快速的更改来发掘更有效的使用方式，所以需要一个更简单的编排工具，其有很多开箱即用的工具，能快速通过连线来构建工作流程。

引用 langchina 的一句话：

> Large language models (LLMs) are emerging as a transformative technology, enabling developers to build applications that they previously could not. But using these LLMs in isolation is often not enough to create a truly powerful app - the real power comes when you can combine them with other sources of computation or knowledge.
>

> 大型语言模型 (LLM) 正在成为一种变革性技术，使开发人员能够构建他们以前无法构建的应用程序。但是单独使用这些 LLM 往往不足以创建一个真正强大的应用程序——当您可以将它们与其他计算或知识来源相结合时，真正的力量就来了。
>

### 数据处理
给定一个对象，依此对单个属性进行判断有效性、处理、解析出其他属性，然后再合并成为一个对象。

## 参考项目

- UI 可以直接抄 [https://github.com/FlowiseAI/Flowise](https://github.com/FlowiseAI/Flowise)，包括他使用的 reactflow 框架。

## 使用方法
暂无文档

## 概念

### Flow
Flow 定义一个工作流，一个 Flow 由多个 Node 组成，Node 之间通过连线来描述输入输出关系。

### Node
Node 是一个节点，多个节点组成 Flow，Node 由 Component 实例化而来，保存了如位置、输入输出等信息。

### Component
Component 保存了名称、描述、Cmd 等信息，Component 是 Node 的模板。

#### Cmd
Cmd 是 Component 的运行命令，支持 Golang 代码、远端。

## 边界

### 错误处理
每个 CMD 都可以返回 error, 有任何一个 CMD 产生 error 都会停止整个流程的调度，你可以通过配置 : retry, 来配置重试策略，默认会重试 3 次。

## 计划
- [ ] 可视化 UI: [writeflow-ui](https://github.com/zbysir/writeflow-ui)
  - [x] 流程配置
  - [x] 运行状态
- [ ] 分布式调度，支持重启恢复，持久化；这不是最优先级的，因为这个项目的编排能力是我最感兴趣的，我要优先实现它。
- 逻辑分支
  - [x] Switch
  - [x] For
- 并行执行
  - [x] 并行执行
- [x] 组件商店
