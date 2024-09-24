<template>
  <div>
    <b-form-group>
      <b-form-checkbox v-model="f.options.presetWithAuthenticated">
        {{ $t('kind.user.presetWithCurrentUser') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      :label="$t('kind.user.roles.label')"
      label-class="text-primary"
    >
      <b-spinner
        v-if="preloadingRoles"
      />

      <c-input-role
        v-else
        v-model="currentRoles"
        :placeholder="$t('kind.user.roles.placeholder')"
        multiple
        @input="f.options.roles = $event.map(r => r.roleID)"
      />
    </b-form-group>

    <template v-if="f.isMulti">
      <b-form-group
        :label="$t('kind.select.optionType.label')"
        label-class="text-primary"
      >
        <b-form-radio-group
          v-model="f.options.selectType"
          :options="selectOptions"
          stacked
          @change="updateIsUniqueMultiValue"
        />
      </b-form-group>

      <b-form-group
        v-if="shouldAllowDuplicates"
      >
        <b-form-checkbox
          v-model="f.options.isUniqueMultiValue"
          :value="false"
          :unchecked-value="true"
        >
          {{ $t('kind.select.allow-duplicates') }}
        </b-form-checkbox>
      </b-form-group>
    </template>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
import base from './base'

const { CInputRole } = components

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    CInputRole,
  },

  extends: base,

  data () {
    return {
      selectOptions: [
        { text: this.$t('kind.select.optionType.default'), value: 'default', allowDuplicates: true },
        { text: this.$t('kind.select.optionType.multiple'), value: 'multiple' },
        { text: this.$t('kind.select.optionType.each'), value: 'each', allowDuplicates: true },
      ],

      preloadingRoles: true,
      currentRoles: [],
    }
  },

  computed: {
    shouldAllowDuplicates () {
      if (!this.f.isMulti) return false

      const { allowDuplicates } = this.selectOptions.find(({ value }) => value === this.f.options.selectType) || {}
      return !!allowDuplicates
    },
  },

  mounted () {
    if (this.f.options.roles.length) {
      this.preloadingRoles = true

      Promise.all(this.f.options.roles.map(roleID => {
        return this.$SystemAPI.roleRead({ roleID }).then(role => {
          this.currentRoles.push(role)
        })
      })).finally(() => {
        this.preloadingRoles = false
      })
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    getOptionKey ({ roleID }) {
      return roleID
    },

    updateIsUniqueMultiValue (value) {
      const { allowDuplicates = false } = this.selectOptions.find(({ value: v }) => v === value) || {}
      if (!allowDuplicates) {
        this.f.options.isUniqueMultiValue = true
      }
    },

    setDefaultValues () {
      this.selectOptions = []
      this.roleOptions = []
    },
  },
}
</script>
