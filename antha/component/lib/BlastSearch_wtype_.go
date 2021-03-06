// Example element demonstrating how to perform a BLAST search using the megablast algorithm

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/blast"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	biogo "github.com/biogo/ncbi/blast"
)

// Input parameters for this protocol

// Data which is returned from this protocol; output data

// Physical inputs to this protocol

// Physical outputs from this protocol

func _BlastSearch_wtypeRequirements() {

}

// Actions to perform before protocol itself
func _BlastSearch_wtypeSetup(_ctx context.Context, _input *BlastSearch_wtypeInput) {

}

// Core process of the protocol: steps to be performed for each input
func _BlastSearch_wtypeSteps(_ctx context.Context, _input *BlastSearch_wtypeInput, _output *BlastSearch_wtypeOutput) {

	var err error
	var hits []biogo.Hit
	/*
		if Querytype == "PROTEIN" {
		hits, err = blast.MegaBlastP(Query)
		if err != nil {
			fmt.Println(err.Error())
		}

		Hits = fmt.Sprintln(blast.HitSummary(hits))


		} else if Querytype == "DNA" {
		hits, err = blast.MegaBlastN(Query)
		if err != nil {
			fmt.Println(err.Error())
		}

		Hits = fmt.Sprintln(blast.HitSummary(hits))
		}
	*/
	_output.AnthaSeq = _input.DNA

	// look for orfs
	orf, orftrue := sequences.FindORF(_output.AnthaSeq.Seq)

	if orftrue == true && len(orf.DNASeq) == len(_output.AnthaSeq.Seq) {
		// if open reading frame is detected, we'll perform a blastP search'
		fmt.Println("ORF detected:", "full sequence length: ", len(_output.AnthaSeq.Seq), "ORF length: ", len(orf.DNASeq))
		hits, err = blast.MegaBlastP(orf.ProtSeq)
	} else {
		// otherwise we'll blast the nucleotide sequence
		hits, err = _output.AnthaSeq.Blast()
	}
	if err != nil {
		fmt.Println(err.Error())

	} //else {

	_output.Hits = fmt.Sprintln(blast.HitSummary(hits))

	// Rename Sequence with ID of top blast hit
	_output.AnthaSeq.Nm = hits[0].Id
	//}

}

// Actions to perform after steps block to analyze data
func _BlastSearch_wtypeAnalysis(_ctx context.Context, _input *BlastSearch_wtypeInput, _output *BlastSearch_wtypeOutput) {

}

func _BlastSearch_wtypeValidation(_ctx context.Context, _input *BlastSearch_wtypeInput, _output *BlastSearch_wtypeOutput) {

}
func _BlastSearch_wtypeRun(_ctx context.Context, input *BlastSearch_wtypeInput) *BlastSearch_wtypeOutput {
	output := &BlastSearch_wtypeOutput{}
	_BlastSearch_wtypeSetup(_ctx, input)
	_BlastSearch_wtypeSteps(_ctx, input, output)
	_BlastSearch_wtypeAnalysis(_ctx, input, output)
	_BlastSearch_wtypeValidation(_ctx, input, output)
	return output
}

func BlastSearch_wtypeRunSteps(_ctx context.Context, input *BlastSearch_wtypeInput) *BlastSearch_wtypeSOutput {
	soutput := &BlastSearch_wtypeSOutput{}
	output := _BlastSearch_wtypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func BlastSearch_wtypeNew() interface{} {
	return &BlastSearch_wtypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &BlastSearch_wtypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _BlastSearch_wtypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &BlastSearch_wtypeInput{},
			Out: &BlastSearch_wtypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type BlastSearch_wtypeElement struct {
	inject.CheckedRunner
}

type BlastSearch_wtypeInput struct {
	DNA wtype.DNASequence
}

type BlastSearch_wtypeOutput struct {
	AnthaSeq wtype.DNASequence
	Hits     string
}

type BlastSearch_wtypeSOutput struct {
	Data struct {
		AnthaSeq wtype.DNASequence
		Hits     string
	}
	Outputs struct {
	}
}

func init() {
	addComponent(Component{Name: "BlastSearch_wtype",
		Constructor: BlastSearch_wtypeNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/BlastSearch/BlastSearch_wtype.an",
			Params: []ParamDesc{
				{Name: "DNA", Desc: "", Kind: "Parameters"},
				{Name: "AnthaSeq", Desc: "", Kind: "Data"},
				{Name: "Hits", Desc: "", Kind: "Data"},
			},
		},
	})
}
