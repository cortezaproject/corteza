<template>
  <div>
    <template v-if="feed.options">
      <b-form-group
        :label-cols="3"
        :label="$t('geometry.recordFeed.moduleLabel')"
        horizontal
        breakpoint="md"
        label-class="text-primary"
      >
        <b-input-group>
          <c-input-select
            v-model="feed.options.moduleID"
            :options="modules"
            :reduce="o => o.moduleID"
            :placeholder="$t('calendar.recordFeed.modulePlaceholder')"
            default-value="0"
            label="name"
          />
        </b-input-group>
      </b-form-group>

      <template v-if="module">
        <b-form-group
          :label-cols="3"
          :label="$t('geometry.recordFeed.geometryFieldLabel')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <b-form-select
            v-model="feed.geometryField"
            :options="geometryFields | optionizeFields"
          >
            <template slot="first">
              <option
                disabled
                value=""
              >
                {{ $t('geometry.recordFeed.geometryFieldPlaceholder') }}
              </option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group
          :label-cols="3"
          :label="$t('geometry.recordFeed.titleLabel')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <c-input-select
            v-model="feed.titleField"
            :options="titleFields | optionizeFields"
            :reduce="o => o.value"
            label="text"
            :placeholder="$t('geometry.recordFeed.titlePlaceholder')"
          />
        </b-form-group>

        <b-form-group
          :label-cols="3"
          :label="$t('calendar.recordFeed.prefilterLabel')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <b-form-textarea
            v-model="feed.options.prefilter"
            :value="true"
            :placeholder="$t('calendar.recordFeed.prefilterPlaceholder')"
          />
        </b-form-group>

        <b-form-group
          :label-cols="3"
          :label="$t('geometry.recordFeed.colorLabel')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <c-input-color-picker
            v-model="feed.options.color"
            :translations="{
              modalTitle: $t('geometry.recordFeed.colorPicker'),
              cancelBtnLabel: $t('general:label.cancel'),
              saveBtnLabel: $t('general:label.saveAndClose')
            }"
          />
        </b-form-group>

        <b-form-group
          :label-cols="3"
          :label="$t('geometry.recordFeed.displayMarker')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <b-form-checkbox
            v-model="feed.displayMarker"
            name="display-marker"
            switch
            size="lg"
          />
        </b-form-group>

        <b-form-group
          :label-cols="3"
          :label="$t('geometry.recordFeed.displayPolygon')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <b-form-checkbox
            v-model="feed.displayPolygon"
            name="display-marker"
            switch
            size="lg"
          />
        </b-form-group>
      </template>
    </template>
  </div>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CInputColorPicker,
  },

  extends: base,

  computed: {
    /**
     * Finds the module, this feed configurator should use
     * @returns {Module|undefined}
     */
    module () {
      if (!(this.feed.options || {}).moduleID) {
        return
      }

      return this.modules.find(({ moduleID }) => moduleID === this.feed.options.moduleID)
    },

    /**
     * Determines if given module has any multi-fields
     * @returns {Boolean}
     */
    hasMultiFields () {
      if (!this.module) {
        return false
      }
      return this.module.fields.reduce((acc, { isMulti }) => acc || isMulti, false)
    },

    /**
     * Determines available title fields based on the given module.
     * @returns {Array}
     */
    titleFields () {
      if (!this.module) {
        return []
      }
      return [...this.module.fields]
        .filter(f => [
          'DateTime',
          'Select',
          'Number',
          'Bool',
          'String',
          'Record',
          'User',
        ].includes(f.kind))
        .sort((a, b) => a.label.localeCompare(b.label))
    },

    /**
     * Determines available geometry fields based on the given module.
     * Currently ignores multi-fields
     * @returns {Array}
     */
    geometryFields () {
      if (!this.module) {
        return []
      }

      const moduleFields = this.module.fields.slice().sort((a, b) => a.label.localeCompare(b.label))

      return [
        ...moduleFields,
        ...this.module.systemFields().map(sf => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        }),
      ].filter(f => f.kind === 'Geometry')
    },
  },
}
</script>
