import { ChartOptions, ChartOptionsRegistry } from './base'
import { FrameDefinition } from '../../frame'
import { Apply } from '../../../../cast'
import { getColorschemeColors } from '../../../../shared'

export class FunnelChartOptions extends ChartOptions {
  public labelColumn = ''
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
    const labels = this.getLabels(dataframes[0])
    const { data = [] } = this.getDatasets(dataframes[0], dataframes) || {}
    const colors = getColorschemeColors(this.colorScheme)

    return {
      animation: !this.noAnimation,
      title: {
        text: this.title,
        left: 'center',
        textStyle: {
          fontSize: 16,
        },
      },
      textStyle: {
        fontFamily: 'Poppins-Regular',
      },
      tooltip: {
        show: true,
        trigger: 'item',
        formatter: '{b} : {c} ({d}%)',
        appendToBody: true,
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
      series: [
        {
          type: 'funnel',
          name: this.labelColumn,
          sort: 'descending',
          width: '90%',
          label: {
            show: this.tooltips.showAlways,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
            formatter: '{c} ({d}%)',
          },
          data: labels.map((name, i) => {
            return { name, value: data[i], itemStyle: { color: colors[i] } }
          }),
          top: this.offset.default ? (this.title ? 60 : 35) : this.offset.top,
          right: this.offset.default ? '5%' : this.offset.right,
          bottom: this.offset.default ? '5%' : this.offset.bottom,
          left: this.offset.default ? '5%' : this.offset.left,
        },
      ],
    }
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

  getDatasets (localDataframe: FrameDefinition, dataframes: Array<FrameDefinition>): any {
    const chartDataset = []

    if (localDataframe && dataframes) {
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
                data[refRowIndex] = parseFloat(df.rows[0][columnIndex] || '0') || 0
              }
            })
          }

          chartDataset.push({
            label: name,
            data,
          })
        }
      }
    }

    return chartDataset[0]
  }
}

ChartOptionsRegistry.set('funnel', FunnelChartOptions)
