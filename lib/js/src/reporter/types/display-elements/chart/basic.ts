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
        show: true,
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
            show: this.tooltips.showAlways,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
            formatter: '{c} ({d}%)',
          },
          itemStyle: {
            borderRadius: 5,
            borderColor: '#FFFFFF',
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
          top: this.offset.default ? undefined : this.offset.top,
          right: this.offset.default ? undefined : this.offset.right,
          bottom: this.offset.default ? undefined : this.offset.bottom,
          left: this.offset.default ? undefined : this.offset.left,
        }
      })
    } else if (['bar', 'line'].includes(this.type)) {
      options.tooltip.trigger = 'axis'

      const {
        label: xLabel,
        type: xType = 'category',
        labelRotation: xLabelRotation = 0
      } = this.xAxis

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
            rotate: xLabelRotation,
          },
        },
      ]

      options.grid = {
        top: this.offset.default ? (this.title ? 70 : 45) : this.offset.top,
        right: this.offset.default ? 30 : this.offset.right,
        bottom: this.offset.default ? (xLabel ? 30 : 25) : this.offset.bottom,
        left: this.offset.default ? 30 : this.offset.left,
        containLabel: true,
      }

      const {
        label: yLabel,
        labelRotation: yLabelRotation = 0,
        type: yType = 'linear',
        position = 'left',
        labelPosition = 'end',
        beginAtZero,
        min,
        max,
      } = this.yAxis

      const tempYAxis = {
        name: yLabel,
        type: yType === 'linear' ? 'value' : 'log',
        position,
        nameGap: labelPosition === 'center' ? 25 : 7,
        nameLocation: labelPosition,
        min: beginAtZero ? 0 : min || undefined,
        max: max || undefined,
        axisLabel: {
          interval: 0,
          overflow: 'truncate',
          hideOverlap: true,
          rotate: yLabelRotation,
        },
        axisLine: {
          show: true,
          onZero: false,
        },
        nameTextStyle: {
          align: labelPosition === 'center' ? 'center' : position,
          padding: labelPosition !== 'center' ? (position === 'left' ? [0, 0, 2, -20] : [0, -20, 2, 0]) : undefined,
        }
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
            show: this.tooltips.showAlways,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
          },
          data: xType === 'time' ? labels.map((name, i) => {
            return [moment(name).valueOf() || undefined, data[i]]
          }) : data,
        }
      })
    }

    return {
      animation: !this.noAnimation,
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
        show: !this.legend.hide,
        type: this.legend.scrollable ? 'scroll' : 'plain',
        top: (this.legend.position.default ? (this.title ? 25 : undefined) : this.legend.position.top) || undefined,
        right: (this.legend.position.default ? undefined : this.legend.position.right) || undefined,
        bottom: (this.legend.position.default ? undefined : this.legend.position.bottom) || undefined,
        left: (this.legend.position.default ? this.legend.align || 'center' : this.legend.position.left) || 'auto',
        orient: this.legend.orientation || 'horizontal'
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
