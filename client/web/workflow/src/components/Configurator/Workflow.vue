<template>
  <div>
    <b-form-group
      :label="$t('label')"
    >
      <b-form-input
        v-model="workflow.meta.name"
        data-test-id="input-label"
        :state="nameState"
        @input="$root.$emit('change-detected')"
      />
    </b-form-group>

    <b-form-group
      :label="$t('general:handle')"
    >
      <b-form-input
        v-model="workflow.handle"
        data-test-id="input-handle"
        :state="handleState"
        :placeholder="$t('workflow.placeholder-handle')"
        @input="$root.$emit('change-detected')"
      />
      <b-form-invalid-feedback
        data-test-id="input-handle-invalid-state"
        :state="handleState"
      >
        {{ $t('workflow.invalid-handle-characters') }}
      </b-form-invalid-feedback>
    </b-form-group>

    <b-form-group
      :label="$t('general:description')"
    >
      <b-form-textarea
        v-model="workflow.meta.description"
        data-test-id="input-description"
        @input="$root.$emit('change-detected')"
      />
    </b-form-group>

    <b-form-group
      :label="$t('workflow.run-as')"
      :description="$t('workflow.not-setup-properly')"
    >
      <vue-select
        :options="user.options"
        data-test-id="select-run-as"
        :get-option-label="getOptionLabel"
        :get-option-key="getOptionKey"
        :value="user.value"
        :calculate-position="calculateDropdownPosition"
        @search="search"
        @input="updateRunAs"
      />
    </b-form-group>

    <b-form-group>
      <b-form-checkbox
        v-model="workflow.enabled"
        data-test-id="checkbox-enable-workflow"
        @change="$root.$emit('change-detected')"
      >
        {{ $t('general:enabled') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      :description="$t('workflow.sub-workflow.description')"
    >
      <b-form-checkbox
        v-model="workflow.meta.subWorkflow"
        data-test-id="checkbox-sub-workflow"
        @change="$root.$emit('change-detected')"
      >
        {{ $t('workflow.sub-workflow.label') }}
      </b-form-checkbox>
    </b-form-group>
  </div>
</template>

<script>
import { debounce } from 'lodash'
import { VueSelect } from 'vue-select'
import { handle } from '@cortezaproject/corteza-vue'

export default {
  i18nOptions: {
    namespaces: 'configurator',
  },

  components: {
    VueSelect,
  },

  props: {
    workflow: {
      type: Object,
      default: () => {},
    },
  },

  data () {
    return {
      user: {
        options: [],
        value: undefined,

        filter: {
          query: null,
          limit: 10,
        },
      },
    }
  },

  computed: {
    nameState () {
      return this.workflow.meta.name ? null : false
    },

    handleState () {
      return handle.handleState(this.workflow.handle)
    },
  },

  created () {
    if (this.workflow.runAs) {
      this.fetchUsers()
      this.getUserByID()
    }
  },

  methods: {
    search: debounce(function (query) {
      if (query !== this.user.filter.query) {
        this.user.filter.query = query
        this.user.filter.page = 1
      }

      if (query) {
        this.fetchUsers()
      }
    }, 300),

    fetchUsers () {
      this.$SystemAPI.userList(this.user.filter)
        .then(({ set }) => {
          this.user.options = set.map(m => Object.freeze(m))
        })
    },

    async getUserByID () {
      if (this.workflow.runAs !== '0') {
        this.$SystemAPI.userRead({ userID: this.workflow.runAs })
          .then(user => {
            this.user.value = user
          }).catch(() => {
            return {}
          })
      }
    },

    updateRunAs (user) {
      if (user && user.userID) {
        this.user.value = user
        this.workflow.runAs = user.userID
      } else {
        this.user.value = null
        this.workflow.runAs = '0'
      }
      this.$root.$emit('change-detected')
    },

    getOptionKey ({ userID }) {
      return userID
    },

    getOptionLabel ({ userID, email, name, username }) {
      return name || username || email || `<@${userID}>`
    },
  },
}
</script>
