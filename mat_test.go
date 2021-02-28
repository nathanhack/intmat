package intmat

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		rows, cols int
		data       []int
		expected   [][]int
	}{
		{1, 1, []int{1}, [][]int{{1}}},
		{2, 2, []int{1, 0, 0, 1}, [][]int{{1, 0}, {0, 1}}},
		{2, 2, []int{}, [][]int{{0, 0}, {0, 0}}},
		{2, 2, nil, [][]int{{0, 0}, {0, 0}}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var m *Matrix
			if test.data != nil {
				m = NewMat(test.rows, test.cols, test.data...)
			} else {
				m = NewMat(test.rows, test.cols)
			}

			for i := 0; i < len(test.expected); i++ {
				for j := 0; j < len(test.expected[i]); j++ {
					expected := test.expected[i][j]
					actual := m.At(i, j)
					if actual != expected {
						t.Fatalf("expected %v but found %v", expected, actual)
					}
				}
			}
		})
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		rows, cols int
		data       []int
		expected   [][]int
	}{
		{1, 1, []int{1}, [][]int{{1}}},
		{2, 2, []int{1, 0, 0, 1}, [][]int{{1, 0}, {0, 1}}},
		{2, 2, []int{}, [][]int{{0, 0}, {0, 0}}},
		{2, 2, nil, [][]int{{0, 0}, {0, 0}}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var m1 *Matrix
			if test.data != nil {
				m1 = NewMat(test.rows, test.cols, test.data...)
			} else {
				m1 = NewMat(test.rows, test.cols)
			}

			m := Copy(m1)

			for i := 0; i < len(test.expected); i++ {
				for j := 0; j < len(test.expected[i]); j++ {
					expected := test.expected[i][j]
					actual := m.At(i, j)
					if actual != expected {
						t.Fatalf("expected %v but found %v", expected, actual)
					}
				}
			}
		})
	}
}

func TestMatrix_Dim(t *testing.T) {
	tests := []struct {
		m            *Matrix
		expectedRows int
		expectedCols int
	}{
		{NewMat(5, 5), 5, 5},
		{NewMat(5, 5).Slice(1, 1, 4, 4), 4, 4},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rows, cols := test.m.Dims()

			if rows != test.expectedRows {
				t.Fatalf("expected %v but found %v", test.expectedRows, rows)
			}
			if cols != test.expectedCols {
				t.Fatalf("expected %v but found %v", test.expectedCols, cols)
			}
		})
	}
}

func TestMatrix_Slice(t *testing.T) {
	tests := []struct {
		sliced   *Matrix
		expected *Matrix
	}{
		{NewMat(2, 2, 1, 0, 0, 1).Slice(0, 0, 2, 1), NewMat(2, 1, 1, 0)},
		{Identity(8).Slice(3, 0, 4, 4), NewMat(4, 4, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{Identity(8).Slice(3, 0, 4, 4).T(), NewMat(4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.sliced.Equals(test.expected) {
				t.Fatalf("expected equality between %v and %v", test.sliced, test.expected)
			}

		})
	}
}

func TestMatrix_Slice2(t *testing.T) {

	orignal := Identity(5)

	slice := orignal.Slice(1, 1, 3, 3).T()

	slice.Zeroize()
	slice.Set(1, 0, 1)
	slice.Set(2, 0, 1)

	expectedSlice := NewMat(3, 3, 0, 0, 0, 1, 0, 0, 1, 0, 0)
	expectedOriginal := NewMat(5, 5, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)

	if !slice.Equals(expectedSlice) {
		t.Fatalf("expcted \n%v\n but found \n%v\n", expectedSlice, slice)
	}
	if !orignal.Equals(expectedOriginal) {
		t.Fatalf("expcted \n%v\n but found \n%v\n", expectedOriginal, orignal)
	}
}

func TestMatrix_Slice3(t *testing.T) {
	x := NewMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	y := NewMat(2, 2, 1, 0, 0, 0)
	result := x.Slice(1, 1, 2, 2)
	s := NewMat(2, 2, 1, 1, 1, 1)
	result.Mul(y, s)
	expected := NewMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)

	if !expected.Equals(x) {
		t.Fatalf("expected %v but found %v", expected, x)
	}
}

func TestMatrix_Equals(t *testing.T) {
	tests := []struct {
		input1, input2 *Matrix
		expected       bool
	}{
		{Identity(3), Identity(3), true},
		{Identity(3).T(), Identity(3), true},
		{Identity(4), Identity(3), false},
		{Identity(4), nil, false},
		{nil, Identity(4), false},
		{nil, nil, true},
		{NewMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0).T(), NewMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), false},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input1.Equals(test.input2)
			if actual != test.expected {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestIdentity(t *testing.T) {
	tests := []struct {
		ident    *Matrix
		expected *Matrix
	}{
		{Identity(3), NewMat(3, 3, 1, 0, 0, 0, 1, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.ident.Equals(test.expected) {
				t.Fatalf("expected equality")
			}
		})
	}
}

func TestMatrix_At(t *testing.T) {
	tests := []struct {
		input    *Matrix
		i, j     int
		expected int
	}{
		{Identity(3), 0, 0, 1},
		{Identity(3), 0, 1, 0},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.At(test.i, test.j)
			if actual != test.expected {
				t.Fatalf("expected %v at (%v,%v) but found %v", test.expected, test.i, test.j, actual)
			}
		})
	}
}

