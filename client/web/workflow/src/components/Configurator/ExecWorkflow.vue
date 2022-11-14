<script>
import Function from './Function'

export default {
  extends: Function,

  data () {
    return {
      showFunctionList: false,
      expressionResults: true,
      functionRef: 'exec-workflow',

      workflowOptions: [],
    }
  },

  mounted () {
    if (!this.workflowOptions.length) {
      this.searchWorkflows('', () => {})
    }
  },

  methods: {
    async getFunctionTypes () {
      this.functions = [
        {
          ref: 'exec-workflow',
          kind: 'function',
          meta: {
            short: 'Execute a workflow',
          },
          parameters: [
            {
              name: 'workflow',
              types: [
                'ID',
                'Handle',
              ],
              required: true,
            },
            {
              name: 'scope',
              types: [
                'Vars',
              ],
              required: true,
            },
          ],
          results: [],
        },
      ]
    },

    searchWorkflows (query = '', loading) {
      loading(true)

      this.$AutomationAPI.workflowList({ query, subWorkflow: 2 })
        .then(({ set }) => {
          this.workflowOptions = set.map(m => Object.freeze(m))
        })
        .finally(() => {
          loading(false)
        })
    },

    getWorkflowLabel ({ workflowID, handle, meta = {} }) {
      return meta.name || handle || workflowID
    },
  },
}
</script>
