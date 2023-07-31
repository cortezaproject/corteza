<template>
  <div>
    <template v-if="feed.options">
      <b-form-group
        :label-cols="3"
        :label="$t('calendar.recordFeed.moduleLabel')"
        horizontal
        breakpoint="md"
        label-class="text-primary"
      >
        <b-input-group>
          <c-input-select
            v-model="feed.options.moduleID"
            :options="modules"
            :reduce="o => o.moduleID"
            default-value="0"
            :placeholder="$t('calendar.recordFeed.modulePlaceholder')"
            label="name"
          />
        </b-input-group>
      </b-form-group>

      <template v-if="module">
        <b-form-group
          :label-cols="3"
          :label="$t('calendar.recordFeed.colorLabel')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <c-input-color-picker
            v-model="feed.options.color"
            :translations="{
              modalTitle: $t('calendar.recordFeed.colorPicker'),
              cancelBtnLabel: $t('general:label.cancel'),
              saveBtnLabel: $t('general:label.saveAndClose')
            }"
          />
        </b-form-group>

        <b-form-group
          :label="$t('calendar.recordFeed.titleLabel')"
          :label-cols="3"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <c-input-select
            v-model="feed.titleField"
            :options="titleFields | optionizeFields"
            :reduce="o => o.value"
            label="text"
            :placeholder="$t('calendar.recordFeed.titlePlaceholder')"
          />
        </b-form-group>

        <b-form-group
          :label-cols="3"
          :label="$t('calendar.recordFeed.eventStartFieldLabel')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <b-form-select
            v-model="feed.startField"
            :options="dateFields | optionizeFields"
          >
            <template slot="first">
              <option
                disabled
                value=""
              >
                {{ $t('calendar.recordFeed.eventStartFieldPlaceholder') }}
              </option>
            </template>
          </b-form-select>

          <b-form-text
            v-if="hasMultiFields"
            class="test-multi-field-ntf"
          >
            {{ $t('calendar.recordFeed.noMultiFields') }}
          </b-form-text>
        </b-form-group>

        <b-form-group
          :label-cols="3"
          :label="$t('calendar.recordFeed.eventEndFieldLabel')"
          horizontal
          breakpoint="md"
          label-class="text-primary"
        >
          <b-form-select
            v-model="feed.endField"
            :options="dateFields | optionizeFields"
          >
            <template slot="first">
              <option value="">
                {{ $t('calendar.recordFeed.eventEndFieldPlaceholder') }}
              </option>
            </template>
          </b-form-select>

          <b-form-text
            v-if="hasMultiFields"
            class="test-multi-field-ntf"
          >
            {{ $t('calendar.recordFeed.noMultiFields') }}
          </b-form-text>

          <b-form-checkbox
            v-model="feed.allDay"
            class="mt-3"
            :value="true"
            :unchecked-value="false"
          >
            {{ $t('calendar.recordFeed.eventAllDay') }}
          </b-form-checkbox>
        </b-form-group>

        <br>

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

    moduleOptions () {
      return [
        { moduleID: '0', name: this.$t('calendar.recordFeed.modulePlaceholder') },
        ...this.modules,
      ]
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
        .filter(f => ['String', 'Email', 'Url'].includes(f.kind))
        .sort((a, b) => a.label.localeCompare(b.label))
    },

    /**
     * Determines available date fields based on the given module.
     * Currently ignores multi-fields
     * @returns {Array}
     */
    dateFields () {
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
      ].filter(f => f.kind === 'DateTime' && !f.isMulti)
    },
  },
}
</script>
