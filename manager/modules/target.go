package modules

type Target struct {
	Pointer			string		`json:"pointer"`
	Weight			int8		`json:"weight"`
	CurrentWeight	int8		`json:"current_weight"`
}
