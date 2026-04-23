# Tenon Development Roadmap

## Completed

### Core Framework
- [x] Component system (Widget / Host dual model)
- [x] Lifecycle hooks (Mount / Update / Unmount)
- [x] Hooks system (UseState / UseEffect / UseMemo / UseRef / UseCallback / UseId / UseContext / UseTransition)
- [x] Yoga layout engine integration (Flexbox)
- [x] Ebiten rendering backend
- [x] Event system (Click / MouseDown / MouseUp / Scroll / MouseMove / KeyDown / KeyUp / FocusIn / FocusOut)
- [x] Focus system (Tab switching, Space/Enter trigger, keyboard dispatch)
- [x] Drag support (MouseDown + MouseMove + MouseUp)
- [x] Host reuse mechanism (preserve Yoga nodes to avoid layout jumps)
- [x] Overlay / Portal floating layer
- [x] Animation system (Tween + easing functions)
- [x] Theme system (global colors/fonts/radius, light/dark mode)

### Built-in Host Components
- [x] View - container (background, border, shadow, radius, clip)
- [x] Text - text (multi-line wrapping, white-space / word-break strategies)
- [x] Button - button (hover / pressed states, focus)
- [x] Image - image
- [x] ScrollView - scroll view (wheel, drag, scrollbar)
- [x] ProgressBar - progress bar
- [x] Checkbox - checkbox
- [x] Slider - slider (drag/click)
- [x] Switch - switch
- [x] Radio - radio button
- [x] Divider - divider
- [x] TextInput - text input (IME, cursor, selection, composition text)
- [x] Menu - menu

### Ant Design Components (Implemented)
- [x] Button - button (primary/default/dashed/text/link, danger, size, loading, disabled)
- [x] Alert - alert (error/warning/success/info, closable, banner, icon, description)
- [x] Badge - badge (count, dot, status, color)
- [x] Card - card (title, extra, shadow, hoverable, children)
- [x] Divider - divider (horizontal/vertical, with text)
- [x] Input - input (placeholder, prefix/suffix, search, password, onChange, onSubmit)
- [x] Table - table (columns, dataSource, header, stripe, hover)
- [x] Tag - tag (color, closable, icon)

### Infrastructure
- [x] Font manager (multi-font family, multi-weight, caching)
- [x] Text layout engine (CJK line breaking, multiple white-space / word-break strategies)
- [x] Anti-aliased circle drawing (vector.Path + Arc)
- [x] Chinese and English bilingual README
- [x] CONTRIBUTING.md
- [x] Core engine unit tests
- [x] Text layout unit tests
- [x] Component layer unit test coverage (Button / Text / ScrollView / TextInput / Checkbox / Radio / Switch / Slider)

---

## Ant Design Components Implementation TODO

Based on Ant Design 5.x official component overview
Implementation strategy: Prefer Widget composition, extend native Host only when necessary.
Complex components should support composition (e.g. Card + Table, Form + Input).

---

## P0 - Core Essentials (High Frequency)

### 1. Typography
Status: pending | Strategy: Widget composition (Text + View)

- [ ] AntTitle - headings h1-h4 (auto fontSize / fontWeight)
- [ ] AntText - text variants (primary, secondary, success, warning, danger, mark, code, keyboard, underline, delete, strong, italic)
- [ ] AntParagraph - paragraph (ellipsis, copyable, editable)
- [ ] AntLink - link (hover underline)

### 2. Space
Status: pending | Strategy: Widget composition (View + Yoga flex gap)

- [ ] AntSpace - spacing container (direction, size, align)
- [ ] AntSpace.Compact - compact mode for form control dense arrangement

### 3. Grid
Status: pending | Strategy: Widget composition (View + Yoga flex)

- [ ] AntRow - row (gutter, justify, align, wrap)
- [ ] AntCol - column (span, offset, push, pull, order, responsive breakpoints)

### 4. Layout
Status: pending | Strategy: Widget composition (View)

- [ ] AntLayout - layout container
- [ ] AntHeader - top navigation
- [ ] AntSider - sidebar (collapsible, collapsedWidth, trigger)
- [ ] AntContent - content area
- [ ] AntFooter - footer

