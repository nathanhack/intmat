package intmat

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Vector struct {
	mat *Matrix
}

type vector struct {
	Mat *Matrix
}

func (vec *Vector) MarshalJSON() ([]byte, error) {
	return json.Marshal(vector{
		Mat: vec.mat,
	})
}

func (vec *Vector) UnmarshalJSON(bytes []byte) error {
	var v vector
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	vec.mat = v.Mat
	return nil
}

func NewVec(length int, values ...int) *Vector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}
	vec := Vector{
		mat: NewMat(1, length, values...),
	}

	return &vec
}

func CopyVec(a *Vector) *Vector {
	return &Vector{
		mat: Copy(a.mat),
	}
}

func (vec *Vector) offset() int {
	return vec.mat.colStart
}

// String returns a string representation of this vector.
func (vec *Vector) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	row := make([]string, vec.Len())
	for i := 0; i < vec.Len(); i++ {
		j := i + vec.offset()
		row[i] = fmt.Sprint(vec.at(j))
	}
	table.Append(row)

	table.Render()
	return buff.String()
}

func (vec *Vector) checkBounds(i int) {
	if i < 0 || i >= vec.Len() {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, vec.Len()-1))
	}
}

// At returns the value at index i.
func (vec *Vector) At(i int) int {
	vec.checkBounds(i)
	j := i + vec.offset()
	return vec.at(j)
}

func (vec *Vector) at(j int) int {
	return vec.mat.rowValues[vec.mat.rowStart][j]
}

// Set sets the value at row index i to value.
func (vec *Vector) Set(i, value int) {
	vec.checkBounds(i)
	j := i + vec.offset()

	vec.set(j, value)
}

func (vec *Vector) set(j, value int) {
	vec.mat.set(vec.mat.rowStart, j, value)
}

// SetVec replaces the values of this vector with the values of from vector a.
func (vec *Vector) SetVec(a *Vector, i int) {
	vec.mat.setMatrix(a.mat, vec.mat.rowStart, i+vec.mat.colStart)
}

func (vec *Vector) Len() int {
	return vec.mat.cols
}

func (vec *Vector) Dot(a *Vector) int {
	m := NewMat(1, 1)
	m.Mul(vec.mat, a.mat.T())
	return m.at(0, 0)
}

func (vec *Vector) NonzeroValues() (indexToValues map[int]int) {
	indexToValues = make(map[int]int)
	end := vec.mat.colStart + vec.mat.cols
	for c, v := range vec.mat.rowValues[vec.mat.rowStart] {
		if c < vec.mat.colStart || end <= c {
			continue
		}
		indexToValues[c] = v
	}
	return
}

func (vec *Vector) T() *TransposedVector {
	return &TransposedVector{
		mat: vec.mat.T(),
	}
}

// Slice creates a slice of the Vector.  The slice will be connected to the original Vector, changes to one
// causes changes in the other.
func (vec *Vector) Slice(i, len int) *Vector {
	if len <= 0 {
		panic("slice len must >0")
	}

	vec.checkBounds(i)
	j := i + vec.offset()

	return &Vector{
		mat: vec.mat.slice(0, j, 1, len),
	}
}

