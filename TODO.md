# Tenon 开发路线图

> 本文档记录 Tenon UI 框架的已完成功能和未来开发计划。

---

## ✅ 已完成功能

### 核心框架
- [x] 组件系统（Widget / Host 双模型）
- [x] 生命周期钩子（Mount / Update / Unmount）
- [x] Hooks 系统（UseState / UseEffect / UseMemo / UseRef / UseCallback / UseId / UseContext / UseTransition）
- [x] Yoga 布局引擎集成（Flexbox）
- [x] Ebiten 渲染后端
- [x] 事件系统（Click / MouseDown / MouseUp / Scroll / MouseMove / KeyDown / KeyUp / FocusIn / FocusOut）
- [x] 焦点系统（Tab 切换、Space/Enter 触发、键盘分发）
- [x] 拖拽支持（MouseDown + MouseMove + MouseUp）
- [x] Host 复用机制（保留 Yoga 节点避免布局跳动）
- [x] Overlay / Portal 浮层

### 内置组件
- [x] View — 容器（背景、边框、阴影、圆角、裁剪）
- [x] Text — 文本（多行换行、white-space / word-break 策略）
- [x] Button — 按钮（hover / pressed 状态、焦点）
- [x] Image — 图片
- [x] ScrollView — 滚动视图（滚轮、拖拽、滚动条）
- [x] ProgressBar — 进度条
- [x] Checkbox — 复选框
- [x] Slider — 滑块（拖拽/点击）
- [x] Switch — 开关
- [x] Radio — 单选框
- [x] Divider — 分割线
- [x] TextInput — 文本输入框（IME、光标、选区、合成文本）

### 基础设施
- [x] 字体管理器（多字体族、多字重、缓存）
- [x] 文本布局引擎（CJK 换行、多种 white-space / word-break 策略）
- [x] 抗锯齿圆形绘制（vector.Path + Arc）

### 文档与测试
- [x] 中英双语 README
- [x] CONTRIBUTING.md
- [x] 核心引擎单元测试
- [x] 文本布局单元测试

---

## 📋 待办事项（TODO）

### P0 — 核心必备（缺少则框架不完整）

#### 1. 动画系统
**优先级**: 🔴 P0  
**状态**: ⏳ 待开始  
**描述**: 现代 UI 离不开动画过渡。当前所有交互都是瞬间跳变，用户体验生硬。

**实现要点**:
- 设计 `Animation` / `Tween` 接口
- 支持缓动函数：linear、ease-in-out、ease-out、bounce、spring
- 在 Engine 的 `Update` 循环中统一更新动画状态
- 支持属性动画：位置(x,y)、尺寸(width,height)、透明度、颜色过渡
- ScrollView 滚动支持惯性动画（fling）
- Switch 滑块切换支持平滑位移动画

**参考 API**:
```go
anim := components.NewTween(duration, ease.EaseInOutQuad).
    OnUpdate(func(progress float32) {
        view.SetWidth(startWidth + (endWidth-startWidth)*progress)
    }).
    OnComplete(func() { ... })
```

---

#### 2. 组件层单元测试覆盖
**优先级**: 🔴 P0  
**状态**: ⏳ 待开始  
**描述**: 目前只有 `pkg/core` 和 `yoga` 有测试，`pkg/components` 零测试。

**测试范围**:
- [ ] Button — 点击事件、hover 状态切换、焦点状态
- [ ] Text — 换行布局（多种 white-space + word-break 组合）
- [ ] ScrollView — 滚动偏移计算、边界限制、最大滚动范围
- [ ] TextInput — IME 状态机、光标位置计算、选区逻辑
- [ ] Checkbox / Radio / Switch — 状态切换、事件回调
- [ ] Slider — 拖拽值计算、边界限制
- [ ] 链式 API 返回值检查（确保所有 SetXxx 返回自身支持链式调用）

---

### P1 — 重要提升（显著改善开发体验）

#### 3. Modal / Dialog 模态弹窗系统
**优先级**: 🟡 P1  
**状态**: ⏳ 待开始  
**描述**: 基于已有的 `MountOverlay` / `UnmountOverlay` 构建模态对话框。

**实现要点**:
- 全局遮罩层（半透明黑色背景 `rgba(0,0,0,0.5)`）
- 点击遮罩关闭
- ESC 键关闭
- 支持自定义内容（传入任意 Component）
- 动画：淡入 + 缩放弹出
- 阻止事件穿透到下层组件

**参考 API**:
```go
modal := components.NewModal().
    SetContent(myDialogView).
    SetClickOutsideToClose(true).
    SetOnClose(func() { ... })
modal.Show()
```

---