### 5. Flex
Status: pending | Strategy: Widget composition (View + Yoga flex)

- [ ] AntFlex - flex layout (vertical, justify, align, gap, wrap)

### 6. Avatar
Status: pending | Strategy: Widget composition (Image / Text + View circular clip)

- [ ] AntAvatar - avatar (size, shape: circle/square, src, icon, alt, gap)
- [ ] AntAvatar.Group - avatar group (maxCount, maxStyle, size)

### 7. Empty
Status: pending | Strategy: Widget composition (View + Text / Image)

- [ ] AntEmpty - empty state (image, description, children custom bottom)

### 8. Tooltip
Status: pending | Strategy: Widget + Overlay system

- [ ] AntTooltip - tooltip (title, placement, color, onOpenChange)
- [ ] Auto position calculation (boundary detection)
- [ ] Small arrow pointing to target

### 9. Popover
Status: pending | Strategy: Widget + Overlay system (extends Tooltip)

- [ ] AntPopover - popover (title, content, placement, trigger)

### 10. Popconfirm
Status: pending | Strategy: Widget + Overlay system (extends Popover)

- [ ] AntPopconfirm - confirmation (title, description, onConfirm, onCancel, okText, cancelText, placement)

### 11. Modal
Status: pending | Strategy: Widget + Overlay system

- [ ] AntModal - modal (title, open, onOk, onCancel, footer, width, centered)
- [ ] Animation: mask fade + content scale pop-in
- [ ] Support Modal.info / warning / success / error / confirm shortcuts

### 12. Drawer
Status: pending | Strategy: Widget + Overlay system

- [ ] AntDrawer - drawer (title, placement, open, onClose, width/height)
- [ ] Animation: slide in/out

### 13. Dropdown
Status: pending | Strategy: Widget + Overlay system

- [ ] AntDropdown - dropdown (menu items, placement, trigger, disabled, arrow)
- [ ] AntDropdown.Button - button with dropdown menu

### 14. Select
Status: pending | Strategy: Widget + Overlay + ScrollView

- [ ] AntSelect - selector (options, value, placeholder, disabled, allowClear, loading)
- [ ] Mode: multiple, tags
- [ ] Dropdown panel via Overlay mount
- [ ] Keyboard navigation (up/down, enter, esc)
- [ ] Multi-select tag display
- [ ] Search filtering

### 15. Checkbox (AntD Wrapper)
Status: pending | Strategy: Widget composition (native Checkbox + Text + View)

- [ ] AntCheckbox - checkbox (checked, disabled, indeterminate, onChange, children label)
- [ ] AntCheckbox.Group - checkbox group (options, value, onChange)

### 16. Radio (AntD Wrapper)
Status: pending | Strategy: Widget composition (native Radio + Text + View)

- [ ] AntRadio - radio (checked, disabled, onChange, children label)
- [ ] AntRadio.Group - radio group (options, value, optionType: default/button, size)
- [ ] AntRadio.Button - button-style radio

### 17. Switch (AntD Wrapper)
Status: pending | Strategy: Widget composition (native Switch + Text)

- [ ] AntSwitch - switch (checked, disabled, loading, size, checkedChildren, unCheckedChildren)

### 18. Slider (AntD Wrapper)
Status: pending | Strategy: Widget composition (native Slider + Tooltip)

- [ ] AntSlider - slider (min, max, step, value, range, vertical, tooltip)

### 19. Progress (AntD Wrapper)
Status: pending | Strategy: Widget composition (native ProgressBar + Text + View)

- [ ] AntProgress - progress (percent, status, showInfo, type: line/circle/dashboard, size)
- [ ] Extend native ProgressBar or draw circle/dashboard in Widget

### 20. Menu (AntD Wrapper)
Status: pending | Strategy: Widget composition (native Menu + View + Text)

- [ ] AntMenu - menu (items, mode: vertical/horizontal/inline, theme, selectedKeys)
- [ ] AntMenu.Item, AntMenu.SubMenu, AntMenu.ItemGroup

