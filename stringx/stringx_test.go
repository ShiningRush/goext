package stringx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasItem(t *testing.T) {
	tests := []struct {
		caseDesc  string
		giveArray []string
		giveItem  string
		wantRet   bool
	}{
		{
			caseDesc:  "existed",
			giveArray: []string{"fine", "bad"},
			giveItem:  "bad",
			wantRet:   true,
		},
		{
			caseDesc:  "not existed",
			giveArray: []string{"fine", "bad"},
			giveItem:  "not",
			wantRet:   false,
		},
		{
			caseDesc:  "nil",
			giveArray: nil,
			giveItem:  "not",
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
		giveSrc  []string
		giveDest []string
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"a", "c"},
			wantRet:  true,
		},
		{
			caseDesc: "not is",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"d", "c"},
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
		giveSrc  []string
		giveDest []string
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []string{"a", "c"},
			giveDest: []string{"a", "b", "c"},
			wantRet:  true,
		},
		{
			caseDesc: "self",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"a", "b", "c"},
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
		giveSrc  []string
		giveDest []string
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []string{"a", "c"},
			giveDest: []string{"a", "b", "c"},
			wantRet:  true,
		},
		{
			caseDesc: "not is",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"a", "b", "c"},
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
		giveA    []string
		giveB    []string
		wantRet  []string
	}{
		{
			caseDesc: "sanity",
			giveA:    []string{"a", "b", "c"},
			giveB:    []string{"a", "c", "e"},
			wantRet:  []string{"a", "c"},
		},
		{
			caseDesc: "a-empty",
			giveA:    nil,
			giveB:    []string{"a", "c"},
			wantRet:  nil,
		},
		{
			caseDesc: "b-empty",
			giveA:    []string{"a", "b", "c"},
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
