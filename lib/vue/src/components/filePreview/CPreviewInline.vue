<template>
  <div class="inline h-100">
    <main class="h-100">
      <component
        :is="previewType"
        v-bind="$attrs"
        :max-pages="1"
        :inline="inline"
        :initial-scale="1.5"
        v-on="$listeners"
      />
    </main>
  </div>
</template>

<script lang="js">
import { InlineIMG as IMG, PDF } from './common/types'
import base from './common/base.vue'

export default {
  components: {
    IMG,
    PDF,
  },

  extends: base,
}
</script>
