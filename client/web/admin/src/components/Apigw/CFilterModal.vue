<template>
  <div>
    <b-modal
      :visible="visible"
      size="lg"
      :title="(filter || {}).label"
      :ok-title="$t('filters.modal.ok')"
      body-class="p-0"
      cancel-variant="link"
      @ok="onSave"
      @hidden="onHidden"
    >
      <div
        v-if="filter"
        class="card-body"
      >
        <c-filter-params
          :filter="filter"
          @update="onUpdate"
        />
        <b-form-checkbox
          v-model="filter.enabled"
          @change="onUpdate"
        >
          {{ $t('filters.enabled') }}
        </b-form-checkbox>
      </div>
    </b-modal>
  </div>
</template>

<script>
import CFilterParams from 'corteza-webapp-admin/src/components/Apigw/CFilterParams'

export default {
  components: {
    CFilterParams,
  },

  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },

    visible: {
      type: Boolean,
      required: false,
      default: false,
    },
  },

  data () {
    return {
      updated: false,

      filteredFields: [],
    }
  },

  methods: {
    onSave () {
      this.$emit('submit', { ...this.filter, updated: this.updated })
      this.updated = false
    },

    onHidden () {
      this.$emit('reset')
    },

    onUpdate () {
      this.updated = true
    },
  },
}
</script>
