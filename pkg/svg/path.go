package svg

import (
	"fmt"
	"math"
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
			rx, ry := args[i]*p.scale, args[i+1]*p.scale
			rot := args[i+2]
			largeArc, sweep := args[i+3] != 0, args[i+4] != 0
			ex, ey := p.apply(args[i+5], args[i+6], isAbs)
			p.arcTo(rx, ry, rot, largeArc, sweep, ex, ey)
		}

	case 'Z', 'z':
		p.path.Close()
		p.px, p.py = p.sx, p.sy
	}

	return nil
}

// arcTo 把 SVG 椭圆弧（endpoint 参数化）转成中心参数化后采样为折线段。
// 处理 rx≠ry 的椭圆与旋转，替代之前的直线近似（否则圆形图标会退化成一条线）。
func (p *PathParser) arcTo(rx, ry, rotDeg float32, largeArc, sweep bool, ex, ey float32) {
	x1, y1 := p.px, p.py
	if rx == 0 || ry == 0 || (x1 == ex && y1 == ey) {
		p.path.LineTo(ex, ey)
		p.px, p.py = ex, ey
		p.updateBounds(ex, ey)
		return
	}
	rxf, ryf := math.Abs(float64(rx)), math.Abs(float64(ry))
	phi := float64(rotDeg) * math.Pi / 180
	cosP, sinP := math.Cos(phi), math.Sin(phi)

	// 端点 -> 中心参数化（SVG 规范 F.6.5）
	dx, dy := float64(x1-ex)/2, float64(y1-ey)/2
	x1p := cosP*dx + sinP*dy
	y1p := -sinP*dx + cosP*dy
	if l := x1p*x1p/(rxf*rxf) + y1p*y1p/(ryf*ryf); l > 1 {
		s := math.Sqrt(l)
		rxf *= s
		ryf *= s
	}
	sign := 1.0
	if largeArc == sweep {
		sign = -1.0
	}
	num := rxf*rxf*ryf*ryf - rxf*rxf*y1p*y1p - ryf*ryf*x1p*x1p
	den := rxf*rxf*y1p*y1p + ryf*ryf*x1p*x1p
	co := 0.0
	if den > 0 && num > 0 {
		co = sign * math.Sqrt(num/den)
	}
	cxp := co * rxf * y1p / ryf
	cyp := -co * ryf * x1p / rxf
	cx := cosP*cxp - sinP*cyp + float64(x1+ex)/2
	cy := sinP*cxp + cosP*cyp + float64(y1+ey)/2

	ang := func(ux, uy, vx, vy float64) float64 {
		d := ux*vx + uy*vy
		l := math.Hypot(ux, uy) * math.Hypot(vx, vy)
		a := math.Acos(math.Max(-1, math.Min(1, d/l)))
		if ux*vy-uy*vx < 0 {
			a = -a
		}
		return a
	}
	t1 := ang(1, 0, (x1p-cxp)/rxf, (y1p-cyp)/ryf)
	dt := ang((x1p-cxp)/rxf, (y1p-cyp)/ryf, (-x1p-cxp)/rxf, (-y1p-cyp)/ryf)
	if !sweep && dt > 0 {
		dt -= 2 * math.Pi
	} else if sweep && dt < 0 {
		dt += 2 * math.Pi
	}
	n := int(math.Ceil(math.Abs(dt) / (math.Pi / 16))) // ~11° 一段
	if n < 1 {
		n = 1
	}
	for k := 1; k <= n; k++ {
		t := t1 + dt*float64(k)/float64(n)
		x := cx + rxf*math.Cos(t)*cosP - ryf*math.Sin(t)*sinP
		y := cy + rxf*math.Cos(t)*sinP + ryf*math.Sin(t)*cosP
		p.path.LineTo(float32(x), float32(y))
		p.updateBounds(float32(x), float32(y))
	}
	p.px, p.py = ex, ey
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
