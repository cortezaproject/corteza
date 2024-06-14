import _ from 'lodash'
import {
  ChartConfig,
  Dimension,
  Metric,
  Report,
  dimensionFunctions,
  makeAlias,
  TemporalDataPoint,
  defFormatData,
} from './util'

import {
  CortezaID,
  NoID,
  ISO8601Date,
  Apply,
} from '../../../cast'

export type PartialChart = Partial<BaseChart>

// The default dataset post processing function to use.
// This one simply returns the current value.
const defaultFx = 'n'

/**
 * BaseChart represents a structure that stores any configuration data.
 * Any display and data rendering operations should be handled by any sub classes.
 */
export class BaseChart {
  public chartID = NoID
  public namespaceID = NoID
  public name = ''
  public handle = ''

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  public canUpdateChart = false
  public canDeleteChart = false
  public canGrant = false

  public config: ChartConfig = {}

  constructor (def: PartialChart = {}) {
    this.merge(def)
  }

  /**
   * The method performs post processing for each value in the given dataset.
   * It works with a simple equation written in javascript (example: n + m).
   * Available variables to use:
   * * n - current value
   * * m - previous value (undefined in case of the first element)
   * * r - entire data array.
   *
   * @param data Array of values in the given data set
   * @param m Metric for the given dataset
   */
  datasetPostProc (data: Array<number|TemporalDataPoint>, m: Metric): Array<number|TemporalDataPoint> {
    // Define a valid function to evaluate
    let fxRaw = (m.fx || defaultFx).trim()
    if (!fxRaw.startsWith('return')) {
      fxRaw = 'return ' + fxRaw
    }
    // eslint-disable-next-line no-new-func
    const fx = new Function('n', 'm', 'r', fxRaw)

    // Define a new array, so we don't alter the original one.
    const r = [...data]

    // Run postprocessing for all data in the given data set
    // There is a slight difference between temporal data points and categorical data points.
    if (data[0] instanceof Object) {
      // Temporal
      for (let i = 0; i < data.length; i++) {
        const a = data[i] as TemporalDataPoint
        const b = data[i - 1] as TemporalDataPoint|undefined

        const n = a.y
        let m: number|undefined
        if (i > 0) {
          m = b?.y
        }

        a.y = fx(n, m, r)
      }
    } else {
      // Categorical
      for (let i = 0; i < data.length; i++) {
        const n = data[i] as number
        let m: number|undefined
        if (i > 0) {
          m = data[i - 1] as number
        }
        data[i] = fx(n, m, r)
      }
    }

    return data
  }

  merge (c: PartialChart) {
    let conf = { ...(c.config || {}) }
    Apply(this, c, CortezaID, 'chartID', 'namespaceID')
    Apply(this, c, String, 'name', 'handle')
    Apply(this, c, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, c, Boolean, 'canUpdateChart', 'canDeleteChart', 'canGrant')
    Apply(this, c, Object, 'config')

    if (typeof c.config === 'object') {
      // Verify & normalize
      const { reports = [], ...rest } = c.config

      conf = { reports: reports || [], ...rest }
    }

    this.config = (conf ? _.merge(this.defConfig(), conf) : false) || this.config || this.defConfig()

    this.config.reports?.forEach(report => {
      const { dimensions = [], metrics = [] } = report || {}

      report.dimensions = dimensions.map(d => {
        // Legacy support
        if (d.modifier === 'auto') {
          d.timeLabels = true
          d.modifier = '(no grouping / buckets)'
        }

        if (d.field === 'created_at') {
          d.field = 'createdAt'
        }

        return _.merge(this.defDimension(), d)
      })

      report.metrics = metrics.map(m => _.merge(this.defMetric(), m))
    })
  }

  /**
   * Checks reports validity.
   * Validates dimensions and metrics.
   * If invalid it throws an error.
   */
  isValid () {
    if (!this.config.reports || !this.config.reports.length) {
      throw new Error('notification.chart.invalidConfig.missingReports')
    }

    this.config.reports.forEach(({ moduleID, dimensions, metrics }) => {
      if (!moduleID) {
        throw new Error('notification.chart.invalidConfig.missingModuleID')
      }

      // Expecting all dimensions to have defined fields
      dimensions?.forEach(this.dimCheck)

      // Expecting all metrics to have defined fields
      metrics?.forEach(this.mtrCheck)
    })

    return true
  }

  /**
   * Checks validity of dimensions.
   * If invalid it throws an error
   */
  dimCheck ({ field, modifier }: Dimension) {
    if (!field) {
      throw new Error('notification.chart.invalidConfig.missingDimensionsField')
    }
    if (!modifier) {
      throw new Error('notification.chart.invalidConfig.missingDimensionsModifier')
    }
  }

  /**
   * Checks validity of metrics.
   * If invalid it throws an error
   */
  mtrCheck ({ field, aggregate, type }: Metric) {
    if (!field) {
      throw new Error('notification.chart.invalidConfig.missingMetricsField')
    }
    if (field !== 'count' && !aggregate) {
      throw new Error('notification.chart.invalidConfig.missingMetricsAggregate')
    }
    if (!type) {
      throw new Error('notification.chart.invalidConfig.missingMetricsType')
    }
  }

