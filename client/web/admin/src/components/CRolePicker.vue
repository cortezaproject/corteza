<template>
  <div
    data-test-id="role-picker"
  >
    <b-input-group>
      <b-input-group-append is-text>
        <font-awesome-icon :icon="['fas', 'search']" />
      </b-input-group-append>
      <b-form-input
        v-model.trim="filter"
        data-test-id="input-role"
      />
      <b-input-group-append
        v-if="filter"
        is-text
      >
        <b-button
          data-test-id="button-clear-role"
          variant="link"
          size="sm"
          class="p-0 m-0"
          @click="filter = ''"
        >
          <font-awesome-icon :icon="['fas', 'times']" />
        </b-button>
      </b-input-group-append>
    </b-input-group>

    <b-container
      v-if="filter && filtered.length > 0"
      class="ml-5 my-2 position-absolute bg-white border results shadow w-50"
    >
      <b-row
        v-for="r in filtered"
        :key="r.roleID"
        data-test-id="filtered-row-list"
        class="filtered-role"
        @click="addRole(r)"
      >
        <b-col class="pt-1">
          {{ r | label }}
          <b-button
            data-test-id="button-add-role"
            variant="link"
            class="float-right"
            @click="addRole(r)"
          >
            <font-awesome-icon :icon="['fas', 'plus']" />
          </b-button>
        </b-col>
      </b-row>
    </b-container>

    <b-form-text
      v-if="$slots['description']"
    >
      <slot name="description" />
    </b-form-text>

    <b-container
      v-if="selected"
      class="p-1"
    >
      <b-row
        v-for="r in selected"
        :key="r.userID"
        data-test-id="selected-row-list"
      >
        <b-col>{{ r | label }}</b-col>
        <b-col class="text-right">
          <b-button
            data-test-id="button-remove-role"
            variant="link"
            @click="removeRole(r)"
          >
            <font-awesome-icon :icon="['far', 'trash-alt']" />
          </b-button>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>

function roleSorter (a, b) {
  return `${a.name} ${a.handle} ${a.roleID}`.localeCompare(`${b.name} ${b.handle} ${b.roleID}`)
}

export default {
  filters: {
    label (r) {
      return r.name || r.handle || r.roleID
    },
  },

  props: {
    label: {
      type: String,
      default: 'count',
    },

    // list of role IDs
    value: {
      type: Array,
      default: () => ([]),
    },
  },

  data () {
    return {
      roles: [],
      filter: '',
    }
  },

  computed: {
    selected () {
      return this.roles
        .filter(({ roleID }) => this.value.includes(roleID))
        .sort(roleSorter)
    },

    filtered () {
      const match = ({ name = '', handle = '', roleID = '' }) => {
        return `${name} ${handle} ${roleID}`.toLocaleLowerCase().indexOf(this.filter.toLocaleLowerCase()) > -1
      }

      const fits = ({ isClosed, meta = {} }) => {
        return !(isClosed || (meta.context && meta.context.resourceTypes))
      }

      return this.roles.filter(r => !this.value.includes(r.roleID) && fits(r) && match(r))
    },
  },

  watch: {
    currentRoles: {
      immediate: true,
      handler () {
        this.filter = ''
      },
    },
  },

  mounted () {
    this.preload()
  },

  methods: {
    addRole (r) {
      if (!this.value.includes(r.roleID)) {
        this.value.push(r.roleID)
      }
    },

    removeRole (r) {
      this.value.splice(this.value.indexOf(r.roleID), 1)
      this.filter = ''
    },

    preload () {
      return this.$SystemAPI.roleList()
        .then(({ set }) => { this.roles = set || [] })
        .catch(this.toastErrorHandler({}))
    },
  },
}
</script>
<style lang="scss">
.results {
  z-index: 100;
  .filtered-role {
    cursor: pointer;
    &:hover {
      background-color: $light;
    }
  }
}

</style>
