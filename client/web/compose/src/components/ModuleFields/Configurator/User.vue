<template>
  <div>
    <b-form-group>
      <b-form-checkbox v-model="f.options.presetWithAuthenticated">
        {{ $t('kind.user.presetWithCurrentUser') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      v-if="f.options.roles"
      :label="$t('kind.user.roles.label')"
    >
      <vue-select
        v-model="f.options.roles"
        :options="roleOptions"
        :reduce="role => role.roleID"
        option-value="roleID"
        option-text="name"
        :close-on-select="false"
        append-to-body
        :placeholder="$t('kind.user.roles.placeholder')"
        multiple
        label="name"
        class="bg-white"
      />
    </b-form-group>

    <b-form-group
      v-if="f.isMulti"
    >
      <label class="d-block">{{ $t('kind.select.optionType.label') }}</label>
      <b-form-radio-group
        v-model="f.options.selectType"
        :options="selectOptions"
        stacked
      />
    </b-form-group>
  </div>
</template>

<script>
import base from './base'
import { VueSelect } from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    VueSelect,
  },

  extends: base,

  data () {
    return {
      selectOptions: [
        { text: this.$t('kind.select.optionType.default'), value: 'default' },
        { text: this.$t('kind.select.optionType.multiple'), value: 'multiple' },
        { text: this.$t('kind.select.optionType.each'), value: 'each' },
      ],
      roleOptions: [],
    }
  },

  mounted () {
    this.$SystemAPI.roleList().then(({ set: roles = [] }) => {
      this.roleOptions = roles
    })
  },
}
</script>
