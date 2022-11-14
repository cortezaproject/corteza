import { ChartOptions, ChartOptionsRegistry } from './base'
import { FrameDefinition } from '../../frame'
import { Apply } from '../../../../cast'
import { getColorschemeColors } from '../../../../shared'
import moment from 'moment'
export class BasicChartOptions extends ChartOptions {
  public labelColumn = ''
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
    const { labels, datasets = [] } = this.getData(dataframes[0], dataframes)

    const options: any = {
      series: [],
      xAxis: [],
      yAxis: [],
      tooltip: {
        show: this.showTooltips,
        appendToBody: true,
      },
    }

    if (['pie', 'doughnut'].includes(this.type)) {
      const startRadius = this.type === 'doughnut' ? 40 : 0
      const endRadius = 80
      const radiusLength = (endRadius - startRadius) / (datasets.length || 1)

      options.tooltip.trigger = 'item'

      options.series = datasets.map(({ label, data }, index) => {
        const sr = startRadius + (index * radiusLength)
        const er = startRadius + ((index + 1) * radiusLength)

        return {
          name: label,
          type: 'pie',
          radius: [`${sr}%`, `${er}%`],
          center: ['50%', '55%'],
          tooltip: {
            formatter: '{a}<br />{b} : {c} ({d}%)',
          },
          label: {
            show: false,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
            fontSize: 14,
          },
          itemStyle: {
            borderRadius: 5,
            borderColor: '#fff',
            borderWidth: 1,
          },
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)',
            },
          },
          data: labels.map((name, i) => {
            return { name, value: data[i] }
          }),
        }
      })
    } else if (['bar', 'line'].includes(this.type)) {
      options.tooltip.trigger = 'axis'

      const { label: xLabel, type: xType = 'category' } = this.xAxis

      options.xAxis = [
        {
          name: xLabel,
          nameLocation: 'center',
          nameGap: 30,
          type: xType,
          data: labels,
          axisLabel: {
            interval: 0,
            overflow: 'truncate',
            hideOverlap: true,
          },
        },
      ]

      options.grid = {
        top: this.title ? 60 : 35,
        bottom: xLabel ? 30 : 20,
        containLabel: true,
      }

      const {
        label: yLabel,
        type: yType = 'linear',
        position = 'left',
        beginAtZero,
        min,
        max,
      } = this.yAxis

      const tempYAxis = {
        name: yLabel,
        type: yType === 'linear' ? 'value' : 'log',
        position,
        nameLocation: 'center',
        nameGap: 30,
        min: beginAtZero ? 0 : min || undefined,
        max: max || undefined,
      }

      // If we provide undefined, log scale breaks
      if (tempYAxis.type === 'log') {
        delete tempYAxis.min
        delete tempYAxis.max
      }

      options.yAxis = [tempYAxis]

      options.series = datasets.map(({ label, data }) => {
        return {
          name: label,
          type: this.type,
          smooth: true,
          areaStyle: {},
          left: 'left',
          label: {
            show: false,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
            fontSize: 14,
          },
          data: xType === 'time' ? labels.map((name, i) => {
            return [moment(name).valueOf() || undefined, data[i]]
          }) : data,
        }
      })
    }

    return {
      title: {
        text: this.title,
        left: 'center',
        textStyle: {
          fontSize: 16,
        },
      },
      color: getColorschemeColors(this.colorScheme),
      textStyle: {
        fontFamily: 'Poppins-Regular',
      },
      legend: {
        show: this.showLegend,
        top: this.title ? 25 : undefined,
        type: 'scroll',
      },
      ...options,
    }
  }

  getColIndex (dataframe: FrameDefinition, col: string) {
    if (!dataframe || !dataframe.columns) return -1

    return dataframe.columns.findIndex(({ name }) => name === col)
  }

  getData (localDataframe: FrameDefinition, dataframes: Array<FrameDefinition>) {
    const datasets: any[] = []
    let labels: string[] = []

    if (localDataframe && dataframes) {
      // Get datasets
      if (this.dataColumns.length && localDataframe.rows) {
        for (const { name } of this.dataColumns) {
          // Assume localDataframe has the dataColumn
          let columnIndex = this.getColIndex(localDataframe, name)

          // If dataColumn is in localDataframe, then set that value
          const data = localDataframe.rows.map(r => {
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
                throw new Error('Local rows not found')
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
            const label = row[columnIndex] || (!this.xAxis.skipMissing ? this.xAxis.defaultValue : undefined)
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
    }

    return { datasets, labels }
  }
}

ChartOptionsRegistry.set('bar', BasicChartOptions)
ChartOptionsRegistry.set('line', BasicChartOptions)
ChartOptionsRegistry.set('pie', BasicChartOptions)
ChartOptionsRegistry.set('doughnut', BasicChartOptions)
