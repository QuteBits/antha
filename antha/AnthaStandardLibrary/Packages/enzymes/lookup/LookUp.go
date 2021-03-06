// Part of the Antha language
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

// Package for looking up restriction enzyme properties
package lookup

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/REBASE"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/asset"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func TypeIIsLookup(name string) (enzyme wtype.TypeIIs, err error) {
	enz := EnzymeLookup(name)

	enzyme, err = wtype.ToTypeIIs(enz)
	return
}

func EnzymeLookup(name string) (enzyme wtype.RestrictionEnzyme) {
	enzymes, err := asset.Asset("rebase/type2.txt")
	if err != nil {
		return
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range rebase.RebaseParse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if strings.ToUpper(record.Name) == strings.ToUpper(name) {
			fmt.Println(record)
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzyme = record
		}

	}
	return enzyme
}

func FindEnzymesofClass(class string) (enzymelist []wtype.RestrictionEnzyme) {
	enzymes, err := asset.Asset("rebase/type2.txt")
	if err != nil {
		return
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range rebase.RebaseParse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if strings.ToUpper(record.Class) == strings.ToUpper(class) {
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzymelist = append(enzymelist, record)
		}
	}
	return enzymelist
}

func FindEnzymeNamesofClass(class string) (enzymelist []string) {
	for _, enzyme := range FindEnzymesofClass(class) {
		enzymelist = append(enzymelist, enzyme.Name)
	}
	return
}
