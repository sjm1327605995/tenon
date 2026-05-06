package render

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle 描述单个粒子状态。
type Particle struct {
	X, Y     float32
	VX, VY   float32 // 速度
	AX, AY   float32 // 加速度
	Life     float32 // 当前生命值（秒）
	MaxLife  float32 // 最大生命值（秒）
	Size     float32 // 半径
	Color    *Color
	Alpha    float32 // 0-1
	Rotation float32 // 弧度
	RotSpeed float32 // 旋转速度
}

// IsAlive 返回粒子是否仍存活。
func (p *Particle) IsAlive() bool { return p.Life > 0 }

// Update 更新粒子状态，dt 为秒。
func (p *Particle) Update(dt float32) {
	p.VX += p.AX * dt
	p.VY += p.AY * dt
	p.X += p.VX * dt
	p.Y += p.VY * dt
	p.Life -= dt
	p.Rotation += p.RotSpeed * dt
	if p.Life < 0 {
		p.Life = 0
	}
	p.Alpha = p.Life / p.MaxLife
	if p.Alpha > 1 {
		p.Alpha = 1
	}
}

// ParticleEmitter 负责生成粒子初始状态。
type ParticleEmitter interface {
	Emit(randSource *rand.Rand) Particle
}

// PointEmitter 从一个点发射。
type PointEmitter struct {
	X, Y float32
}

func (e PointEmitter) Emit(r *rand.Rand) Particle {
	return Particle{X: e.X, Y: e.Y}
}

// RectEmitter 从矩形区域内均匀发射。
type RectEmitter struct {
	X, Y, W, H float32
}

func (e RectEmitter) Emit(r *rand.Rand) Particle {
	return Particle{
		X: e.X + float32(r.Float64())*e.W,
		Y: e.Y + float32(r.Float64())*e.H,
	}
}

// CircleEmitter 从圆形区域内发射。
type CircleEmitter struct {
	CX, CY float32
	Radius float32
}

func (e CircleEmitter) Emit(r *rand.Rand) Particle {
	angle := r.Float64() * 2 * math.Pi
	dist := r.Float64() * float64(e.Radius)
	return Particle{
		X: e.CX + float32(math.Cos(angle)*dist),
		Y: e.CY + float32(math.Sin(angle)*dist),
	}
}

// ParticleBehavior 配置粒子的速度、加速度、颜色、大小等变化。
type ParticleBehavior struct {
	MinSpeed, MaxSpeed       float32
	MinAngle, MaxAngle       float32 // 弧度
	MinSize, MaxSize         float32
	MinLife, MaxLife         float32 // 秒
	AX, AY                   float32 // 全局加速度（重力）
	MinRotSpeed, MaxRotSpeed float32
	Colors                   []*Color // 随机颜色池，nil 则使用默认白色
	SizeOverLife             bool     // 大小随生命周期衰减
}

// Apply 将行为应用到新生成的粒子上。
func (b ParticleBehavior) Apply(p *Particle, r *rand.Rand) {
	speed := b.MinSpeed + float32(r.Float64())*float32(b.MaxSpeed-b.MinSpeed)
	angle := b.MinAngle + float32(r.Float64())*float32(b.MaxAngle-b.MinAngle)
	p.VX = speed * float32(math.Cos(float64(angle)))
	p.VY = speed * float32(math.Sin(float64(angle)))
	p.AX = b.AX
	p.AY = b.AY
	p.Size = b.MinSize + float32(r.Float64())*float32(b.MaxSize-b.MinSize)
	p.MaxLife = b.MinLife + float32(r.Float64())*float32(b.MaxLife-b.MinLife)
	p.Life = p.MaxLife
	p.RotSpeed = b.MinRotSpeed + float32(r.Float64())*float32(b.MaxRotSpeed-b.MinRotSpeed)
	if len(b.Colors) > 0 {
		p.Color = b.Colors[r.Intn(len(b.Colors))]
	} else {
		p.Color = NewColor(255, 255, 255, 255)
	}
	p.Alpha = 1
}

// ParticleSystem 管理粒子数组的更新和绘制。
type ParticleSystem struct {
	particles  []Particle
	emitter    ParticleEmitter
	behavior   ParticleBehavior
	randSource *rand.Rand

	EmissionRate float32 // 每秒发射数量（0 表示只 burst）
	BurstCount   int     // 初始 burst 数量
	emitTimer    float32 // 累积时间

	Loop      bool    // 死亡后是否自动重新 burst
	LoopDelay float32 // 循环间隔（秒）
	loopTimer float32

	Running bool
}

// NewParticleSystem 创建粒子系统。
func NewParticleSystem(emitter ParticleEmitter, behavior ParticleBehavior) *ParticleSystem {
	ps := &ParticleSystem{
		particles:    make([]Particle, 0, 64),
		emitter:      emitter,
		behavior:     behavior,
		randSource:   rand.New(rand.NewSource(time.Now().UnixNano())),
		EmissionRate: 0,
		BurstCount:   30,
		Loop:         false,
		LoopDelay:    0,
		Running:      true,
	}
	return ps
}

// SetEmitter 设置发射器。
func (ps *ParticleSystem) SetEmitter(e ParticleEmitter) { ps.emitter = e }

