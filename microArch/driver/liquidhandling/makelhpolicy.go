// /anthalib/driver/liquidhandling/makelhpolicy.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	antha "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/internal/github.com/ghodss/yaml"
)

var DOEliquidhandlingFile = "8run4cpFactorial.xlsx" //"FullFactorial.xlsx" // "ScreenLHPolicyDOE2.xlsx"
var DXORJMP = "JMP"                                 //"DX"

func MakePolicies() map[string]LHPolicy {
	pols := make(map[string]LHPolicy)

	// what policies do we need?
	pols["water"] = MakeWaterPolicy()
	pols["culture"] = MakeCulturePolicy()
	pols["culturereuse"] = MakeCultureReusePolicy()
	pols["glycerol"] = MakeGlycerolPolicy()
	pols["solvent"] = MakeSolventPolicy()
	pols["default"] = MakeDefaultPolicy()
	pols["foamy"] = MakeFoamyPolicy()
	pols["dna"] = MakeDNAPolicy()
	pols["DoNotMix"] = MakeDefaultPolicy()
	pols["NeedToMix"] = MakeNeedToMixPolicy()
	pols["viscous"] = MakeViscousPolicy()
	pols["Paint"] = MakePaintPolicy()

	//      pols["lysate"] = MakeLysatePolicy()
	pols["protein"] = MakeProteinPolicy()
	pols["detergent"] = MakeDetergentPolicy()
	pols["load"] = MakeLoadPolicy()
	pols["loadlow"] = MakeLoadPolicy()
	pols["loadwater"] = MakeLoadWaterPolicy()
	pols["DispenseAboveLiquid"] = MakeDispenseAboveLiquidPolicy()
	pols["PEG"] = MakePEGPolicy()
	pols["Protoplasts"] = MakeProtoplastPolicy()
	pols["dna_mix"] = MakeDNAMixPolicy()
	//      pols["lysate"] = MakeLysatePolicy()

	/*policies, names := PolicyMaker(Allpairs, "DOE_run", false)
	for i, policy := range policies {
		pols[names[i]] = policy
	}
	*/

	if antha.Anthafileexists(DOEliquidhandlingFile) {
		fmt.Println("found lhpolicy doe file", DOEliquidhandlingFile)
		policies, names, _, err := PolicyMakerfromDesign(DXORJMP, DOEliquidhandlingFile, "DOE_run")

		for i, policy := range policies {
			pols[names[i]] = policy
		}
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("no lhpolicy doe file found named: ", DOEliquidhandlingFile)
	}
	return pols

}

func PolicyMakerfromDesign(DXORJMP string, dxdesignfilename string, prepend string) (policies []LHPolicy, names []string, runs []Run, err error) {

	if DXORJMP == "DX" {

		runs, err = RunsFromDXDesign(filepath.Join(antha.Dirpath(), dxdesignfilename), []string{"Pre_MIX", "POST_MIX"})
		if err != nil {
			return policies, names, runs, err
		}

	} else if DXORJMP == "JMP" {

		patterncolumn := 0
		factorcolumns := []int{1, 2, 3, 4, 5}
		responsecolumns := []int{6, 7, 8, 9}

		runs, err = RunsFromJMPDesign(filepath.Join(antha.Dirpath(), dxdesignfilename), patterncolumn, factorcolumns, responsecolumns, []string{"PRE_MIX", "POST_MIX"})
		if err != nil {
			return policies, names, runs, err
		}
	} else {
		return policies, names, runs, fmt.Errorf("only JMP or DX allowed as valid inputs for DXORJMP variable")
	}
	policies, names = PolicyMakerfromRuns(runs, prepend, false)
	return policies, names, runs, err
}

func PolicyMaker(factors []DOEPair, nameprepend string, concatfactorlevelsinname bool) (policies []LHPolicy, names []string) {

	runs := AllCombinations(factors)

	policies, names = PolicyMakerfromRuns(runs, nameprepend, concatfactorlevelsinname)

	return
}

