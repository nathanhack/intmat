package intmat

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type BigIntVector struct {
	mat *BigIntMatrix
}

type bigintvector struct {
	Mat *BigIntMatrix
}

func (vec *BigIntVector) MarshalJSON() ([]byte, error) {
	return json.Marshal(bigintvector{
		Mat: vec.mat,
	})
}

func (vec *BigIntVector) UnmarshalJSON(bytes []byte) error {
	var v bigintvector
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	vec.mat = v.Mat
	return nil
}

func NewBigIntVec(length int, values ...*big.Int) *BigIntVector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}
	vec := BigIntVector{
		mat: NewBigIntMat(1, length, values...),
	}

	return &vec
}

func CopyBigIntVec(a *BigIntVector) *BigIntVector {
	return &BigIntVector{
		mat: BigIntCopy(a.mat), // Assumes intmat.Copy for BigIntMatrix
	}
}

func (vec *BigIntVector) offset() int {
	return vec.mat.colStart
}

// String returns a string representation of this vector.
func (vec *BigIntVector) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	row := make([]string, vec.Len())
	for i := 0; i < vec.Len(); i++ {
		// Use the public At method for safety and consistency
		row[i] = fmt.Sprint(vec.At(i))
	}
	table.Append(row)

	table.Render()
	return buff.String()
}

func (vec *BigIntVector) checkBounds(i int) {
	if i < 0 || i >= vec.Len() {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, vec.Len()-1))
	}
}

// At returns the value at index i.
func (vec *BigIntVector) At(i int) *big.Int {
	vec.checkBounds(i)
	// vec.mat is a 1xN view. Its row is 0, its col is i (0-indexed for the vector).
	// The underlying BigIntMatrix's At method handles offsets.
	ret := vec.at(i)
	if ret == nil {
		return big.NewInt(0)
	}
	return ret
}

// internal at method, takes absolute column index
func (vec *BigIntVector) at(absColIdx int) *big.Int {
	// This replicates BigIntMatrix.at(r, c) for r = vec.mat.rowStart, c = absColIdx
	// or rather, it's a direct access for String() if that pattern is kept.
	// For consistency with BigIntMatrix's sparse nature:
	if rowMap, ok := vec.mat.rowValues[vec.mat.rowStart]; ok {
		if val, ok2 := rowMap[absColIdx]; ok2 {
			return val
		}
	}
	return nil // Default to zero if not found
}

// Set sets the value at index i to value.
func (vec *BigIntVector) Set(i int, value *big.Int) {
	vec.checkBounds(i)
	j := i + vec.offset()
	vec.set(j, value)
}

// internal set method
func (vec *BigIntVector) set(j int, value *big.Int) {
	vec.mat.set(vec.mat.rowStart, j, value)
}

// SetVec replaces the values of this vector with the values of from vector a, starting at index i of this vector.
func (vec *BigIntVector) SetVec(a *BigIntVector, i int) {
	// This implies setting a portion of vec with a.
	// vec.mat is the destination matrix (1xN)
	// a.mat is the source matrix (1xM)
	// The operation should be like: vec.mat's submatrix starting at (0,i) gets values from a.mat
	if a == nil {
		panic("source vector a cannot be nil")
	}
	if i < 0 || i+a.Len() > vec.Len() {
		panic(fmt.Sprintf("SetVec out of bounds: trying to set %v elements at index %v in a vector of length %v", a.Len(), i, vec.Len()))
	}

	// vec.mat.setMatrix takes absolute offsets in the underlying shared map.
	// We want to set a sub-vector within vec.
	// The target location in vec.mat is (row 0, column i) from vec's perspective.
	// This corresponds to (vec.mat.rowStart, i + vec.mat.colStart) in absolute map coordinates.
	vec.mat.setMatrix(a.mat, vec.mat.rowStart, i+vec.mat.colStart)
}

func (vec *BigIntVector) Len() int {
	if vec.mat == nil {
		return 0
	}
	return vec.mat.cols
}

func (vec *BigIntVector) Dot(a *BigIntVector) *big.Int {
	if vec.Len() != a.Len() {
		panic(fmt.Sprintf("Dot product vectors must have the same length: %v != %v", vec.Len(), a.Len()))
	}
	m := NewBigIntMat(1, 1) // Result is 1x1 matrix
	m.Mul(vec.mat, a.mat.T())
	return m.At(0, 0) // Use At method of BigIntMatrix
}

