package intmat

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type BigIntMatrix struct {
	rowValues map[int]map[int]*big.Int //hold rowValues for (X,Y)
	colValues map[int]map[int]*big.Int //easy access to (Y,X)
	rows      int                      // total number rows available to this matrix
	rowStart  int                      // [rowStart,rowEnd)
	cols      int                      // total number cols available to this matrix
	colStart  int                      // [colStart,colEnd)

}

type bigintmatrix struct {
	RowValues map[int]map[int]*big.Int //hold rowValues for (X,Y)
	ColValues map[int]map[int]*big.Int //easy access to (Y,X)
	Rows      int                      // total number rows available to this matrix
	RowStart  int                      // [rowStart,rowEnd)
	Cols      int                      // total number cols available to this matrix
	ColStart  int                      // [colStart,colEnd)

}

func (mat *BigIntMatrix) MarshalJSON() ([]byte, error) {
	return json.Marshal(bigintmatrix{
		RowValues: mat.rowValues,
		ColValues: mat.colValues,
		Rows:      mat.rows,
		RowStart:  mat.rowStart,
		Cols:      mat.cols,
		ColStart:  mat.colStart,
	})
}

func (mat *BigIntMatrix) UnmarshalJSON(bytes []byte) error {
	var m bigintmatrix
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}
	mat.rowValues = m.RowValues
	mat.colValues = m.ColValues
	mat.rows = m.Rows
	mat.rowStart = m.RowStart
	mat.cols = m.Cols
	mat.colStart = m.ColStart
	return nil
}

// NewMat creates a new matrix with the specified number of rows and cols.
// If values is empty, the matrix will be zeroized.
// If values are not empty it must have rows*cols items.  The values are expected to
// be big.NewInt(0)/nil or big.NewInt(1) for logical operations, but can be any *big.Int for arithmetic.
// Note value refs passed in will NOT be copied into new big.ints
func NewBigIntMat(rows, cols int, values ...*big.Int) *BigIntMatrix {
	if len(values) != 0 && len(values) != rows*cols {
		panic(fmt.Sprintf("matrix data length (%v) to size mismatch expected %v", len(values), rows*cols))
	}

	mat := BigIntMatrix{
		rowValues: map[int]map[int]*big.Int{},
		colValues: map[int]map[int]*big.Int{},
		rows:      rows,
		rowStart:  0,
		cols:      cols,
		colStart:  0,
	}

	if len(values) > 0 {
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				index := i*cols + j
				if values[index] != nil {
					mat.set(i, j, values[index])
				}
			}
		}
	}

	return &mat
}

func NewBigIntMatFromVec(vec *BigIntVector) *BigIntMatrix {
	return BigIntCopy(vec.mat)
}

// Identity create an identity matrix (one's on the diagonal).
func BigIntIdentity(size int) *BigIntMatrix {
	mat := BigIntMatrix{
		rowValues: map[int]map[int]*big.Int{},
		colValues: map[int]map[int]*big.Int{},
		rows:      size,
		rowStart:  0,
		cols:      size,
		colStart:  0,
	}

	for i := 0; i < size; i++ {
		mat.set(i, i, big.NewInt(1))
	}

	return &mat
}

// Copy will create a NEW matrix that will have all the same values as m.
func BigIntCopy(m *BigIntMatrix) *BigIntMatrix {
	mat := BigIntMatrix{
		rowValues: make(map[int]map[int]*big.Int),
		colValues: make(map[int]map[int]*big.Int),
		rows:      m.rows,
		rowStart:  0,
		cols:      m.cols,
		colStart:  0,
	}

	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.set(i, j, m.At(i, j))
		}
	}

	return &mat
}

// Slice creates a slice of the matrix.  The slice will be connected to the original matrix, changes to one
// causes changes in the other.
func (mat *BigIntMatrix) Slice(i, j, rows, cols int) *BigIntMatrix {
	if rows <= 0 || cols <= 0 {
		panic("slice rows and cols must >= 1")
	}

	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	r := i + mat.rowStart
	c := j + mat.colStart

	if r+rows-1 > mat.rows || c+cols-1 > mat.cols {
		panic("slice rows and cols must be in bounds of matrix")
	}
	mat.checkRowBounds(i + rows - 1)
	mat.checkColBounds(j + cols - 1)

	return mat.slice(r, c, rows, cols)
}

