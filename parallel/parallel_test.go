package parallel

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	tests := []struct {
		caseDesc    string
		giveWork    Work
		giveWorkNum int
		wantErrs    []error
	}{
		{
			caseDesc: "sanity-default-worker",
			giveWork: func(workerIdx int) error {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return nil
			},
		},
		{
			caseDesc: "sanity-1-worker",
			giveWork: func(workerIdx int) error {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return nil
			},
			giveWorkNum: 1,
		},
		{
			caseDesc: "sanity-100-worker",
			giveWork: func(workerIdx int) error {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return nil
			},
			giveWorkNum: 100,
		},
		{
			caseDesc: "all-failed",
			giveWork: func(workerIdx int) error {
				return errors.New("expected err")
			},
			giveWorkNum: 2,
			wantErrs: []error{
				errors.New("expected err"),
				errors.New("expected err"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			var ops []OptionOp
			if tc.giveWorkNum > 0 {
				ops = append(ops, WorkerNumber(tc.giveWorkNum))
			}

			ret := Do(tc.giveWork, ops...)

			assert.Equal(t, tc.wantErrs, ret)
		})
	}
}

func TestStreamDo(t *testing.T) {
	tests := []struct {
		caseDesc         string
		giveInputs       []interface{}
		giveWork         StreamWork
		giveWorkNum      int
		giveIgnoreResult bool
		wantRets         []*StreamPayload
	}{
		{
			caseDesc: "sanity-default-worker",
			giveInputs: []interface{}{
				"item1",
				"item2",
				"item3",
			},
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload{
				{Result: "item1"},
				{Result: "item2"},
				{Result: "item3"},
			},
		},
		{
			caseDesc: "sanity-default-worker-ignore-result",
			giveInputs: []interface{}{
				"item1",
				"item2",
				"item3",
			},
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			giveIgnoreResult: true,
		},
		{
			caseDesc: "sanity-1-worker",
			giveInputs: []interface{}{
				"item1",
				"item2",
			},
			giveWorkNum: 1,
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload{
				{Result: "item1"},
				{Result: "item2"},
			},
		},
		{
			caseDesc: "sanity-100-worker",
			giveInputs: []interface{}{
				"item1",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload{
				{Result: "item1"},
			},
		},
		{
			caseDesc: "all-failed",
			giveInputs: []interface{}{
				"item1",
				"item2",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return nil, errors.New("expected err")
			},
			wantRets: []*StreamPayload{
				{Result: nil, Err: errors.New("expected err")},
				{Result: nil, Err: errors.New("expected err")},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			var ops []OptionOp
			if tc.giveIgnoreResult {
				ops = append(ops, IgnoreResult())
			}
			if tc.giveWorkNum > 0 {
				ops = append(ops, WorkerNumber(tc.giveWorkNum))
			}

			session := StreamDo(tc.giveWork, ops...)

			for _, v := range tc.giveInputs {
				session.Send(v)
			}

			rets := session.ReceivedPayloads()
			assert.ElementsMatch(t, rets, tc.wantRets)
		})
	}
}

func TestStreamDo_ReceiveFromChan(t *testing.T) {
	tests := []struct {
		caseDesc         string
		giveInputs       []interface{}
		giveWork         StreamWork
		giveWorkNum      int
		giveIgnoreResult bool
		wantRets         []*StreamPayload
	}{
		{
			caseDesc: "sanity-default-worker",
			giveInputs: []interface{}{
				"item1",
				"item2",
				"item3",
			},
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload{
				{Result: "item1"},
				{Result: "item2"},
				{Result: "item3"},
			},
		},
		{
			caseDesc: "sanity-default-worker-ignore-result",
			giveInputs: []interface{}{
				"item1",
				"item2",
				"item3",
			},
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			giveIgnoreResult: true,
		},
		{
			caseDesc: "sanity-1-worker",
			giveInputs: []interface{}{
				"item1",
				"item2",
			},
			giveWorkNum: 1,
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload{
				{Result: "item1"},
				{Result: "item2"},
			},
		},
		{
			caseDesc: "sanity-100-worker",
			giveInputs: []interface{}{
				"item1",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload{
				{Result: "item1"},
			},
		},
		{
			caseDesc: "all-failed",
			giveInputs: []interface{}{
				"item1",
				"item2",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item interface{}) (ret interface{}, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return nil, errors.New("expected err")
			},
			wantRets: []*StreamPayload{
				{Result: nil, Err: errors.New("expected err")},
				{Result: nil, Err: errors.New("expected err")},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			var ops []OptionOp
			ops = append(ops, ReceiveDataFromChan())
			if tc.giveIgnoreResult {
				ops = append(ops, IgnoreResult())
			}
			if tc.giveWorkNum > 0 {
				ops = append(ops, WorkerNumber(tc.giveWorkNum))
			}

			session := StreamDo(tc.giveWork, ops...)

			go func() {
				for _, v := range tc.giveInputs {
					session.Send(v)
				}
				session.CompleteSend()
			}()

			var rets []*StreamPayload
			for v := range session.ReceiveChan() {
				rets = append(rets, v)
			}
			assert.ElementsMatch(t, rets, tc.wantRets)
		})
	}
}