func (vec *BigIntVector) NonzeroValues() (indexToValues map[int]*big.Int) {
	indexToValues = make(map[int]*big.Int)
	if vec.mat == nil || vec.mat.rowValues == nil {
		return
	}

	rowMap, ok := vec.mat.rowValues[vec.mat.rowStart]
	if !ok {
		return
	}

	// Iterate through the columns relevant to this vector slice
	for absColIdx, val := range rowMap {
		// Check if this absolute column index falls within the vector's range
		if absColIdx >= vec.mat.colStart && absColIdx < vec.mat.colStart+vec.mat.cols {
			if val.Cmp(big.NewInt(0)) != 0 { // Only include non-zero values
				relativeIndex := absColIdx - vec.mat.colStart
				indexToValues[relativeIndex] = val
			}
		}
	}
	return
}

func (vec *BigIntVector) T() *TransposedBigIntVector {
	return &TransposedBigIntVector{
		mat: vec.mat.T(),
	}
}

// Slice creates a slice of the Vector.  The slice will be connected to the original Vector, changes to one
// causes changes in the other.
func (vec *BigIntVector) Slice(i, len int) *BigIntVector {
	if len <= 0 {
		panic("slice length must be > 0")
	}
	vec.checkBounds(i) // Checks i
	j := i + vec.offset()

	return &BigIntVector{
		mat: vec.mat.slice(0, j, 1, len),
	}
}

// Add sets vec equal to the sum of a and b.
func (vec *BigIntVector) Add(a, b *BigIntVector) {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}
	if vec == a || vec == b {
		panic("addition self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic("adding vectors must have the same length")
	}
	if vec.Len() != a.Len() {
		panic("adding vectors, destination must have the same length")
	}

	vec.mat.add(a.mat, b.mat)
}

func (tvec *TransposedBigIntVector) Add(a, b *TransposedBigIntVector) {
	if a == nil || b == nil {
		panic("addition input was found to be nil")
	}
	if tvec == a || tvec == b {
		panic("addition self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic("adding transposed vectors must have the same length")
	}
	if tvec.Len() != a.Len() {
		panic("adding transposed vectors, destination must have the same length")
	}

	tvec.mat.add(a.mat, b.mat)
}
func (vec *BigIntVector) Equals(v *BigIntVector) bool {
	if vec == v {
		return true
	}

	return vec.mat.Equals(v.mat)
}

func (vec *BigIntVector) Mul(vec2 *BigIntVector, matInput *BigIntMatrix) {
	if vec == nil || vec.mat == nil || vec2 == nil || vec2.mat == nil || matInput == nil {
		panic("vector multiply input (vec, vec2, or matInput) or their underlying matrices were found to be nil")
	}

	if vec == vec2 || vec.mat == matInput || vec.mat == vec2.mat { // Check underlying matrices too
		panic("vector multiply self assignment not allowed (receiver, vec2, or matInput share memory)")
	}

	if vec2.mat.cols != matInput.rows {
		panic(fmt.Sprintf("multiply shape misalignment: can't vector-matrix multiply dims (%v)x(%v,%v). vec2.Len()=%v, matInput.rows=%v", vec2.mat.cols, matInput.rows, matInput.cols, vec2.Len(), matInput.rows))
	}

	_, matColsResult := matInput.Dims()
	if vec.Len() != matColsResult {
		panic(fmt.Sprintf("vector (receiver) not long enough to hold result, actual length:%v required:%v", vec.Len(), matColsResult))
	}

	vec.mat.mul(vec2.mat, matInput)
}

func (vec *BigIntVector) Negate() {
	if vec.mat == nil {
		return
	}
	vec.mat.Negate()
}

// TransposedVector in this file context uses BigIntMatrix
type TransposedBigIntVector struct {
	mat *BigIntMatrix
}

type transposedBigIntVector struct { // for JSON marshalling
	Mat *BigIntMatrix
}

func (tvec *TransposedBigIntVector) MarshalJSON() ([]byte, error) {
	return json.Marshal(transposedBigIntVector{
		Mat: tvec.mat,
	})
}

func (tvec *TransposedBigIntVector) UnmarshalJSON(bytes []byte) error {
	var v transposedBigIntVector
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	tvec.mat = v.Mat
	return nil
}

func NewTBigIntVec(length int, values ...*big.Int) *TransposedBigIntVector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}
	// Create a BigIntVector (row vector) first, then transpose it.
	// NewBigIntVec handles creating a 1xlength BigIntMatrix with the values.
	tempRowVec := NewBigIntVec(length, values...)
	return tempRowVec.T() // .T() on BigIntVector returns *TransposedVector (this file's type)
}