func PolicyMakerfromRuns(runs []Run, nameprepend string, concatfactorlevelsinname bool) (policies []LHPolicy, names []string) {

	names = make([]string, 0)
	policies = make([]LHPolicy, 0)

	//policy := make(LHPolicy, 0)
	policy := MakeDefaultPolicy()
	for _, run := range runs {
		for j, desc := range run.Factordescriptors {
			policy[desc] = run.Setpoints[j]
		}

		// raising runtime error when using concat == true
		if concatfactorlevelsinname {
			name := nameprepend
			for key, value := range policy {
				name = fmt.Sprint(name, "_", key, ":", value)

			}
			//	fmt.Println(name)
		} else {
			names = append(names, nameprepend+strconv.Itoa(run.RunNumber))
		}
		policies = append(policies, policy)
		//fmt.Println("len policy = ", len(policy))
		policy = MakeDefaultPolicy()
	}

	return
}

//func MakeLysatePolicy() LHPolicy {
//        lysatepolicy := make(LHPolicy, 6)
//        lysatepolicy["ASPSPEED"] = 1.0
//        lysatepolicy["DSPSPEED"] = 1.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["DSP_WAIT"] = 2.0
//        lysatepolicy["PRE_MIX"] = 5
//        lysatepolicy["CAN_MSA"]= false
//        return lysatepolicy
//}
//func MakeProteinPolicy() LHPolicy {
//        proteinpolicy := make(LHPolicy, 4)
//        proteinpolicy["DSPREFERENCE"] = 2
//        proteinpolicy["CAN_MULTI"] = true
//        proteinpolicy["PRE_MIX"] = 3
//        proteinpolicy["CAN_MSA"] = false
//        return proteinpolicy
//}

func GetPolicyByName(policyname string) (lhpolicy LHPolicy, policypresent bool) {
	policymap := MakePolicies()

	lhpolicy, policypresent = policymap[policyname]
	return
}

/*
Available policy field names and policy types to use:

Here is a list of everything currently implemented in the liquid handling policy framework

ASPENTRYSPEED,                    ,float64,      ,allows slow moves into liquids
ASPSPEED,                                ,float64,     ,aspirate pipetting rate
ASPZOFFSET,                           ,float64,      ,mm above well bottom when aspirating
ASP_WAIT,                                   ,float64,     ,wait time in seconds post aspirate
BLOWOUTOFFSET,                    ,float64,     ,mm above BLOWOUTREFERENCE
BLOWOUTREFERENCE,          ,int,             ,where to be when blowing out: 0 well bottom, 1 well top
BLOWOUTVOLUME,                ,float64,      ,how much to blow out
CAN_MULTI,                              ,bool,         ,is multichannel operation allowed?
DSPENTRYSPEED,                    ,float64,     ,allows slow moves into liquids
DSPREFERENCE,                      ,int,            ,where to be when dispensing: 0 well bottom, 1 well top
DSPSPEED,                              ,float64,       ,dispense pipetting rate
DSPZOFFSET,                         ,float64,          ,mm above DSPREFERENCE
DSP_WAIT,                               ,float64,        ,wait time in seconds post dispense
EXTRA_ASP_VOLUME,            ,wunit.Volume,       ,additional volume to take up when aspirating
EXTRA_DISP_VOLUME,           ,wunit.Volume,       ,additional volume to dispense
JUSTBLOWOUT,                      ,bool,            ,shortcut to get single transfer
POST_MIX,                               ,int,               ,number of mix cycles to do after dispense
POST_MIX_RATE,                    ,float64,          ,pipetting rate when post mixing
POST_MIX_VOL,                      ,float64,          ,volume to post mix (ul)
POST_MIX_X,                          ,float64,           ,x offset from centre of well (mm) when post-mixing
POST_MIX_Y,                          ,float64,           ,y offset from centre of well (mm) when post-mixing
POST_MIX_Z,                          ,float64,           ,z offset from centre of well (mm) when post-mixing
PRE_MIX,                                ,int,               ,number of mix cycles to do before aspirating
PRE_MIX_RATE,                     ,float64,           ,pipetting rate when pre mixing
PRE_MIX_VOL,                       ,float64,           ,volume to pre mix (ul)
PRE_MIX_X,                              ,float64,          ,x offset from centre of well (mm) when pre-mixing
PRE_MIX_Y,                              ,float64,           ,y offset from centre of well (mm) when pre-mixing
PRE_MIX_Z,                              ,float64,           ,z offset from centre of well (mm) when pre-mixing
TIP_REUSE_LIMIT,                    ,int,                ,number of times tips can be reused for asp/dsp cycles
TOUCHOFF,                              ,bool,             ,whether to move to TOUCHOFFSET after dispense
TOUCHOFFSET,                         ,float64,          ,mm above wb to touch off at


*/

