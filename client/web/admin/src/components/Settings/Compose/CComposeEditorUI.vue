<template>
  <b-card
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

    <b-form
      @submit.prevent="$emit('submit', settings)"
    >
      <b-form-group
        :label="$t('sidebar.title')"
        label-class="text-primary"
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
        label-class="text-primary"
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

export default {
  name: 'CComposeEditorUI',

  i18nOptions: {
    namespaces: 'compose.settings',
    keyPrefix: 'editor.ui',
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
