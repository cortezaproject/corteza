<script>
import Function from './Function'
import { components } from '@cortezaproject/corteza-vue'

export default {
  extends: Function,

  methods: {
    async getFunctionTypes () {
      this.functions = components.promptDefinitions || []
    },
  },
}
</script>