func (mat *BigIntMatrix) slice(r, c, rows, cols int) *BigIntMatrix {
	return &BigIntMatrix{
		rowValues: mat.rowValues,
		rows:      rows,
		rowStart:  r,
		colValues: mat.colValues,
		cols:      cols,
		colStart:  c,
	}
}

func (mat *BigIntMatrix) checkRowBounds(i int) {
	if i < 0 || i >= mat.rows {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, mat.rows-1))
	}
}

func (mat *BigIntMatrix) checkColBounds(j int) {
	if j < 0 || j >= mat.cols {
		panic(fmt.Sprintf("%v out of range: [0-%v]", j, mat.cols-1))
	}
}

// Dims returns the dimensions of the matrix.
func (mat *BigIntMatrix) Dims() (int, int) {
	return mat.rows, mat.cols
}

// At returns the value at row index i and column index j.
func (mat *BigIntMatrix) At(i, j int) *big.Int {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	r := i + mat.rowStart
	c := j + mat.colStart

	ret := mat.at(r, c)
	if ret == nil {
		return big.NewInt(0)
	}
	return ret
}

func (mat *BigIntMatrix) at(r, c int) *big.Int {
	ys, ok := mat.rowValues[r]
	if !ok {
		return nil
	}
	v, ok := ys[c]
	if !ok {
		return nil
	}
	return v
}

// Set sets the value at row index i and column index j to value.
func (mat *BigIntMatrix) Set(i, j int, value *big.Int) {
	mat.checkRowBounds(i)
	mat.checkColBounds(j)
	r := i + mat.rowStart
	c := j + mat.colStart

	mat.set(r, c, value)
}

func (mat *BigIntMatrix) set(r, c int, value *big.Int) {

	// Treat nil as zero for convenience
	if value == nil || value.Sign() == 0 {
		ys, ok := mat.rowValues[r]
		if !ok {
			return
		}

		_, ok = ys[c]
		if !ok {
			return
		}

		delete(ys, c)
		if len(mat.rowValues[r]) == 0 {
			delete(mat.rowValues, r)
		}

		delete(mat.colValues[c], r)
		if len(mat.colValues[c]) == 0 {
			delete(mat.colValues, c)
		}

		return
	}

	ys, ok := mat.rowValues[r]
	if !ok {
		ys = make(map[int]*big.Int)
		mat.rowValues[r] = ys
	}
	ys[c] = value

	xs, ok := mat.colValues[c]
	if !ok {
		xs = make(map[int]*big.Int)
		mat.colValues[c] = xs
	}
	xs[r] = value
}

// T returns a matrix that is the transpose of the underlying matrix. Note the transpose
// is connected to matrix it is a transpose of, and changes made to one affect the other.
func (mat *BigIntMatrix) T() *BigIntMatrix {
	return &BigIntMatrix{
		rowValues: mat.colValues,
		rows:      mat.cols,
		rowStart:  mat.colStart,
		colValues: mat.rowValues,
		cols:      mat.rows,
		colStart:  mat.rowStart,
	}
}

// Zeroize take the current matrix sets all values to 0.
func (mat *BigIntMatrix) Zeroize() {
	mat.zeroize(mat.rowStart, mat.colStart, mat.rows, mat.cols)
}

// ZeroizeRange take the current matrix sets values inside the range to zero.
func (mat *BigIntMatrix) ZeroizeRange(i, j, rows, cols int) {
	if i < 0 || j < 0 || rows < 0 || cols < 0 {
		panic("zeroize must have positive values")
	}
	if mat.rows < i+rows || mat.cols < j+cols {
		panic(fmt.Sprintf("zeroize bounds check failed can't zeroize shape (%v,%v) on a (%v,%v) matrix", i+rows, j+cols, mat.rows, mat.cols))
	}

	r := i + mat.rowStart
	c := j + mat.colStart

	mat.zeroize(r, c, rows, cols)
}

func (mat *BigIntMatrix) zeroize(r, c, rows, col int) {
	for rv, cs := range mat.rowValues {
		if rv < r || r+rows <= rv {
			continue
		}
		for cv, _ := range cs {
			if cv < c || c+col <= cv {
				continue
			}
			mat.set(rv, cv, nil)
		}
	}
}

