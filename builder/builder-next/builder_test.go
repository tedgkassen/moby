package buildkit

import (
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"gotest.tools/assert"
)

func makeTestBuildCachePruneOptions(untilDuration string) types.BuildCachePruneOptions {
	arg := filters.Arg("until", untilDuration)
	return types.BuildCachePruneOptions{
		false, 50000, filters.NewArgs(arg),
	}
}


func TestToBuildKitPruneInfoUntilFormats(t *testing.T) {
	now := time.Now().UTC()
	tenHours, _ := time.ParseDuration("24h")
	inTenHoursTs := now.Add(tenHours).Unix()
	inTenHoursTime := time.Unix(inTenHoursTs, 0)
	testCases := []string{
		"24h",
		strconv.Itoa(int(inTenHoursTs)),
		inTenHoursTime.Format(time.RFC3339),
		inTenHoursTime.Format(time.RFC3339Nano),
		inTenHoursTime.Format("2006-01-02T15:04:05"),
		inTenHoursTime.Format("2006-01-02T15:04:05.999999999"),
	}

	getCaseResult := func (testCase string) time.Duration {
		pruneOptions := makeTestBuildCachePruneOptions(testCase)
		result, err := toBuildkitPruneInfo(pruneOptions)
		if err != nil {
			t.Fatalf("An until argument of format %s should be accepted", testCase)
		}
		return result.KeepDuration
	}


	for i, testCase := range(testCases) {
		result := getCaseResult(testCase)
		resultSecs := result.Seconds()
		assert.Assert(t, resultSecs >= 84399 && resultSecs <= 86400, testCase, resultSecs, i)
	}

	testCases = []string{
		inTenHoursTime.Format("2006-01-02Z07:00"),
		inTenHoursTime.Format("2006-01-02"),
	}

	for i, testCase := range(testCases) {
		result := getCaseResult(testCase)
		resultHours := int(result.Hours())
		truncated := 24 - now.Hour()
		assert.Assert(t, resultHours >= (truncated - 1) && resultHours <= truncated, testCase, resultHours, i)
	}
}