### 21. Tabs
Status: pending | Strategy: Widget composition (View + Text + Button)

- [ ] AntTabs - tabs (items, activeKey, type: line/card, size, onChange)
- [ ] AntTabs.TabPane - tab pane content
- [ ] Top/bottom/left/right tab bar placement
- [ ] Tab switching animation

### 22. Breadcrumb
Status: pending | Strategy: Widget composition (Text + View + separator)

- [ ] AntBreadcrumb - breadcrumb (items, separator)
- [ ] AntBreadcrumb.Item - breadcrumb item
- [ ] AntBreadcrumb.Separator - custom separator

### 23. Pagination
Status: pending | Strategy: Widget composition (Button + Text + View)

- [ ] AntPagination - pagination (current, total, pageSize, showSizeChanger, showQuickJumper)

### 24. Steps
Status: pending | Strategy: Widget composition (View + Text)

- [ ] AntSteps - steps (current, direction, size, status, items)
- [ ] AntSteps.Step - step item (title, description, icon)

### 25. Form
Status: pending | Strategy: Widget composition + state management

- [ ] AntForm - form (initialValues, onFinish, layout, labelCol, wrapperCol)
- [ ] AntForm.Item - form field (name, label, rules, validateStatus, help, required)
- [ ] Validation rules (required, min, max, len, pattern, validator)
- [ ] Error display and form submission

### 26. InputNumber
Status: pending | Strategy: Extend native TextInput or Widget composition

- [ ] AntInputNumber - number input (min, max, step, precision, formatter, parser)
- [ ] Step buttons (+/-) and keyboard up/down

### 27. Rate
Status: pending | Strategy: Widget composition (Text/Button stars)

- [ ] AntRate - rate (value, count, allowHalf, allowClear, character)
- [ ] Half star support, hover preview

### 28. Segmented
Status: pending | Strategy: Widget composition (View + Button/Text)

- [ ] AntSegmented - segmented control (options, value, disabled, size, block)

### 29. Descriptions
Status: pending | Strategy: Widget composition (View + Text + Grid)

- [ ] AntDescriptions - descriptions (title, bordered, column, layout)
- [ ] AntDescriptions.Item - description item (label, span)

### 30. Statistic
Status: pending | Strategy: Widget composition (Text + View)

- [ ] AntStatistic - statistic (title, value, precision, prefix, suffix)
- [ ] AntStatistic.Countdown - countdown (value, format, onFinish)

### 31. Timeline
Status: pending | Strategy: Widget composition (View + Text)

- [ ] AntTimeline - timeline (mode, pending, reverse)
- [ ] AntTimeline.Item - timeline item (color, dot, label)

### 32. List
Status: pending | Strategy: Widget composition (ScrollView + View + Text)

- [ ] AntList - list (dataSource, renderItem, bordered, split, loading)
- [ ] AntList.Item - list item (actions, extra)
- [ ] AntList.Item.Meta - list meta (avatar, title, description)

### 33. Collapse
Status: pending | Strategy: Widget composition (View + Text + Button)

- [ ] AntCollapse - collapse (items, activeKey, accordion, bordered, ghost)
- [ ] AntCollapse.Panel - panel (header, key, extra, disabled)


## P1 - Important Components (Medium Frequency)

### 34. Anchor
Status: pending | Strategy: Widget composition

- [ ] AntAnchor - anchor (affix, bounds, onChange)
- [ ] AntAnchor.Link - anchor link (href, title)

### 35. Affix
Status: pending | Strategy: Widget composition

- [ ] AntAffix - affix (offsetTop, offsetBottom, onChange)

### 36. BackTop
Status: pending | Strategy: Widget + Overlay

- [ ] AntBackTop - back to top (visibilityHeight, onClick)

### 37. Carousel
Status: pending | Strategy: Widget composition

- [ ] AntCarousel - carousel (autoplay, dotPosition, effect)
- [ ] Slide switching, indicator dots

### 38. Image (AntD Wrapper)
Status: pending | Strategy: Widget composition + Overlay

- [ ] AntImage - image (src, alt, width, height, fallback, placeholder, preview)
- [ ] Preview mode (overlay zoom, drag, close)