// SetBehavior 设置粒子行为。
func (ps *ParticleSystem) SetBehavior(b ParticleBehavior) { ps.behavior = b }

// Start 启动（burst 一次）。
func (ps *ParticleSystem) Start() {
	ps.Running = true
	ps.emitBurst()
}

// Stop 停止生成新粒子。
func (ps *ParticleSystem) Stop() {
	ps.Running = false
}

// Reset 清空所有粒子并重新 burst。
func (ps *ParticleSystem) Reset() {
	ps.particles = ps.particles[:0]
	ps.emitBurst()
}

func (ps *ParticleSystem) emitBurst() {
	for i := 0; i < ps.BurstCount; i++ {
		p := ps.emitter.Emit(ps.randSource)
		ps.behavior.Apply(&p, ps.randSource)
		ps.particles = append(ps.particles, p)
	}
}

// Update 更新所有粒子，dt 为秒。
func (ps *ParticleSystem) Update(dt float32) {
	if !ps.Running && len(ps.particles) == 0 {
		return
	}

	if ps.Running && ps.EmissionRate > 0 {
		ps.emitTimer += dt
		interval := 1.0 / ps.EmissionRate
		for ps.emitTimer >= interval {
			ps.emitTimer -= interval
			p := ps.emitter.Emit(ps.randSource)
			ps.behavior.Apply(&p, ps.randSource)
			ps.particles = append(ps.particles, p)
		}
	}

	aliveIdx := 0
	for i := range ps.particles {
		ps.particles[i].Update(dt)
		if ps.particles[i].IsAlive() {
			ps.particles[aliveIdx] = ps.particles[i]
			aliveIdx++
		}
	}
	ps.particles = ps.particles[:aliveIdx]

	if ps.Loop && len(ps.particles) == 0 {
		ps.loopTimer += dt
		if ps.loopTimer >= ps.LoopDelay {
			ps.loopTimer = 0
			ps.emitBurst()
		}
	}
}

// HasAliveParticles 返回是否还有存活粒子。
func (ps *ParticleSystem) HasAliveParticles() bool {
	return len(ps.particles) > 0
}

// Paint 绘制所有粒子到 screen，offset 为 RenderObject 在屏幕上的偏移。
func (ps *ParticleSystem) Paint(screen *ebiten.Image, offset Offset) {
	if len(ps.particles) == 0 {
		return
	}
	for i := range ps.particles {
		p := &ps.particles[i]
		size := p.Size
		if ps.behavior.SizeOverLife {
			size *= p.Alpha
		}
		if size <= 0 {
			continue
		}
		cx := float32(offset.X) + p.X
		cy := float32(offset.Y) + p.Y
		c := p.Color
		if c == nil {
			c = NewColor(255, 255, 255, 255)
		}
		alpha := uint8(float32(c.A) * p.Alpha)
		drawColor := color.RGBA{R: c.R, G: c.G, B: c.B, A: alpha}
		vector.DrawFilledCircle(screen, cx, cy, size, drawColor, true)
	}
}

// ParticleSystemRenderObject 是一个特殊的 RenderObject，仅用于绘制粒子。
type ParticleSystemRenderObject struct {
	BaseRenderObject
	ps *ParticleSystem
}

// NewParticleSystemRenderObject 创建粒子渲染对象。
func NewParticleSystemRenderObject(emitter ParticleEmitter, behavior ParticleBehavior) *ParticleSystemRenderObject {
	ro := &ParticleSystemRenderObject{}
	ro.BaseRenderObject.Init(ro)
	ro.ps = NewParticleSystem(emitter, behavior)
	return ro
}

// Paint 重写以绘制粒子，同时更新粒子状态。
func (r *ParticleSystemRenderObject) Paint(screen *ebiten.Image, offset Offset) {
	if r.ps == nil {
		return
	}
	dt := float32(1.0 / 60.0)
	r.ps.Update(dt)
	r.ps.Paint(screen, offset)
	if r.ps.HasAliveParticles() || r.ps.Running {
		r.MarkNeedsPaint()
	}
}

func (r *ParticleSystemRenderObject) Start() {
	if r.ps != nil {
		r.ps.Start()
		r.MarkNeedsPaint()
	}
}

func (r *ParticleSystemRenderObject) Stop() {
	if r.ps != nil {
		r.ps.Stop()
	}
}

func (r *ParticleSystemRenderObject) SetEmitter(e ParticleEmitter) {
	if r.ps != nil {
		r.ps.SetEmitter(e)
	}
}

func (r *ParticleSystemRenderObject) SetBehavior(b ParticleBehavior) {
	if r.ps != nil {
		r.ps.SetBehavior(b)
	}
}

func (r *ParticleSystemRenderObject) SetLoop(v bool) {
	if r.ps != nil {
		r.ps.Loop = v
	}
}

func (r *ParticleSystemRenderObject) SetLoopDelay(v float32) {
	if r.ps != nil {
		r.ps.LoopDelay = v
	}
}

func (r *ParticleSystemRenderObject) SetBurstCount(n int) {
	if r.ps != nil {
		r.ps.BurstCount = n
	}
}
