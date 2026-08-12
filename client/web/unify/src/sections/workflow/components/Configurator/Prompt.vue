<script>
import Function from './Function.vue'
import { components } from '@planetcrust/human-vue'

export default {
  extends: Function,

  methods: {
    async getFunctionTypes() {
      this.functions = components.promptDefinitions || []
    },
  },
}
</script>
