package wtype

import (
	"strings"
)

//  instruction to a liquid handler
type LHInstruction struct {
	ID               string
	ProductID        string
	BlockID          BlockID
	SName            string
	Order            int
	Components       []*LHComponent
	ContainerType    string
	Welladdress      string
	Plateaddress     string
	plateID          string
	Platetype        string
	Vol              float64
	Type             string
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Result           *LHComponent
	gen              int
}

func NewLHInstruction() *LHInstruction {
	var lhi LHInstruction
	lhi.ID = GetUUID()
	lhi.Majorlayoutgroup = -1
	return &lhi
}
func (inst *LHInstruction) AddProduct(cmp *LHComponent) {
	inst.Result = cmp
	inst.ProductID = cmp.ID
}

func (inst *LHInstruction) AddComponent(cmp *LHComponent) {
	if inst == nil {
		return
	}

	inst.Components = append(inst.Components, cmp)
}

func (ins *LHInstruction) Generation() int {
	return ins.gen
}
func (ins *LHInstruction) SetGeneration(i int) {
	ins.gen = i
}

func (ins *LHInstruction) PlateID() string {
	return ins.plateID
}

func (ins *LHInstruction) SetPlateID(pid string) {
	ins.plateID = pid
}

func (ins *LHInstruction) IsMixInPlace() bool {
	smp := ins.Components[0].IsSample()
	return !smp
}

func (ins *LHInstruction) HasAnyParent() bool {
	for _, v := range ins.Components {
		if v.HasAnyParent() {
			return true
		}
	}

	return false
}

func (ins *LHInstruction) HasParent(id string) bool {
	for _, v := range ins.Components {
		if v.HasParent(id) {
			return true
		}
	}
	return false
}

func (ins *LHInstruction) ParentString() string {
	if ins == nil {
		return ""
	}

	tx := make([]string, 0, 1)

	for _, v := range ins.Components {
		//s += v.ParentID + "_"

		pid := v.ParentID

		if pid != "" {
			tx = append(tx, pid)
		}
	}

	if len(tx) == 0 {
		return ""
	} else {
		return strings.Join(tx, "_")
	}

}

func (ins *LHInstruction) ComponentsMoving() string {
	sa := make([]string, 0, 1)
	for i, v := range ins.Components {
		// ignore component 1 if this is a mix-in-place
		if i == 0 && !v.IsSample() {
			continue
		}
		sa = append(sa, v.CName)
	}
	return strings.Join(sa, "+")
}
