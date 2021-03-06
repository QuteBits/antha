// wtype/wtype_test.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wtype

import (
	"fmt"
	"sort"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

/*
func testBS(bs BioSequence) {
	fmt.Println(bs.Sequence())
}

func TestOne(*testing.T) {
	dna := DNASequence{"test", "ACCACACATAGCTAGCTAGCTAG", false, false, Overhang{}, Overhang{}, ""}
	testBS(&dna)
}

func ExampleOne() {
	dna := DNASequence{"test", "ACCACACATAGCTAGCTAGCTAG", false, false, Overhang{}, Overhang{}, ""}
	testBS(&dna)
	// Output:
	// ACCACACATAGCTAGCTAGCTAG
}

func TestLocations(*testing.T) {
	nl := NewLocation("liquidhandler", 9, NewShape("box", "", 0, 0, 0))
	nl2 := NewLocation("anotherliquidhandler", 9, NewShape("box", "", 0, 0, 0))
	fmt.Println("Location ", nl.Location_Name(), " ", nl.Location_ID(), " and location ", nl2.Location_Name(), " ", nl2.Location_ID(), " are the same? ", SameLocation(nl, nl2, 0))

	fmt.Println("Location ", nl.Positions()[0].Location_Name(), " and location ", nl.Positions()[1].Location_Name(), " are the same? ", SameLocation(nl.Positions()[0], nl.Positions()[1], 0), " share a parent? ", SameLocation(nl.Positions()[0], nl.Positions()[1], 1))

	fmt.Println("Locations ", nl.Location_Name(), " and ", nl.Positions()[0].Location_Name(), " share a parent? ", SameLocation(nl, nl.Positions()[0], 1))
}

func TestWellCoords(*testing.T) {
	fmt.Println("Testing Well Coords")
	wc := MakeWellCoordsA1("A1")
	fmt.Println(wc.FormatA1())
	fmt.Println(wc.Format1A())
	fmt.Println(wc.FormatXY())
	fmt.Println(wc.X, " ", wc.Y)
	wc = MakeWellCoordsXYsep("X1", "Y1")
	fmt.Println(wc.FormatA1())
	fmt.Println(wc.Format1A())
	fmt.Println(wc.FormatXY())
	fmt.Println(wc.X, " ", wc.Y)
	wc = MakeWellCoordsXY("X1Y1")
	fmt.Println(wc.FormatA1())
	fmt.Println(wc.Format1A())
	fmt.Println(wc.FormatXY())
	fmt.Println(wc.X, " ", wc.Y)
	fmt.Println("Finished Testing Well Coords")
}
*/

func TestWellCoords(t *testing.T) {
	wc := MakeWellCoordsA1("A1")

	if wc.X != 0 || wc.Y != 0 {
		t.Fatal(fmt.Sprint("Well Coords A1 expected {0,0} got ", wc))
	}

	if wc.FormatA1() != "A1" {
		t.Fatal(fmt.Sprint("Well coords A1 expected formatA1 to return A1, instead got ", wc.FormatA1()))
	}

	if wc.Format1A() != "1A" {
		t.Fatal(fmt.Sprint("Well coords A1 expected format1A to return 1A, instead got ", wc.FormatA1()))
	}

	wc = MakeWellCoords1A("1A")
	if wc.X != 0 || wc.Y != 0 {
		t.Fatal(fmt.Sprint("Well Coords 1A expected {0,0} got ", wc))
	}

	wc = MakeWellCoordsA1("AA1")

	if wc.X != 0 || wc.Y != 26 {
		t.Fatal(fmt.Sprint("Well Coords AA1 expected {0,26} got ", wc))
	}

	wc = MakeWellCoords1A("1AA")
	if wc.X != 0 || wc.Y != 26 {
		t.Fatal(fmt.Sprint("Well Coords 1AA expected {0,26} got ", wc))
	}

	wc = MakeWellCoordsA1("AAA1")

	if wc.X != 0 || wc.Y != 702 {
		t.Fatal(fmt.Sprint("Well Coords AAA1 expected {0,702} got ", wc))
	}

	wc = MakeWellCoords1A("1AAA")
	if wc.X != 0 || wc.Y != 702 {
		t.Fatal(fmt.Sprint("Well Coords AAA1 expected {0,702} got ", wc))
	}

}

func TestLHComponentSampleStuff(t *testing.T) {
	var c LHComponent

	faux := c.IsSample()

	if faux {
		t.Fatal("IsSample() must return false on new components")
	}

	c.SetSample(true)

	vrai := c.IsSample()

	if !vrai {
		t.Fatal("IsSample() must return true following SetIsSample(true)")
	}

	c.SetSample(false)

	faux = c.IsSample()

	if faux {
		t.Fatal("IsSample() must return false following SetIsSample(false)")
	}

	// now the same from NewLHComponent

	c2 := NewLHComponent()

	faux = c2.IsSample()

	if faux {
		t.Fatal("IsSample() must return false on new components")
	}

	c2.SetSample(true)

	vrai = c2.IsSample()

	if !vrai {
		t.Fatal("IsSample() must return true following SetIsSample(true)")
	}

	c2.SetSample(false)

	faux = c2.IsSample()

	if faux {
		t.Fatal("IsSample() must return false following SetIsSample(false)")
	}

	// finally need to make sure sample works
	// grrr import cycle not allowed: honestly I think Sample should just be an
	// instance method of LHComponent now anyway
	/*

		c.CName = "YOMAMMA"
		s := mixer.Sample(c, wunit.NewVolume(10.0, "ul"))

		vrai = s.IsSample()

		if !vrai {
			t.Fatal("IsSample() must return true for results of mixer.Sample()")
		}
		s = mixer.SampleForConcentration(c, wunit.NewConcentration(10.0, "mol/l"))

		vrai = s.IsSample()

		if !vrai {
			t.Fatal("IsSample() must return true for results of mixer.SampleForConcentration()")
		}
		s = mixer.SampleForTotalVolume(c, wunit.NewVolume(10.0, "ul"))

		vrai = s.IsSample()

		if !vrai {
			t.Fatal("IsSample() must return true for results of mixer.SampleForTotalVolume()")
		}
	*/
}

type testpair struct {
	ltstring string
	ltint    int
	err      error
}

var lts []testpair = []testpair{testpair{ltstring: "170516CCFDesign_noTouchoff_noBlowout2", ltint: 102}, testpair{ltstring: "190516OnePolicy0", ltint: 3000}, testpair{ltstring: "dna_mix", ltint: LTDNAMIX}, testpair{ltstring: "PreMix", ltint: LTPreMix} /*testpair{ltstring: "InvalidEntry", ltint: LTWater, err: fmt.Errorf("!")}*/}

func TestLiquidTypeFromString(t *testing.T) {

	for _, lt := range lts {

		ltnum, err := LiquidTypeFromString(lt.ltstring)
		if int(ltnum) != lt.ltint {
			t.Error("running LiquidTypeFromString on ", lt.ltstring, "expected", lt.ltint, "got", ltnum)
		}
		if err != nil {
			if err != lt.err {
				t.Error("running LiquidTypeFromString on ", lt.ltstring, "expected err:", lt.err.Error(), "got", err.Error())
			}
		}
	}
}

func TestLiquidTypeName(t *testing.T) {

	for _, lt := range lts {

		ltstr := LiquidTypeName(LiquidType(lt.ltint))
		if ltstr != lt.ltstring {
			t.Error("running LiquidTypeName on ", lt.ltint, "expected", lt.ltstring, "got", ltstr)
		}

	}
}

func TestParent(t *testing.T) {
	c := NewLHComponent()

	c.AddParent("A")
	c.AddParent("B")
	c.AddParent("C")

	vrai := c.HasParent("A")

	if !vrai {
		t.Error("LHComponent.HasParent() must return true for values set with AddParent")
	}

	vrai = c.HasParent("B")

	if !vrai {
		t.Error("LHComponent.HasParent() must return true for values set with AddParent")
	}

	faux := c.HasParent("D")

	if faux {
		t.Error("LHComponent.HasParent() must return false for values not set")
	}

}

func testLHCP() LHChannelParameter {
	return LHChannelParameter{
		ID:          "dummydummy",
		Name:        "mrdummy",
		Minvol:      wunit.NewVolume(1.0, "ul"),
		Maxvol:      wunit.NewVolume(1.0, "ul"),
		Minspd:      wunit.NewFlowRate(0.5, "ml/min"),
		Maxspd:      wunit.NewFlowRate(0.6, "ml/min"),
		Multi:       8,
		Independent: false,
		Orientation: LHVChannel,
		Head:        0,
	}
}

func TestLHMultiConstraint(t *testing.T) {
	params := testLHCP()

	cnst := params.GetConstraint(8)

	expected := LHMultiChannelConstraint{0, 1, 8}

	if !cnst.Equals(expected) {
		t.Fatal(fmt.Sprint("Expected: ", expected, " GOT: ", cnst))
	}

}

func TestWCSorting(t *testing.T) {
	v := make([]WellCoords, 0, 1)

	v = append(v, WellCoords{0, 2})
	v = append(v, WellCoords{4, 2})
	v = append(v, WellCoords{0, 1})
	v = append(v, WellCoords{8, 9})
	v = append(v, WellCoords{1, 3})
	v = append(v, WellCoords{3, 6})
	v = append(v, WellCoords{8, 0})

	sort.Sort(WellCoordArrayRow(v))

	if v[0].FormatA1() != "B1" {
		t.Fatal(fmt.Sprint("Row-first sort incorrect: expected B1 first, got ", v[0].FormatA1()))
	}

	sort.Sort(WellCoordArrayCol(v))

	if v[0].FormatA1() != "A9" {
		t.Fatal("Col-first sort incorrect: expected A9 first, got ", v[0].FormatA1())
	}
}