func TestMatrix_Mul(t *testing.T) {
	tests := []struct {
		m1, m2, result, expected *Matrix
	}{
		0: {NewMat(1, 4, 1, 0, 1, 0), NewMat(4, 1, 1, 0, 1, 0), NewMat(1, 1), NewMat(1, 1, 2)},
		1: {NewMat(1, 4, 1, 0, 1, 0), NewMat(4, 1, 1, 0, 0, 0), NewMat(1, 1), NewMat(1, 1, 1)},
		2: {NewMat(1, 4, 1, 1, 1, 1), NewMat(4, 1, 1, 1, 1, 0), NewMat(1, 1), NewMat(1, 1, 3)},
		3: {Identity(3), Identity(3), NewMat(3, 3), Identity(3)},
		4: {Identity(3), NewMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), NewMat(3, 3), NewMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		5: {NewMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), Identity(3), NewMat(3, 3), NewMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0)},
		6: {NewMat(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T(), Identity(4), NewMat(3, 4), NewMat(4, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1).T()},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Mul(test.m1, test.m2)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v equality but found %v", test.expected, test.result)
			}
		})
	}
}

func TestMatrix_Zeroize(t *testing.T) {
	tests := []struct {
		original *Matrix
		expected *Matrix
	}{
		{Identity(3), NewMat(3, 3)},
		{NewMat(3, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1), NewMat(3, 3)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.Zeroize()
			if !test.original.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original)
			}
		})
	}
}

func TestMatrix_ZeroizeRange(t *testing.T) {
	tests := []struct {
		original         *Matrix
		i, j, rows, cols int
		expected         *Matrix
	}{
		{NewMat(4, 4, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 1, 1, 2, 2, NewMat(4, 4, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.ZeroizeRange(test.i, test.j, test.rows, test.cols)
			if !test.original.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original)
			}
		})
	}
}

func TestMatrix_T(t *testing.T) {
	tests := []struct {
		original *Matrix
		expected *Matrix
	}{
		{NewMat(3, 3, 0, 1, 1, 0, 1, 1, 0, 0, 0), NewMat(3, 3, 0, 0, 0, 1, 1, 0, 1, 1, 0)},
		{NewMat(4, 2, 0, 1, 0, 0, 0, 0, 1, 0), NewMat(2, 4, 0, 0, 0, 1, 1, 0, 0, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !test.original.T().Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.original.T())
			}
		})
	}
}

func TestMatrix_Add(t *testing.T) {
	tests := []struct {
		a, b, result *Matrix
		expected     *Matrix
	}{
		{Identity(3), Identity(3), NewMat(3, 3), NewMat(3, 3, 2, 0, 0, 0, 2, 0, 0, 0, 2)},
		{Identity(3), NewMat(3, 3), NewMat(3, 3), Identity(3)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Add(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.result)
			}
		})
	}
}

