package mixer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	planner "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

var (
	_ target.Device = &Mixer{}
)

type Mixer struct {
	driver     driver.ExtendedLiquidhandlingDriver
	properties driver.LHProperties
	opt        Opt
}

func (a *Mixer) String() string {
	return "Mixer"
}

func (a *Mixer) CanCompile(req ast.Request) bool {
	// TODO: remove when mixers have wait instruction
	if req.Time != nil {
		return false
	}
	if req.Temp != nil {
		return false
	}
	if req.Move != nil {
		return false
	}
	// TODO: Add specific volume constraints
	return req.MixVol != nil
}

func (a *Mixer) MoveCost(from target.Device) int {
	if from == a {
		return 0
	}
	return human.HumanByXCost + 1
}

type lhreq struct {
	*planner.LHRequest     // A request
	*planner.Liquidhandler // ... and its associated planner
}

func (a *Mixer) makeLhreq() (*lhreq, error) {
	addPlate := func(req *planner.LHRequest, ip *wtype.LHPlate) error {
		if _, seen := req.Input_plates[ip.ID]; seen {
			return fmt.Errorf("plate %q already added", ip.ID)
		} else {
			req.Input_plates[ip.ID] = ip
			return nil
		}
	}

	req := planner.NewLHRequest()
	pols, err := driver.GetLHPolicyForTest()

	if err != nil {
		return nil, err
	}
	req.Policies = pols

	plan := planner.Init(&a.properties)

	if p := a.opt.MaxPlates; p != nil {
		req.Input_setup_weights["MAX_N_PLATES"] = *p
	}

	if p := a.opt.MaxWells; p != nil {
		req.Input_setup_weights["MAX_N_WELLS"] = *p
	}

	if p := a.opt.ResidualVolumeWeight; p != nil {
		req.Input_setup_weights["RESIDUAL_VOLUME_WEIGHT"] = *p
	}

	// TODO -- error check here to prevent nil values

	if p := a.opt.InputPlateType; len(p) != 0 {
		for _, v := range p {
			req.Input_platetypes = append(req.Input_platetypes, factory.GetPlateByType(v))
		}
	}

	if p := a.opt.OutputPlateType; len(p) != 0 {
		for _, v := range p {
			req.Output_platetypes = append(req.Output_platetypes, factory.GetPlateByType(v))
		}
	}

	if p := a.opt.TipType; len(p) != 0 {
		for _, v := range p {
			req.Tips = append(req.Tips, factory.GetTipByType(v))
		}
	}

	if p := a.opt.InputPlateFiles; len(p) != 0 {
		for _, filename := range p {
			ip, err := parseInputPlateFile(filename)
			if err != nil {
				return nil, fmt.Errorf("cannot parse file %q: %s", filename, err)
			}
			if err := addPlate(req, ip); err != nil {
				return nil, err
			}
		}
	}

	if p := a.opt.InputPlateData; len(p) != 0 {
		for idx, bs := range p {
			buf := bytes.NewBuffer(bs)
			ip, err := parseInputPlateData(buf)
			if err != nil {
				return nil, fmt.Errorf("cannot parse data at idx %d: %s", idx, err)
			}
			if err := addPlate(req, ip); err != nil {
				return nil, err
			}
		}
	}

	if ips := a.opt.InputPlates; len(ips) != 0 {
		for _, ip := range ips {
			if err := addPlate(req, ip); err != nil {
				return nil, err
			}
		}
	}

	err = req.ConfigureYourself()
	if err != nil {
		return nil, err
	}

	return &lhreq{
		LHRequest:     req,
		Liquidhandler: plan,
	}, nil
}

func (a *Mixer) Compile(cmds []ast.Command) ([]target.Inst, error) {
	var mixes []*wtype.LHInstruction
	for _, c := range cmds {
		if m, ok := c.(*ast.Mix); !ok {
			return nil, fmt.Errorf("cannot compile %T", c)
		} else {
			mixes = append(mixes, m.Inst)
		}
	}
	if inst, err := a.makeMix(mixes); err != nil {
		return nil, err
	} else {
		return []target.Inst{inst}, nil
	}
}

func (a *Mixer) saveFile(name string) ([]byte, error) {
	data, status := a.driver.GetOutputFile()
	if !status.OK {
		return nil, fmt.Errorf("%d: %s", status.Errorcode, status.Msg)
	} else if len(data) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	bs := []byte(data)

	if err := tw.WriteHeader(&tar.Header{
		Name:    name,
		Mode:    0644,
		Size:    int64(len(bs)),
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	} else if _, err := tw.Write(bs); err != nil {
		return nil, err
	} else if err := tw.Close(); err != nil {
		return nil, err
	} else if err := gw.Close(); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (a *Mixer) makeMix(mixes []*wtype.LHInstruction) (target.Inst, error) {
	hasPlate := func(plates []*wtype.LHPlate, typ, id string) bool {
		for _, p := range plates {
			if p.Type == typ && (id == "" || p.ID == id) {
				return true
			}
		}
		return false
	}

	getId := func(mixes []*wtype.LHInstruction) (r wtype.BlockID) {
		m := make(map[wtype.BlockID]bool)
		for _, mix := range mixes {
			m[mix.BlockID] = true
		}
		for k := range m {
			r = k
			break
		}
		return
	}

	r, err := a.makeLhreq()
	if err != nil {
		return nil, err
	}
	r.LHRequest.BlockID = getId(mixes)

	for _, mix := range mixes {
		if len(mix.Platetype) != 0 && !hasPlate(r.LHRequest.Output_platetypes, mix.Platetype, mix.PlateID()) {
			p := factory.GetPlateByType(mix.Platetype)
			p.ID = mix.PlateID()
			r.LHRequest.Output_platetypes = append(r.LHRequest.Output_platetypes, p)
		}
		r.LHRequest.Add_instruction(mix)
	}

	err = r.Liquidhandler.MakeSolutions(r.LHRequest)

	if err != nil {
		// depending on what went wrong we might error out or return
		// an error instruction

		if wtype.LHErrorIsInternal(err) {
			return nil, err
		} else {
			return &target.CmpError{Err: err, Dev: a}, nil
		}
	}

	// TODO: Desired filename not exposed in current driver interface, so pick
	// a name. So far, at least Gilson software cares what the filename is, so
	// use .sqlite for compatibility
	tarball, err := a.saveFile("input.sqlite")
	if err != nil {
		return nil, err
	}

	var ftype string
	if a.properties.Mnfr != "" {
		ftype = fmt.Sprintf("application/%s", strings.ToLower(a.properties.Mnfr))
	}
	return &target.Mix{
		Dev:        a,
		Request:    r.LHRequest,
		Properties: a.properties,
		Files: target.Files{
			Tarball: tarball,
			Type:    ftype,
		},
	}, nil
}

func New(opt Opt, d driver.ExtendedLiquidhandlingDriver) (*Mixer, error) {
	p, status := d.GetCapabilities()
	if !status.OK {
		return nil, fmt.Errorf("cannot get capabilities: %s", status.Msg)
	}
	p.Driver = d
	return &Mixer{driver: d, properties: p, opt: opt}, nil
}
