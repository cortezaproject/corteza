import { ChartOptions, ChartOptionsRegistry } from './base'
import { FrameDefinition } from '../../frame'
import { Apply } from '../../../../cast'
import moment from 'moment'

export class BasicChartOptions extends ChartOptions {
  public labelColumn: string = ''
  public dataColumns: Array<{ name: string; label?: string }> = []

  constructor (o?: BasicChartOptions | Partial<BasicChartOptions>) {
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
      data: this.getData(dataframes[0], dataframes),
      options: {
        title: {
          display: !!this.title,
          text: this.title,
        },
        legend: {
          display: this.showLegend,
          labels: {
            // This more specific font property overrides the global property
            fontFamily: "'Poppins-Regular'",
          },
        },
        responsive: true,
        maintainAspectRatio: false,
        scales: {},
        tooltips: {
          enabled: this.showTooltips,
          displayColors: false,
          intersect: !['bar', 'line'].includes(this.type),
          callbacks: {},
          titleFontFamily: "'Poppins-Regular'",
          bodyFontFamily: "'Poppins-Regular'",
          footerFontFamily: "'Poppins-Regular'",
        },
        plugins: {
          colorschemes: {
            scheme: this.colorScheme,
            reverse: true,
          },
        },
      }
    }

    if (['bar', 'line'].includes(this.type)) {
      const {
        label: xLabel,
        type: xType,
        unit,
      } = this.xAxis

      const {
        label: yLabel,
        type: yType = 'linear',
        position = 'left',
        beginAtZero = true,
        stepSize,
        min,
        max,
      } = this.yAxis

      config.options.scales = {
        xAxes: [{
          type: xType || undefined,
          offset: true,
          time: {
            unit,
            round: true,
            minUnit: 'day',
          },
          scaleLabel: {
            display: !!xLabel,
            labelString: xLabel,
          },
          ticks: {
            autoSkip: false,
          }
        }],

        yAxes: [{
          display: true,
          type: yType,
          position,
          scaleLabel: {
            display: !!yLabel,
            labelString: yLabel,
          },
          ticks: {
            beginAtZero,
            stepSize: stepSize ? parseFloat(stepSize) : undefined,
            min: min ? parseFloat(min) : undefined,
            max: max ? parseFloat(max) : undefined,
          },
        }],
      }
    } else {
      config.options.tooltips.callbacks = {
        label: this.makeLabel,
      }
    }


    return config
  }

  getColIndex (dataframe: FrameDefinition, col: string) {
    if (!dataframe || !dataframe.columns) return -1

    return dataframe.columns.findIndex(({ name }) => name === col)
  }

  makeLabel ({ datasetIndex, index }: any, { datasets, labels }: any): string {
    const dataset = datasets[datasetIndex]
    const total = dataset.data.reduce((acc: string, v: string) => {
      return parseFloat(v) ? acc + parseFloat(v) : acc
    }, 0)

    let suffix = `(${total.toFixed(2)})%`
    if (total) {
      suffix = `(${((dataset.data[index] * 100) / total).toFixed(2)}%)`
    }

    return `${labels[index]}: ${dataset.data[index]} ${suffix}`
  }

  getData (localDataframe: FrameDefinition, dataframes: Array<FrameDefinition>) {
    const datasets: any[] = []
    let labels: any[] = []

    // Get datasets
    if (this.dataColumns.length && localDataframe.rows) {
      for (const { name } of this.dataColumns) {
        // Assume localDataframe has the dataColumn
        let columnIndex = this.getColIndex(localDataframe, name)

        // If dataColumn is in localDataframe, then set that value
        let data = localDataframe.rows.map(r => {
          return columnIndex < 0 ? undefined : r[columnIndex]
        })

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
              data[refRowIndex] = df.rows[0][columnIndex]
            }
          })
        }

        datasets.push({
          label: name,
          data,
        })
      }
    }

    // Get labels, if dimensions type is not time
    if (this.labelColumn && localDataframe) {
      const columnIndex = this.getColIndex(localDataframe, this.labelColumn)
      if (columnIndex < 0) {
        throw new Error(`Column ${this.labelColumn} not found`)
      }

      if (localDataframe.rows) {
        for (const row of localDataframe.rows) {
          let label = row[columnIndex] || (!this.xAxis.skipMissing ? this.xAxis.defaultValue : undefined)

          if (this.xAxis.type === 'time') {
            label = label ? moment(label).toDate() : undefined
          }

          labels.push(label)
        }
      }
    }

    if (this.xAxis.skipMissing) {
      labels.forEach((label, index) => {
        if (!label) {
          datasets.forEach(ds => {
            ds.data.splice(index, 1)
          })
        }
      })

      labels = labels.filter(label => label)
    }

    return { datasets, labels }
  }
}

ChartOptionsRegistry.set('bar', BasicChartOptions)
ChartOptionsRegistry.set('line', BasicChartOptions)
ChartOptionsRegistry.set('pie', BasicChartOptions)
ChartOptionsRegistry.set('doughnut', BasicChartOptions)

