Work your way through the following folders in order.

Each shows two key files, both in json format:

1. A worflow definition file
2. A parameters file 

Workflow:
The workflow file specifies a set of Processes which call Antha elements 
(components) which are to be run. 
This could be: 
(A) a single element run once 
(B) parallel copies of a single element run in parallel, for example multiple runs of the same protocol for different samples or with different conditions.
(C) multiple different elements run at the same time
(C) multiple elements which may be connected; i.e. one or more outputs (ports) from a source element (src) may feed in as inputs (also ports) into the downstream target element (tgt).


Parameters:
The parameters file assigns parameters for each of the processes specified in the workflow file

i.e. the parameters file is used to set the values for the input parameters.

The example parameters files in these folders show how to set variables specified in the parameters file to the actual values we want to assign to them.
One of the key variables you'll likely want to set are the liquid handling components (wtype.LHComponent) 


LHComponent:

One of the key antha types which will typically be specified in the parameters file is the wtype.LHComponent

LHComponents can be accessed in the parameters.yml file in the following way:

These are written as a string: e.g. 
"Diluent":water,
“dnastock”:”gfpstock”,

Before a component can be used, currently, the concept of that component needs to be added to the factory.
i.e. When we say the concept of a component we don't mean a specific sample of water, which would be called from an inventory instead, but any sample of water, i.e. which has the liquidhandling properties of water.

The factory is located in the following path:

$GOPATH/src/github.com/antha-lang/antha/microArch/factory/make_component_library.go

Open the file and add the component to the list within the body of the func makeComponentLibrary()

e.g.

A = wtype.NewLHComponent()
	A.CName = "tartrazine"
	A.Type = wtype.LTWater // or could use wtype.LiquidTypeFromString("water")
	A.Smax = 9999
	cmap[A.CName] = A
therefore a new component would be specified as follows:

A = wtype.NewLHComponent()
    A.CName = "mynewviscouscomponent"
    A.Type = wtype.LTVISCOUS
    A.Smax = 9999
    cmap[A.CName] = A


LiquidTypes:
	
You may want to change the .Type to something else as this will determine how the liquid type is pipetted. 
Currently this consists of:

	LTWater
	LTGlycerol
	LTEthanol
	LTDetergent
	LTCulture
	LTProtein
	LTDNA
	LTload
	LTDoNotMix
	LTloadwater
	LTNeedToMix
	LTPostMix
	LTPreMix
	LTVISCOUS
	LTPAINT
	LTDISPENSEABOVE
	LTPEG
	LTProtoplasts
	LTCulutureReuse
	LTDNAMIX
	
The full list can be found in

$GOPATH/src/github.com/antha-lang/antha/antha/anthalib/wtype/LiquidType.go
The details of these policies can be found in

$GOPATH/src/github.com/antha-lang/antha/microArch/driver/liquidhandling/makelhpolicy.go
