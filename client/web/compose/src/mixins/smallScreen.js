// Tracks whether the viewport is narrower than Bootstrap's md breakpoint,
// where blocks with many horizontal items switch to a compact layout
const query = '(max-width: 767.98px)'

export default {
  data () {
    return {
      isSmallScreen: false,
    }
  },

  mounted () {
    this.smallScreenQuery = window.matchMedia(query)
    this.isSmallScreen = this.smallScreenQuery.matches
    this.smallScreenQuery.addEventListener('change', this.onSmallScreenChange)
  },

  beforeDestroy () {
    this.smallScreenQuery.removeEventListener('change', this.onSmallScreenChange)
  },

  methods: {
    onSmallScreenChange ({ matches }) {
      this.isSmallScreen = matches
    },
  },
}
