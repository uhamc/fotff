package rec

import (
	"fotff/tester"
	"os"
	"strconv"
	"testing"
)

type FotffMocker struct {
	StepsNum   int
	FirstFail  int
	runningPkg string
}

func TestMain(m *testing.M) {
	defer os.RemoveAll(".fotff")
	defer os.RemoveAll("logs")
	m.Run()
}

func (f *FotffMocker) TaskName() string {
	return "mocker"
}

func (f *FotffMocker) DoTestTask() ([]tester.Result, error) {
	return []tester.Result{{TestCaseName: f.TestCaseName(), Status: tester.ResultFail}}, nil
}

func (f *FotffMocker) DoTestCase(testCase string) (tester.Result, error) {
	running, _ := strconv.Atoi(f.runningPkg)
	if running >= f.FirstFail {
		return tester.Result{TestCaseName: testCase, Status: tester.ResultFail}, nil
	}
	return tester.Result{TestCaseName: testCase, Status: tester.ResultPass}, nil
}

func (f *FotffMocker) Flash(pkg string) error {
	f.runningPkg = pkg
	return nil
}

func (f *FotffMocker) LastIssue(pkg string) (string, error) {
	return "issue" + pkg, nil
}

func (f *FotffMocker) Steps(from, to string) (ret []string, err error) {
	fromIndex, _ := strconv.Atoi(from)
	toIndex, _ := strconv.Atoi(to)
	for i := fromIndex + 1; i <= toIndex; i++ {
		ret = append(ret, strconv.Itoa(i))
	}
	return
}

func (f *FotffMocker) GetNewer(cur string) (string, error) {
	return strconv.Itoa(f.StepsNum), nil
}

func (f *FotffMocker) TestCaseName() string {
	return "MOCK_FAILED_TEST_CASE"
}

func (f *FotffMocker) First() string {
	return "0"
}

func (f *FotffMocker) Last() string {
	return strconv.Itoa(f.StepsNum)
}

func TestFindOutTheFirstFail(t *testing.T) {
	tests := []struct {
		name   string
		mocker *FotffMocker
	}{
		{
			name:   "0-1(X)",
			mocker: &FotffMocker{StepsNum: 1, FirstFail: 1},
		},
		{
			name:   "0-1(X)-2-3-4",
			mocker: &FotffMocker{StepsNum: 4, FirstFail: 1},
		},
		{
			name:   "0-1-2(X)-3-4",
			mocker: &FotffMocker{StepsNum: 4, FirstFail: 2},
		},
		{
			name:   "0-1-2-3(X)-4",
			mocker: &FotffMocker{StepsNum: 4, FirstFail: 3},
		},
		{
			name:   "0-1-2-3-4(X)",
			mocker: &FotffMocker{StepsNum: 4, FirstFail: 4},
		},
		{
			name:   "0-1(X)-2-3-4-5",
			mocker: &FotffMocker{StepsNum: 5, FirstFail: 1},
		},
		{
			name:   "0-1-2(X)-3-4-5",
			mocker: &FotffMocker{StepsNum: 5, FirstFail: 2},
		},
		{
			name:   "0-1-2-3(X)-4-5",
			mocker: &FotffMocker{StepsNum: 5, FirstFail: 3},
		},
		{
			name:   "0-1-2-3-4(X)-5",
			mocker: &FotffMocker{StepsNum: 5, FirstFail: 4},
		},
		{
			name:   "0-1-2-3-4-5(X)",
			mocker: &FotffMocker{StepsNum: 5, FirstFail: 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err := FindOutTheFirstFail(tt.mocker, tt.mocker, tt.mocker.TestCaseName(), tt.mocker.First(), tt.mocker.Last())
			if err != nil {
				t.Errorf("err: expcect: <nil>, actual: %v", err)
			}
			expectIssue, _ := tt.mocker.LastIssue(strconv.Itoa(tt.mocker.FirstFail))
			if ret != expectIssue {
				t.Errorf("fotff result: expect: %s, actual: %s", expectIssue, ret)
			}
		})
	}
}
