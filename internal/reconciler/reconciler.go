package reconciler

import (
	"github.com/sjm1327605995/tenon/internal/layout"
	"github.com/sjm1327605995/tenon/pkg/types"
)

type Reconciler struct {
	rootUI      types.UI
	rootElement types.Element
	needsUpdate bool
}

func NewReconciler(rootUI types.UI) *Reconciler {
	reconciler := &Reconciler{
		rootUI:      rootUI,
		needsUpdate: true,
	}

	reconciler.reconcile()
	return reconciler
}

func (r *Reconciler) reconcile() {
	if r.needsUpdate {
		r.rootElement = r.rootUI.Render()
		layout.CreateElement(r.rootElement)
		r.calculateLayout()
		r.needsUpdate = false
	}
}

func (r *Reconciler) calculateLayout() {
	if r.rootElement == nil {
		return
	}

	layout.CalculateLayout(r.rootElement, 1000, 1000)
	layout.UpdateElementLayout(r.rootElement)
}

func (r *Reconciler) GetRootElement() types.Element {
	r.reconcile()
	return r.rootElement
}

func (r *Reconciler) Update() {
	r.needsUpdate = true
}