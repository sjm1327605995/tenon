# Documentation / 文档索引

Start with the [English README](../README.md) or [中文 README](../README.zh-CN.md) for installation, a complete application, and runnable examples.

## Reading paths / 阅读路径

| Goal / 目标 | Read / 入口 |
| --- | --- |
| Build your first app / 编写第一个应用 | [Quick start](../README.md#quick-start) · [快速开始](../README.zh-CN.md#快速开始) |
| Compose components and use hooks / 组件与 Hooks | [UI API](../pkg/ui/README.md) · [Component example](../example/component) |
| Bind data to controls / 数据绑定 | [Value / Bind](../pkg/value/README.md) |
| Use themed widgets / 主题控件 | [shadcn components](../pkg/shadcn/README.md) · [Showcase](../example/accordion) |
| Use struct declarations / 可选结构体声明 | [Declaration API](../pkg/decl/README.md) |
| Add navigation / 路由导航 | [Router](../pkg/router) · [Example](../example/router) |
| Understand or change the engine / 理解或修改引擎 | [Architecture](../ARCHITECTURE.md) · [Contributing](../CONTRIBUTING.md) |
| Check current scope / 查看当前能力与计划 | [Roadmap](../ROADMAP.md) |

## Current design notes / 当前设计说明

- [API ergonomics](API_ERGONOMICS.md) — component composition, fluent methods, optional struct declarations, and Value / Bind. / 组件组合、链式方法、可选结构体声明与数据协议。
- [Screenshots](screenshots/README.md) — visual reference and capture instructions. / 界面截图与生成方式。

Use the root architecture document and package guides above as the reference for the current implementation.

当前实现以仓库根目录架构文档、各包指南和源码为准。

## Historical material / 历史资料

The numbered notes (`01-research.md` through `10-implementation-roadmap.md`), older topic designs in this directory, and [archive](archive/) record earlier explorations. They may describe the former Flutter/SwiftUI-style API, older rendering backends, or features that were only proposed. Treat them as background, not current usage instructions.

编号文档（`01-research.md` 至 `10-implementation-roadmap.md`）、本目录较早的专题设计与 [archive](archive/) 保留了早期探索，其中可能出现旧的 Flutter/SwiftUI 风格 API、旧渲染后端或尚未落地的方案，阅读时请与当前实现区分。

- [Early architecture](ARCHITECTURE.md)
- [Research and constraints](01-research.md)
- [Original API exploration](04-declarative-api.md)
- [Original implementation plan](10-implementation-roadmap.md)
