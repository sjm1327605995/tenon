package ui

import (
	"reflect"
	"testing"
)

// 窗口设置必须真的抵达后端。标题曾被写死在 Run 里（"Tenon UI"），
// 任何只检查 winCfg 字段的测试都抓不到那种 bug —— 得看后端收到了什么。
func TestWindowOptionsReachBackend(t *testing.T) {
	old, oldCfg := backendRun, winCfg
	t.Cleanup(func() { backendRun, winCfg = old, oldCfg })

	var got windowConfig
	backendRun = func(root *Node, cfg windowConfig) { got = cfg }

	winCfg = windowConfig{w: 800, h: 600, title: "Tenon UI"} // 默认
	WindowSize(1280, 720)
	WindowTitle("我的应用")
	WindowMinSize(640, 480)
	WindowMaxSize(1920, 1080)
	WindowFullscreen(true)
	Run(Text("x"))

	want := windowConfig{
		w: 1280, h: 720, title: "我的应用",
		minW: 640, minH: 480, maxW: 1920, maxH: 1080,
		fullscreen: true, sync: FrameSync,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("后端收到 %+v\nwant %+v", got, want)
	}
}

// 默认值：不设任何东西时后端该收到什么。
func TestWindowDefaults(t *testing.T) {
	old, oldCfg := backendRun, winCfg
	t.Cleanup(func() { backendRun, winCfg = old, oldCfg })

	var got windowConfig
	backendRun = func(root *Node, cfg windowConfig) { got = cfg }
	winCfg = windowConfig{w: 800, h: 600, title: "Tenon UI"}
	Run(Text("x"))

	if got.w != 800 || got.h != 600 {
		t.Errorf("默认尺寸 %dx%d，want 800x600", got.w, got.h)
	}
	if got.title != "Tenon UI" {
		t.Errorf("默认标题 %q，want %q", got.title, "Tenon UI")
	}
	if got.minW != 0 || got.maxW != 0 || got.fullscreen {
		t.Errorf("默认不该有 min/max/全屏限制，got %+v", got)
	}
}

// WindowSize 拒绝非正尺寸（沿用既有行为）。
func TestWindowSizeRejectsNonPositive(t *testing.T) {
	oldCfg := winCfg
	t.Cleanup(func() { winCfg = oldCfg })
	winCfg = windowConfig{w: 800, h: 600}
	WindowSize(0, -5)
	if winCfg.w != 800 || winCfg.h != 600 {
		t.Errorf("非正尺寸把配置改成了 %dx%d，应当忽略", winCfg.w, winCfg.h)
	}
}

// gio 侧的映射：app.Size 会把窗口模式重置为 Windowed，所以全屏必须排在它之后，
// 否则设了全屏也不生效。选项是不透明的 func，只能校验数量与顺序约束。
func TestGioWindowOptionCount(t *testing.T) {
	for _, tc := range []struct {
		name string
		cfg  windowConfig
		want int
	}{
		{"仅标题+尺寸", windowConfig{w: 800, h: 600, title: "t"}, 2},
		{"加最小尺寸", windowConfig{w: 800, h: 600, minW: 640, minH: 480}, 3},
		{"加最大尺寸", windowConfig{w: 800, h: 600, maxW: 1920, maxH: 1080}, 3},
		{"全部", windowConfig{w: 800, h: 600, minW: 1, minH: 1, maxW: 2, maxH: 2, fullscreen: true}, 5},
	} {
		if n := len(gioWindowOptions(tc.cfg)); n != tc.want {
			t.Errorf("%s: 生成 %d 个选项，want %d", tc.name, n, tc.want)
		}
	}
}
