import { makeDataLabel } from '../util'
interface PluginOptions {
  [_: string]: any;
}

export function makeTipper (Tooltip: any, options: PluginOptions = {}) {
  return {
    id: 'tipper',

    beforeRender: (chart: any) => {
      chart.options.tooltips.enabled = true
      // build a tooltip for each data set data node
      const tipperTips = chart.config.data.datasets.map((_: any, i: number) => {
        // ignore if dataset is hidden
        if (!chart.isDatasetVisible(i)) {
          return undefined
        }

        // getDatasetMeta provides element selector; dataset does not
        return chart.getDatasetMeta(i).data.map((sc: any) => {
          const opts: PluginOptions = {
            ...chart.options.tooltips,
            // force these values
            displayColors: false,
            backgroundColor: 'rgba(255,255,255,0.5)',
            ...options,
          }
          opts.callbacks = {
            ...opts.callbacks,
            labelTextColor: function () {
              return '#000000'
            },
            // We should only show the value when doing this
            beforeTitle: () => '',
            title: () => '',
            afterTitle: () => '',
            beforeBody: () => '',
            beforeLabel: () => '',
            afterLabel: () => '',
            afterBody: () => '',
            beforeFooter: () => '',
            footer: () => '',
            afterFooter: () => '',
          }

          return new Tooltip({
            _chart: chart.chart,
            _chartInstance: chart,
            _data: chart.data,
            _options: opts,
            _active: [sc],
          }, chart)
        })
      })

      chart.config.tipperTips = tipperTips.filter((tt: any) => tt)
      chart.options.tooltips.enabled = false
    },

    afterDraw: (chart: any, easing: any) => {
      // enable for drawing window
      chart.options.tooltips.enabled = true
      chart.config.tipperTips.forEach((ds: Array<any>) => {
        ds.forEach(tt => {
          tt.initialize()
          tt.update()
          tt.transition(easing).draw()
        })
      })
      chart.options.tooltips.enabled = false
    },
  }
}