#### 4. Tooltip 气泡提示
**优先级**: 🟡 P1  
**状态**: ⏳ 待开始  
**描述**: 鼠标悬停延迟显示提示文本。

**实现要点**:
- 悬停延迟 500ms 显示
- 自动计算位置（上下左右自适应边界）
- 使用 Overlay 层挂载，避免被父组件裁剪
- 小三角箭头指向目标组件
- 支持自定义内容（不限于文本）

---

#### 5. Select / Dropdown 下拉选择
**优先级**: 🟡 P1  
**状态**: ⏳ 待开始  
**描述**: 下拉列表选择组件。

**实现要点**:
- 点击展开下拉列表（基于 Overlay）
- 支持单选、选项高亮、滚动
- 键盘导航（↑↓ 选择、Enter 确认、ESC 关闭）
- 支持禁用某一项
- 支持分组（optgroup）

---

#### 6. 主题系统
**优先级**: 🟡 P1  
**状态**: ⏳ 待开始  
**描述**: 全局配色/字体/圆角配置，支持一键换肤。

**实现要点**:
- `Theme` 全局配置对象：
  - `PrimaryColor` — 主题色（按钮、滑块、开关）
  - `BorderColor` / `FocusBorderColor`
  - `BorderRadius` — 全局圆角
  - `FontSize` 层级（sm / base / lg / xl）
  - `BackgroundColor` / `TextColor`
- 所有组件默认从 Theme 读取配色，支持单独覆盖
- 支持明暗模式切换
- Demo 中展示一键换肤

---

### P2 — 锦上添花（高级功能与性能）

#### 7. Tabs 标签页组件
**优先级**: 🟢 P2  
**状态**: ⏳ 待开始  

**实现要点**:
- 顶部/底部标签栏切换
- 内容区随标签切换重新 Render
- 支持禁用某项
- 支持徽标小红点
- 切换动画（下划线滑动或内容淡入淡出）

---

#### 8. Loading / Spinner 加载动画
**优先级**: 🟢 P2  
**状态**: ⏳ 待开始  

**实现要点**:
- 旋转圆圈动画（基于角度属性的 Tween）
- 骨架屏 Skeleton（灰色闪烁块，支持矩形/圆形/文本行）
- 支持全屏遮罩加载和局部加载

---

#### 9. 虚拟列表 VirtualList
**优先级**: 🟢 P2  
**状态**: ⏳ 待开始  
**描述**: 只渲染可视区域内的列表项，适用于大数据量场景。

**实现要点**:
- 计算可视范围（基于 ScrollView 的 scrollOffset 和 viewport 高度）
- 预估每项高度（或支持动态高度）
- 只渲染 visibleStart ~ visibleEnd 范围内的项
- 配合 ScrollView 的滚动偏移更新可见范围
- 适用于聊天记录、日志列表、长表格

---

#### 10. 手势系统
**优先级**: 🟢 P2  
**状态**: ⏳ 待开始  

**实现要点**:
- 双指缩放（Pinch）—— 用于 Image 查看器
- 惯性滚动（Fling）—— ScrollView 松手后继续滑动
- 边缘回弹（Bounce）—— 滚动到边界时的弹性效果
- 需要扩展 Engine 事件系统支持多点触控

---

#### 11. 性能优化
**优先级**: 🟢 P2  
**状态**: ⏳ 待开始  

**实现要点**:
- **脏矩形渲染**：只重绘发生变化的区域，而非全屏
- **绘制缓存**：静态内容（如背景、边框）缓存到 offscreen Image
- **Yoga 布局缓存**：避免每帧重复计算未变化的节点布局
- **组件级缓存**：给 View 添加 `cacheAsBitmap` 选项

---

## 🗓️ 建议实施顺序

```
第一阶段（P0）
  ├── 动画系统基础（Tween + 缓动函数）
  ├── 组件层单元测试覆盖
  └── Switch 滑块平滑动画（验证动画系统）

第二阶段（P1）
  ├── Modal / Dialog 模态弹窗
  ├── Tooltip 气泡提示
  ├── Select / Dropdown 下拉选择
  └── 主题系统 + Demo 一键换肤

第三阶段（P2）
  ├── Tabs 标签页
  ├── Loading / Spinner
  ├── 虚拟列表 VirtualList
  ├── 手势系统（惯性滚动、双指缩放）
  └── 性能优化（脏矩形、绘制缓存）
```

---

## 📝 备注

- **P0** 项完成后，Tenon 将具备完整的 GUI 开发能力，可用于实际项目。
- **P1** 项完成后，开发体验将接近现代 Web 前端框架。
- **P2** 项完成后，性能和高级交互将接近原生桌面应用。
- 如需调整优先级或添加新需求，请直接修改本文件。
