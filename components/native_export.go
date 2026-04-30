package components

import (
	"github.com/sjm1327605995/tenon/internal/native"
)

// Native element type aliases exported for public use.
// These are the same underlying types as internal/native; users can reference
// them in their own code without importing internal packages.
type (
	View       = native.View
	Text       = native.Text
	Image      = native.Image
	Divider    = native.Divider
	SVGIcon    = native.SVGIcon
	ShadowText = native.ShadowText
)

// NewView creates a new native View container.
func NewView() *View { return native.NewView() }

// NewText creates a new native Text element.
func NewText(content string) *Text { return native.NewText(content) }

// NewImage creates a new native Image element.
func NewImage() *Image { return native.NewImage() }

// NewDivider creates a new native Divider element.
func NewDivider() *Divider { return native.NewDivider() }

// NewSVGIcon creates a new native SVGIcon element.
func NewSVGIcon(pathData string) *SVGIcon { return native.NewSVGIcon(pathData) }

// NewShadowText creates a new native ShadowText element.
func NewShadowText(content string) *ShadowText { return native.NewShadowText(content) }
