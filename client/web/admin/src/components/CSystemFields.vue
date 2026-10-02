<template>
  <b-row>
    <b-col
      v-if="showDivider"
      cols="12"
    >
      <hr>
    </b-col>

    <b-col
      v-if="hasID"
      cols="12"
    >
      <b-form-group
        :label="$t('id')"
        label-class="text-primary"
      >
        {{ id }}
      </b-form-group>
    </b-col>

    <b-col
      v-for="(f, i) in systemFields"
      :key="i"
      cols="12"
    >
      <b-form-group
        v-if="getFieldValue(f) !== '0'"
        :label="$t(`general:system.field.${f}`)"
        label-class="text-primary"
        :data-test-id="`input-${generateTestID(f)}`"
      >
        <c-user-label
          v-if="isUserField(f)"
          :user="getUserForField(f)"
        />
        <template v-else>
          {{ getFieldValue(f) }}
        </template>
      </b-form-group>
    </b-col>
    <slot name="custom-field" />
  </b-row>
</template>

<script>
import Vue from 'vue'
import { components } from '@cortezaproject/corteza-vue'
import { getSystemFields, kebabize, userFields } from 'corteza-webapp-admin/src/lib/sysFields'

const { CUserLabel } = components

const userCache = new Map()
const inflight = new Map()

export default {
  name: 'CSystemFields',

  components: {
    CUserLabel,
  },

  props: {
    resource: {
      type: Object,
      required: true,
    },

    id: {
      type: String,
      default: '',
    },
  },

  data () {
    return {
      resolvedUsers: {},
    }
  },

  computed: {
    systemFields () {
      return getSystemFields(this.resource)
    },

    hasID () {
      return this.id && this.id !== '0'
    },

    hasVisibleSystemFields () {
      return this.systemFields.some(f => this.getFieldValue(f) !== '0')
    },

    showDivider () {
      return this.hasID || this.hasVisibleSystemFields
    },
  },

  watch: {
    resource: {
      immediate: true,
      deep: true,
      handler () {
        userFields.forEach(f => {
          const userID = this.resource[f]
          if (userID && userID !== '0') {
            this.resolveUser(userID)
          }
        })
      },
    },
  },

  methods: {
    generateTestID (field) {
      return kebabize(field)
    },

    isUserField (field) {
      return userFields.includes(field)
    },

    getFieldValue (field) {
      const isTimeValue = field.substring(field.length - 2) === 'At'
      return isTimeValue ? this.$options.filters.locFullDateTime(this.resource[field]) : this.resource[field]
    },

    getUserForField (field) {
      const userID = this.resource[field]
      if (!userID || userID === '0') return {}
      return this.resolvedUsers[userID] || { userID }
    },

    resolveUser (userID) {
      if (this.resolvedUsers[userID]) return
      if (userCache.has(userID)) {
        Vue.set(this.resolvedUsers, userID, userCache.get(userID))
        return
      }
      const pending = inflight.get(userID) || this.$SystemAPI.userRead({ userID })
        .then(user => {
          userCache.set(userID, user)
          return user
        })
        .finally(() => {
          inflight.delete(userID)
        })
      inflight.set(userID, pending)
      pending.then(user => {
        Vue.set(this.resolvedUsers, userID, user)
      }).catch(() => {
        Vue.set(this.resolvedUsers, userID, { userID })
      })
    },
  },
}
</script>
