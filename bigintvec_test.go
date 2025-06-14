package intmat

import (
	"encoding/json"
	"math/big"
	"reflect"
	"strconv"
	"testing"
)

func TestBigIntVector_Add(t *testing.T) {
	tests := []struct {
		a, b, result *BigIntVector
		expected     *BigIntVector
	}{
		{NewBigIntVec(3, intsToBigInts([]int{0, 1, 0})...), NewBigIntVec(3, intsToBigInts([]int{1, 0, 0})...), NewBigIntVec(3), NewBigIntVec(3, intsToBigInts([]int{1, 1, 0})...)},
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

func TestBigIntVector_Dot(t *testing.T) {
	tests := []struct {
		a, b     *BigIntVector
		expected *big.Int
	}{
		{NewBigIntVec(3, intsToBigInts([]int{1, 1, 1})...), NewBigIntVec(3, intsToBigInts([]int{1, 1, 1})...), big.NewInt(3)}, // 1*1 + 1*1 + 1*1 = 3
		{NewBigIntVec(3, intsToBigInts([]int{1, 0, 1})...), NewBigIntVec(3, intsToBigInts([]int{1, 1, 1})...), big.NewInt(2)}, // 1*1 + 0*1 + 1*1 = 2
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.a.Dot(test.b)
			if test.expected.Cmp(actual) != 0 {
				t.Fatalf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestBigIntVector_Mul(t *testing.T) {
	tests := []struct {
		a        *BigIntVector
		b        *BigIntMatrix
		result   *BigIntVector
		expected *BigIntVector
	}{
		{NewBigIntVec(3, intsToBigInts([]int{1, 0, 1})...), BigIntIdentity(3), NewBigIntVec(3), NewBigIntVec(3, intsToBigInts([]int{1, 0, 1})...)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Mul(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expected, test.result)
			}
		})
	}
}

func TestBigIntVector_Equals(t *testing.T) {
	tests := []struct {
		a, b     *BigIntVector
		expected bool
	}{
		{NewBigIntVec(3), NewBigIntVec(3), true},
		{NewBigIntVec(3), NewBigIntVec(4), false},
		{NewBigIntVec(3, intsToBigInts([]int{1, 0, 1})...), NewBigIntVec(3), false},
		{NewBigIntVec(3, intsToBigInts([]int{1, 0, 1})...), NewBigIntVec(3, intsToBigInts([]int{1, 0, 1})...), true},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.a.Equals(test.b) != test.expected {
				t.Fatalf("expected %v for %v == %v", test.expected, test.a, test.b)
			}
		})
	}
}

func TestBigIntVector_Set(t *testing.T) {
	tests := []struct {
		source, result *BigIntVector
	}{
		// {NewBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...), NewBigIntVec(5)},
		{BigIntIdentity(5).Row(2), NewBigIntVec(5)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for i := 0; i < test.source.Len(); i++ {
				test.result.Set(i, new(big.Int).Set(test.source.At(i))) // Use new(big.Int).Set to avoid aliasing
			}

			if !test.source.Equals(test.result) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.source, test.result)
			}
		})
	}
}

func TestBigIntVector_Set2(t *testing.T) {
	m := NewBigIntMat(5, 5)

	for i := 0; i < 5; i++ {
		row := m.Row(i)
		row.Set(i, big.NewInt(1))
	}

	expected := NewBigIntMat(5, 5, intsToBigInts([]int{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1})...)

	if !expected.Equals(m) {
		t.Fatalf("expected\n%v\nbut found\n%v", expected, m)
	}
}