func MakePEGPolicy() LHPolicy {
	policy := make(LHPolicy, 9)
	policy["ASP_SPEED"] = 1.5
	policy["DSP_SPEED"] = 1.5
	policy["ASP_WAIT"] = 2.0
	policy["DSP_WAIT"] = 2.0
	policy["ASPZOFFSET"] = 2.5
	policy["DSPZOFFSET"] = 2.5
	policy["POST_MIX"] = 3
	policy["POST_MIX_Z"] = 3.5
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = true
	policy["CAN_MULTI"] = false
	return policy
}

func MakeProtoplastPolicy() LHPolicy {
	policy := make(LHPolicy, 7)
	policy["ASP_SPEED"] = 0.15
	policy["DSP_SPEED"] = 0.15
	policy["ASPZOFFSET"] = 1.5
	policy["DSPZOFFSET"] = 1.5
	//policy["BLOWOUTVOLUME"] = 0.0
	//policy["BLOWOUTVOLUMEUNIT"] = "ul"
	//policy["TOUCHOFF"] = true
	policy["TIP_REUSE_LIMIT"] = 5
	policy["CAN_MULTI"] = false
	return policy
}

func MakePaintPolicy() LHPolicy {

	policy := make(LHPolicy, 13)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["ASP_SPEED"] = 1.5
	policy["DSP_SPEED"] = 1.5
	policy["ASP_WAIT"] = 1.0
	policy["DSP_WAIT"] = 1.0
	policy["PRE_MIX"] = 3
	policy["POST_MIX"] = 3
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = true
	policy["CAN_MULTI"] = false

	return policy
}

func MakeDispenseAboveLiquidPolicy() LHPolicy {

	policy := make(LHPolicy, 7)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	policy["ASP_SPEED"] = 3.0
	policy["DSP_SPEED"] = 3.0
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false

	return policy
}

func MakeColonyPolicy() LHPolicy {

	policy := make(LHPolicy, 10)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.0
	policy["ASP_SPEED"] = 3.0
	policy["DSP_SPEED"] = 3.0
	policy["ASP_WAIT"] = 1.0
	policy["POST_MIX"] = 3
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = true
	policy["CAN_MULTI"] = false

	return policy
}

func MakeWaterPolicy() LHPolicy {
	waterpolicy := make(LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 0
	//waterpolicy["CAN_MULTI"] = true
	waterpolicy["CAN_MULTI"] = false
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["DSPZOFFSET"] = 0.5
	return waterpolicy
}
func MakeFoamyPolicy() LHPolicy {
	foamypolicy := make(LHPolicy, 5)
	foamypolicy["DSPREFERENCE"] = 0
	foamypolicy["TOUCHOFF"] = true
	foamypolicy["CAN_MULTI"] = false
	foamypolicy["CAN_MSA"] = false
	foamypolicy["CAN_SDD"] = true
	return foamypolicy
}
func MakeCulturePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy["PRE_MIX"] = 2
	culturepolicy["ASPSPEED"] = 2.0
	culturepolicy["DSPSPEED"] = 2.0
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = false
	culturepolicy["CAN_SDD"] = false
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["TIP_REUSE_LIMIT"] = 0
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["BLOWOUTVOLUME"] = 0.0
	culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

	return culturepolicy
}

func MakePlateOutPolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 17)
	culturepolicy["PRE_MIX"] = 2
	culturepolicy["PRE_MIX_VOLUME"] = 50
	culturepolicy["PRE_MIX_Z"] = 2.0
	culturepolicy["PRE_MIX_RATE"] = 4.0
	culturepolicy["ASPSPEED"] = 4.0
	culturepolicy["ASPZOFFSET"] = 2.0
	culturepolicy["DSPSPEED"] = 4.0
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = false
	culturepolicy["CAN_SDD"] = false
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["TIP_REUSE_LIMIT"] = 0
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["BLOWOUTVOLUME"] = 0.0
	culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

	return culturepolicy
}

func MakeCultureReusePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy["PRE_MIX"] = 2
	culturepolicy["ASPSPEED"] = 2.0
	culturepolicy["DSPSPEED"] = 2.0
	//culturepolicy["CAN_MULTI"] = true
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = true
	culturepolicy["CAN_SDD"] = true
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["BLOWOUTVOLUME"] = 0.0
	culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

	return culturepolicy
}

func MakeGlycerolPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 6)
	glycerolpolicy["ASP_SPEED"] = 1.5
	glycerolpolicy["DSP_SPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	glycerolpolicy["CAN_MULTI"] = false
	return glycerolpolicy
}

func MakeViscousPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 4)
	glycerolpolicy["ASP_SPEED"] = 1.5
	glycerolpolicy["DSP_SPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	//glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	return glycerolpolicy
}
func MakeSolventPolicy() LHPolicy {
	solventpolicy := make(LHPolicy, 5)
	solventpolicy["PRE_MIX"] = 3
	solventpolicy["DSPREFERENCE"] = 0
	solventpolicy["DSPZOFFSET"] = 0.5
	solventpolicy["NO_AIR_DISPENSE"] = true
	solventpolicy["CAN_MULTI"] = false
	return solventpolicy
}

func MakeDNAPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 10)
	dnapolicy["ASPSPEED"] = 2.0
	dnapolicy["DSPSPEED"] = 2.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy
}

func MakeDNAMixPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 50
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	return dnapolicy
}

func MakeDetergentPolicy() LHPolicy {
	detergentpolicy := make(LHPolicy, 9)
	//        detergentpolicy["POST_MIX"] = 3
	detergentpolicy["ASPSPEED"] = 1.0
	detergentpolicy["DSPSPEED"] = 1.0
	detergentpolicy["CAN_MULTI"] = false
	detergentpolicy["CAN_MSA"] = false
	detergentpolicy["CAN_SDD"] = false
	detergentpolicy["DSPREFERENCE"] = 0
	detergentpolicy["DSPZOFFSET"] = 0.5
	detergentpolicy["TIP_REUSE_LIMIT"] = 8
	detergentpolicy["NO_AIR_DISPENSE"] = true
	return detergentpolicy
}
func MakeProteinPolicy() LHPolicy {
	proteinpolicy := make(LHPolicy, 10)
	proteinpolicy["POST_MIX"] = 5
	proteinpolicy["POST_MIX_VOLUME"] = 50
	proteinpolicy["ASPSPEED"] = 2.0
	proteinpolicy["DSPSPEED"] = 2.0
	proteinpolicy["CAN_MULTI"] = false
	proteinpolicy["CAN_MSA"] = false
	proteinpolicy["CAN_SDD"] = false
	proteinpolicy["DSPREFERENCE"] = 0
	proteinpolicy["DSPZOFFSET"] = 0.5
	proteinpolicy["TIP_REUSE_LIMIT"] = 0
	proteinpolicy["NO_AIR_DISPENSE"] = true
	return proteinpolicy
}
func MakeLoadPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 0
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	return loadpolicy
}

func MakeLoadWaterPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 100
	return loadpolicy
}

func MakeLoadlowPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 1.0
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 0
	loadpolicy["DSPZOFFSET"] = 0.5
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	return loadpolicy
}

func MakeNeedToMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 10)
	dnapolicy["POST_MIX"] = 4
	dnapolicy["POST_MIX_VOLUME"] = 75
	dnapolicy["ASPSPEED"] = 4.0
	dnapolicy["DSPSPEED"] = 4.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy

}

func MakeDefaultPolicy() LHPolicy {
	defaultpolicy := make(LHPolicy, 21)
	// don't set this here -- use defaultpipette speed or there will be inconsistencies
	// defaultpolicy["ASP_SPEED"] = 3.0
	// defaultpolicy["DSP_SPEED"] = 3.0
	defaultpolicy["TOUCHOFF"] = false
	defaultpolicy["TOUCHOFFSET"] = 0.5
	defaultpolicy["ASPREFERENCE"] = 0
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPREFERENCE"] = 0
	defaultpolicy["DSPZOFFSET"] = 0.5
	defaultpolicy["CAN_MULTI"] = false
	defaultpolicy["CAN_MSA"] = false
	defaultpolicy["CAN_SDD"] = true
	defaultpolicy["TIP_REUSE_LIMIT"] = 100
	defaultpolicy["BLOWOUTREFERENCE"] = 1
	defaultpolicy["BLOWOUTOFFSET"] = -0.5
	defaultpolicy["BLOWOUTVOLUME"] = 200.0
	defaultpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	defaultpolicy["PTZREFERENCE"] = 1
	defaultpolicy["PTZOFFSET"] = -0.5
	defaultpolicy["NO_AIR_DISPENSE"] = false
	defaultpolicy["DEFAULTPIPETTESPEED"] = 3.0
	defaultpolicy["MANUALPTZ"] = false
	defaultpolicy["JUSTBLOWOUT"] = false
	defaultpolicy["DONT_BE_DIRTY"] = true
	return defaultpolicy
}

