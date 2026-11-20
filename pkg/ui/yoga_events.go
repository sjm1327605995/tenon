package ui

import (
	"fmt"
	"os"
	"sync"

	"github.com/sjm1327605995/tenon/yoga"
)

var debugYogaLayout = os.Getenv("TENON_DEBUG_LAYOUT") != ""
var debugYogaLayoutOnce sync.Once

func enableYogaLayoutDebug() {
	if !debugYogaLayout {
		return
	}
	debugYogaLayoutOnce.Do(func() {
		yoga.Subscribe(yoga.EventLayoutPassEnd, func(event yoga.Event) {
			data := event.LayoutData
			fmt.Printf("YOGA layouts=%d measures=%d cached_layouts=%d cached_measures=%d measure_callbacks=%d max_measure_cache=%d\n",
				data.Layouts(), data.Measures(), data.CachedLayouts(), data.CachedMeasures(),
				data.MeasureCallbacks(), data.MaxMeasureCache())
		})
	})
}
