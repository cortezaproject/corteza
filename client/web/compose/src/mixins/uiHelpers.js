export default {
  methods: {
    getBreakpoint () {
      let breakpoint
      const width = window.innerWidth

      if (width < 576) {
        breakpoint = 'xs'
      } else if (width >= 576 && width < 768) {
        breakpoint = 'sm'
      } else if (width >= 768 && width < 992) {
        breakpoint = 'md'
      } else if (width >= 992 && width < 1200) {
        breakpoint = 'lg'
      } else if (width >= 1200 && width < 1400) {
        breakpoint = 'xl'
      } else {
        breakpoint = 'xxl'
      }

      return breakpoint
    },
  },
}
