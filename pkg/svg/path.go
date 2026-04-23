package svg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PathParser 解析 SVG path 的 d 属性。
type PathParser struct {
	path        *vector.Path
	px, py      float32 // 当前点
	sx, sy      float32 // 子路径起点（用于 Z 命令）
	scale       float32 // 缩放比例
	shiftX      float32 // X 平移
	shiftY      float32 // Y 平移
	trackBounds bool
	minX, minY  float32
	maxX, maxY  float32
}

// ParsePath 将 SVG path 数据解析为 Ebiten vector.Path。
func ParsePath(d string) (*vector.Path, error) {
	return ParsePathScaled(d, 1)
}

// ParsePathScaled 解析 SVG path 并按指定比例缩放。
func ParsePathScaled(d string, scale float32) (*vector.Path, error) {
	return ParsePathScaledAndShifted(d, scale, 0, 0)
}

// ParsePathScaledAndShifted 解析 SVG path，缩放并平移。
func ParsePathScaledAndShifted(d string, scale, shiftX, shiftY float32) (*vector.Path, error) {
	p := &PathParser{path: &vector.Path{}, scale: scale, shiftX: shiftX, shiftY: shiftY}
	if err := p.parse(d); err != nil {
		return nil, err
	}
	return p.path, nil
}

// ParsePathBounds 解析 SVG path 并返回其 bounding box。
func ParsePathBounds(d string) (minX, minY, maxX, maxY float32, err error) {
	p := &PathParser{path: &vector.Path{}, scale: 1, trackBounds: true, minX: 1e9, minY: 1e9, maxX: -1e9, maxY: -1e9}
	if err = p.parse(d); err != nil {
		return 0, 0, 0, 0, err
	}
	return p.minX, p.minY, p.maxX, p.maxY, nil
}

func (p *PathParser) parse(d string) error {
	d = strings.TrimSpace(d)
	if d == "" {
		return nil
	}

	// 标准化：在命令字母前后加空格，方便分割
	var sb strings.Builder
	for i, r := range d {
		if isCommand(r) {
			if i > 0 && d[i-1] != ' ' && d[i-1] != ',' {
				sb.WriteByte(' ')
			}
			sb.WriteRune(r)
			if i+1 < len(d) && d[i+1] != ' ' && d[i+1] != ',' && !isCommand(rune(d[i+1])) {
				sb.WriteByte(' ')
			}
		} else if r == ',' {
			sb.WriteByte(' ')
		} else if r == '-' && i > 0 && !isCommand(rune(d[i-1])) && d[i-1] != ' ' && d[i-1] != ',' {
			// 负数前面加空格
			sb.WriteByte(' ')
			sb.WriteRune(r)
		} else {
			sb.WriteRune(r)
		}
	}

	tokens := strings.Fields(sb.String())
	i := 0
	var lastCmd byte

	for i < len(tokens) {
		tok := tokens[i]
		cmd := lastCmd
		if isCommandByte(tok[0]) {
			cmd = tok[0]
			lastCmd = cmd
			i++
			if i >= len(tokens) {
				break
			}
		}

		// 读取后续的数字参数
		start := i
		for i < len(tokens) && !isCommandByte(tokens[i][0]) {
			i++
		}
		args, err := parseFloats(tokens[start:i])
		if err != nil {
			return fmt.Errorf("svg path parse error at cmd %c: %w", cmd, err)
		}

		if err := p.executeCmd(cmd, args); err != nil {
			return err
		}
	}

	return nil
}

func isCommand(r rune) bool {
	return strings.ContainsRune("MmLlHhVvCcSsQqTtAaZz", r)
}

func isCommandByte(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func parseFloats(tokens []string) ([]float32, error) {
	res := make([]float32, len(tokens))
	for i, t := range tokens {
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return nil, err
		}
		res[i] = float32(f)
	}
	return res, nil
}

