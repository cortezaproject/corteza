<template>
  <div>
    <div :class="classes">
      <template v-if="field.isMulti">
        <template v-for="(u, index) in resolvedUsers">
          <c-user-label
            v-if="u"
            :key="`u-${index}`"
            :user="u"
          />
          <span
            v-else
            :key="`na-${index}`"
          >{{ $t('kind.user.na') }}</span>
          <span
            v-if="index !== resolvedUsers.length - 1"
            :key="`d-${index}`"
          >{{ field.options.multiDelimiter }}</span>
        </template>
      </template>

      <template v-else>
        <c-user-label
          v-if="singleUser"
          :user="singleUser"
        />
        <template v-else>
          {{ $t('kind.user.na') }}
        </template>
      </template>
    </div>
  </div>
</template>

<script>
import base from './base'
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
const { CUserLabel } = components

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    CUserLabel,
  },

  extends: base,

  computed: {
    ...mapGetters({
      findByID: 'user/findByID',
    }),

    resolvedUsers () {
      const resolve = (u) => {
        if (!u || u === NoID) return undefined
        if (typeof u === 'string') return this.findByID(u)
        return u
      }

      return (this.value || []).map(resolve)
    },

    singleUser () {
      const v = this.value
      if (!v || v === NoID) return undefined
      if (typeof v === 'string') return this.findByID(v)
      return v
    },
  },
}
</script>
