package components

import (
	"github.com/sjm1327605995/tenon/yoga"
)

// SheetSide defines which side the sheet appears from.
type SheetSide int

const (
	SheetTop SheetSide = iota
	SheetRight
	SheetBottom
	SheetLeft
)

// Sheet is a side panel overlay (similar to Drawer but with shadcn/ui naming).
type Sheet struct {
	*Drawer
}

// NewSheet creates a sheet.
func NewSheet(side SheetSide) *Sheet {
	var drawerSide DrawerSide
	switch side {
	case SheetLeft:
		drawerSide = DrawerLeft
	case SheetRight:
		drawerSide = DrawerRight
	case SheetTop:
		drawerSide = DrawerTop
	case SheetBottom:
		drawerSide = DrawerBottom
	}
	s := &Sheet{Drawer: NewDrawer(drawerSide)}
	s.panel.SetPadding(yoga.EdgeAll, 24)
	s.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	s.panel.SetGap(yoga.GutterAll, 16)
	return s
}

// ElementType returns type identifier.
func (s *Sheet) ElementType() string { return "Sheet" }
