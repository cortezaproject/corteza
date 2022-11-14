<template>
  <b-row
    no-gutters
    class="mb-lg-1"
  >
    <b-col
      cols="12"
      lg="6"
      order="2"
      order-lg="1"
    >
      <access
        :access="access"
        :current="current"
        :enabled="enabled"
        @update="onUpdate"
      />

      <b-button
        v-show="isChanged"
        data-test-id="button-reset-back"
        variant="link"
        class="pl-0"
        @click="onReset"
      >
        {{ $t('ui.reset-back', {current}) }}
      </b-button>
    </b-col>
    <b-col
      cols="12"
      lg="6"
      class="mb-lg-4 text-wrap"
      order="1"
      order-lg="2"
    >
      <b v-html="title || `${operation} on ${resource}`" />
      <div>{{ description || '&nbsp;' }}</div>
    </b-col>
  </b-row>
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
