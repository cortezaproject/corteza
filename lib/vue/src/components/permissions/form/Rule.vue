<template>
  <div
    :data-test-id="title || `${operation} on ${resource}`"
  >
    <p
      :title="title || `${operation} on ${resource}`"
      class="mb-1"
    >
      {{ title || `${operation} on ${resource}` }}
    </p>

    <access
      :access="access"
      :current="current"
      :enabled="enabled"
      class="w-100"
      @update="onUpdate"
    />
  </div>
</template>
<script lang="js">
import Access from './Access.vue'

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    Access,
  },

  props: {
    resource: {
      type: String,
      required: true,
    },

    operation: {
      type: String,
      required: true,
    },

    title: {
      type: String,
      required: false,
      default: undefined,
    },

    description: {
      type: String,
      required: false,
      default: undefined,
    },

    enabled: {
      type: Boolean,
      default: true,
    },

    access: {
      type: String,
      required: false,
      default: 'inherit',
    },

    current: {
      type: String,
      required: false,
      default: 'inherit',
    },
  },

  computed: {
    isChanged () {
      return this.access !== this.current
    },
  },

  methods: {
    onUpdate (access) {
      this.emit(access)
    },

    onReset () {
      this.emit(this.current)
    },

    emit (access) {
      this.$emit('update', {
        resource: this.resource,
        operation: this.operation,
        access,
      })
    },
  },
}
</script>