func TestBigIntVector_SetVec(t *testing.T) {
	tests := []struct {
		original         *BigIntVector
		setToSlice       *BigIntVector
		index            int
		expectedOriginal *BigIntVector
	}{
		{NewBigIntVec(5, intsToBigInts([]int{1, 1, 1, 1, 1})...), NewBigIntVec(3, intsToBigInts([]int{0, 1, 0})...), 1, NewBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.SetVec(test.setToSlice, test.index)

			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestBigIntVector_NonzeroValues(t *testing.T) {
	tests := []struct {
		input    *BigIntVector
		expected map[int]*big.Int
	}{
		{BigIntIdentity(4).Row(2), map[int]*big.Int{2: big.NewInt(1)}},
		{NewBigIntMat(4, 6, intsToBigInts([]int{1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1})...).Row(0), map[int]*big.Int{0: big.NewInt(1), 1: big.NewInt(1), 3: big.NewInt(1)}},
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

func TestBigIntVector_Slice(t *testing.T) {
	tests := []struct {
		original         *BigIntVector
		i, len           int
		addToSlice       *BigIntVector
		expectedOriginal *BigIntVector
		expectedSlice    *BigIntVector
	}{
		{
			NewBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...),
			1,
			3,
			NewBigIntVec(3, intsToBigInts([]int{1, 1, 1})...),
			NewBigIntVec(5, intsToBigInts([]int{1, 1, 2, 1, 1})...),
			NewBigIntVec(3, intsToBigInts([]int{1, 2, 1})...),
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.len)
			sl.Add(CopyBigIntVec(sl), test.addToSlice)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("Slice test - expected slice:\n%v\nbut found:\n%v", test.expectedSlice, sl)
			}
			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("Slice test - expected original:\n%v\nbut found:\n%v", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestBigIntVector_Slice2(t *testing.T) {
	tests := []struct {
		original      *BigIntVector
		i1, len1      int
		i2, len2      int
		expectedSlice *BigIntVector
	}{
		{NewBigIntVec(7, intsToBigInts([]int{0, 0, 1, 0, 1, 0, 0})...), 1, 5, 1, 3, NewBigIntVec(3, intsToBigInts([]int{1, 0, 1})...)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i1, test.len1).Slice(test.i2, test.len2)

			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expectedSlice, sl)
			}

		})
	}
}

func TestBigIntVector_Negate(t *testing.T) {
	tests := []struct {
		x, expected *BigIntVector
	}{
		{NewBigIntVec(4, intsToBigInts([]int{0, 1, 0, 1})...), NewBigIntVec(4, intsToBigInts([]int{0, -1, 0, -1})...)},
		{NewBigIntVec(4, intsToBigInts([]int{0, 0, 1, 1})...), NewBigIntVec(4, intsToBigInts([]int{0, 0, -1, -1})...)},
		{NewBigIntVec(4, intsToBigInts([]int{2, -3, 0, 5})...), NewBigIntVec(4, intsToBigInts([]int{-2, 3, 0, -5})...)},
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

func TestTransposedBigIntVector_Set(t *testing.T) {
	tests := []struct {
		source, result *TransposedBigIntVector
	}{
		{NewTBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...), NewTBigIntVec(5)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for i := 0; i < test.source.Len(); i++ {
				test.result.Set(i, new(big.Int).Set(test.source.At(i))) // Use new(big.Int).Set to avoid aliasing
			}

			if !test.source.Equals(test.result) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.source, test.result)
			}
		})
	}
}

func TestTransposedBigIntVector_Set2(t *testing.T) {
	m := NewBigIntMat(5, 5)
	for j := 0; j < 5; j++ {
		col := m.Column(j)
		col.Set(j, big.NewInt(1))
	}

	expected := NewBigIntMat(5, 5, intsToBigInts([]int{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1})...)

	if !expected.Equals(m) {
		t.Fatalf("expected\n%v\nbut found\n%v", expected, m)
	}
}

func TestTransposedBigIntVector_MulVec(t *testing.T) {
	tests := []struct {
		a        *BigIntMatrix
		b        *BigIntVector
		result   *TransposedBigIntVector
		expected *BigIntVector // Expected is a normal vector because we transpose the result for comparison
	}{
		{BigIntIdentity(3), NewBigIntVec(3, intsToBigInts([]int{0, 1, 0})...), NewTBigIntVec(3), NewBigIntVec(3, intsToBigInts([]int{0, 1, 0})...)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.MulVec(test.a, test.b.T())
			if !test.result.T().Equals(test.expected) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expected, test.result.T())
			}
		})
	}
}

