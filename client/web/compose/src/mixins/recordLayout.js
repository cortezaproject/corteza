export default {
  data () {
    return {
      resizeObserver: null,
      columnWrapClass: '',
    }
  },

  methods: {
    initializeResizeObserver (el) {
      this.resizeObserver = new ResizeObserver((entries) => {
        for (let entry of entries) {
          // Handle the resize event here
          this.applyColumnClasses(entry.contentRect.width)
        }
      })

      this.resizeObserver.observe(el)
    },

    applyColumnClasses (width) {
      const breakpoints = {
        xs: 576,
        md: 768,
        lg: 992,
        xl: 1200,
      }

      const columnClasses = {
        xs: 'col-12',
        md: 'col-6',
        lg: 'col-4',
        xl: 'col-3',
      }

      let columnClass

      switch (true) {
        case width <= breakpoints.xs:
          columnClass = columnClasses.xs
          break
        case width > breakpoints.xs && width <= breakpoints.md:
          columnClass = columnClasses.md
          break
        case width > breakpoints.md && width <= breakpoints.lg:
          columnClass = columnClasses.lg
          break
        default:
          columnClass = columnClasses.xl
          break
      }

      this.columnWrapClass = `field-col ${columnClass}`
    },

    destroyEvents (el) {
      if (this.resizeObserver) {
        this.resizeObserver.unobserve(el)
        this.resizeObserver.disconnect()
      }
    },
  },
}