### 39. Calendar
Status: pending | Strategy: Widget composition

- [ ] AntCalendar - calendar (value, mode, fullscreen, dateCellRender, monthCellRender)
- [ ] Month/year view, date selection

### 40. Comment
Status: pending | Strategy: Widget composition

- [ ] AntComment - comment (author, avatar, content, datetime, actions)

### 41. Result
Status: pending | Strategy: Widget composition

- [ ] AntResult - result (status, title, subTitle, extra, icon)

### 42. Skeleton
Status: pending | Strategy: Widget composition (animated placeholder)

- [ ] AntSkeleton - skeleton (active, avatar, paragraph, title, loading)
- [ ] AntSkeleton.Button, AntSkeleton.Input, AntSkeleton.Avatar, AntSkeleton.Image

### 43. Spin
Status: pending | Strategy: Widget composition

- [ ] AntSpin - spinner (spinning, size, tip, delay, indicator)
- [ ] Rotating loader animation

### 44. Transfer
Status: pending | Strategy: Widget composition

- [ ] AntTransfer - transfer (dataSource, targetKeys, selectedKeys, titles, onChange)
- [ ] Two-panel shuttle selection

### 45. Tree
Status: pending | Strategy: Widget composition

- [ ] AntTree - tree (treeData, checkable, defaultExpandAll, onSelect, onCheck)
- [ ] TreeNode expand/collapse, checkbox selection
- [ ] Virtual scrolling for large trees

### 46. TreeSelect
Status: pending | Strategy: Widget + Overlay

- [ ] AntTreeSelect - tree select (treeData, value, multiple, treeCheckable, onChange)
- [ ] Dropdown tree panel via Overlay

### 47. Upload
Status: pending | Strategy: Widget composition

- [ ] AntUpload - upload (action, accept, multiple, directory, maxCount, fileList, onChange)
- [ ] File list display, upload progress, remove
- [ ] Drag and drop area (Dragger)

### 48. Watermark
Status: pending | Strategy: Widget composition

- [ ] AntWatermark - watermark (content, rotate, gap, offset, font, zIndex)

### 49. QRCode
Status: pending | Strategy: Widget composition

- [ ] AntQRCode - QR code (value, type, icon, size, color, bgColor, status)
- [ ] Error correction level, logo in center

### 50. FloatButton
Status: pending | Strategy: Widget composition

- [ ] AntFloatButton - float button (icon, type, shape, description, tooltip)
- [ ] AntFloatButton.Group - float button group
- [ ] AntFloatButton.BackTop - back to top variant

### 51. ColorPicker
Status: pending | Strategy: Widget + Overlay

- [ ] AntColorPicker - color picker (value, disabled, showText, format, onChange)
- [ ] Color panel with hue/saturation/brightness selection
- [ ] Preset colors

### 52. DatePicker
Status: pending | Strategy: Widget + Overlay

- [ ] AntDatePicker - date picker (value, format, placeholder, disabled, allowClear, showTime, onChange)
- [ ] AntRangePicker - range picker (start/end date)
- [ ] Year/month/day selection panel

### 53. TimePicker
Status: pending | Strategy: Widget + Overlay

- [ ] AntTimePicker - time picker (value, format, placeholder, disabled, onChange)
- [ ] Hour/minute/second column selection

### 54. Cascader
Status: pending | Strategy: Widget + Overlay

- [ ] AntCascader - cascader (options, value, placeholder, disabled, onChange)
- [ ] Multi-level cascading selection panel
- [ ] Search filtering

### 55. AutoComplete
Status: pending | Strategy: Widget + Overlay

- [ ] AntAutoComplete - autocomplete (options, value, placeholder, onChange, onSelect)
- [ ] Input with dropdown suggestions

### 56. Mentions
Status: pending | Strategy: Widget + Overlay

- [ ] AntMentions - mentions (options, value, prefix, placeholder, onChange, onSelect)
- [ ] Trigger dropdown on prefix character (e.g. @)

### 57. Tour
Status: pending | Strategy: Widget + Overlay

