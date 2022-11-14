<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <b-form
      @submit.prevent="$emit('submit', settings)"
    >
      <b-form-group
        :label="$t('sidebar.title')"
      >
        <b-form-checkbox
          v-model="sidebar.hideNamespaceList"
        >
          {{ $t('sidebar.hide-namespace-list') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="sidebar.hideNamespaceListLink"
        >
          {{ $t('sidebar.hide-namespace-list-link') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        :label="$t('record-toolbar.title')"
        class="mb-0"
      >
        <b-form-checkbox
          v-model="recordToolbar.hideSubmit"
        >
          {{ $t('record-toolbar.hide-submit') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="recordToolbar.hideDelete"
        >
          {{ $t('record-toolbar.hide-delete') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="recordToolbar.hideEdit"
        >
          {{ $t('record-toolbar.hide-edit') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="recordToolbar.hideNew"
        >
          {{ $t('record-toolbar.hide-new') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="recordToolbar.hideClone"
        >
          {{ $t('record-toolbar.hide-clone') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="recordToolbar.hideBack"
        >
          {{ $t('record-toolbar.hide-back') }}
        </b-form-checkbox>
      </b-form-group>
    </b-form>

    <template #footer>
      <c-submit-button
        class="float-right"
        :disabled="!canManage"
        :processing="processing"
        :success="success"
        @submit="onSubmit"
      />
    </template>
  </b-card>
</template>

<script>
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CComposeEditorUI',

  i18nOptions: {
    namespaces: 'compose.settings',
    keyPrefix: 'editor.ui',
  },

  components: {
    CSubmitButton,
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
      sidebar: {},
      recordToolbar: {},
    }
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        this.sidebar = settings['compose.ui.sidebar'] || {}
        this.recordToolbar = settings['compose.ui.record-toolbar'] || {}
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', {
        'compose.ui.sidebar': this.sidebar,
        'compose.ui.record-toolbar': this.recordToolbar,
      })
    },
  },
}
</script>
