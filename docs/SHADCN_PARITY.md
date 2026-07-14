# shadcn/ui Parity Plan

Goal: implement every [shadcn/ui component](https://ui.shadcn.com/docs/components) in `pkg/shadcn`, one by one, with full property support. Tracked here; checked off as landed.

Legend: ✅ done · 🟡 partial (needs fuller props/variants) · ⬜ missing

## Core components

| Component | Status | Notes |
|---|---|---|
| Accordion | ✅ | |
| Alert | ✅ | |
| Alert Dialog | ✅ | title + desc + cancel/action (destructive) |
| Aspect Ratio | ✅ | |
| Avatar | ✅ | image fallback? initials only |
| Badge | ✅ | |
| Breadcrumb | ✅ | |
| Button | ✅ | 6 variants / 4 sizes |
| Button Group | ✅ | segmented, active highlight + dividers |
| Calendar | ✅ | |
| Card | ✅ | |
| Carousel | ✅ | |
| Chart | ✅ | Bar/Line/Area/Pie (Vector primitive) |
| Checkbox | ✅ | |
| Collapsible | ✅ | |
| Combobox | ✅ | |
| Command | ✅ | |
| Context Menu | ✅ | right-click (engine OnContextMenu) at cursor |
| Data Table | ✅ | search + sortable headers + pagination |
| Date Picker | ✅ | trigger + Calendar floatPanel |
| Dialog | ✅ | |
| Drawer | ✅ | bottom slide-up + grab handle |
| Dropdown Menu | ✅ | |
| Empty | ✅ | icon + title + description + actions |
| Field | ✅ | label + control + description/error |
| Form | 🟡 | audit props |
| Hover Card | ✅ | |
| Input | ✅ | |
| Input Group | ✅ | leading/trailing addons, borderless input |
| Input OTP | ✅ | segmented digit slots, active highlight, digit filter |
| Item | ✅ | media + title/description + trailing |
| Kbd | ✅ | `Kbd` + `KbdGroup` |
| Label | ✅ | |
| Menubar | ✅ | bar of DropdownMenus |
| Native Select | ✅ | covered by Select (no OS-native select on Ebiten) |
| Navigation Menu | ✅ | |
| Pagination | ✅ | |
| Popover | ✅ | |
| Progress | ✅ | |
| Radio Group | ✅ | |
| Resizable | ✅ | |
| Scroll Area | ✅ | |
| Select | ✅ | |
| Separator | ✅ | |
| Sheet | ✅ | |
| Sidebar | ✅ | groups + items + collapse + divider |
| Skeleton | ✅ | |
| Slider | ✅ | |
| Sonner / Toast | ✅ | Toast + Toaster |
| Spinner | ✅ | themed wrapper of `ui.Spinner` |
| Switch | ✅ | |
| Table | ✅ | |
| Tabs | ✅ | |
| Textarea | ✅ | |
| Toggle | ✅ | |
| Toggle Group | ✅ | |
| Tooltip | ✅ | |
| Typography | ✅ | H1–H4/P/Lead/Muted/Large/Small/InlineCode/Blockquote |

Out of scope (AI/chatbot blocks, not core UI): Attachment, Bubble, Direction, Marker, Message, Message Scroller.

## Implementation order

1. **Primitives** — Kbd, Spinner, Empty, Typography.
2. **Form/input family** — Native Select, Input OTP, Input Group, Button Group, Field, Item.
3. **Overlays/menus** — Alert Dialog, Context Menu, Menubar, Drawer, Date Picker.
4. **Data/large** — Data Table, Sidebar, Chart (Line/Area/Pie).
5. **Audit pass** — fill missing props/variants on existing components.
