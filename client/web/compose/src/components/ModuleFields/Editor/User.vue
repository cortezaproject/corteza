<template>
  <b-form-group
    label-class="text-primary"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <div
        class="d-flex align-items-top"
      >
        <label
          class="mb-0"
        >
          {{ label }}
        </label>

        <hint
          :id="field.fieldID"
          :text="hint"
        />
      </div>
      <small
        class="text-muted"
      >
        {{ description }}
      </small>
    </template>

    <multi
      v-if="field.isMulti"
      :value.sync="value"
      :errors="errors"
      :single-input="field.options.selectType !== 'each'"
      :removable="field.options.selectType !== 'multiple'"
    >
      <template v-slot:single>
        <vue-select
          v-if="field.options.selectType === 'default'"
          ref="singleSelect"
          :placeholder="$t('kind.user.suggestionPlaceholder')"
          :options="options"
          :get-option-label="getOptionLabel"
          :get-option-key="getOptionKey"
          :append-to-body="appendToBody"
          :calculate-position="calculatePosition"
          :clearable="false"
          :filterable="false"
          :selectable="option => option.selectable"
          :loading="processing"
          class="bg-white w-100"
          @search="search"
          @input="updateValue($event)"
        >
          <pagination
            v-if="showPagination"
            slot="list-footer"
            :has-prev-page="hasPrevPage"
            :has-next-page="hasNextPage"
            @prev="goToPage(false)"
            @next="goToPage(true)"
          />
        </vue-select>
        <vue-select
          v-else-if="field.options.selectType === 'multiple'"
          v-model="multipleSelected"
          :placeholder="$t('kind.user.suggestionPlaceholder')"
          :options="options"
          :get-option-label="getOptionLabel"
          :get-option-key="getOptionKey"
          :append-to-body="appendToBody"
          :calculate-position="calculatePosition"
          :filterable="false"
          :selectable="option => option.selectable"
          :loading="processing"
          multiple
          class="bg-white w-100"
          @search="search"
        >
          <pagination
            v-if="showPagination"
            slot="list-footer"
            :has-prev-page="hasPrevPage"
            :has-next-page="hasNextPage"
            @prev="goToPage(false)"
            @next="goToPage(true)"
          />
        </vue-select>
      </template>
      <template v-slot:default="ctx">
        <vue-select
          v-if="field.options.selectType === 'each'"
          :placeholder="$t('kind.user.suggestionPlaceholder')"
          :options="options"
          :get-option-label="getOptionLabel"
          :get-option-key="getOptionKey"
          :value="getUserByIndex(ctx.index)"
          :append-to-body="appendToBody"
          :calculate-position="calculatePosition"
          :clearable="false"
          :filterable="false"
          :selectable="option => option.selectable"
          :loading="processing"
          class="bg-white w-100"
          @search="search"
          @input="updateValue($event, ctx.index)"
        >
          <pagination
            v-if="showPagination"
            slot="list-footer"
            :has-prev-page="hasPrevPage"
            :has-next-page="hasNextPage"
            @prev="goToPage(false)"
            @next="goToPage(true)"
          />
        </vue-select>
        <span v-else>{{ getOptionLabel(getUserByIndex(ctx.index)) }}</span>
      </template>
    </multi>
    <template
      v-else
    >
      <vue-select
        :placeholder="$t('kind.user.suggestionPlaceholder')"
        :options="options"
        :get-option-label="getOptionLabel"
        :get-option-key="getOptionKey"
        :value="getUserByIndex()"
        :append-to-body="appendToBody"
        :calculate-position="calculatePosition"
        :filterable="false"
        :selectable="option => option.selectable"
        :loading="processing"
        class="bg-white w-100"
        @input="updateValue($event)"
        @search="search"
      >
        <pagination
          v-if="showPagination"
          slot="list-footer"
          :has-prev-page="hasPrevPage"
          :has-next-page="hasNextPage"
          @prev="goToPage(false)"
          @next="goToPage(true)"
        />
      </vue-select>
      <errors :errors="errors" />
    </template>
  </b-form-group>
