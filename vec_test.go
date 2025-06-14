package intmat

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

func TestVector_Add(t *testing.T) {
	tests := []struct {
		a, b, result *Vector
		expected     *Vector
	}{
		{NewVec(3, 0, 1, 0), NewVec(3, 1, 0, 0), NewVec(3), NewVec(3, 1, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Add(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestVector_Dot(t *testing.T) {
	tests := []struct {
		a, b     *Vector
		expected int
	}{
		{NewVec(3, 1, 1, 1), NewVec(3, 1, 1, 1), 3},
		{NewVec(3, 1, 0, 1), NewVec(3, 1, 1, 1), 2},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.a.Dot(test.b)
			if test.expected != actual {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestVector_Mul(t *testing.T) {
	tests := []struct {
		a        *Vector
		b        *Matrix
		result   *Vector
		expected *Vector
	}{
		{NewVec(3, 1, 0, 1), Identity(3), NewVec(3), NewVec(3, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Mul(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestVector_Equals(t *testing.T) {
	tests := []struct {
		a, b     *Vector
		expected bool
	}{
		{NewVec(3), NewVec(3), true},
		{NewVec(3), NewVec(4), false},
		{NewVec(3, 1, 0, 1), NewVec(3), false},
		{NewVec(3, 1, 0, 1), NewVec(3, 1, 0, 1), true},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.a.Equals(test.b) != test.expected {
				t.Fatalf("expected %v for %v == %v", test.expected, test.a, test.b)
			}
		})
	}
}

func TestVector_Set(t *testing.T) {
	tests := []struct {
		source, result *Vector
	}{
		{NewVec(5, 1, 0, 1, 0, 1), NewVec(5)},
		{Identity(5).Row(2), NewVec(5)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for i := 0; i < test.source.Len(); i++ {
				test.result.Set(i, test.source.At(i))
			}

			if !test.source.Equals(test.result) {
				t.Fatalf("expected %v but found %v", test.source, test.result)
			}
		})
	}
}

func TestVector_Set2(t *testing.T) {
	m := NewMat(5, 5)
	for i := 0; i < 5; i++ {
		row := m.Row(i)
		row.Set(i, 1)
	}

	expected := NewMat(5, 5, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1)

	if !expected.Equals(m) {
		t.Fatalf("expected %v but found %v", expected, m)
	}
}

func TestVector_SetVec(t *testing.T) {
	tests := []struct {
		original         *Vector
		setToSlice       *Vector
		index            int
		expectedOriginal *Vector
	}{
		{NewVec(5, 1, 1, 1, 1, 1), NewVec(3, 0, 1, 0), 1, NewVec(5, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.SetVec(test.setToSlice, test.index)

			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestVector_NonzeroValues(t *testing.T) {
	tests := []struct {
		input    *Vector
		expected map[int]int
	}{
		{Identity(4).Row(2), map[int]int{2: 1}},
		{NewMat(4, 6, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1).Row(0), map[int]int{0: 1, 1: 1, 3: 1}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.NonzeroValues()
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestVector_Slice(t *testing.T) {
	tests := []struct {
		original         *Vector
		i, len           int
		addToSlice       *Vector
		expectedOriginal *Vector
		expectedSlice    *Vector
	}{
		{
			NewVec(5, 1, 0, 1, 0, 1),
			1,
			3,
			NewVec(3, 1, 1, 1),
			NewVec(5, 1, 1, 2, 1, 1),
			NewVec(3, 1, 2, 1),
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.len)
			sl.Add(CopyVec(sl), test.addToSlice)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestVector_Slice2(t *testing.T) {
	tests := []struct {
		original      *Vector
		i1, len1      int
		i2, len2      int
		expectedSlice *Vector
	}{
		{NewVec(7, 0, 0, 1, 0, 1, 0, 0), 1, 5, 1, 3, NewVec(3, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i1, test.len1).Slice(test.i2, test.len2)

			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}

		})
	}
}

func TestVector_And(t *testing.T) {
	tests := []struct {
		x, y, result, expected *Vector
	}{
		{NewVec(4, 0, 1, 0, 1), NewVec(4, 0, 0, 1, 1), NewVec(4), NewVec(4, 0, 0, 0, 1)},
		{NewVec(4, 0, 0, 1, 1), NewVec(4, 0, 1, 0, 1), NewVec(4), NewVec(4, 0, 0, 0, 1)},
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

func TestVector_Or(t *testing.T) {
	tests := []struct {
		x, y, result, expected *Vector
	}{
		{NewVec(4, 0, 1, 0, 1), NewVec(4, 0, 0, 1, 1), NewVec(4), NewVec(4, 0, 1, 1, 1)},
		{NewVec(4, 0, 0, 1, 1), NewVec(4, 0, 1, 0, 1), NewVec(4), NewVec(4, 0, 1, 1, 1)},
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

func TestVector_XOr(t *testing.T) {
	tests := []struct {
		x, y, result, expected *Vector
	}{
		{NewVec(4, 0, 1, 0, 1), NewVec(4, 0, 0, 1, 1), NewVec(4), NewVec(4, 0, 1, 1, 0)},
		{NewVec(4, 0, 0, 1, 1), NewVec(4, 0, 1, 0, 1), NewVec(4), NewVec(4, 0, 1, 1, 0)},
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

func TestVector_Negate(t *testing.T) {
	tests := []struct {
		x, expected *Vector
	}{
		{NewVec(4, 0, 1, -2, 3), NewVec(4, 0, -1, 2, -3)},
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

func TestTransposedVector_Set(t *testing.T) {
	tests := []struct {
		source, result *TransposedVector
	}{
		{NewTVec(5, 1, 0, 1, 0, 1), NewTVec(5)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for i := 0; i < test.source.Len(); i++ {
				test.result.Set(i, test.source.At(i))
			}

			if !test.source.Equals(test.result) {
				t.Fatalf("expected %v but found %v", test.source, test.result)
			}
		})
	}
}

func TestTransposedVector_Set2(t *testing.T) {
	m := NewMat(5, 5)
	for j := 0; j < 5; j++ {
		col := m.Column(j)
		col.Set(j, 1)
	}

	expected := NewMat(5, 5, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1)

	if !expected.Equals(m) {
		t.Fatalf("expected %v but found %v", expected, m)
	}
}

func TestTransposedVector_MulVec(t *testing.T) {
	tests := []struct {
		a        *Matrix
		b        *Vector
		result   *TransposedVector
		expected *Vector
	}{
		{Identity(3), NewVec(3, 0, 1, 0), NewTVec(3), NewVec(3, 0, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.MulVec(test.a, test.b.T())
			if !test.result.T().Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result.T())
			}
		})
	}
}

func TestTransposedVector_Add(t *testing.T) {
	tests := []struct {
		a, b, result *TransposedVector
		expected     *TransposedVector
	}{
		{NewTVec(3, 0, 1, 0), NewTVec(3, 1, 0, 0), NewTVec(3), NewTVec(3, 1, 1, 0)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Add(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, test.result)
			}
		})
	}
}

func TestTransposedVector_NonzeroValues(t *testing.T) {
	tests := []struct {
		input    *TransposedVector
		expected map[int]int
	}{
		{Identity(4).Column(2), map[int]int{2: 1}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.NonzeroValues()
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestTransposedVector_Slice(t *testing.T) {
	tests := []struct {
		original         *TransposedVector
		i, len           int
		addToSlice       *TransposedVector
		expectedOriginal *TransposedVector
		expectedSlice    *TransposedVector
	}{
		{
			NewTVec(5, 1, 0, 1, 0, 1),
			1,
			3,
			NewTVec(3, 1, 1, 1),
			NewTVec(5, 1, 1, 2, 1, 1),
			NewTVec(3, 1, 2, 1),
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.len)
			sl.Add(CopyTVec(sl), test.addToSlice)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestTransposedVector_Slice2(t *testing.T) {
	tests := []struct {
		original      *TransposedVector
		i1, len1      int
		i2, len2      int
		expectedSlice *TransposedVector
	}{
		{NewTVec(7, 0, 0, 1, 0, 1, 0, 0), 1, 5, 1, 3, NewTVec(3, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i1, test.len1).Slice(test.i2, test.len2)

			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}

		})
	}
}

func TestTransposedVector_SetVec(t *testing.T) {
	tests := []struct {
		original         *TransposedVector
		setToSlice       *TransposedVector
		index            int
		expectedOriginal *TransposedVector
	}{
		{NewTVec(5, 1, 1, 1, 1, 1), NewTVec(3, 0, 1, 0), 1, NewTVec(5, 1, 0, 1, 0, 1)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.SetVec(test.setToSlice, test.index)

			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestVector_JSON(t *testing.T) {
	v := NewVec(5, 1, 0, 1, 0, 1)

	bs, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual Vector
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !v.Equals(&actual) {
		t.Fatalf("expected %v but found %v", v, actual)
	}
}

func TestTransposedVector_JSON(t *testing.T) {
	v := NewTVec(5, 1, 0, 1, 0, 1)

	bs, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual TransposedVector
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !v.Equals(&actual) {
		t.Fatalf("expected %v but found %v", v, actual)
	}
}

func TestTransposedVector_And(t *testing.T) {
	tests := []struct {
		x, y, result, expected *TransposedVector
	}{
		{NewTVec(4, 0, 1, 0, 1), NewTVec(4, 0, 0, 1, 1), NewTVec(4), NewTVec(4, 0, 0, 0, 1)},
		{NewTVec(4, 0, 0, 1, 1), NewTVec(4, 0, 1, 0, 1), NewTVec(4), NewTVec(4, 0, 0, 0, 1)},
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

func TestTransposedVector_Or(t *testing.T) {
	tests := []struct {
		x, y, result, expected *TransposedVector
	}{
		{NewTVec(4, 0, 1, 0, 1), NewTVec(4, 0, 0, 1, 1), NewTVec(4), NewTVec(4, 0, 1, 1, 1)},
		{NewTVec(4, 0, 0, 1, 1), NewTVec(4, 0, 1, 0, 1), NewTVec(4), NewTVec(4, 0, 1, 1, 1)},
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

func TestTransposedVector_XOr(t *testing.T) {
	tests := []struct {
		x, y, result, expected *TransposedVector
	}{
		{NewTVec(4, 0, 1, 0, 1), NewTVec(4, 0, 0, 1, 1), NewTVec(4), NewTVec(4, 0, 1, 1, 0)},
		{NewTVec(4, 0, 0, 1, 1), NewTVec(4, 0, 1, 0, 1), NewTVec(4), NewTVec(4, 0, 1, 1, 0)},
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

func TestTransposedVector_Negate(t *testing.T) {
	tests := []struct {
		x, expected *TransposedVector
	}{
		{NewTVec(4, 0, 1, 0, 1), NewTVec(4, 0, -1, 0, -1)},
		{NewTVec(4, 0, 0, 1, 1), NewTVec(4, 0, 0, -1, -1)},
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
