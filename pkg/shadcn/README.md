# tenon/pkg/shadcn

A [shadcn/ui](https://ui.shadcn.com)-style component library for [tenon/pkg/ui](../ui). Each component is a base `ui` primitive restyled with theme tokens (`ui.UseTheme`) and interaction state (`ui.UseInteraction`). shadcn components **mix freely with base `ui` nodes**.

## Usage

Wrap your app in `ui.ThemeProvider` (light/dark) so components pick up tokens:

```go
ui.ThemeProvider(ui.LightTheme,
    shadcn.Card(
        shadcn.CardHeader(shadcn.CardTitle("Sign in")),
        shadcn.CardContent(
            shadcn.Label("Email"),
            shadcn.Input(shadcn.InputProps{Value: email, OnChange: setEmail}),
        ),
        shadcn.CardFooter(
            shadcn.Button(shadcn.ButtonProps{OnClick: submit}, ui.Text("Continue")),
        ),
    ),
)
```

## Components

| Component | Notes |
|---|---|
| `Button` | 6 variants (Default/Secondary/Destructive/Outline/Ghost/Link), 4 sizes; hover + press feedback |
| `Badge` | Default/Secondary/Destructive/Outline |
| `Card` | + `CardHeader`/`CardTitle`/`CardDescription`/`CardContent`/`CardFooter` |
| `Input` | controlled text field |
| `Label` | form label |
| `Checkbox` | controlled |
| `Switch` | controlled, animated thumb |
| `RadioGroup` | controlled, from `Options []string` |
| `Slider` | controlled, draggable |
| `Progress` | 0..1 |
| `Separator` | horizontal/vertical |
| `Skeleton` | pulsing placeholder |
| `Avatar` | initials |
| `Alert` | Default/Destructive + `AlertTitle`/`AlertDescription` |
| `Tabs` | segmented tab bar |
| `Toggle` | pressable toggle |
| `Dialog` | Portal modal + `DialogTitle`/`DialogDescription`, enter/exit transition |
| `Textarea` | multi-line input (wrapping, Enter=newline, grows) |
| `Popover` | anchored floating panel, click-outside to close |
| `Tooltip` | hover-anchored label above the trigger |
| `DropdownMenu` | anchored menu (`[]MenuItem`) |
| `Select` | anchored options dropdown (controlled value) |
| `Combobox` | searchable dropdown: type-to-filter options + check on selected (Select × Command) |
| `Table` | `TableRow`/`TableHead`/`TableCell` (equal columns) |
| `Accordion` | collapsible sections, single-open, height animation |
| `Toast` / `Toaster` | global notifications, auto-dismiss (mount `Toaster()` at root, call `Toast(...)` anywhere) |
| `Sheet` | edge drawer (left/right), slide-in transition |
| `Command` | command palette: search input + filtered list |
| `Calendar` | month date picker (prev/next, select) |
| `NavigationMenu` | horizontal nav bar with dropdowns |
| `Breadcrumb` | path navigation |
| `Pagination` | page numbers + prev/next |
| `Collapsible` | single collapsible, height animation |
| `ToggleGroup` | single-select segmented group |
| `AspectRatio` | fixed-ratio container |
| `ScrollArea` | fixed-height scroll region |
| `HoverCard` | hover-triggered info card (stays open over the card) |
| `Carousel` | sliding slideshow with prev/next + dots |
| `Resizable` | two panels + draggable divider |
| `BarChart` | simple bar chart |

Controlled components take value in / `OnChange` out — hold the state with `ui.UseState`.

See `example/accordion` for a full showcase — a shadcn/ui docs-page re-creation (with a light/dark toggle).

Anchored overlays (`Popover`/`Tooltip`/`DropdownMenu`/`Select`) use `ui.UseMeasure()` to read the trigger's on-screen rect and position a `ui.Portal` panel at it; the panel measures itself and **flips above the trigger when it would overflow the bottom** (via `ui.Viewport()`). `Textarea` uses the base `ui.Multiline()` input. `Accordion` measures content height (`UseMeasure`) and animates it with `UseTween`.

`Sheet` slides via `UseTransition` + `TranslateXY`; `Command` filters items by the search query; `Calendar` uses Go `time` for the month grid; `NavigationMenu` reuses `DropdownMenu` per item. All overlays (Dialog/Sheet/Command/Popover/Select/DropdownMenu) close on **Escape** (base `ui.UseEscape`, topmost-first).

## Coverage

~42 components across form / display / feedback / navigation / overlay / data — the full shadcn/ui core plus common extras — all on the same theme + interaction + measure + animation foundation.
