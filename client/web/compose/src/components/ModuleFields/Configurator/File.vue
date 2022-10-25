<template>
  <div>
    <b-form-group
      :label="$t('kind.file.view.maxSizeLabel')"
    >
      <b-input-group>
        <b-form-input
          v-model="f.options.maxSize"
          type="number"
          number
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
      v-if="enablePreviewStyling"
      class="mb-0"
    >
      <b-form-checkbox
        v-model="f.options.hideFileName"
      >
        {{ $t('kind.file.view.showName') }}
      </b-form-checkbox>
    </b-form-group>

    <template v-if="enablePreviewStyling">
      <hr>

      <h5 class="mb-2">
        {{ $t('kind.file.view.previewStyle') }}
      </h5>

      <b-row
        align-v="center"
        class="mb-2 mt-2"
      >
        <b-col
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.height')"
          >
            <b-input-group
              :append="$t('kind.file.view.px')"
            >
              <b-form-input
                v-model="f.options.height"
                type="number"
                number
                :placeholder="$t('kind.file.view.height')"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.width')"
          >
            <b-input-group
              :append="$t('kind.file.view.px')"
            >
              <b-form-input
                v-model="f.options.width"
                type="number"
                number
                :placeholder="$t('kind.file.view.width')"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.maxHeight')"
          >
            <b-input-group
              :append="$t('kind.file.view.px')"
            >
              <b-form-input
                v-model="f.options.maxHeight"
                type="number"
                number
                :placeholder="$t('kind.file.view.maxHeight')"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.maxWidth')"
          >
            <b-input-group
              :append="$t('kind.file.view.px')"
            >
              <b-form-input
                v-model="f.options.maxWidth"
                type="number"
                number
                :placeholder="$t('kind.file.view.maxWidth')"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.borderRadius')"
          >
            <b-input-group
              :append="$t('kind.file.view.px')"
            >
              <b-form-input
                v-model="f.options.borderRadius"
                type="number"
                number
                :placeholder="$t('kind.file.view.borderRadius')"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.background')"
          >
            <b-form-input
              v-model="f.options.backgroundColor"
              type="color"
              debounce="300"
            />
          </b-form-group>
        </b-col>

        <b-col
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.margin')"
          >
            <b-input-group
              :append="$t('kind.file.view.px')"
            >
              <b-form-input
                v-model="f.options.margin"
                type="number"
                number
                :placeholder="$t('kind.file.view.margin')"
              />
            </b-input-group>
          </b-form-group>
        </b-col>
      </b-row>
    </template>
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

    enablePreviewStyling () {
      const { mode } = this.f.options
      return (mode === 'single') || (mode === 'gallery')
    },
  },
}
</script>
