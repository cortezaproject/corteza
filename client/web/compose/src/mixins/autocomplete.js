import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  props: {
    record: {
      type: compose.Record,
      required: false,
      default: undefined,
    },

    page: {
      type: compose.Page,
      required: true,
    },
  },

  computed: {
    isRecordPage () {
      return this.page && this.page.moduleID !== NoID
    },

    recordAutoCompleteParams () {
      return this.processRecordAutoCompleteParams({})
    },

    visibilityAutoCompleteParams () {
      return this.processVisibilityAutoCompleteParams({})
    },
  },

  methods: {
    processRecordAutoCompleteParams ({ module = this.module, operators = false }) {
      const { fields = [] } = module || {}
      const moduleFields = fields.map(({ name }) => name)
      const userProperties = this.$auth.user.properties() || []

      const recordSuggestions = this.isRecordPage
        ? [
            ...(['ownerID', 'recordID'].map(value => ({ interpolate: true, value }))),
            {
              interpolate: true,
              value: 'record',
              properties: [
                ...(this.record.properties || []),
                { value: 'values', properties: Object.keys(this.record.values) || [] },
              ],
            },
          ]
        : []

      return [
        ...moduleFields,
        ...recordSuggestions,
        ...(operators ? ['AND', 'OR'] : []),
        { interpolate: true, value: 'userID' },
        { interpolate: true, value: 'user', properties: userProperties },
      ]
    },

    processVisibilityAutoCompleteParams ({ module = this.module }) {
      const { fields = [] } = module || {}
      const moduleFields = fields.map(({ name }) => name)
      const userProperties = this.$auth.user.properties() || []

      const recordSuggestions = this.isRecordPage
        ? [
            {
              value: 'record',
              properties: [
                ...(this.record.properties || []),
                { value: 'values', properties: Object.keys(this.record.values) || [] },
              ],
            },
          ]
        : []

      return [
        ...moduleFields,
        ...recordSuggestions,
        { value: 'user', properties: userProperties },
        { value: 'screen', properties: ['width', 'height', 'userAgent', 'breakpoint'] },
      ]
    },
  },
}