</template>
<script>
import { debounce } from 'lodash'
import base from './base'
import { VueSelect } from 'vue-select'
import { mapActions, mapGetters } from 'vuex'
import calculatePosition from 'corteza-webapp-compose/src/mixins/vue-select-position'
import Pagination from '../Common/Pagination.vue'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    VueSelect,
    Pagination,
  },

  extends: base,

  mixins: [
    calculatePosition,
  ],

  data () {
    return {
      processing: false,

      users: [],

      filter: {
        query: null,
        limit: 10,
        pageCursor: '',
        prevPage: '',
        nextPage: '',
        roles: [],
      },
    }
  },

  computed: {
    ...mapGetters({
      resolved: 'user/set',
      findByID: 'user/findByID',
    }),

    options () {
      return this.users.map(u => {
        return { ...u, selectable: this.field.isMulti ? !(this.value || []).includes(u.userID) : this.value !== u.userID }
      })
    },

    // This is used in the case of using the multiple select option
    multipleSelected: {
      get () {
        const map = userID => {
          return this.findByID(userID) || { userID }
        }

        return this.field.isMulti ? this.value.map(map) : map(this.value)
      },

      set (users) {
        if (users && Array.isArray(users)) {
          // When adding/removing items from vue-selects[multiple],
          // we get array of options back

          this.addUserToResolved(users)
          this.value = users.map(({ userID }) => userID)
        }
      },
    },

    showPagination () {
      return this.hasPrevPage || this.hasNextPage
    },

    hasPrevPage () {
      return !!this.filter.prevPage
    },

    hasNextPage () {
      return !!this.filter.nextPage
    },
  },

  watch: {
    value: {
      async handler (value) {
        value = this.field.isMulti ? [...value] : [value]
        if (value) {
          await this.resolveUsers(value)
        }
      },
    },

    'filter.pageCursor': {
      handler (pageCursor) {
        if (pageCursor) {
          this.fetchUsers()
        }
      },
    },
  },

  created () {
    // Prefill value with current user
    if ((!this.value || this.value.length === 0) && this.field.options.presetWithAuthenticated) {
      this.updateValue(this.$auth.user)
    }

    this.fetchUsers()
  },

  methods: {
    ...mapActions({
      resolveUsers: 'user/fetchUsers',
      addUserToResolved: 'user/push',
    }),

    getOptionKey ({ userID }) {
      return userID
    },

    getOptionLabel ({ userID, email, name, username }) {
      return name || username || email || `<@${userID}>`
    },

    /**
     * Updates record value with user
     *
     * Handles single & multi value fields
     */
    updateValue (user, index = -1) {
      // reset singleSelect value for better value presentation
      if (this.$refs.singleSelect) {
        this.$refs.singleSelect._data._value = undefined
      }

      if (user) {
        // update list of resolved users for every item we add
        this.addUserToResolved({ ...user })

        // update valie on record
        const { userID } = user
        if (this.field.isMulti) {
          if (index >= 0) {
            this.value[index] = userID
          } else {
            // <0, assume we're appending
            this.value.push(userID)
          }
        } else {
          this.value = userID
        }
      } else {
        if (index >= 0) {
          this.value.splice(index, 1)
        } else {
          this.value = undefined
        }
      }
    },

    /**
     * Retrives user (via value) from record field
     * Handles single & multi value fields
     */
    getUserByIndex (index = 0) {
      const userID = this.field.isMulti ? this.value[index] : this.value
      if (userID) {
        return this.findByID(userID) || {}
      }
    },

    search: debounce(function (query = '') {
      if (query !== this.filter.query) {
        this.filter.query = query
        this.filter.pageCursor = undefined
      }

      this.fetchUsers()
    }, 300),

    fetchUsers () {
      this.processing = true

      const roleID = this.field.options.roles || []

      this.$SystemAPI.userList({ ...this.filter, roleID })
        .then(({ filter, set }) => {
          this.filter = { ...this.filter, ...filter }
          this.filter.nextPage = filter.nextPage
          this.filter.prevPage = filter.prevPage
          this.users = set.map(m => Object.freeze(m))
          return { filter, set }
        })
        .finally(() => {
          this.processing = false
        })
    },

    goToPage (next = true) {
      this.filter.pageCursor = next ? this.filter.nextPage : this.filter.prevPage
    },
  },
}
</script>
