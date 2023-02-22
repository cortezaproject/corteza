// vue mixin for conditional fields
import { NoID } from '@cortezaproject/corteza-js'

export default {
  data () {
    return {
      conditions: [],
      evaluating: false,
    }
  },

  methods: {
    async evaluateExpressions () {
      if (!this.block.options.fieldConditions.length || this.$route.name === 'admin.pages.builder') return
      const { expressions, variables } = this.prepareFieldConditionsData()

      return this.$SystemAPI
        .expressionEvaluate({ variables, expressions })
        .then(res => {
          this.conditions = []

          Object.keys(res).forEach(v => {
            if (!res[v]) this.conditions.push(v)
          })
        }).catch(this.toastErrorHandler(this.$t('record.fieldConditions.error')))
    },

    prepareFieldConditionsData () {
      const expressions = {}
      const variables = { record: this.record.serialize() }
      this.block.options.fieldConditions.forEach(({ field, condition }) => {
        if (field && condition) {
          expressions[field] = condition
        }
      })
      return { expressions, variables }
    },

    canDisplay ({ fieldID, name }) {
      return !this.conditions.includes(fieldID !== NoID ? fieldID : name)
    },
  },
}