func CopyTBigIntVec(a *TransposedBigIntVector) *TransposedBigIntVector {
	if a == nil || a.mat == nil {
		panic("cannot copy nil TransposedVector or one with nil matrix")
	}
	return &TransposedBigIntVector{
		mat: BigIntCopy(a.mat), // Assumes intmat.Copy for BigIntMatrix
	}
}

func (tvec *TransposedBigIntVector) offset() int {
	// For a transposed vector (column vector), the "offset" refers to the rowStart of its underlying matrix.
	return tvec.mat.rowStart
}

func (tvec *TransposedBigIntVector) T() *BigIntVector {
	return &BigIntVector{
		mat: tvec.mat.T(),
	}
}

func (tvec *TransposedBigIntVector) Len() int {
	if tvec.mat == nil {
		return 0
	}
	return tvec.mat.rows // Length of a column vector is its number of rows
}

func (tvec *TransposedBigIntVector) checkBounds(i int) {
	if i < 0 || i >= tvec.Len() {
		panic(fmt.Sprintf("%v out of range for TransposedVector: [0-%v]", i, tvec.Len()-1))
	}
}

// At returns the value at index i (row i for the column vector).
func (tvec *TransposedBigIntVector) At(i int) *big.Int {
	tvec.checkBounds(i)
	// tvec.mat is an Nx1 view. Its row is i, its col is 0.
	return tvec.mat.At(i, 0)
}

// internal at method, takes absolute row index
func (tvec *TransposedBigIntVector) at(absRowIdx int) *big.Int {
	// This replicates BigIntMatrix.at(r, c) for r = absRowIdx, c = tvec.mat.colStart
	if rowMap, ok := tvec.mat.rowValues[absRowIdx]; ok {
		if val, ok2 := rowMap[tvec.mat.colStart]; ok2 {
			return val
		}
	}
	return big.NewInt(0)
}

// Set sets the value at index i (row i for the column vector) to value.
func (tvec *TransposedBigIntVector) Set(i int, value *big.Int) {
	tvec.checkBounds(i)
	// tvec.mat is an Nx1 view. Its row is i, its col is 0.
	tvec.mat.Set(i, 0, value)
}

// internal set method, takes absolute row index
func (tvec *TransposedBigIntVector) set(absRowIdx int, value *big.Int) {
	tvec.mat.set(absRowIdx, tvec.mat.colStart, value)
}

// SetVec replaces the values of this transposed vector with the values of from transposed vector a, starting at index j of this vector.
func (tvec *TransposedBigIntVector) SetVec(a *TransposedBigIntVector, j int) {
	if a == nil {
		panic("source transposed vector a cannot be nil")
	}
	if j < 0 || j+a.Len() > tvec.Len() {
		panic(fmt.Sprintf("SetVec out of bounds: trying to set %v elements at index %v in a transposed vector of length %v", a.Len(), j, tvec.Len()))
	}
	// tvec.mat is the destination matrix (Nx1)
	// a.mat is the source matrix (Mx1)
	// Target location in tvec.mat is (row j, column 0) from tvec's perspective.
	// This corresponds to (j + tvec.mat.rowStart, tvec.mat.colStart) in absolute map coordinates.
	tvec.mat.setMatrix(a.mat, j+tvec.mat.rowStart, tvec.mat.colStart)
}