// Add sets vec equal to the sum of a and b.
func (vec *Vector) Add(a, b *Vector) {
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

func (vec *Vector) Equals(v *Vector) bool {
	if vec == v {
		return true
	}
	return vec.mat.Equals(v.mat)
}

func (vec *Vector) Mul(vec2 *Vector, mat *Matrix) {
	if vec == nil || mat == nil {
		panic("vector multiply input was found to be nil")
	}

	if vec == vec2 || vec.mat == mat {
		panic("vector multiply self assignment not allowed")
	}

	if vec2.mat.cols != mat.rows {
		panic(fmt.Sprintf("multiply shape misalignment can't vector-matrix multiply dims: (%v)x(%v,%v)", vec2.mat.cols, mat.rows, mat.cols))
	}

	_, matCols := mat.Dims()
	if vec.mat.cols != matCols {
		panic(fmt.Sprintf("vector not long enough to hold result, actual length:%v required:%v", vec.Len(), mat.cols))
	}

	vec.mat.mul(vec2.mat, mat)
}

func (vec *Vector) And(a, b *Vector) {
	if a == nil || b == nil {
		panic("AND input was found to be nil")
	}

	if vec == a || vec == b {
		panic("AND self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("AND shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if vec.Len() != a.Len() {
		panic(fmt.Sprintf("vec len:%v does not match expected %v", vec.Len(), a.Len()))
	}
	vec.mat.and(a.mat, b.mat)
}

func (vec *Vector) Or(a, b *Vector) {
	if a == nil || b == nil {
		panic("OR input was found to be nil")
	}

	if vec == a || vec == b {
		panic("OR self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("OR shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if vec.Len() != a.Len() {
		panic(fmt.Sprintf("vec len:%v does not match expected %v", vec.Len(), a.Len()))
	}
	vec.mat.or(a.mat, b.mat)
}

func (vec *Vector) XOr(a, b *Vector) {
	if a == nil || b == nil {
		panic("XOR input was found to be nil")
	}

	if vec == a || vec == b {
		panic("XOR self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("XOR shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if vec.Len() != a.Len() {
		panic(fmt.Sprintf("vec len:%v does not match expected %v", vec.Len(), a.Len()))
	}
	vec.mat.xor(a.mat, b.mat)
}

func (vec *Vector) Negate() {
	vec.mat.Negate()
}

type TransposedVector struct {
	mat *Matrix
}

type transposedVector struct {
	Mat *Matrix
}

func (tvec *TransposedVector) MarshalJSON() ([]byte, error) {
	return json.Marshal(transposedVector{
		Mat: tvec.mat,
	})
}

func (tvec *TransposedVector) UnmarshalJSON(bytes []byte) error {
	var v transposedVector
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	tvec.mat = v.Mat
	return nil
}

func NewTVec(length int, values ...int) *TransposedVector {
	if len(values) != 0 {
		if length != len(values) {
			panic("length and number of values must be equal")
		}
	}
	vec := Vector{
		mat: NewMat(1, length, values...),
	}

	return vec.T()
}

func CopyTVec(a *TransposedVector) *TransposedVector {
	return &TransposedVector{
		mat: Copy(a.mat),
	}
}

func (tvec *TransposedVector) checkBounds(i int) {
	if i < 0 || i >= tvec.Len() {
		panic(fmt.Sprintf("%v out of range: [0-%v]", i, tvec.Len()-1))
	}
}

func (tvec *TransposedVector) offset() int {
	return tvec.mat.rowStart
}

func (tvec *TransposedVector) T() *Vector {
	return &Vector{
		mat: tvec.mat.T(),
	}
}

func (tvec *TransposedVector) Len() int {
	return tvec.mat.rows
}

func (tvec *TransposedVector) MulVec(a *Matrix, b *TransposedVector) {
	if a == nil || b == nil {
		panic("multiply input was found to be nil")
	}

	if tvec == b || tvec.mat == a {
		panic("multiply self assignment not allowed")
	}

	if a.cols != b.mat.rows {
		panic(fmt.Sprintf("multiply shape misalignment can't matrix-vector multiply (%v,%v)x(%v,1)", a.rows, a.cols, b.mat.rows))
	}

	if tvec.Len() != b.Len() {
		panic(fmt.Sprintf("transposed vector length (%v) does not match expected (%v)", tvec.Len(), b.Len()))
	}

	tvec.mat.mul(a, b.mat)
}

func (tvec *TransposedVector) Add(a, b *TransposedVector) {
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

// At returns the value at index i.
func (tvec *TransposedVector) At(j int) int {
	tvec.checkBounds(j)
	i := j + tvec.offset()
	return tvec.at(i)
}

func (tvec *TransposedVector) at(i int) int {
	return tvec.mat.rowValues[i][tvec.mat.colStart]
}

// Set sets the value at row index i and column index j to value.
func (tvec *TransposedVector) Set(j, value int) {
	tvec.checkBounds(j)
	i := j + tvec.offset()

	tvec.set(i, value)
}

// SetVec replaces the values of this vector with the values of from vector a.
func (tvec *TransposedVector) SetVec(a *TransposedVector, j int) {
	tvec.mat.setMatrix(a.mat, j+tvec.mat.rowStart, tvec.mat.colStart)
}

// Slice creates a slice of the TransposedVector.  The slice will be connected to the original TransposedVector, changes to one
// causes changes in the other.
func (tvec *TransposedVector) Slice(j, len int) *TransposedVector {
	if len <= 0 {
		panic("slice len must >0")
	}

	tvec.checkBounds(j)
	i := j + tvec.offset()

	return &TransposedVector{
		mat: tvec.mat.slice(i, 0, len, 1),
	}
}

func (tvec *TransposedVector) set(i, value int) {
	tvec.mat.set(i, tvec.mat.colStart, value)
}

func (tvec *TransposedVector) Equals(v *TransposedVector) bool {
	return tvec.mat.Equals(v.mat)
}

func (tvec *TransposedVector) NonzeroValues() (indexToValues map[int]int) {
	indexToValues = make(map[int]int)
	end := tvec.mat.rowStart + tvec.mat.rows
	for r, v := range tvec.mat.colValues[tvec.mat.colStart] {
		if r < tvec.mat.rowStart || end <= r {
			continue
		}
		indexToValues[r] = v
	}
	return
}

func (tvec *TransposedVector) String() string {
	buff := &strings.Builder{}
	table := tablewriter.NewWriter(buff)

	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	row := make([]string, tvec.Len())
	for i := 0; i < tvec.Len(); i++ {
		j := i + tvec.offset()
		row[i] = fmt.Sprint(tvec.at(j))
	}
	table.Append(row)

	table.Render()
	return buff.String()
}

func (tvec *TransposedVector) And(a, b *TransposedVector) {
	if a == nil || b == nil {
		panic("AND input was found to be nil")
	}

	if tvec == a || tvec == b {
		panic("AND self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("AND shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if tvec.Len() != a.Len() {
		panic(fmt.Sprintf("tvec len:%v does not match expected %v", tvec.Len(), a.Len()))
	}
	tvec.mat.and(a.mat, b.mat)
}

func (tvec *TransposedVector) Or(a, b *TransposedVector) {
	if a == nil || b == nil {
		panic("OR input was found to be nil")
	}

	if tvec == a || tvec == b {
		panic("OR self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("OR shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if tvec.Len() != a.Len() {
		panic(fmt.Sprintf("tvec len:%v does not match expected %v", tvec.Len(), a.Len()))
	}
	tvec.mat.or(a.mat, b.mat)
}

func (tvec *TransposedVector) XOr(a, b *TransposedVector) {
	if a == nil || b == nil {
		panic("XOR input was found to be nil")
	}

	if tvec == a || tvec == b {
		panic("XOR self assignment not allowed")
	}

	if a.Len() != b.Len() {
		panic(fmt.Sprintf("XOR shape misalignment both inputs must be equal length found  %v and %v", a.Len(), b.Len()))
	}

	if tvec.Len() != a.Len() {
		panic(fmt.Sprintf("tvec len:%v does not match expected %v", tvec.Len(), a.Len()))
	}
	tvec.mat.xor(a.mat, b.mat)
}

//Negate 
func (tvec *TransposedVector) Negate() {
	tvec.mat.Negate()
}
