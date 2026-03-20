<script>
import Function from './Function.vue'
import { components } from '@cortezaproject/corteza-vue-next'

export default {
  extends: Function,

  methods: {
    async getFunctionTypes () {
      this.functions = components.promptDefinitions || []
    },
  },
}
</script>
