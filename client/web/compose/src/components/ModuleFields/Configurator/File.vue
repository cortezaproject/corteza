<template>
  <div>
    <b-form-group
      :label="$t('kind.file.view.maxSizeLabel')"
    >
      <b-input-group>
        <b-form-input
          v-model="f.options.maxSize"
          type="number"
        />
      </b-input-group>
    </b-form-group>

    <b-form-group
      :label="$t('kind.file.view.mimetypesLabel')"
      :description="$t('kind.file.view.mimetypesFootnote')"
      class="mt-2"
    >
      <b-form-input
        v-model="f.options.mimetypes"
      />
    </b-form-group>

    <b-form-group
      :description="$t('kind.file.view.modeFootnote')"
      :label="$t('kind.file.view.modeLabel')"
    >
      <b-form-radio-group
        v-model="f.options.mode"
        buttons
        button-variant="outline-secondary"
        size="sm"
        name="buttons2"
        :options="modes"
      />
    </b-form-group>

    <b-form-group
      v-if="enableFileNameHiding"
      class="mb-0"
    >
      <b-form-checkbox
        v-model="f.options.hideFileName"
      >
        {{ $t('kind.file.view.showName') }}
      </b-form-checkbox>
    </b-form-group>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  computed: {
    modes () {
      return [
        { value: 'list', text: this.$t('kind.file.view.list') },
        { value: 'grid', text: this.$t('kind.file.view.grid') },
        { value: 'single', text: this.$t('kind.file.view.single') },
        { value: 'gallery', text: this.$t('kind.file.view.gallery') },
      ]
    },

    enableFileNameHiding () {
      return (this.f.options.mode === 'single') || (this.f.options.mode === 'gallery')
    },
  },
}
</script>
