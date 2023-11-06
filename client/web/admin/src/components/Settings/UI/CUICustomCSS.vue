<template>
  <b-card
    body-class="p-0"
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <c-ace-editor
      v-model="customCSS"
      lang="css"
      height="300px"
      font-size="14px"
      show-line-numbers
      :border="false"
      :show-popout="true"
      @open="openEditorModal"
    />

    <b-modal
      id="custom-css-editor"
      v-model="showEditorModal"
      :title="$t('modal.editor')"
      cancel-variant="link"
      size="lg"
      :ok-title="$t('general:label.saveAndClose')"
      :cancel-title="$t('general:label.cancel')"
      body-class="p-0"
      @ok="saveCustomCSSInput"
      @hidden="resetCustomCSSInput"
    >
      <c-ace-editor
        v-model="modalCSSInput"
        lang="scss"
        height="500px"
        font-size="14px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />
    </b-modal>

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
  name: 'CUIEditorCustomCSS',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.custom-css',
  },

  components: {
    CAceEditor,
  },

  props: {
    settings: {
      type: Object,
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

  data () {
    return {
      customCSS: '',
      modalCSSInput: undefined,
      showEditorModal: false,
    }
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        this.customCSS = settings['ui.custom-css'] || ''
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', { 'ui.custom-css': this.customCSS })
    },

    openEditorModal () {
      this.modalCSSInput = this.customCSS
      this.showEditorModal = true
    },

    saveCustomCSSInput () {
      this.customCSS = this.modalCSSInput
    },

    resetCustomCSSInput () {
      this.modalCSSInput = undefined
    },
  },

}
</script>
