package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"regexp"
	"strconv"
	"strings"
)

// convenience structure for handling well coordinates
type WellCoords struct {
	X int
	Y int
}

func ZeroWellCoords() WellCoords {
	return WellCoords{-1, -1}
}
func (wc WellCoords) IsZero() bool {
	if wc.Equals(ZeroWellCoords()) {
		return true
	}

	return false
}

func MatchString(s1, s2 string) bool {
	m, _ := regexp.MatchString(s1, s2)
	return m
}

func (wc WellCoords) Equals(w2 WellCoords) bool {
	if wc.X == w2.X && wc.Y == w2.Y {
		return true
	}

	return false
}

func MakeWellCoords(wc string) WellCoords {
	// try each one in turn

	r := MakeWellCoordsA1(wc)

	zero := WellCoords{-1, -1}

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoords1A(wc)

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoordsXY(wc)

	return r
}

// make well coordinates in the "A1" convention
func MakeWellCoordsA1(a1 string) WellCoords {
	if !MatchString("[A-Z]{1,}[0-9]{1,2}", a1) {
		return WellCoords{-1, -1}
	}
	re, _ := regexp.Compile("[A-Z]{1,}")
	ix := re.FindIndex([]byte(a1))
	endC := ix[1]

	X := wutil.ParseInt(a1[endC:len(a1)]) - 1
	Y := AlphaToNum(string(a1[0:endC])) - 1
	return WellCoords{X, Y}
}

// make well coordinates in the "1A" convention
func MakeWellCoords1A(a1 string) WellCoords {

	if !MatchString("[0-9]{1,2}[A-Z]{1,}", a1) {
		return WellCoords{-1, -1}
	}
	re, _ := regexp.Compile("[A-Z]{1,}")
	ix := re.FindIndex([]byte(a1))
	startC := ix[0]

	Y := AlphaToNum(string(a1[startC:len(a1)])) - 1
	X := wutil.ParseInt(a1[0:startC]) - 1
	return WellCoords{X, Y}
}

// make well coordinates in a manner compatble with "X1,Y1" etc.
func MakeWellCoordsXYsep(x, y string) WellCoords {
	r := WellCoords{wutil.ParseInt(y[1:len(y)]) - 1, wutil.ParseInt(x[1:len(x)]) - 1}

	if r.X < 0 || r.Y < 0 {
		return WellCoords{-1, -1}
	}

	return r
}

func MakeWellCoordsXY(xy string) WellCoords {
	tx := strings.Split(xy, "Y")
	if tx == nil || len(tx) != 2 || len(tx[0]) == 0 || len(tx[1]) == 0 {
		return WellCoords{-1, -1}
	}
	x := wutil.ParseInt(tx[0][1:len(tx[0])]) - 1
	y := wutil.ParseInt(tx[1]) - 1
	return WellCoords{x, y}
}

// return well coordinates in "X1Y1" format
func (wc WellCoords) FormatXY() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return "X" + strconv.Itoa(wc.X+1) + "Y" + strconv.Itoa(wc.Y+1)
}
func (wc WellCoords) Format1A() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X+1) + NumToAlpha(wc.Y+1)
}
func (wc WellCoords) FormatA1() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return NumToAlpha(wc.Y+1) + strconv.Itoa(wc.X+1)
}
func (wc WellCoords) WellNumber() int {
	if wc.X < 0 || wc.Y < 0 {
		return -1
	}
	return (8*(wc.X-1) + wc.Y)
}

func (wc WellCoords) ColNumString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X + 1)
}
func (wc WellCoords) RowLettString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return NumToAlpha(wc.Y + 1)
}

// comparison operators

func (wc WellCoords) RowLessThan(wc2 WellCoords) bool {
	if wc.Y == wc2.Y {
		return wc.X < wc2.Y
	}
	return wc.Y < wc2.Y
}

func (wc WellCoords) ColLessThan(wc2 WellCoords) bool {
	if wc.X == wc2.X {
		return wc.Y < wc2.Y
	}
	return wc.X < wc2.X
}

// convenience structure to allow sorting

type WellCoordArrayCol []WellCoords
type WellCoordArrayRow []WellCoords

func (wca WellCoordArrayCol) Len() int           { return len(wca) }
func (wca WellCoordArrayCol) Swap(i, j int)      { t := wca[i]; wca[i] = wca[j]; wca[j] = t }
func (wca WellCoordArrayCol) Less(i, j int) bool { return wca[i].RowLessThan(wca[j]) }

func (wca WellCoordArrayRow) Len() int           { return len(wca) }
func (wca WellCoordArrayRow) Swap(i, j int)      { t := wca[i]; wca[i] = wca[j]; wca[j] = t }
func (wca WellCoordArrayRow) Less(i, j int) bool { return wca[i].ColLessThan(wca[j]) }