// Slice creates a slice of the TransposedVector. The slice will be connected to the original.
func (tvec *TransposedBigIntVector) Slice(j, length int) *TransposedBigIntVector {
	if length <= 0 {
		panic("slice length must be > 0")
	}
	tvec.checkBounds(j) // Checks j
	if j+length > tvec.Len() {
		panic(fmt.Sprintf("slice [%v:%v] out of bounds for transposed vector of length %v", j, j+length, tvec.Len()))
	}
	// tvec.mat.slice takes absolute start row, absolute start col, num rows, num cols
	// For a TransposedVector, it's always 1 column.
	// The absolute start row for the slice is tvec.mat.rowStart + j.
	absoluteStartRow := tvec.mat.rowStart + j
	return &TransposedBigIntVector{
		mat: tvec.mat.slice(absoluteStartRow, tvec.mat.colStart, length, 1),
	}
}

func (tvec *TransposedBigIntVector) MulVec(matInput *BigIntMatrix, b *TransposedBigIntVector) {
	// This method's name tvec.MulVec(a,b) implies tvec = a * b where a is Matrix, b is TransposedVector
	// So, tvec.mat = matInput * b.mat
	if tvec == nil || tvec.mat == nil || matInput == nil || b == nil || b.mat == nil {
		panic("multiply input (tvec, matInput, or b) or their underlying matrices were found to be nil")
	}
	if tvec.mat == matInput || tvec.mat == b.mat {
		panic("multiply self assignment not allowed (receiver, matInput, or b share memory)")
	}

	// matInput (A) is M x K, b.mat (B_col) is K x 1. Result (tvec.mat) should be M x 1.
	if matInput.cols != b.mat.rows { // K_A != K_B (rows of B_col)
		panic(fmt.Sprintf("multiply shape misalignment: can't matrix-vector multiply (%v,%v)x(%v,1). matInput.cols=%v, b.Len()=%v", matInput.rows, matInput.cols, b.mat.rows, matInput.cols, b.Len()))
	}

	if tvec.Len() != matInput.rows { // M_tvec != M_A
		panic(fmt.Sprintf("transposed vector (receiver) length (%v) does not match expected matrix rows (%v)", tvec.Len(), matInput.rows))
	}
	if tvec.mat.cols != 1 { // Ensure receiver is a column vector
		panic(fmt.Sprintf("receiver transposed vector must have 1 column, found %v", tvec.mat.cols))
	}

	tvec.mat.mul(matInput, b.mat)
}

func (tvec *TransposedBigIntVector) Equals(v *TransposedBigIntVector) bool {
	if tvec == v {
		return true
	}
	if tvec == nil || v == nil || tvec.mat == nil || v.mat == nil {
		return false
	}
	return tvec.mat.Equals(v.mat)
}

func (tvec *TransposedBigIntVector) NonzeroValues() (indexToValues map[int]*big.Int) {
	indexToValues = make(map[int]*big.Int)
	if tvec.mat == nil || tvec.mat.colValues == nil { // Transposed vector uses colValues of its matrix if it's a direct column view
		return
	}

	// A transposed vector is represented by an N x 1 matrix.
	// We are interested in the non-zero values in its single column (tvec.mat.colStart).
	// The keys of indexToValues are the row indices (0 to Len-1).
	colMap, ok := tvec.mat.colValues[tvec.mat.colStart]
	if !ok {
		return
	}

	for absRowIdx, val := range colMap {
		// Check if this absolute row index falls within the vector's range
		if absRowIdx >= tvec.mat.rowStart && absRowIdx < tvec.mat.rowStart+tvec.mat.rows {
			if val.Cmp(big.NewInt(0)) != 0 {
				relativeIndex := absRowIdx - tvec.mat.rowStart
				indexToValues[relativeIndex] = val
			}
		}
	}
	return
}

func (tvec *TransposedBigIntVector) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	// Even though it's a column vector, the original String method printed it as a row.
	// To print as a column, each element would be its own row in the table.
	// Sticking to original behavior: print as a single row.
	rowStrings := make([]string, tvec.Len())
	for i := 0; i < tvec.Len(); i++ {
		rowStrings[i] = fmt.Sprint(tvec.At(i))
	}
	table.Append(rowStrings)

	table.Render()
	return buff.String()
}

func (tvec *TransposedBigIntVector) Negate() {
	if tvec.mat == nil {
		return
	}
	tvec.mat.Negate()
}
