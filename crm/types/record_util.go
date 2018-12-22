package types

import (
	"regexp"
	"strings"
)

type (
	RecordReport struct {
		Metrics    []RecordReportMetric
		Dimensions []RecordReportDimensions
	}

	RecordReportMetric struct {
		Alias      string
		Expression string
	}

	RecordReportDimensions struct {
		Alias     string
		Field     string
		Modifiers []string
	}
)

var (
	recordReportMetricScanRE    *regexp.Regexp
	recordReportDimensionScanRE *regexp.Regexp
)

func init() {
	recordReportMetricScanRE = regexp.MustCompile("^(?:(\\w+):)?(.+)$")
	recordReportDimensionScanRE = regexp.MustCompile("^(?:(\\w+):)?(\\w+)((?:\\|?\\w+)+)?$")
}

func (r *RecordReport) ScanMetrics(metrics ...string) {
	r.Metrics = make([]RecordReportMetric, len(metrics))
	for i := 0; i < len(metrics); i++ {
		r.Metrics[i].Scan(metrics[i])
	}
}

func (r *RecordReport) ScanDimensions(dimensions ...string) {
	r.Dimensions = make([]RecordReportDimensions, len(dimensions))
	for i := 0; i < len(dimensions); i++ {
		r.Dimensions[i].Scan(dimensions[i])
	}
}

func (m *RecordReportMetric) Scan(metric string) {
	if match := recordReportMetricScanRE.FindStringSubmatch(metric); len(match) == 3 {
		m.Alias = match[1]
		m.Expression = match[2]
	}
}

func (d *RecordReportDimensions) Scan(dimension string) {
	if match := recordReportDimensionScanRE.FindStringSubmatch(dimension); len(match) == 4 {
		d.Alias = match[1]
		d.Field = match[2]

		if len(match[3]) > 0 {
			d.Modifiers = strings.Split(match[3][1:], "|")
		}
	}
}
