package plot

import (
	"fmt"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/internal/code.google.com/p/plotinum/plot"
	"github.com/antha-lang/antha/internal/code.google.com/p/plotinum/plotter"
	"github.com/antha-lang/antha/internal/code.google.com/p/plotinum/plotutil"
	"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
)

var (
	alphabet string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func Export(plt *plot.Plot, filename string) {
	plt.Save(4, 4, filename)
}

func Plot(Xvalues []float64, Yvaluearray [][]float64) (plt *plot.Plot) {
	// now plot the graph

	// the data points
	pts := make([]plotter.XYer, 0) //len(Xdatarange))

	for ptsindex := 0; ptsindex < len(Xvalues); ptsindex++ {

		// each specific set for each datapoint

		for Xdatarangeindex, xfloat := range Xvalues {

			xys := make(plotter.XYs, len(Yvaluearray))

			yfloats := make([]float64, 0)
			for _, yvalues := range Yvaluearray {
				yfloat := yvalues[Xdatarangeindex]
				yfloats = append(yfloats, yfloat)
			}

			for j := range xys {
				xys[j].X = xfloat
				xys[j].Y = yfloats[j]

			}
			//fmt.Println(xys)
			pts = append(pts, xys) //
		}

	}
	plt, err := plot.New()

	if err != nil {
		panic(err)
	}

	// Create two lines connecting points and error bars. For
	// the first, each point is the mean x and y value and the
	// error bars give the 95% confidence intervals.  For the
	// second, each point is the median x and y value with the
	// error bars showing the minimum and maximum values.
	/*
	   	fmt.Println("pts", pts)
	   	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	   	if err != nil {
	   		panic(err)
	   	}
	   	//medMinMax, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	   //	if err != nil {
	   //		panic(err)
	   //	}
	   	plotutil.AddLinePoints(plt,
	   		"mean and 95% confidence", mean95,
	   	) //	"median and minimum and maximum", medMinMax)
	   	//plotutil.AddErrorBars(plt, mean95, medMinMax)

	   	// Add the points that are summarized by the error points.


	*/

	fmt.Println(len(pts))
	plotutil.AddScattersXYer(plt, pts)
	return
}

func PlotfromMinMaxpairs(sheet *xlsx.Sheet, Xminmax []string, Yminmaxarray [][]string, Exportedfilename string) {
	Xdatarange, err := spreadsheet.ConvertMinMaxtoArray(Xminmax)
	if err != nil {
		fmt.Println(Xminmax, Xdatarange)
		panic(err)
	}
	fmt.Println(Xdatarange)

	Ydatarangearray := make([][]string, 0)
	for i, Yminmax := range Yminmaxarray {
		Ydatarange, err := spreadsheet.ConvertMinMaxtoArray(Yminmax)
		if err != nil {
			panic(err)
		}
		if len(Xdatarange) != len(Ydatarange) {
			panicmessage := fmt.Errorf("for index", i, "of array", "len(Xdatarange) != len(Ydatarange)")
			panic(panicmessage.Error())
		}
		Ydatarangearray = append(Ydatarangearray, Ydatarange)
		fmt.Println(Ydatarange)
	}
	Plotfromspreadsheet(sheet, Xdatarange, Ydatarangearray, Exportedfilename)
}

// Xdatarange would consist of an array of
func Plotfromspreadsheet(sheet *xlsx.Sheet, Xdatarange []string, Ydatarangearray [][]string, Exportedfilename string) {

	// now plot the graph

	// the data points
	pts := make([]plotter.XYer, 0) //len(Xdatarange))

	for ptsindex := 0; ptsindex < len(Xdatarange); ptsindex++ {

		// each specific set for each datapoint

		for Xdatarangeindex, Xdatapoint := range Xdatarange {

			xys := make(plotter.XYs, len(Ydatarangearray))

			fmt.Println("going here3")
			fmt.Println("Xdatapoint= ", Xdatapoint)

			xrow, xcol, err := spreadsheet.A1formattorowcolumn(Xdatapoint)
			if err != nil {
				panic(err)
			}
			fmt.Println("row, col line 155:", xrow, xcol)
			xpoint := sheet.Rows[xcol].Cells[xrow]
			fmt.Println("datapoint", Xdatarangeindex, Xdatapoint, "xpoint = ", xpoint)

			// get each y point and work out average

			//yvalues := make([]xlsx.Cell, 0)
			yfloats := make([]float64, 0)
			for ydatarangearrayindex, Ydatarange := range Ydatarangearray {
				yrow, ycol, err := spreadsheet.A1formattorowcolumn(Ydatarange[Xdatarangeindex])
				if err != nil {
					panic(err)
				}
				fmt.Println("row, col line 148:", yrow, ycol)
				//ypoint := sheet.Cell(yrow, ycol)
				ypoint := sheet.Rows[ycol].Cells[yrow]
				//yvalues = append(yvalues, ypoint)
				yfloat, err := ypoint.Float()
				if err != nil {
					panic(err)
				}
				yfloats = append(yfloats, yfloat)
				fmt.Println("datapoint", Xdatarangeindex, ydatarangearrayindex, Ydatarange[ydatarangearrayindex], "Ycol", ycol, "yrow", yrow, "ypoint = ", ypoint)

				//n, m := 5, 10
				//pts := make([]plotter.XYer, len(Xdatarange))
				//for i := range pts {

				//pts[i] =

			}

			/*ymean, err := stats.Mean(yfloats)
			if err != nil {
				panic(err)
			}*/
			xfloat, err := xpoint.Float()
			if err != nil {
				panic(err)
			}

			//type XYs []struct{ X, Y float64 }
			//pts[ptsindex] = &MyXYs{xfloat, ymean}

			if err != nil {
				panic(err)
			}

			fmt.Println("datapoint", Xdatarangeindex, Xdatapoint, "xpoint = ", xpoint)
			if err != nil {
				panic(err)
			}
			fmt.Println("going here2")
			//center := float64(i)
			for j := range xys {
				fmt.Println("going here")
				fmt.Println(ptsindex)
				//x, _ := pts[l].XY(j)
				xys[j].X = xfloat
				fmt.Println("xfloat", j, xfloat)
				xys[j].Y = yfloats[j]
				fmt.Println("yfloats[j]", j, yfloats[j])
				fmt.Println("xysssssssssx", Xdatarangeindex, j, xys)
			}
			fmt.Println(xys)
			pts = append(pts, xys) //
			//pts[Xdatarangeindex] = xys
			fmt.Println("hello:", pts[Xdatarangeindex])
			fmt.Println("hello again", pts)
		}

		fmt.Println("len(pts)", len(pts))
	}
	plt, err := plot.New()

	if err != nil {
		panic(err)
	}

	// Create two lines connecting points and error bars. For
	// the first, each point is the mean x and y value and the
	// error bars give the 95% confidence intervals.  For the
	// second, each point is the median x and y value with the
	// error bars showing the minimum and maximum values.

	fmt.Println("pts", pts)
	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	if err != nil {
		panic(err)
	}
	/*medMinMax, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	if err != nil {
		panic(err)
	}*/
	plotutil.AddLinePoints(plt,
		"mean and 95% confidence", mean95,
	) //	"median and minimum and maximum", medMinMax)
	//plotutil.AddErrorBars(plt, mean95, medMinMax)

	// Add the points that are summarized by the error points.
	fmt.Println(len(pts))
	plotutil.AddScatters(plt, pts[0], pts[1], pts[2], pts[3], pts[4])

	plt.Save(4, 4, Exportedfilename)

}
