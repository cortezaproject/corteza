import { ChartOptions, ChartOptionsRegistry } from './base'
import { FrameDefinition } from '../../frame'
import { makeDataLabel } from '../../../../compose/types/chart/util'
import { makeTipper } from '../../../../compose/types/chart/chartjs/plugins'
import { Apply } from '../../../../cast'
const ChartJS = require('chart.js')

export class FunnelChartOptions extends ChartOptions {
  public labelColumn: string = ''
  public dataColumns: Array<{ name: string; label?: string }> = []

  constructor (o?: FunnelChartOptions | Partial<FunnelChartOptions>) {
    super(o)

    if (!o) return

    Apply(this, o, String, 'labelColumn')

    if (o.dataColumns) {
      this.dataColumns = o.dataColumns || []
    }
  }

  getChartConfiguration (dataframes: Array<FrameDefinition>) {
    const config = {
      type: this.type,
      data: {
        labels: this.getLabels(dataframes[0]),
        datasets: this.getDatasets(dataframes[0], dataframes),
      },
      options: {
        title: {
          display: !!this.title,
          text: this.title,
          labels: {
            // This more specific font property overrides the global property
            fontFamily: "'Poppins-Regular'",
          },
        },
        legend: {
          display: this.showLegend
        },
        responsive: true,
        maintainAspectRatio: false,
        sort: 'desc',
        tooltips: {
          enabled: true,
          displayColors: false,
          callbacks: {
            label: this.makeLabel,
          },
          titleFontFamily: "'Poppins-Regular'",
          bodyFontFamily: "'Poppins-Regular'",
          footerFontFamily: "'Poppins-Regular'",
        },
        plugins: {},
      },
    }

    config.options.plugins = {
      colorschemes: {
        scheme: this.colorScheme,
        custom: (e: Array<string>) => {
          config.data.datasets[0].backgroundColor = config.data.labels.map((label, index) => {
            return e[index]
          })

          return e
        },
      },
      tipper: makeTipper(ChartJS.Tooltip, {}),
      datalabels: {
        display: false,
      },
    }


    return config
  }

  getColIndex (dataframe: FrameDefinition, col: string) {
    if (!dataframe || !dataframe.columns) return -1

    return dataframe.columns.findIndex(({ name }) => name === col)
  }

  getLabels (localDataframe: FrameDefinition) {
    const labels = []

    if (this.labelColumn && localDataframe) {
      const columnIndex = this.getColIndex(localDataframe, this.labelColumn)
      if (columnIndex < 0) {
        throw new Error(`Column ${this.labelColumn} not found`)
      }

      if (localDataframe.rows) {
        for (const row of localDataframe.rows) {
          labels.push(row[columnIndex])
        }
      }
    }

    return labels
  }

  makeLabel ({ datasetIndex, index }: any, { datasets, labels }: any): string {
    const dataset = datasets[datasetIndex]
    const total = dataset.data.reduce((acc: number, v: string) => acc + parseFloat(v), 0)

    let suffix = `(${total.toFixed(2)})%`
    if (total) {
      suffix = `(${((dataset.data[index] * 100) / total).toFixed(2)}%)`
    }

    return `${labels[index]}: ${dataset.data[index]} ${suffix}`
  }

  getDatasets (localDataframe: FrameDefinition, dataframes: Array<FrameDefinition>) {
    const chartDataset = []

    if (this.dataColumns.length && localDataframe.rows) {
      // Create dataset for each dataColumn
      for (const { name } of this.dataColumns) {
        // Assume localDataframe has the dataColumn
        let columnIndex = this.getColIndex(localDataframe, name)

        // If dataColumn is in localDataframe, then set that value
        const data: any = localDataframe.rows.map(r => {
          return columnIndex < 0 ? undefined : parseFloat(r[columnIndex] || '0') || 0
        })

        // Otherwise check other dataframes for that columnn
        if (columnIndex < 0) {
          dataframes.slice(1).forEach(df => {
            const { relColumn = '', refValue = '' } = df

            // Get column that is referenced by relColumn
            const relColumnIndex = this.getColIndex(localDataframe, relColumn)
            if (relColumnIndex < 0) {
              throw new Error(`Column ${relColumn} not found`)
            }

            if (!localDataframe.rows) {
              throw new Error(`Local rows not found`)
            }

            // Get row index that matches refValue
            const refRowIndex = localDataframe.rows.findIndex(row => row[relColumnIndex] === refValue)
            if (refRowIndex < 0) {
              throw new Error(`Row that matches refRowIndex ${refValue} not found`)
            }

            columnIndex = this.getColIndex(df, name)
            if (columnIndex < 0) {
              throw new Error(`Column ${name} not found`)
            } else if (df.rows) {
              data[refRowIndex] = parseFloat(df.rows[0][columnIndex] || '0') || 0
            }
          })
        }

        const backgroundColor: string[] = []

        chartDataset.push({
          label: name,
          data,
          backgroundColor,
        })
      }
    }

    return chartDataset
  }
}

ChartOptionsRegistry.set('funnel', FunnelChartOptions)