  /**
   * Prepares params that the reporter can use for querying.
   */
  formatReporterParams ({ moduleID, metrics, dimensions, filter }: Report) {
    return {
      moduleID,
      filter,

      // Remove count (we'll get it anyway) and construct FUNC(ARG) params
      metrics: metrics?.filter((m: Metric) => m.field !== 'count').map((m: Metric) => `${m.aggregate}(${m.field}) AS ${makeAlias(m)}`).join(','),

      // Construct dimensions \w modifiers...
      dimensions: dimensions?.map(d => ({ field: 'createdAt', ...d })).map((d: Dimension) => dimensionFunctions.convert(d))[0],
    }
  }

  /**
   * Fetcher reports defined in the given configuration with the help of the provided
   * reporter.
   */
  async fetchReports ({ reporter }: { reporter(p: any): Promise<any> }) {
    const out: Array<any> = []

    // Prepare params & filter out invalid combos (formatReporterParams will return null on invalid params)
    const reports: any = this.config.reports?.map(this.formatReporterParams)
      // Send requests to reporter (API caller)
      .map(params => reporter(params))
      // Process each result
      .map((p: any, index: number) => p.then((results: any) => {
        results = results || []
        out[index] = this.processReporterResults(results, (this.config.reports || [])[index])
      }))

    // Wait for all requests to finish and return new promise, with results
    return Promise.all(reports).then(() => new Promise(resolve => {
      resolve(out)
    }))
  }

  /**
   * Processes provided report with it's results:
   * * skip missing values, if so requested,
   * * generate labels,
   * * creates dataset for the chart.
   */
  private processReporterResults (results: Array<object> = [], report: Report): object {
    const dLabel = 'dimension_0'
    const { dimensions: [dimension] = [] } = report
    let labels: Array<string> = []

    // helper to choose between eight the provided value, default value or a generic 'undefined'
    const pickValue = (val: unknown, { default: dDft }: Dimension): unknown => {
      return val || val === 0 ? val : dDft || 'undefined'
    }

    // Skip missing values; if so requested
    if (dimension.skipMissing) {
      results = results.filter((r: any) => r[dLabel] || r[dLabel] === 0)
    }

    // Not a time dimensions, build set of labels
    labels = results.map((r: any) => pickValue(r[dLabel], dimension)) as Array<string>

    // Build data sets
    const datasets = report.metrics?.map(m => {
      const alias = makeAlias({ field: m.field, aggregate: m.aggregate })
      const data = results.map((r: any) => {
        return pickValue(r[m.field === 'count' ? m.field : alias], dimension)
      })

      // Any sub class has the ability to define how the dataset looks like.
      // this comes in handy when we want to support charts with different definitions.
      return this.makeDataset(m, dimension, data, alias)
    })

    return {
      labels: this.processLabels(labels, dimension),
      datasets,
      dimension,
    }
  }

  processLabels (ll: Array<string>, d: Dimension) {
    return ll
  }

  makeDataset (m: Metric, d: Dimension, data: Array<number|any>, alias: string) {
    throw new Error('method.makeDataset.notImplemented')
  }

  makeOptions (data?: any) {
    throw new Error('method.makeOptions.notImplemented')
  }

  plugins (mm: Array<Metric>) {
    throw new Error('method.plugins.notImplemented')
  }

  baseChartType (datasets: Array<any>) {
    throw new Error('method.baseChartType.notImplemented')
  }

  /**
   * Performs chart export; used by exporter feature.
   */
  async export (findModuleByID: ({ namespaceID, moduleID }: { namespaceID: string; moduleID: string }) => Promise<any>) {
    const { namespaceID } = this
    const copy = new BaseChart(this)
    if (copy.config?.reports) {
      await Promise.all(copy.config.reports.map(async (r: any) => {
        const { moduleID } = r
        if (moduleID) {
          const module = await findModuleByID({ namespaceID, moduleID })
          r.moduleID = module.name
          return r
        } else {
          return null
        }
      })).then((a: any) => {
        return a
      })
    }
    return copy
  }

  /**
   * Performs import; used by importer feature
   */
  import (getModuleID: (moduleID: string) => string) {
    const copy = new BaseChart(this)
    copy.config.reports = copy.config?.reports?.map(r => {
      const { moduleID } = r
      if (moduleID) {
        r.moduleID = getModuleID(moduleID)
      }
      return r
    })
    return copy
  }

  defDimension (): Dimension {
    return Object.assign({}, {
      conditions: {},
      meta: {},
      rotateLabel: 0,
    })
  }

  defMetric (): Metric {
    return Object.assign({}, {
      formatting: defFormatData(),
    })
  }

  defReport (): Report {
    return Object.assign({}, {
      moduleID: undefined,
      filter: '',
      dimensions: [this.defDimension()],
      metrics: [this.defMetric()],
      yAxis: {
        axisType: 'linear',
        axisPosition: 'left',
        labelPosition: 'end',
        rotateLabel: 0,
        formatting: defFormatData(),
      },
      tooltip: {},
      legend: {
        isScrollable: true,
        orientation: 'horizontal',
        align: 'center',
        position: {
          top: undefined,
          right: undefined,
          bottom: undefined,
          left: undefined,
          isDefault: true,
        },
      },
      offset: {
        top: '50',
        right: '30',
        bottom: '20',
        left: '30',
        isDefault: true,
      },
    })
  }

  defConfig (): ChartConfig {
    return Object.assign({}, {
      colorScheme: '',
      reports: [this.defReport()],
      noAnimation: false,
      toolbox: {
        saveAsImage: false,
        timeline: '',
      },
    })
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'compose:chart'
  }

  clone (): BaseChart {
    return new BaseChart(JSON.parse(JSON.stringify(this)))
  }
}

export { chartUtil } from './util'
