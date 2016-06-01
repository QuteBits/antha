// anthalib//liquidhandling/input_plate_setup.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
)

type InputSorter struct {
	Ordered []string
	Values  map[string]wunit.Volume
}

// @implement sort.Interface
func (is InputSorter) Len() int {
	return len(is.Ordered)
}

func (is InputSorter) Swap(i, j int) {
	s := is.Ordered[i]
	is.Ordered[i] = is.Ordered[j]
	is.Ordered[j] = s
}

func (is InputSorter) Less(i, j int) bool {
	vv1 := is.Values[is.Ordered[i]]
	vv2 := is.Values[is.Ordered[j]]

	v1 := vv1.SIValue()
	v2 := vv2.SIValue()

	// we want ascending sort here
	if v1 < v2 {
		return false
	} else if v1 > v2 {
		return true
	}

	// volumes are equal

	ss := sort.StringSlice(is.Ordered)

	return ss.Less(i, j)
}

//  TASK: 	Map inputs to input plates
// INPUT: 	"input_platetype", "inputs"
//OUTPUT: 	"input_plates"      -- these each have components in wells
//		"input_assignments" -- map with arrays of assignment strings, i.e. {tea: [plate1:A:1, plate1:A:2...] }etc.
func input_plate_setup(request *LHRequest) (*LHRequest, error) {
	// I think this might need moving too
	logger.Debug("in input plate setup")
	input_platetypes := (*request).Input_platetypes
	if input_platetypes == nil || len(input_platetypes) == 0 {
		// XXX this is dangerous... until input_plate_linear is replaced we will hit big problems here
		// this configuration needs to happen outside but for now...
		list := factory.GetPlateList()
		input_platetypes = make([]*wtype.LHPlate, len(list))
		for i, platetype := range list {
			input_platetypes[i] = factory.GetPlateByType(platetype)
		}
		(*request).Input_platetypes = input_platetypes
		//debug
	}

	// we assume that input_plates is set if any locs are set
	input_plates := (*request).Input_plates

	if len(input_plates) == 0 {
		input_plates = make(map[string]*wtype.LHPlate, 3)
	}

	// need to fill each plate type

	var curr_plate *wtype.LHPlate

	inputs := (*request).Input_solutions
	//	input_order := (*request).Input_order

	input_order := make([]string, len((*request).Input_order))
	for i, v := range (*request).Input_order {
		input_order[i] = v
	}

	// this needs to be passed in via the request... must specify how much of inputs cannot
	// be satisfied by what's already passed in
	//	input_volumes := make(map[string]wunit.Volume, len(inputs))
	input_volumes := request.Input_vols_wanting

	// sort to make deterministic
	// we sort by a) volume (descending) b) name (alphabetically)

	isrt := InputSorter{input_order, input_volumes}

	sort.Sort(isrt)

	input_order = isrt.Ordered

	weights_constraints := request.Input_setup_weights

	// get the assignment

	//well_count_assignments := choose_plate_assignments(input_volumes, input_platetypes, weights_constraints)

	var well_count_assignments map[string]map[*wtype.LHPlate]int

	if len(input_volumes) != 0 {
		well_count_assignments = choose_plate_assignments(input_volumes, input_platetypes, weights_constraints)
	}

	input_assignments := make(map[string][]string, len(well_count_assignments))

	plates_in_play := make(map[string]*wtype.LHPlate)

	curplaten := 1
	for _, cname := range input_order {
		volume, ok := input_volumes[cname]

		if !ok {
			continue
		}

		component := inputs[cname][0]
		//logger.Debug(fmt.Sprintln("Plate_setup - component", cname, ":"))

		well_assignments, ok := well_count_assignments[cname]

		if !ok {
			continue
		}

		//logger.Debug(fmt.Sprintln("Well assignments: ", well_assignments))

		var curr_well *wtype.LHWell
		ass := make([]string, 0, 3)

		// best hack so far: add an extra well of everything
		// except we should do this at the end
		for platetype, nwells := range well_assignments {
			for i := 0; i < nwells+1; i++ {
				curr_plate = plates_in_play[platetype.Type]

				if curr_plate == nil {
					plates_in_play[platetype.Type] = factory.GetPlateByType(platetype.Type)
					curr_plate = plates_in_play[platetype.Type]
					platename := fmt.Sprintf("Input_plate_%d", curplaten)
					curr_plate.PlateName = platename
					curplaten += 1
					curr_plate.DeclareTemporary()
				}

				// find somewhere to put it
				curr_well, ok = wtype.Get_Next_Well(curr_plate, component, curr_well)

				if !ok {
					// if no space, reset
					plates_in_play[platetype.Type] = nil
					curr_plate = nil
					curr_well = nil
					i -= 1
					continue
				}

				// now put it there

				location := curr_plate.ID + ":" + curr_well.Crds
				ass = append(ass, location)

				// make a duplicate of this component to stick in the well
				// wait wait wait is this right?
				newcomponent := component.Dup()
				newcomponent.Vol = curr_well.MaxVol
				volume.Subtract(curr_well.WorkingVolume())

<<<<<<< HEAD
				fmt.Println("ADDING component ", component.CName, " to ", location, " AT VOLUME ", newcomponent.Vol)
=======
				fmt.Println("ADDING component ", component.CName, " to ", location)
>>>>>>> origin/master

				curr_well.Add(newcomponent)
				input_plates[curr_plate.ID] = curr_plate
			}
		}

		input_assignments[cname] = ass
	}

	// add any remaining assignments

	for _, v := range inputs {
		for _, vv := range v {
			if vv.Loc != "" {
				// append it

				input_assignments[vv.CName] = append(input_assignments[vv.CName], vv.Loc)
			}
		}
	}

	(*request).Input_plates = input_plates
	(*request).Input_assignments = input_assignments
	//return input_plates, input_assignments
	return request, nil
}