// Pow raises the matrix to the power of k using exponentiation by squaring.
// The matrix must be square.
func (mat *BigIntMatrix) Pow(k int) *BigIntMatrix {
	if mat.rows != mat.cols {
		panic(fmt.Sprintf("matrix must be square to raise to a power, got %dx%d", mat.rows, mat.cols))
	}

	if k < 0 {
		panic("power k must be non-negative")
	}

	if k == 0 {
		return BigIntIdentity(mat.rows)
	}

	result := BigIntIdentity(mat.rows)
	currentPower := BigIntCopy(mat) // Use a copy to avoid modifying the original matrix if k=1

	for k > 0 {
		if k%2 == 1 {
			temp := NewBigIntMat(mat.rows, mat.cols)
			temp.Mul(result, currentPower)
			result = temp
		}
		k /= 2
		if k > 0 { // Avoid unnecessary multiplication if k becomes 0
			temp := NewBigIntMat(mat.rows, mat.cols)
			temp.Mul(currentPower, currentPower)
			currentPower = temp
		}
	}
	return result
}

// Mul multiplies two matrices and stores the values in this matrix.
func (mat *BigIntMatrix) Mul(a, b *BigIntMatrix) {
	if a == nil || b == nil {
		panic("multiply input was found to be nil")
	}

	if mat == a || mat == b {
		panic("multiply self assignment not allowed")
	}

	if a.cols != b.rows {
		panic(fmt.Sprintf("multiply shape misalignment can't multiply (%v,%v)x(%v,%v)", a.rows, a.cols, b.rows, b.cols))
	}

	mRows, mCols := mat.Dims()
	aRows, _ := a.Dims()
	_, bCols := b.Dims()
	if mRows != aRows || mCols != bCols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, a.rows, b.cols))
	}

	mat.mul(a, b)
}

func (mat *BigIntMatrix) mul(a, b *BigIntMatrix) {
	//first we need to clear mat
	mat.zeroize(mat.rowStart, mat.colStart, mat.rows, mat.cols)

	for r, cs := range a.rowValues {
		if r < a.rowStart || a.rowStart+a.rows <= r {
			continue
		}
		i := r - a.rowStart

		for c, rs := range b.colValues {
			if c < b.colStart || b.colStart+b.cols <= c {
				continue
			}
			j := c - b.colStart
			value := big.NewInt(0)
			for ics, v1 := range cs {
				ci := ics - a.colStart

				v2, ok := rs[ci+b.rowStart]
				if ok {
					prod := new(big.Int).Mul(v1, v2)
					value.Add(value, prod)
				}
			}

			mat.Set(i, j, value)
		}
	}
}

// Add stores the addition of a and b in this matrix.
func (mat *BigIntMatrix) Add(a, b *BigIntMatrix) {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}
	if mat == a || mat == b {
		panic("addition self assignment not allowed")
	}

	if a.rows != b.rows || a.cols != b.cols {
		panic(fmt.Sprintf("addition input mat shapes do not match a=(%v,%v) b=(%v,%v)", a.rows, a.cols, b.rows, b.cols))
	}
	if mat.rows != a.rows || mat.cols != a.cols {
		panic(fmt.Sprintf("mat shape (%v,%v) does not match expected (%v,%v)", mat.rows, mat.cols, a.rows, a.cols))
	}

	mat.add(a, b)
}

func (mat *BigIntMatrix) add(a, b *BigIntMatrix) {
	//first we need to clear mat
	mat.setMatrix(a, mat.rowStart, mat.colStart)

	for r, cs := range b.rowValues {
		i := r - b.rowStart
		mr := i + mat.rowStart
		for c, v := range cs {
			j := c - b.colStart
			mc := j + mat.colStart
			currentVal := mat.at(mr, mc)
			if currentVal == nil {
				mat.set(mr, mc, new(big.Int).Set(v))
			} else {
				mat.set(mr, mc, new(big.Int).Add(currentVal, v))
			}
		}
	}
}

// Column returns a map containing the non zero row indices as the keys and it's associated values.
func (mat *BigIntMatrix) Column(j int) *TransposedBigIntVector {
	mat.checkColBounds(j)

	return &TransposedBigIntVector{
		mat: mat.Slice(0, j, mat.rows, 1),
	}
}

// SetColumn sets the values in column j. The values' keys are expected to be row indices.
func (mat *BigIntMatrix) SetColumn(j int, vec *TransposedBigIntVector) {
	mat.checkColBounds(j)

	if mat.rows != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	c := j + mat.colStart

	//first we'll zeroize
	rs := mat.colValues[c]
	for r := range rs {
		mat.set(r, c, nil)
	}

	//now set the new values
	for i, v := range vec.mat.colValues[vec.mat.colStart] {
		r := i + mat.rowStart
		mat.set(r, c, v)
	}
}

