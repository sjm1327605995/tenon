package ui

import (
	"math"
	"strconv"
	"strings"

	"gioui.org/f32"
	"gioui.org/op/clip"
)

// ---- gio 后端：SVG path 解析（图标 / 图表矢量）----
//
// 解析 SVG path 的 d 属性并按 scale 缩放，发到一个中立 pathSink（可脱离 gio 单测）。
// 椭圆弧 A 采样成折线；其余命令直译。gioPathSink 把它写进 gio 的 clip.Path。

type pathSink interface {
	moveTo(x, y float32)
	lineTo(x, y float32)
	cubeTo(x1, y1, x2, y2, x, y float32)
	quadTo(x1, y1, x, y float32)
	closePath()
}

// gioPathSink 把解析结果写入 clip.Path，并整体平移 (ox,oy)。
type gioPathSink struct {
	p      *clip.Path
	ox, oy float32
}

func (s *gioPathSink) pt(x, y float32) f32.Point { return f32.Pt(x+s.ox, y+s.oy) }
func (s *gioPathSink) moveTo(x, y float32)       { s.p.MoveTo(s.pt(x, y)) }
func (s *gioPathSink) lineTo(x, y float32)       { s.p.LineTo(s.pt(x, y)) }
func (s *gioPathSink) cubeTo(x1, y1, x2, y2, x, y float32) {
	s.p.CubeTo(s.pt(x1, y1), s.pt(x2, y2), s.pt(x, y))
}
func (s *gioPathSink) quadTo(x1, y1, x, y float32) { s.p.QuadTo(s.pt(x1, y1), s.pt(x, y)) }
func (s *gioPathSink) closePath()                  { s.p.Close() }

// svgParser 解析 d 属性，按 scale 缩放后驱动 sink。
type svgParser struct {
	sink   pathSink
	px, py float32 // 当前点
	sx, sy float32 // 子路径起点
	scale  float32
}

func parseSVGInto(d string, scale float32, sink pathSink) {
	(&svgParser{sink: sink, scale: scale}).parse(d)
}

func svgIsCommand(r rune) bool  { return strings.ContainsRune("MmLlHhVvCcSsQqTtAaZz", r) }
func svgIsCommandB(b byte) bool { return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' }

func (p *svgParser) parse(d string) {
	d = strings.TrimSpace(d)
	if d == "" {
		return
	}
	// 标准化：命令字母、逗号、负号前后补空格，便于分割。
	var sb strings.Builder
	for i, r := range d {
		switch {
		case svgIsCommand(r):
			if i > 0 && d[i-1] != ' ' && d[i-1] != ',' {
				sb.WriteByte(' ')
			}
			sb.WriteRune(r)
			if i+1 < len(d) && d[i+1] != ' ' && d[i+1] != ',' && !svgIsCommand(rune(d[i+1])) {
				sb.WriteByte(' ')
			}
		case r == ',':
			sb.WriteByte(' ')
		case r == '-' && i > 0 && !svgIsCommand(rune(d[i-1])) && d[i-1] != ' ' && d[i-1] != ',':
			sb.WriteByte(' ')
			sb.WriteRune(r)
		default:
			sb.WriteRune(r)
		}
	}
	tokens := strings.Fields(sb.String())
	i := 0
	var lastCmd byte
	for i < len(tokens) {
		tok := tokens[i]
		cmd := lastCmd
		if svgIsCommandB(tok[0]) {
			cmd = tok[0]
			lastCmd = cmd
			i++ // 注意：不能在此处 i>=len 就 break，否则末尾无参命令（如 Z）会漏掉
		}
		start := i
		for i < len(tokens) && !svgIsCommandB(tokens[i][0]) {
			i++
		}
		args := parseFloats(tokens[start:i])
		p.exec(cmd, args)
	}
}

func parseFloats(toks []string) []float32 {
	res := make([]float32, 0, len(toks))
	for _, t := range toks {
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			res = append(res, float32(f))
		}
	}
	return res
}

func (p *svgParser) apply(x, y float32, abs bool) (float32, float32) {
	if abs {
		return x * p.scale, y * p.scale
	}
	return p.px + x*p.scale, p.py + y*p.scale
}

