<script>
import Function from './Function.vue'

export default {
  extends: Function,

  data () {
    return {
      showFunctionList: false,
      functionRef: 'error-handler',
    }
  },

  methods: {
    async getFunctionTypes () {
      this.functions = [
        {
          ref: 'error-handler',
          kind: 'function',
          meta: {
            short: 'Handle error',
          },
          parameters: [],
          results: [
            {
              name: 'error',
              types: [
                'Any',
              ],
            },
            {
              name: 'errorMessage',
              types: [
                'String',
              ],
            },
            {
              name: 'errorStepID',
              types: [
                'Integer',
              ],
            },
          ],
        },
      ]
    },
  },
}
</script>
