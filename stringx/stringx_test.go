package stringx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHaveItem(t *testing.T) {
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
			ret := HaveItem(tc.giveArray, tc.giveItem)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}
