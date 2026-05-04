package ui

import (
	"fmt"
	"image/color"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// Inspector 是运行时元素树检查器。
// 按 F12 切换显示，鼠标悬停高亮节点，显示类型/bounds/属性。
type Inspector struct {
	visible    bool
	engine     *Engine
	hoverNode  render.RenderObject
	selectedRO render.RenderObject

	// 样式
	highlightColor color.Color
	selectedColor  color.Color
	bgColor        color.Color
	textColor      color.Color
}

// NewInspector 创建元素树检查器。
func NewInspector(engine *Engine) *Inspector {
	return &Inspector{
		visible:        false,
		engine:         engine,
		highlightColor: color.RGBA{R: 255, G: 0, B: 0, A: 100}, // 红色半透明
		selectedColor:  color.RGBA{R: 0, G: 120, B: 255, A: 100}, // 蓝色半透明
		bgColor:        color.RGBA{A: 200},
		textColor:      color.RGBA{R: 0, G: 255, B: 0, A: 255},  // 绿色
	}
}

// Toggle 切换检查器可见性。
func (ins *Inspector) Toggle() {
	ins.visible = !ins.visible
}

// SetVisible 设置可见性。
func (ins *Inspector) SetVisible(v bool) {
	ins.visible = v
}

// IsVisible 返回是否可见。
func (ins *Inspector) IsVisible() bool {
	return ins.visible
}

// Update 更新检查器状态（鼠标悬停检测）。
func (ins *Inspector) Update() {
	if !ins.visible || ins.engine == nil {
		return
	}
	mx, my := ebiten.CursorPosition()
	if ins.engine.rootRenderObject != nil {
		ins.hoverNode = ins.engine.hitTest(ins.engine.rootRenderObject, float32(mx), float32(my))
	}
}

// Paint 绘制检查器覆盖层。
func (ins *Inspector) Paint(screen *ebiten.Image, offset render.Offset) {
	if !ins.visible || ins.engine == nil {
		return
	}

	// 高亮悬停节点
	if ins.hoverNode != nil {
		ins.drawHighlight(screen, ins.hoverNode, ins.highlightColor, offset)
	}

	// 高亮选中节点
	if ins.selectedRO != nil {
		ins.drawHighlight(screen, ins.selectedRO, ins.selectedColor, offset)
	}

	// 绘制信息面板
	ins.drawInfoPanel(screen, offset)
}

// SelectNode 选中一个节点进行详细查看。
func (ins *Inspector) SelectNode(ro render.RenderObject) {
	ins.selectedRO = ro
}

// GetSelected 返回当前选中的节点。
func (ins *Inspector) GetSelected() render.RenderObject {
	return ins.selectedRO
}

// GetHoverNode 返回当前鼠标悬停的节点。
func (ins *Inspector) GetHoverNode() render.RenderObject {
	return ins.hoverNode
}

// DrawTree 绘制元素树的文本表示（用于调试）。
func (ins *Inspector) DrawTree() string {
	if ins.engine == nil || ins.engine.rootRenderObject == nil {
		return "<empty>"
	}
	return ins.renderTreeText(ins.engine.rootRenderObject, 0)
}

func (ins *Inspector) renderTreeText(ro render.RenderObject, depth int) string {
	if ro == nil {
		return ""
	}
	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "  "
	}
	b := ro.GetBounds()
	line := fmt.Sprintf("%s%s [%.0f,%.0f %.0fx%.0f]\n",
		prefix, typeName(ro), b.X, b.Y, b.Width, b.Height)

	for _, child := range ro.GetChildren() {
		line += ins.renderTreeText(child, depth+1)
	}
	return line
}

func (ins *Inspector) drawHighlight(screen *ebiten.Image, ro render.RenderObject, c color.Color, offset render.Offset) {
	b := ro.GetBounds()
	x := b.X + offset.X
	y := b.Y + offset.Y
	w := b.Width
	h := b.Height
	if w <= 0 || h <= 0 {
		return
	}

	// 画四条边（2px 宽）
	r, g, bb, a := c.RGBA()
	clr := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(bb >> 8), A: uint8(a >> 8)}
	drawBorder(screen, x, y, w, h, 2, clr)
}

func (ins *Inspector) drawInfoPanel(screen *ebiten.Image, offset render.Offset) {
	target := ins.selectedRO
	if target == nil {
		target = ins.hoverNode
	}
	if target == nil {
		return
	}

	b := target.GetBounds()
	info := []string{
		fmt.Sprintf("Type: %s", typeName(target)),
		fmt.Sprintf("Bounds: %.0f,%.0f %.0fx%.0f", b.X, b.Y, b.Width, b.Height),
		fmt.Sprintf("Visible: %v", target.IsVisible()),
	}

	// 如果有 transform 信息
	t := target.GetTransform()
	if !t.IsIdentity() {
		info = append(info, fmt.Sprintf("Alpha: %.2f", t.Alpha))
		if t.Rotation != 0 {
			info = append(info, fmt.Sprintf("Rotation: %.1f°", t.Rotation))
		}
	}

	// 绘制背景面板
	panelW := float32(250)
	panelH := float32(len(info)*18 + 16)
	panelX := float32(10)
	panelY := float32(10)
	drawRect(screen, panelX+offset.X, panelY+offset.Y, panelW, panelH, ins.bgColor)

	// 绘制文本
	for i, line := range info {
		drawDebugText(screen, line, panelX+offset.X+8, panelY+offset.Y+float32(i*18)+8, ins.textColor)
	}
}

func typeName(ro render.RenderObject) string {
	if ro == nil {
		return "nil"
	}
	// 使用 reflect 获取类型名
	t := reflect.TypeOf(ro)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func drawBorder(screen *ebiten.Image, x, y, w, h, thickness float32, c color.Color) {
	// 上边
	drawRect(screen, x, y, w, thickness, c)
	// 下边
	drawRect(screen, x, y+h-thickness, w, thickness, c)
	// 左边
	drawRect(screen, x, y, thickness, h, c)
	// 右边
	drawRect(screen, x+w-thickness, y, thickness, h, c)
}