func MakeJBPolicy() LHPolicy {
	jbp := make(LHPolicy, 1)
	jbp["JUSTBLOWOUT"] = true
	//jbp["TOUCHOFF"] = true
	return jbp
}

func MakeTOPolicy() LHPolicy {
	top := make(LHPolicy, 1)
	top["TOUCHOFF"] = true
	return top
}

func MakeLVExtraPolicy() LHPolicy {
	lvep := make(LHPolicy, 2)
	lvep["EXTRA_ASP_VOLUME"] = wunit.NewVolume(0.5, "ul")
	lvep["EXTRA_DISP_VOLUME"] = wunit.NewVolume(0.5, "ul")
	return lvep
}

func MakeLVOffsetPolicy() LHPolicy {
	lvop := make(LHPolicy, 2)
	lvop["ASPZOFFSET"] = 0.0
	lvop["DSPZOFFSET"] = 0.0
	lvop["POST_MIX_Z"] = 0.0
	lvop["PRE_MIX_Z"] = 0.0
	lvop["DSPREFERENCE"] = 0
	lvop["ASPREFERENCE"] = 0
	return lvop
}

func GetLHPolicyForTest() (*LHPolicyRuleSet, error) {

	// make some policies

	policies := MakePolicies()

	// now make rules

	lhpr := NewLHPolicyRuleSet()

	for name, policy := range policies {
		rule := NewLHPolicyRule(name)
		err := rule.AddCategoryConditionOn("LIQUIDCLASS", name)

		if err != nil {
			return nil, err
		}
		lhpr.AddRule(rule, policy)
	}

	// add a specific case for transfers of water to dry wells
	// nb for this to really work I think we still need to make sure well volumes
	// are being properly kept in sync

	rule := NewLHPolicyRule("BlowOutToEmptyWells")
	err := rule.AddNumericConditionOn("WELLTOVOLUME", 0.0, 1.0)

	if err != nil {
		return nil, err
	}

	err = rule.AddCategoryConditionOn("LIQUIDCLASS", "water")
	if err != nil {
		return nil, err
	}
	pol := MakeJBPolicy()
	lhpr.AddRule(rule, pol)

	// a further refinement: for low volumes we need to add extra volume
	// for aspirate and dispense
	/*
		rule = NewLHPolicyRule("ExtraVolumeForLV")
		rule.AddNumericConditionOn("VOLUME", 0.0, 20.0)
		pol = MakeLVExtraPolicy()
		lhpr.AddRule(rule, pol)
	*/

	// hack to fix plate type problems
	rule = NewLHPolicyRule("LVOffsetFix")
	rule.AddNumericConditionOn("VOLUME", 0.0, 20.0)
	rule.AddCategoryConditionOn("FROMPLATETYPE", "pcrplate_skirted_riser")
	pol = MakeLVOffsetPolicy()
	lhpr.AddRule(rule, pol)

	rule = NewLHPolicyRule("LVOffsetFix2")
	rule.AddNumericConditionOn("VOLUME", 0.0, 20.0)
	rule.AddCategoryConditionOn("TOPLATETYPE", "pcrplate_skirted_riser")
	pol = MakeLVOffsetPolicy()

	lhpr.AddRule(rule, pol)

	return lhpr, nil
}

func LoadLHPoliciesFromFile() (*LHPolicyRuleSet, error) {
	lhPoliciesFileName := os.Getenv("ANTHA_LHPOLICIES_FILE")
	if lhPoliciesFileName == "" {
		return nil, fmt.Errorf("Env variable ANTHA_LHPOLICIES_FILE not set")
	}
	contents, err := ioutil.ReadFile(lhPoliciesFileName)
	if err != nil {
		return nil, err
	}
	lhprs := NewLHPolicyRuleSet()
	lhprs.Policies = make(map[string]LHPolicy)
	lhprs.Rules = make(map[string]LHPolicyRule)
	//	err = readYAML(contents, lhprs)
	err = readJSON(contents, lhprs)
	if err != nil {
		return nil, err
	}
	return lhprs, nil
}

func readYAML(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := yaml.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}

func readJSON(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := json.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}
