package intx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasItem(t *testing.T) {
	tests := []struct {
		caseDesc  string
		giveArray []int
		giveItem  int
		wantRet   bool
	}{
		{
			caseDesc:  "existed",
			giveArray: []int{1, 2},
			giveItem:  1,
			wantRet:   true,
		},
		{
			caseDesc:  "not existed",
			giveArray: []int{1, 2},
			giveItem:  3,
			wantRet:   false,
		},
		{
			caseDesc:  "nil",
			giveArray: nil,
			giveItem:  1,
			wantRet:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := HasItem(tc.giveArray, tc.giveItem)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIsSuperset(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveSrc  []int
		giveDest []int
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []int{1, 2, 3},
			giveDest: []int{1, 3},
			wantRet:  true,
		},
		{
			caseDesc: "not is",
			giveSrc:  []int{1, 2, 3},
			giveDest: []int{3, 4},
			wantRet:  false,
		},
		{
			caseDesc: "empty",
			giveSrc:  nil,
			giveDest: nil,
			wantRet:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := IsSuperset(tc.giveSrc, tc.giveDest)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIsSubset(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveSrc  []int
		giveDest []int
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []int{1, 3},
			giveDest: []int{1, 2, 3},
			wantRet:  true,
		},
		{
			caseDesc: "self",
			giveSrc:  []int{1, 2, 3},
			giveDest: []int{1, 2, 3},
			wantRet:  true,
		},
		{
			caseDesc: "empty",
			giveSrc:  nil,
			giveDest: nil,
			wantRet:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := IsSubset(tc.giveSrc, tc.giveDest)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIsProperSubset(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveSrc  []int
		giveDest []int
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []int{1, 3},
			giveDest: []int{1, 2, 3},
			wantRet:  true,
		},
		{
			caseDesc: "not is",
			giveSrc:  []int{1, 2, 3},
			giveDest: []int{1, 2, 3},
			wantRet:  false,
		},
		{
			caseDesc: "empty",
			giveSrc:  nil,
			giveDest: nil,
			wantRet:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := IsProperSubset(tc.giveSrc, tc.giveDest)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveA    []int
		giveB    []int
		wantRet  []int
	}{
		{
			caseDesc: "sanity",
			giveA:    []int{1, 2, 3},
			giveB:    []int{1, 2, 5},
			wantRet:  []int{1, 2},
		},
		{
			caseDesc: "a-empty",
			giveA:    nil,
			giveB:    []int{1, 3},
			wantRet:  nil,
		},
		{
			caseDesc: "b-empty",
			giveA:    []int{1, 2, 3},
			giveB:    nil,
			wantRet:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := Intersect(tc.giveA, tc.giveB)
			assert.ElementsMatch(t, tc.wantRet, ret)
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveA    []int
		giveB    []int
		wantRet  []int
	}{
		{
			caseDesc: "sanity",
			giveA:    []int{1, 2, 3},
			giveB:    []int{1, 2, 5},
			wantRet:  []int{3},
		},
		{
			caseDesc: "a-empty",
			giveA:    nil,
			giveB:    []int{1, 3},
			wantRet:  nil,
		},
		{
			caseDesc: "b-empty",
			giveA:    []int{1, 2, 3},
			giveB:    nil,
			wantRet:  []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := Diff(tc.giveA, tc.giveB)
			assert.ElementsMatch(t, tc.wantRet, ret)
		})
	}
}
