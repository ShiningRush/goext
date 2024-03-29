package parallel

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		giveInputs       []string
		giveWork         StreamWork[string, string]
		giveWorkNum      int
		giveIgnoreResult bool
		wantRets         []*StreamPayload[string]
	}{
		{
			caseDesc: "sanity-default-worker",
			giveInputs: []string{
				"item1",
				"item2",
				"item3",
			},
			giveWork: func(workerIdx int, item string) (ret string, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload[string]{
				{Result: "item1"},
				{Result: "item2"},
				{Result: "item3"},
			},
		},
		{
			caseDesc: "sanity-default-worker-ignore-result",
			giveInputs: []string{
				"item1",
				"item2",
				"item3",
			},
			giveWork: func(workerIdx int, item string) (ret string, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			giveIgnoreResult: true,
		},
		{
			caseDesc: "sanity-1-worker",
			giveInputs: []string{
				"item1",
				"item2",
			},
			giveWorkNum: 1,
			giveWork: func(workerIdx int, item string) (ret string, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload[string]{
				{Result: "item1"},
				{Result: "item2"},
			},
		},
		{
			caseDesc: "sanity-100-worker",
			giveInputs: []string{
				"item1",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item string) (ret string, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item, nil
			},
			wantRets: []*StreamPayload[string]{
				{Result: "item1"},
			},
		},
		{
			caseDesc: "all-failed",
			giveInputs: []string{
				"item1",
				"item2",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item string) (ret string, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return "", errors.New("expected err")
			},
			wantRets: []*StreamPayload[string]{
				{Result: "", Err: errors.New("expected err")},
				{Result: "", Err: errors.New("expected err")},
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

type HashItem string

func (s HashItem) GetKey() string {
	return string(s)
}

func TestStreamDo_ReceiveFromChan(t *testing.T) {
	tests := []struct {
		caseDesc         string
		giveInputs       []interface{}
		giveWork         StreamWork[interface{}, string]
		giveWorkNum      int
		giveIgnoreResult bool
		wantRets         []*StreamPayload[string]
	}{
		{
			caseDesc: "sanity-default-worker",
			giveInputs: []interface{}{
				"item1",
				"item2",
				"item3",
			},
			giveWork: func(workerIdx int, item interface{}) (ret string, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item.(string), nil
			},
			wantRets: []*StreamPayload[string]{
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
			giveWork: func(workerIdx int, item interface{}) (ret string, err error) {
				assert.GreaterOrEqual(t, workerIdx, 0)

				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item.(string), nil
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
			giveWork: func(workerIdx int, item interface{}) (ret string, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return item.(string), nil
			},
			wantRets: []*StreamPayload[string]{
				{Result: "item1"},
				{Result: "item2"},
			},
		},
		{
			caseDesc: "sanity-100-worker",
			giveInputs: []interface{}{
				"item1",
				"item2",
				"item3",
				"item4",
				"item1",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item interface{}) (ret string, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return fmt.Sprintf("%s-%d", item, workerIdx), nil
			},
			wantRets: []*StreamPayload[string]{
				{Result: "item1-0"},
				{Result: "item2-1"},
				{Result: "item3-2"},
				{Result: "item4-3"},
				{Result: "item1-4"},
			},
		},
		{
			caseDesc: "all-failed",
			giveInputs: []interface{}{
				"item1",
				"item2",
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item interface{}) (ret string, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return "", errors.New("expected err")
			},
			wantRets: []*StreamPayload[string]{
				{Result: "", Err: errors.New("expected err")},
				{Result: "", Err: errors.New("expected err")},
			},
		},
		{
			caseDesc: "sanity-100-hash-worker",
			giveInputs: []interface{}{
				HashItem("item1"),
				HashItem("item2"),
				HashItem("item3"),
				HashItem("item4"),
				HashItem("item1"),
			},
			giveWorkNum: 100,
			giveWork: func(workerIdx int, item interface{}) (ret string, err error) {
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				return fmt.Sprintf("%s-%d", item, workerIdx), nil
			},
			wantRets: []*StreamPayload[string]{
				{Result: "item1-55"},
				{Result: "item2-17"},
				{Result: "item3-7"},
				{Result: "item4-36"},
				{Result: "item1-55"},
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

			var rets []*StreamPayload[string]
			for v := range session.ReceiveChan() {
				rets = append(rets, v)
			}
			assert.ElementsMatch(t, rets, tc.wantRets)
		})
	}
}
