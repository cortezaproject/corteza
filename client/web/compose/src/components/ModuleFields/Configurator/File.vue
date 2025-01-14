<template>
  <div>
    <b-form-group
      :label="$t('kind.file.view.maxSizeLabel')"
      label-class="text-primary"
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
      label-class="text-primary"
      class="mt-2"
    >
      <b-form-input
        v-model="f.options.mimetypes"
      />
    </b-form-group>

    <b-form-group
      :description="$t('kind.file.view.webcam.enable.footnote')"
      class="mt-2"
    >
      <b-form-checkbox
        v-model="f.options.enableWebcam"
      >
        {{ $t('kind.file.view.webcam.enable.label') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      :description="$t('kind.file.view.modeFootnote')"
      :label="$t('kind.file.view.modeLabel')"
      label-class="text-primary"
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

    <b-form-group>
      <b-form-checkbox
        v-if="enablePreviewStyling"
        v-model="f.options.hideFileName"
      >
        {{ $t('kind.file.view.showName') }}
      </b-form-checkbox>

      <b-form-checkbox
        v-model="f.options.clickToView"
      >
        {{ $t('kind.file.view.clickToView') }}
      </b-form-checkbox>

      <b-form-checkbox
        v-model="f.options.enableDownload"
      >
        {{ $t('kind.file.view.enableDownload') }}
      </b-form-checkbox>
    </b-form-group>

    <template v-if="enablePreviewStyling">
      <hr>

      <h5 class="mb-2">
        {{ $t('kind.file.view.previewStyle') }}
      </h5>

      <small>{{ $t('kind.file.view.description' ) }}</small>

      <b-row
        align-v="center"
        class="mb-2 mt-2"
      >
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.height')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.height"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.width')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.width"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.maxHeight')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.maxHeight"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.maxWidth')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.maxWidth"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.borderRadius')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.borderRadius"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.margin')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="f.options.margin"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.background')"
            label-class="text-primary"
          >
            <c-input-color-picker
              v-model="f.options.backgroundColor"
              :translations="{
                modalTitle: $t('kind.file.view.colorPicker'),
                light: $t('general:themes.labels.light'),
                dark: $t('general:themes.labels.dark'),
                cancelBtnLabel: $t('general:label.cancel'),
                saveBtnLabel: $t('general:label.saveAndClose')
              }"
              :theme-settings="themeSettings"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>
  </div>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    CInputColorPicker,
  },

  extends: base,

  computed: {
    modes () {
      return [
        { value: 'list', text: this.$t('kind.file.view.list') },
        { value: 'gallery', text: this.$t('kind.file.view.gallery') },
      ]
    },

    enablePreviewStyling () {
      const { mode } = this.f.options
      return mode === 'gallery'
    },

    themeSettings () {
      return this.$Settings.get('ui.studio.themes', [])
    },
  },
}
</script>
