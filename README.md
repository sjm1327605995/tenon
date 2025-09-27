# Goui

A lightweight UI library built on top of [Ebiten](https://github.com/hajimehoshi/ebiten), designed for creating simple and efficient user interfaces in Go.

## Features
- **Lightweight**: Minimal dependencies and fast rendering.
- **Flexible**: Supports custom widgets and layouts.
- **Cross-platform**: Works on Windows, macOS, and Linux.

## Installation
```bash
go get github.com/m/ebitenimgui
```

## Quick Start
```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m/ebitenimgui"
)

func main() {
	ui := ebitenimgui.NewUI()
	ebiten.RunGame(ui)
}
```

## Documentation
For detailed documentation, visit the [GitHub repository](https://github.com/m/ebitenimgui).

---

# EbitenImgUI

基于 [Ebiten](https://github.com/hajimehoshi/ebiten) 构建的轻量级 UI 库，用于在 Go 中创建简单高效的用户界面。

## 特性
- **轻量级**：依赖少，渲染快。
- **灵活**：支持自定义控件和布局。
- **跨平台**：支持 Windows、macOS 和 Linux。

## 安装
```bash
go get github.com/m/ebitenimgui
```

## 快速开始
```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m/ebitenimgui"
)

func main() {
	ui := ebitenimgui.NewUI()
	ebiten.RunGame(ui)
}
```

## 文档
详细文档请访问 [GitHub 仓库](https://github.com/m/ebitenimgui)。