func TestMatrix_Add2(t *testing.T) {
	tests := []struct {
		original         *Matrix
		i, j, rows, cols int
		addToSlice       *Matrix
		expectedOriginal *Matrix
		expectedSlice    *Matrix
	}{
		{
			NewMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
			1, 1, 3, 3,
			Identity(3),
			NewMat(5, 5, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1),
			NewMat(3, 3, 2, 1, 1, 1, 2, 1, 1, 1, 2),
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.j, test.rows, test.cols)
			c := Copy(sl)
			sl.Add(c, test.addToSlice)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestMatrix_Column(t *testing.T) {
	tests := []struct {
		m        *Matrix
		j        int //column
		expected *TransposedVector
	}{
		{Identity(3), 1, NewTVec(3, 0, 1, 0)},
		{Identity(3), 0, NewTVec(3, 1, 0, 0)},
		{NewMat(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 0, NewTVec(2, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.m.Column(test.j)

			if !actual.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestMatrix_Row(t *testing.T) {
	tests := []struct {
		m        *Matrix
		i        int //row index
		expected *Vector
	}{
		{Identity(3), 1, NewVec(3, 0, 1, 0)},
		{Identity(3), 0, NewVec(3, 1, 0, 0)},
		{NewMat(4, 4, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0).Slice(1, 1, 2, 2).T(), 1, NewVec(2, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.m.Row(test.i)

			if !actual.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestMatrix_SetColumn(t *testing.T) {
	tests := []struct {
		m        *Matrix
		j        int //column to change
		vec      *TransposedVector
		expected *Matrix
	}{
		{Identity(3), 0, NewTVec(3, 0, 1, 0), NewMat(3, 3, 0, 0, 0, 1, 1, 0, 0, 0, 1)},
		{Identity(3), 1, Identity(3).Column(2), NewMat(3, 3, 1, 0, 0, 0, 0, 0, 0, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.m.SetColumn(test.j, test.vec)
			if !test.m.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.m)
			}
		})
	}
}

func TestMatrix_SetRow(t *testing.T) {
	tests := []struct {
		m        *Matrix
		i        int //row to change
		vec      *Vector
		expected *Matrix
	}{
		{Identity(3), 0, NewVec(3, 0, 1, 0), NewMat(3, 3, 0, 1, 0, 0, 1, 0, 0, 0, 1)},
		{Identity(3), 1, Identity(3).Row(2), NewMat(3, 3, 1, 0, 0, 0, 0, 1, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.m.SetRow(test.i, test.vec)
			if !test.m.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.m)
			}
		})
	}
}

func TestMatrix_SetMatrix(t *testing.T) {
	tests := []struct {
		dest             *Matrix
		source           *Matrix
		iOffset, jOffset int
		expected         *Matrix
	}{
		{NewMat(3, 3), Identity(3), 0, 0, Identity(3)},
		{NewMat(4, 4), NewMat(2, 2, 1, 1, 1, 1), 1, 1, NewMat(4, 4, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0)},
		{NewMat(4, 4), NewMat(2, 2, 0, 1, 0, 0).T(), 1, 1, NewMat(4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0)},
		{Identity(4), NewMat(2, 2, 1, 1, 1, 1), 1, 1, NewMat(4, 4, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1)},
		{NewMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), Identity(3), 1, 1, NewMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.dest.SetMatrix(test.source, test.iOffset, test.jOffset)
			if !test.dest.Equals(test.expected) {
				t.Fatalf("expcted \n%v\n but found \n%v\n", test.expected, test.dest)
			}
		})
	}
}

func TestMatrix_SetMatrix2(t *testing.T) {
	tests := []struct {
		original         *Matrix
		i, j, rows, cols int
		source           *Matrix
		iOffset, jOffset int
		expectedOriginal *Matrix
		expectedSlice    *Matrix
	}{
		{NewMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
			1, 1, 3, 3,
			Identity(2),
			1, 1,
			NewMat(5, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1),
			NewMat(3, 3, 1, 1, 1, 1, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.j, test.rows, test.cols)
			sl.SetMatrix(test.source, test.iOffset, test.jOffset)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestMatrix_JSON(t *testing.T) {
	m := Identity(3)

	bs, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual Matrix
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !m.Equals(&actual) {
		t.Fatalf("expected %v but found %v", m, actual)
	}
}

func TestMatrix_And(t *testing.T) {
	tests := []struct {
		x, y, result, expected *Matrix
	}{
		{NewMat(2, 2, 0, 1, 0, 1), NewMat(2, 2, 0, 0, 1, 1), NewMat(2, 2), NewMat(2, 2, 0, 0, 0, 1)},
		{NewMat(2, 2, 0, 0, 1, 1), NewMat(2, 2, 0, 1, 0, 1), NewMat(2, 2), NewMat(2, 2, 0, 0, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.And(test.x, test.y)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestMatrix_Or(t *testing.T) {
	tests := []struct {
		x, y, result, expected *Matrix
	}{
		{NewMat(2, 2, 0, 1, 0, 1), NewMat(2, 2, 0, 0, 1, 1), NewMat(2, 2), NewMat(2, 2, 0, 1, 1, 1)},
		{NewMat(2, 2, 0, 0, 1, 1), NewMat(2, 2, 0, 1, 0, 1), NewMat(2, 2), NewMat(2, 2, 0, 1, 1, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Or(test.x, test.y)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestMatrix_XOr(t *testing.T) {
	tests := []struct {
		x, y, result, expected *Matrix
	}{
		{NewMat(2, 2, 0, 1, 0, 1), NewMat(2, 2, 0, 0, 1, 1), NewMat(2, 2), NewMat(2, 2, 0, 1, 1, 0)},
		{NewMat(2, 2, 0, 0, 1, 1), NewMat(2, 2, 0, 1, 0, 1), NewMat(2, 2), NewMat(2, 2, 0, 1, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.XOr(test.x, test.y)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestMatrix_Negate(t *testing.T) {
	tests := []struct {
		x, expected *Matrix
	}{
		{Identity(3), NewMat(3, 3, -1, 0, 0, 0, -1, 0, 0, 0, -1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.x.Negate()

			if !test.x.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.x)
			}
		})
	}
}