func TestTransposedBigIntVector_Add(t *testing.T) {
	tests := []struct {
		a, b, result *TransposedBigIntVector
		expected     *TransposedBigIntVector
	}{
		{NewTBigIntVec(3, intsToBigInts([]int{0, 1, 0})...), NewTBigIntVec(3, intsToBigInts([]int{1, 0, 0})...), NewTBigIntVec(3), NewTBigIntVec(3, intsToBigInts([]int{1, 1, 0})...)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.result.Add(test.a, test.b)
			if !test.result.Equals(test.expected) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expected, test.result)
			}
		})
	}
}

func TestTransposedBigIntVector_NonzeroValues(t *testing.T) {
	tests := []struct {
		input    *TransposedBigIntVector
		expected map[int]*big.Int
	}{
		{BigIntIdentity(4).Column(2), map[int]*big.Int{2: big.NewInt(1)}},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := test.input.NonzeroValues()
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expected, actual)
			}
		})
	}
}

func TestTransposedBigIntVector_Slice(t *testing.T) {
	tests := []struct {
		original         *TransposedBigIntVector
		i, len           int
		addToSlice       *TransposedBigIntVector
		expectedOriginal *TransposedBigIntVector
		expectedSlice    *TransposedBigIntVector
	}{
		{
			NewTBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...),
			1,
			3,
			NewTBigIntVec(3, intsToBigInts([]int{1, 1, 1})...),
			NewTBigIntVec(5, intsToBigInts([]int{1, 1, 2, 1, 1})...),
			NewTBigIntVec(3, intsToBigInts([]int{1, 2, 1})...),
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i, test.len)
			sl.Add(CopyTBigIntVec(sl), test.addToSlice)
			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedSlice, sl)
			}
			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected \n%v\n but found \n%v\n", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestTransposedBigIntVector_Slice2(t *testing.T) {
	tests := []struct {
		original      *TransposedBigIntVector
		i1, len1      int
		i2, len2      int
		expectedSlice *TransposedBigIntVector
	}{
		{NewTBigIntVec(7, intsToBigInts([]int{0, 0, 1, 0, 1, 0, 0})...), 1, 5, 1, 3, NewTBigIntVec(3, intsToBigInts([]int{1, 0, 1})...)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sl := test.original.Slice(test.i1, test.len1).Slice(test.i2, test.len2)

			if !sl.Equals(test.expectedSlice) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expectedSlice, sl)
			}

		})
	}
}

func TestTransposedBigIntVector_SetVec(t *testing.T) {
	tests := []struct {
		original         *TransposedBigIntVector
		setToSlice       *TransposedBigIntVector
		index            int
		expectedOriginal *TransposedBigIntVector
	}{
		{NewTBigIntVec(5, intsToBigInts([]int{1, 1, 1, 1, 1})...), NewTBigIntVec(3, intsToBigInts([]int{0, 1, 0})...), 1, NewTBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...)},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.original.SetVec(test.setToSlice, test.index)

			if !test.original.Equals(test.expectedOriginal) {
				t.Fatalf("expected\n%v\nbut found\n%v", test.expectedOriginal, test.original)
			}
		})
	}
}

func TestBigIntVector_JSON(t *testing.T) {
	v := NewBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...)

	bs, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual BigIntVector
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !v.Equals(&actual) {
		t.Fatalf("expected\n%v\nbut found\n%v", v, actual)
	}
}

func TestTransposedBigIntVector_JSON(t *testing.T) {
	v := NewTBigIntVec(5, intsToBigInts([]int{1, 0, 1, 0, 1})...)

	bs, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}

	var actual TransposedBigIntVector
	err = json.Unmarshal(bs, &actual)
	if err != nil {
		t.Fatalf("expected no error found:%v", err)
	}
	if !v.Equals(&actual) {
		t.Fatalf("expected\n%v\nbut found\n%v", v, actual)
	}
}

func TestTransposedBigIntVector_Negate(t *testing.T) {
	tests := []struct {
		x, expected *TransposedBigIntVector
	}{
		{NewTBigIntVec(4, intsToBigInts([]int{0, 1, -2, 3})...), NewTBigIntVec(4, intsToBigInts([]int{0, -1, 2, -3})...)},
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
