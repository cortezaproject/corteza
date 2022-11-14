<template>
  <div
    class="d-flex flex-column h-100"
  >
    <portal to="topbar-title">
      {{ $t(`request:kind.${kind}`) }}
    </portal>

    <component
      :is="kind"
      v-on="$listeners"
    />
  </div>
</template>

<script>
import Correct from './Correct'
import Delete from './Delete'
import Export from './Export'

export default {
  i18nOptions: {
    namespaces: 'request',
    keyPrefix: 'edit',
  },

  components: {
    Correct,
    Delete,
    Export,
  },

  props: {
    kind: {
      type: String,
      required: true,
    },
  },
}
</script>
