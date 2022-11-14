import _ from 'lodash'
import moment from 'moment'
import {
  ChartConfig,
  ChartRenderer,
  Dimension,
  Metric,
  Report,
  dimensionFunctions,
  makeAlias,
  isRadialChart,
} from './util'

import {
  CortezaID,
  NoID,
  ISO8601Date,
  Apply,
} from '../../../cast'

export type PartialChart = Partial<BaseChart>

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

  merge (c: PartialChart) {
    let conf = { ...(c.config || {}) }
    Apply(this, c, CortezaID, 'chartID', 'namespaceID')
    Apply(this, c, String, 'name', 'handle')
    Apply(this, c, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, c, Boolean, 'canUpdateChart', 'canDeleteChart', 'canGrant')
    Apply(this, c, Object, 'config')

    if (typeof c.config === 'object') {
      // Verify & normalize
      let { renderer, reports, ...rest } = c.config

      if (renderer) {
        const { version } = renderer || {}

        if (version !== 'chart.js') {
          throw Error('notification.chart.unsupportedRenderer')
        }
      } else {
        renderer = { version: ChartRenderer.chartJS }
      }

      conf = { renderer, reports: reports || [], ...rest }
    }

    this.config = (conf ? _.merge(this.defConfig(), conf) : false) || this.config || this.defConfig()
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
    this.config.reports.map(({ moduleID, dimensions, metrics }) => {
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
      dimensions: dimensions?.map((d: Dimension) => dimensionFunctions.convert(d))[0],
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
  private processReporterResults (results: Array<object>, report: Report) {
    const dLabel = 'dimension_0'
    const { dimensions: [dimension] = [] } = report
    const isTimeDimension = !!(dimensionFunctions.lookup(dimension) || {}).time
    const hasRadialChart = report.metrics?.find(isRadialChart)
    let labels: Array<string> = []

    // helper to choose between eight the provided value, default value or a generic 'undefined'
    const pickValue = (val: unknown, { default: dDft }: Dimension): unknown => {
      if (val !== undefined && val !== null) return val
      if (dDft !== undefined && dDft !== null) return dDft
      return 'undefined'
    }

    // Skip missing values; if so requested
    if (dimension.skipMissing) {
      results = results.filter((r: any) => r[dLabel] !== null)
    }

    // Not a time dimensions, build set of labels
    if (!isTimeDimension || hasRadialChart) {
      labels = results.map((r: any) => pickValue(r[dLabel], dimension)) as Array<string>
    }

    // Build data sets
    const datasets = report.metrics?.map(m => {
      const alias = makeAlias({ field: m.field, aggregate: m.aggregate })
      const data = results.map((r: any) => {
        const y: any = r[m.field === 'count' ? m.field : alias]
        if (!isTimeDimension || hasRadialChart) {
          return pickValue(y, dimension)
        }
        return { y, t: moment(pickValue(r[dLabel], dimension) as string).toDate() }
      })

      // Any sub class has the ability to define how the dataset looks like.
      // this comes in handy when we want to support charts with different definitions.
      return this.makeDataset(m, dimension, data, alias)
    })

    return {
      labels: this.processLabels(labels, dimension),
      datasets,
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
    copy.config?.reports?.map(r => {
      const { moduleID } = r
      if (moduleID) {
        r.moduleID = getModuleID(moduleID)
      }
      return r
    })
    return copy
  }

  defDimension (): Dimension {
    return Object.assign({}, { conditions: {}, meta: {} })
  }

  defMetrics (): Metric {
    return Object.assign({}, {})
  }

  defReport (): Report {
    return Object.assign({}, {
      moduleID: null,
      filter: null,
      dimensions: [this.defDimension()],
      metrics: [this.defMetrics()],
      yAxis: {},
    })
  }

  defConfig (): ChartConfig {
    return Object.assign({}, {
      colorScheme: undefined,
      reports: [this.defReport()],
      renderer: {
        version: ChartRenderer.chartJS,
      },
    })
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'compose:chart'
  }
}

export { chartUtil } from './util'
