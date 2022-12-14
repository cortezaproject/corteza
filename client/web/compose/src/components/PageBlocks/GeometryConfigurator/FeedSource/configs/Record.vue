<template>
  <div>
    <template v-if="feed.options">
      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('geometry.recordFeed.moduleLabel')"
      >
        <b-input-group>
          <b-form-select
            v-model="feed.options.moduleID"
            :options="modules"
            value-field="moduleID"
            text-field="name"
          >
            <template slot="first">
              <option value="0">
                {{ $t('geometry.recordFeed.modulePlaceholder') }}
              </option>
            </template>
          </b-form-select>
        </b-input-group>
      </b-form-group>

      <template v-if="module">
        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('geometry.recordFeed.colorLabel')"
        >
          <b-input-group>
            <b-form-input
              v-model="feed.options.color"
              style="max-width: 50px;"
              type="color"
              debounce="300"
            />
          </b-input-group>
        </b-form-group>

        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('geometry.recordFeed.titleLabel')"
        >
          <b-form-select
            v-model="feed.titleField"
            :options="titleFields | optionizeFields"
          >
            <template slot="first">
              <option
                disabled
                value=""
              >
                {{ $t('geometry.recordFeed.titlePlaceholder') }}
              </option>
            </template>
          </b-form-select>
        </b-form-group>

        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('geometry.recordFeed.displayMarker')"
        >
          <b-form-checkbox
            v-model="feed.displayMarker"
            name="display-marker"
            switch
            size="lg"
          />
        </b-form-group>

        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('geometry.recordFeed.displayPolygon')"
        >
          <b-form-checkbox
            v-model="feed.displayPolygon"
            name="display-marker"
            switch
            size="lg"
          />
        </b-form-group>

        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('geometry.recordFeed.geometryFieldLabel')"
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

        <br>

        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('calendar.recordFeed.prefilterLabel')"
        >
          <b-form-textarea
            v-model="feed.options.prefilter"
            :value="true"
            :placeholder="$t('calendar.recordFeed.prefilterPlaceholder')"
          />
        </b-form-group>
      </template>
    </template>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'block',
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
