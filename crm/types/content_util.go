package types

import (
	"regexp"
	"strings"
)

type (
	ContentReport struct {
		Metrics    []ContentReportMetric
		Dimensions []ContentReportDimensions
	}

	ContentReportMetric struct {
		Alias      string
		Expression string
	}

	ContentReportDimensions struct {
		Alias     string
		Field     string
		Modifiers []string
	}
)

var (
	contentReportMetricScanRE    *regexp.Regexp
	contentReportDimensionScanRE *regexp.Regexp
)

func init() {
	contentReportMetricScanRE = regexp.MustCompile("^(?:(\\w+):)?(\\w+)$")
	contentReportDimensionScanRE = regexp.MustCompile("^(?:(\\w+):)?(\\w+)((?:\\|?\\w+)+)?$")
}

func (r *ContentReport) ScanMetrics(metrics ...string) {
	r.Metrics = make([]ContentReportMetric, len(metrics))
	for i := 0; i < len(metrics); i++ {
		r.Metrics[i].Scan(metrics[i])
	}
}

func (r *ContentReport) ScanDimensions(dimensions ...string) {
	r.Dimensions = make([]ContentReportDimensions, len(dimensions))
	for i := 0; i < len(dimensions); i++ {
		r.Dimensions[i].Scan(dimensions[i])
	}
}

func (m *ContentReportMetric) Scan(metric string) {
	if match := contentReportMetricScanRE.FindStringSubmatch(metric); len(match) == 3 {
		m.Alias = match[1]
		m.Expression = match[2]
	}
}

func (d *ContentReportDimensions) Scan(dimension string) {
	if match := contentReportDimensionScanRE.FindStringSubmatch(dimension); len(match) == 4 {
		d.Alias = match[1]
		d.Field = match[2]

		if len(match[3]) > 0 {
			d.Modifiers = strings.Split(match[3][1:], "|")
		}
	}
}