func (p *svgParser) exec(cmd byte, a []float32) {
	abs := cmd >= 'A' && cmd <= 'Z'
	switch cmd {
	case 'M', 'm':
		if len(a) < 2 {
			return
		}
		p.px, p.py = p.apply(a[0], a[1], abs)
		p.sink.moveTo(p.px, p.py)
		p.sx, p.sy = p.px, p.py
		for i := 2; i+1 < len(a); i += 2 { // 后续坐标对视为 L
			p.px, p.py = p.apply(a[i], a[i+1], abs)
			p.sink.lineTo(p.px, p.py)
		}
	case 'L', 'l':
		for i := 0; i+1 < len(a); i += 2 {
			p.px, p.py = p.apply(a[i], a[i+1], abs)
			p.sink.lineTo(p.px, p.py)
		}
	case 'H', 'h':
		for _, x := range a {
			if abs {
				p.px = x * p.scale
			} else {
				p.px += x * p.scale
			}
			p.sink.lineTo(p.px, p.py)
		}
	case 'V', 'v':
		for _, y := range a {
			if abs {
				p.py = y * p.scale
			} else {
				p.py += y * p.scale
			}
			p.sink.lineTo(p.px, p.py)
		}
	case 'C', 'c':
		for i := 0; i+5 < len(a); i += 6 {
			x1, y1 := p.apply(a[i], a[i+1], abs)
			x2, y2 := p.apply(a[i+2], a[i+3], abs)
			x, y := p.apply(a[i+4], a[i+5], abs)
			p.sink.cubeTo(x1, y1, x2, y2, x, y)
			p.px, p.py = x, y
		}
	case 'S', 's':
		for i := 0; i+3 < len(a); i += 4 {
			x2, y2 := p.apply(a[i], a[i+1], abs)
			x, y := p.apply(a[i+2], a[i+3], abs)
			p.sink.cubeTo(p.px, p.py, x2, y2, x, y)
			p.px, p.py = x, y
		}
	case 'Q', 'q':
		for i := 0; i+3 < len(a); i += 4 {
			x1, y1 := p.apply(a[i], a[i+1], abs)
			x, y := p.apply(a[i+2], a[i+3], abs)
			p.sink.quadTo(x1, y1, x, y)
			p.px, p.py = x, y
		}
	case 'T', 't':
		for i := 0; i+1 < len(a); i += 2 {
			x, y := p.apply(a[i], a[i+1], abs)
			p.sink.quadTo(p.px, p.py, x, y)
			p.px, p.py = x, y
		}
	case 'A', 'a':
		for i := 0; i+6 < len(a); i += 7 {
			rx, ry := a[i]*p.scale, a[i+1]*p.scale
			ex, ey := p.apply(a[i+5], a[i+6], abs)
			p.arcTo(rx, ry, a[i+2], a[i+3] != 0, a[i+4] != 0, ex, ey)
		}
	case 'Z', 'z':
		p.sink.closePath()
		p.px, p.py = p.sx, p.sy
	}
}

// arcTo 把 SVG 椭圆弧（endpoint 参数化）转中心参数化后采样成折线（rx≠ry、旋转均支持）。
func (p *svgParser) arcTo(rx, ry, rotDeg float32, largeArc, sweep bool, ex, ey float32) {
	x1, y1 := p.px, p.py
	if rx == 0 || ry == 0 || (x1 == ex && y1 == ey) {
		p.sink.lineTo(ex, ey)
		p.px, p.py = ex, ey
		return
	}
	rxf, ryf := math.Abs(float64(rx)), math.Abs(float64(ry))
	phi := float64(rotDeg) * math.Pi / 180
	cosP, sinP := math.Cos(phi), math.Sin(phi)
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
	n := int(math.Ceil(math.Abs(dt) / (math.Pi / 16)))
	if n < 1 {
		n = 1
	}
	for k := 1; k <= n; k++ {
		t := t1 + dt*float64(k)/float64(n)
		x := cx + rxf*math.Cos(t)*cosP - ryf*math.Sin(t)*sinP
		y := cy + rxf*math.Cos(t)*sinP + ryf*math.Sin(t)*cosP
		p.sink.lineTo(float32(x), float32(y))
	}
	p.px, p.py = ex, ey
}
