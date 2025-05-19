package meterial

type MaterialAttrib struct {
	Index      int     `json:"index"`
	Threshold  float64 `json:"threshold"`
	MaterialId int     `json:"materialId"`
	Name       string  `json:"name"`
	Required   *bool   `json:"required"`
}

var RefFalse = false

func (this *MaterialAttrib) SetDefaults() {
	if this.Required == nil {
		this.Required = &RefFalse
	}
}
