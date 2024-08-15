<template>
  <div>
    <b-form-group
      :label="$t('provider.label')"
      label-class="text-primary"
    >
      <b-input-group>
        <b-form-input v-model="value.name" />
      </b-input-group>
    </b-form-group>

    <div>
      <div class="mb-2">
        <h5>
          {{ $t('value.label') }}
        </h5>
        <span class="text-muted">
          {{ $t('value.description') }}
        </span>
      </div>

      <c-ace-editor
        v-model="value.cdnScript"
        lang="javascript"
        height="500px"
        font-size="14px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />
    </div>
  </div>
</template>
<script>
import { components } from '@cortezaproject/corteza-vue'
const { CAceEditor } = components

export default {
  name: 'CSystemCdnProviderEditor',

  i18nOptions: {
    namespaces: 'system.cdns',
    keyPrefix: 'editor.cdns-provider.form',
  },

  components: {
    CAceEditor,
  },

  props: {
    value: {
      type: Object,
      required: true,
      default: () => ({
        'name': '',
        'cdnScript': '',
      }),
    },
  },

}
</script>
<style scoped lang="scss">
</style>
