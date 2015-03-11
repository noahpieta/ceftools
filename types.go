package ceftools

type Attribute struct {
	Name   string
	Values []string
}

type Header struct {
	Name  string
	Value string
}

const (
	Transposed = 1 << iota
)

const (
	MajorVersion = 0
	MinorVersion = 1
	MagicCEB     = 0x09424543
	MagicCEF     = 0x09464543
)

type Cef struct {
	MajorVersion     int32
	MinorVersion     int32
	NumRows          int64
	NumColumns       int64
	Headers          []Header
	Flags            int64
	RowAttributes    []Attribute
	ColumnAttributes []Attribute
	Matrix           []float32
}

func (cef Cef) Get(col int64, row int64) float32 {
	return cef.Matrix[col+row*cef.NumColumns]
}

func (cef Cef) Set(col int64, row int64, val float32) {
	cef.Matrix[col+row*cef.NumColumns] = val
}