// Row returns a map containing the non zero column indices as the keys and it's associated values.
func (mat *BigIntMatrix) Row(i int) *BigIntVector {
	mat.checkRowBounds(i)

	return &BigIntVector{
		mat: mat.Slice(i, 0, 1, mat.cols),
	}
}

// SetRow sets the values in row i. The values' keys are expected to be column indices.
func (mat *BigIntMatrix) SetRow(i int, vec *BigIntVector) {
	mat.checkRowBounds(i)

	if mat.cols != vec.Len() {
		panic("matrix number of columns must equal length of vector")
	}

	r := i + mat.rowStart

	//first we'll zeroize
	cs := mat.rowValues[r]
	for c := range cs {
		mat.set(r, c, nil)
	}

	//now set the new values
	for j, v := range vec.mat.rowValues[vec.mat.rowStart] {
		c := j + mat.colStart
		mat.set(r, c, v)
	}
}

// Equals return true if the m matrix has the same shape and values as this matrix.
func (mat *BigIntMatrix) Equals(m *BigIntMatrix) bool {
	if mat == m {
		return true
	}

	if mat == nil || m == nil {
		return false
	}

	if mat.rows != m.rows || mat.cols != m.cols {
		return false
	}

	for i := 0; i < mat.rows; i++ {
		r := i + mat.rowStart
		cs, ok1 := mat.rowValues[r]
		ar := i + m.rowStart
		acs, ok2 := m.rowValues[ar]

		if !ok1 && !ok2 {
			continue
		}

		for j := 0; j < mat.cols; j++ {
			c := j + mat.colStart
			v1, ok1Inner := cs[c]
			ac := j + m.colStart
			v2, ok2Inner := acs[ac]

			// A missing value mean a value of zero
			// If both are missing then they are equal
			if !ok1Inner && !ok2Inner {
				continue
			}

			// if one is missing
			if ok1Inner != ok2Inner {
				return false
			}

			if v1.Cmp(v2) != 0 {
				return false
			}
		}
	}
	return true
}

// String returns a string representation of this matrix.
func (mat BigIntMatrix) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	for i := 0; i < mat.rows; i++ {
		row := make([]string, mat.cols)
		for j := 0; j < mat.cols; j++ {
			row[j] = fmt.Sprint(mat.At(i, j)) // mat.At(i,j) now returns *big.Int, fmt.Sprint handles it.
		}
		table.Append(row)
	}

	table.Render()
	return buff.String()
}

// SetMatrix replaces the values of this matrix with the values of from matrix a. The shape of 'a' must be less than or equal mat.
// If the 'a' shape is less then iOffset and jOffset can be used to place 'a' matrix in a specific location.
func (mat *BigIntMatrix) SetMatrix(a *BigIntMatrix, iOffset, jOffset int) {
	if iOffset < 0 || jOffset < 0 {
		panic("offsets must be positive values [0,+)")
	}
	if mat.rows < iOffset+a.rows || mat.cols < jOffset+a.cols {
		panic(fmt.Sprintf("set matrix have equal or smaller shape (%v,%v), found a=(%v,%v)", mat.rows, mat.cols, iOffset+a.rows, jOffset+a.cols))
	}

	mat.setMatrix(a, iOffset+mat.rowStart, jOffset+mat.colStart)
}

func (mat *BigIntMatrix) setMatrix(a *BigIntMatrix, rOffset, cOffset int) {
	mat.zeroize(rOffset, cOffset, a.rows, a.cols)

	for r, cs := range a.rowValues {
		i := r - a.rowStart
		mr := i + rOffset
		for c, v := range cs {
			j := c - a.colStart
			mc := j + cOffset
			mat.set(mr, mc, v)
		}
	}
}

// Negate performs an inplace piecewise negation. For big.Int, this means changing the sign (flip 1 to -1, 0 to 0).
func (mat *BigIntMatrix) Negate() {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			r := i + mat.rowStart
			c := j + mat.colStart

			currentVal := mat.at(r, c)
			if currentVal != nil {
				mat.set(r, c, currentVal.Neg(currentVal))
			}
		}
	}
}
