<template>
  <b-card
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    body-class="p-0"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('custom-js.title') }}
      </h4>
    </template>

    <c-ace-editor
      v-model="settings[0].value"
      lang="javascript"
      height="500px"
      font-size="14px"
      show-line-numbers
      :border="false"
      :show-popout="false"
    />

    <template #footer>
      <c-button-submit
        v-if="canManage"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="onSubmit"
      />
    </template>
  </b-card>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CAceEditor } = components

export default {
  name: 'cui-cdns-editor',

  components: {
    CAceEditor,
  },

  props: {
    settings: {
      type: Array,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canManage: {
      type: Boolean,
      required: true,
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', { 'ui.cdn-scripts': this.settings[0].value })
    },
  },
}
</script>