func (p *PathParser) executeCmd(cmd byte, args []float32) error {
	isAbs := cmd >= 'A' && cmd <= 'Z'

	switch cmd {
	case 'M', 'm':
		if len(args) < 2 {
			return fmt.Errorf("M/m requires 2 args, got %d", len(args))
		}
		p.px, p.py = p.apply(args[0], args[1], isAbs)
		p.path.MoveTo(p.px, p.py)
		p.updateBounds(p.px, p.py)
		p.sx, p.sy = p.px, p.py
		// 后续坐标对视为 L/l
		for i := 2; i+1 < len(args); i += 2 {
			p.px, p.py = p.apply(args[i], args[i+1], isAbs)
			p.path.LineTo(p.px, p.py)
			p.updateBounds(p.px, p.py)
		}

	case 'L', 'l':
		for i := 0; i+1 < len(args); i += 2 {
			p.px, p.py = p.apply(args[i], args[i+1], isAbs)
			p.path.LineTo(p.px, p.py)
			p.updateBounds(p.px, p.py)
		}

	case 'H', 'h':
		for _, x := range args {
			if isAbs {
				p.px = x*p.scale + p.shiftX
			} else {
				p.px += x * p.scale
			}
			p.path.LineTo(p.px, p.py)
			p.updateBounds(p.px, p.py)
		}

	case 'V', 'v':
		for _, y := range args {
			if isAbs {
				p.py = y*p.scale + p.shiftY
			} else {
				p.py += y * p.scale
			}
			p.path.LineTo(p.px, p.py)
			p.updateBounds(p.px, p.py)
		}

	case 'C', 'c':
		for i := 0; i+5 < len(args); i += 6 {
			x1, y1 := p.apply(args[i], args[i+1], isAbs)
			x2, y2 := p.apply(args[i+2], args[i+3], isAbs)
			x, y := p.apply(args[i+4], args[i+5], isAbs)
			p.path.CubicTo(x1, y1, x2, y2, x, y)
			p.px, p.py = x, y
			p.updateBounds(x1, y1)
			p.updateBounds(x2, y2)
			p.updateBounds(x, y)
		}

	case 'S', 's':
		for i := 0; i+3 < len(args); i += 4 {
			x2, y2 := p.apply(args[i], args[i+1], isAbs)
			x, y := p.apply(args[i+2], args[i+3], isAbs)
			p.path.CubicTo(p.px, p.py, x2, y2, x, y)
			p.px, p.py = x, y
			p.updateBounds(x2, y2)
			p.updateBounds(x, y)
		}

	case 'Q', 'q':
		for i := 0; i+3 < len(args); i += 4 {
			x1, y1 := p.apply(args[i], args[i+1], isAbs)
			x, y := p.apply(args[i+2], args[i+3], isAbs)
			p.path.QuadTo(x1, y1, x, y)
			p.px, p.py = x, y
			p.updateBounds(x1, y1)
			p.updateBounds(x, y)
		}

	case 'T', 't':
		for i := 0; i+1 < len(args); i += 2 {
			x, y := p.apply(args[i], args[i+1], isAbs)
			p.path.QuadTo(p.px, p.py, x, y)
			p.px, p.py = x, y
			p.updateBounds(x, y)
		}

	case 'A', 'a':
		for i := 0; i+6 < len(args); i += 7 {
			_ = args[i+2]
			_ = args[i+3] != 0
			_ = args[i+4] != 0
			x, y := p.apply(args[i+5], args[i+6], isAbs)
			// Ebiten vector.Path.Arc 只支持圆弧，不支持椭圆弧。
			// 这里用 LineTo 近似（大多数图标场景下 A 命令很少见）
			p.path.LineTo(x, y)
			p.px, p.py = x, y
			p.updateBounds(x, y)
		}

	case 'Z', 'z':
		p.path.Close()
		p.px, p.py = p.sx, p.sy
	}

	return nil
}

func (p *PathParser) apply(x, y float32, isAbs bool) (float32, float32) {
	if isAbs {
		return x*p.scale + p.shiftX, y*p.scale + p.shiftY
	}
	return p.px + x*p.scale, p.py + y*p.scale
}

func (p *PathParser) updateBounds(x, y float32) {
	if !p.trackBounds {
		return
	}
	if x < p.minX {
		p.minX = x
	}
	if y < p.minY {
		p.minY = y
	}
	if x > p.maxX {
		p.maxX = x
	}
	if y > p.maxY {
		p.maxY = y
	}
}
