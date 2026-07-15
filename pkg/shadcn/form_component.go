package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// FormField 描述表单里的一个字段：标签、说明、校验、以及如何渲染其控件。
type FormField struct {
	Name        string
	Label       string
	Description string
	Initial     string
	// Validate 返回错误信息（""=通过）。提交时对每个字段调用。
	Validate func(value string) string
	// Control 渲染绑定到 value/onChange 的控件，如 shadcn.Input。
	Control func(value string, onChange func(string)) *ui.Node
}

type FormProps struct {
	Fields      []FormField
	SubmitLabel string // 默认「提交」
	OnSubmit    func(values map[string]string)
}

// Form 是带校验的受控表单：管理各字段值、提交时逐个校验，全部通过才回调 OnSubmit，
// 否则在对应字段下显示错误（复用 Field + Button）。
func Form(p FormProps) *ui.Node { return ui.Use(form, p) }

func form(p FormProps) *ui.Node {
	values, setValues := ui.UseState(initialFormValues(p.Fields))
	errors, setErrors := ui.UseState(map[string]string{})

	setField := func(name, val string) {
		nv := make(map[string]string, len(values))
		for k, v := range values {
			nv[k] = v
		}
		nv[name] = val
		setValues(nv)
	}
	submit := func() {
		errs := map[string]string{}
		for _, f := range p.Fields {
			if f.Validate != nil {
				if msg := f.Validate(values[f.Name]); msg != "" {
					errs[f.Name] = msg
				}
			}
		}
		setErrors(errs)
		if len(errs) == 0 && p.OnSubmit != nil {
			p.OnSubmit(values)
		}
	}

	kids := []*ui.Node{ui.Style(ui.Column, ui.Gap(16))}
	for _, f := range p.Fields {
		field := f
		var control *ui.Node
		if field.Control != nil {
			control = field.Control(values[field.Name], func(v string) { setField(field.Name, v) })
		}
		kids = append(kids, Field(FieldProps{
			Label: field.Label, Description: field.Description, Error: errors[field.Name]}, control))
	}
	sl := p.SubmitLabel
	if sl == "" {
		sl = "提交"
	}
	kids = append(kids, Button(ButtonProps{OnClick: submit}, ui.Text(sl)))
	return ui.Div(kids...)
}

func initialFormValues(fields []FormField) map[string]string {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f.Name] = f.Initial
	}
	return m
}
