package react

import (
	"github.com/sjm1327605995/tenon/core/vdom"
)

// Renderer 渲染器接口，类似于React Native的UIManager
type Renderer interface {
	// CreateNode 创建原生节点
	CreateNode(element *vdom.Element) NativeComponent
	// UpdateNode 更新原生节点属性
	UpdateNode(component NativeComponent, element *vdom.Element)
	// RemoveNode 移除原生节点
	RemoveNode(component NativeComponent)
	// AppendChild 添加子节点
	AppendChild(parent, child NativeComponent)
	// RemoveChild 移除子节点
	RemoveChild(parent, child NativeComponent)
	// ReplaceChild 替换子节点
	ReplaceChild(parent, oldChild, newChild NativeComponent)
	// InsertBefore 在指定节点前插入
	InsertBefore(parent, newChild, referenceChild NativeComponent)
}

// NativeRenderer 原生渲染器实现
type NativeRenderer struct {
	rootComponents map[vdom.Node]NativeComponent // vdom节点到原生组件的映射
	manager        *NativeRendererManager       // 原生渲染管理器
}

// NewNativeRenderer 创建新的原生渲染器
func NewNativeRenderer() *NativeRenderer {
	return &NativeRenderer{
		rootComponents: make(map[vdom.Node]NativeComponent),
		manager:        NewNativeRendererManager(),
	}
}

// Render 渲染vdom树到原生组件
func (r *NativeRenderer) Render(tree *vdom.Tree) NativeComponent {
	// 初始化原生组件注册
	InitializeNativeComponents()
	
	var rootComponent NativeComponent
	
	// 处理根节点
	for _, rootNode := range tree.Children {
		component := r.renderNode(rootNode, nil)
		if rootComponent == nil {
			rootComponent = component
		} else {
			// 如果有多个根节点，创建一个容器组件
			container := NewNativeView()
			r.AppendChild(container, rootComponent)
			r.AppendChild(container, component)
			rootComponent = container
		}
	}
	
	// 挂载根组件
	if rootComponent != nil {
		r.manager.MountComponent(rootComponent)
	}
	
	return rootComponent
}

// Update 更新vdom树
func (r *NativeRenderer) Update(oldTree, newTree *vdom.Tree) {
	// 计算差异
	patches, err := vdom.Diff(oldTree, newTree)
	if err != nil {
		// 错误处理
		return
	}
	
	// 应用差异
	// 注意：这里需要使用自定义的差异应用方法，而不是DOM的方法
	r.applyPatches(patches)
}

// renderNode 递归渲染节点
func (r *NativeRenderer) renderNode(node vdom.Node, parent NativeComponent) NativeComponent {
	switch n := node.(type) {
	case *vdom.Element:
		// 创建原生组件
		component := r.CreateNode(n)
		
		// 保存映射关系
		r.rootComponents[n] = component
		
		// 如果有父组件，添加为子组件
		if parent != nil {
			r.AppendChild(parent, component)
		}
		
		// 递归渲染子节点
		for _, child := range n.Children() {
			r.renderNode(child, component)
		}
		
		return component
		
	case *vdom.Text:
		// 创建文本组件
		textComponent := CreateNativeComponent("Text", map[string]interface{}{
			"value": string(n.Value),
		})
		
		// 保存映射关系
		r.rootComponents[n] = textComponent
		
		// 如果有父组件，添加为子组件
		if parent != nil {
			r.AppendChild(parent, textComponent)
		}
		
		return textComponent
	}
	
	return nil
}

// CreateNode 创建原生节点
func (r *NativeRenderer) CreateNode(element *vdom.Element) NativeComponent {
	// 将Attrs转换为Props
	props := make(map[string]interface{})
	for _, attr := range element.Attrs {
		props[attr.Name] = attr.Value
	}
	
	// 创建原生组件
	return CreateNativeComponent(element.Name, props)
}

// UpdateNode 更新原生节点属性
func (r *NativeRenderer) UpdateNode(component NativeComponent, element *vdom.Element) {
	// 将Attrs转换为Props
	props := make(map[string]interface{})
	for _, attr := range element.Attrs {
		props[attr.Name] = attr.Value
	}
	
	// 更新组件属性
	component.UpdateProps(props)
}

// RemoveNode 移除原生节点
func (r *NativeRenderer) RemoveNode(component NativeComponent) {
	// 卸载组件
	r.manager.UnmountComponent(component)
}

// AppendChild 添加子节点
func (r *NativeRenderer) AppendChild(parent, child NativeComponent) {
	// 调用父组件的AddChild方法
	if view, ok := parent.(*NativeView); ok {
		view.AddChild(child)
	}
}

// RemoveChild 移除子节点
func (r *NativeRenderer) RemoveChild(parent, child NativeComponent) {
	// 调用父组件的RemoveChild方法
	if view, ok := parent.(*NativeView); ok {
		view.RemoveChild(child)
	}
}

// ReplaceChild 替换子节点
func (r *NativeRenderer) ReplaceChild(parent, oldChild, newChild NativeComponent) {
	// 先移除旧子节点
	r.RemoveChild(parent, oldChild)
	// 再添加新子节点
	r.AppendChild(parent, newChild)
}

// InsertBefore 在指定节点前插入
func (r *NativeRenderer) InsertBefore(parent, newChild, referenceChild NativeComponent) {
	// 简化实现：先移除referenceChild，添加newChild，再添加referenceChild
	r.RemoveChild(parent, referenceChild)
	r.AppendChild(parent, newChild)
	r.AppendChild(parent, referenceChild)
}

// applyPatches 应用差异补丁
func (r *NativeRenderer) applyPatches(patches []vdom.Patcher) {
	for _, patch := range patches {
		switch p := patch.(type) {
		case *vdom.Append:
			// 实现添加节点
			parentComponent := r.rootComponents[p.Parent]
			childComponent := r.renderNode(p.Child, parentComponent)
			
		case *vdom.Remove:
			// 实现移除节点
			component := r.rootComponents[p.Node]
			r.RemoveNode(component)
			delete(r.rootComponents, p.Node)
			
		case *vdom.Replace:
			// 实现替换节点
			oldComponent := r.rootComponents[p.Old]
			newComponent := r.renderNode(p.New, nil)
			
			// 找到父组件
			var parent NativeComponent
			for _, comp := range r.rootComponents {
				if view, ok := comp.(*NativeView); ok {
					// 检查是否包含oldComponent
					// 简化实现，实际需要遍历子组件
					parent = view
					break
				}
			}
			
			if parent != nil {
				r.ReplaceChild(parent, oldComponent, newComponent)
			}
			
			delete(r.rootComponents, p.Old)
			r.rootComponents[p.New] = newComponent
			
		case *vdom.SetAttr:
			// 实现设置属性
			component := r.rootComponents[p.Node]
			r.UpdateNode(component, p.Node.(*vdom.Element))
			
		case *vdom.RemoveAttr:
			// 实现移除属性
			component := r.rootComponents[p.Node]
			r.UpdateNode(component, p.Node.(*vdom.Element))
		}
	}
}