- [ ] AntTour - tour (steps, current, open, onClose, onFinish)
- [ ] Guided product tour with step highlighting

### 58. App
Status: pending | Strategy: Context provider wrapper

- [ ] AntApp - app container (message, notification, modal global config)
- [ ] Provides global static methods for message/notification/modal

### 59. ConfigProvider
Status: pending | Strategy: Context-based theme injection

- [ ] AntConfigProvider - config provider (theme, locale, componentSize)
- [ ] Global configuration for all child AntD components

---

## P2 - Advanced / Complex Components (Lower Frequency)

### 60. Message
Status: pending | Strategy: Widget + Overlay (global singleton)

- [ ] Message.info / success / error / warning / loading (content, duration, onClose)
- [ ] Auto-dismiss after duration
- [ ] Stacking multiple messages vertically

### 61. Notification
Status: pending | Strategy: Widget + Overlay (global singleton)

- [ ] Notification.open / success / info / warning / error (message, description, duration, placement)
- [ ] Top-left / top-right / bottom-left / bottom-right placement
- [ ] Auto-dismiss, close button

### 62. Icon
Status: pending | Strategy: Text-based icons or custom drawing

- [ ] AntIcon - icon component (type, spin, rotate, style)
- [ ] Built-in common icon set (close, search, check, loading, etc.)
- [ ] Support custom icon rendering

### 63. VirtualList
Status: pending | Strategy: Extend ScrollView or Widget composition

- [ ] AntVirtualList - virtual list (data, itemHeight, renderItem, height)
- [ ] Only render visible items
- [ ] Calculate visible range from scrollOffset and viewport height
- [ ] Dynamic height support (optional)

---

## Implementation Notes

### Component Composition Patterns
Complex components should be built by composing simpler ones:
- Form = View + Input/Radio/Checkbox/Select + Label + Error Text
- Table = View + ScrollView + Row (header + data) + Cell
- Modal = Overlay + View (container) + Text (title) + View (content) + Button (actions)
- Dropdown = Button/Input + Overlay + ScrollView + List of selectable items
- Select = Input + Overlay + ScrollView + Checkbox/Radio items
- DatePicker = Input + Overlay + Calendar panel
- TreeSelect = Input + Overlay + Tree panel
- Upload = View (drop zone) + Button + List (file list) + ProgressBar
- Transfer = List (left) + Button (move) + List (right)

### Required Native Host Extensions
Some components need new Host primitives or extensions:
- ProgressBar: add circle and dashboard drawing modes
- TextInput: add numeric input mode (for InputNumber)
- ScrollView: virtual scrolling support (for VirtualList and large Tree/Table)
- View: add gradient background support (for Skeleton shimmer effect)

### Overlay-Dependent Components
These require the existing Overlay system:
Tooltip, Popover, Popconfirm, Modal, Drawer, Dropdown, Select, DatePicker, TimePicker, Cascader, ColorPicker, AutoComplete, Mentions, TreeSelect, Message, Notification, Tour, Image preview

### Recommended Implementation Order
Phase 1 (Layout & Typography): Typography, Space, Grid, Layout, Flex
Phase 2 (Basic Data Display): Avatar, Empty, Descriptions, Statistic, Timeline, List, Collapse, Breadcrumb, Steps
Phase 3 (Feedback & Navigation): Tooltip, Popover, Popconfirm, Modal, Drawer, Dropdown, Tabs, Pagination, Menu, BackTop, Affix, Anchor
Phase 4 (Form Controls): Checkbox, Radio, Switch, Slider, Select, Form, InputNumber, Rate, Segmented, AutoComplete, Mentions, Cascader, DatePicker, TimePicker, ColorPicker, TreeSelect
Phase 5 (Data Display Advanced): Table enhancements, Tree, Transfer, Calendar, Carousel, Image, Comment, Result
Phase 6 (Feedback Advanced): Progress (circle), Spin, Skeleton, Message, Notification, Tour
Phase 7 (Utilities): Icon, QRCode, Watermark, FloatButton, Upload, VirtualList, App, ConfigProvider

---

End of TODO
