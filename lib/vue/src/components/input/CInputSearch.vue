<template>
  <b-input-group
    :size="size"
  >
    <b-input
      ref="searchInput"
      data-test-id="input-search"
      :type="inputType"
      name="search"
      :value="value"
      :debounce="debounce"
      :disabled="disabled"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      class="h-100 pr-0 border-light border-right-0 text-truncate bg-white"
      @update="search"
      @keyup.enter="submitQuery"
    />
    <b-input-group-append
      :class="{ 'border-left': showSubmittable }"
      class="bg-white border-light rounded-right append-group-border"
    >
      <b-button
        v-if="showSubmittable"
        variant="link"
        :disabled="disabled"
        class="py-0"
        :class="{
          'search-icon-border': showSubmittableAndClearable,
          'cursor-default': !isSubmittable
        }"
        @[isSubmittable]="submitQuery"
      >
        <font-awesome-icon
          :icon="['fas', 'search']"
          class="align-middle"
        />
      </b-button>
    </b-input-group-append>
  </b-input-group>
</template>

<script>
export default {
  name: 'CInputSearch',

  props: {
    value: {
      type: String,
      default: '',
    },

    placeholder: {
      type: String,
      default: '',
    },

    size: {
      type: String,
    },

    disabled: {
      type: Boolean,
    },

    clearable: {
      type: Boolean,
      default: true,
    },

    submittable: {
      type: Boolean,
    },

    autocomplete: {
      type: String,
      default: 'on',
    },

    debounce: {
      type: Number,
      default: 0,
    },
  },

  computed: {
    inputType () {
      return this.clearable ? 'search' : 'text'
    },

    showSubmittable () {
      return !this.value || (this.value && this.showSubmittableAndClearable)
    },

    isSubmittable () {
      return this.submittable && !this.disabled ? 'click' : null
    },

    showSubmittableAndClearable () {
      return this.clearable && this.submittable
    },
  },

  methods: {
    search (e) {
      if (!this.submittable) {
        this.$emit('input', e)
      }
    },

    submitQuery () {
      if (this.submittable) {
        this.$emit('search', this.$refs.searchInput.localValue)
      }
    },

    clearQuery () {
      this.$refs.searchInput.focus()
      this.$emit('input', '')
    },
  },
}
</script>
<style lang="scss" scoped>
$border-color: 2px solid #F3F5F7;

input:focus::placeholder {
  color: transparent;
}

input[type="search"]::-webkit-search-cancel-button {
  height: 13px;
  width: 13px;
  padding-left: 12px;
  cursor: pointer;
  background: url("data:image/svg+xml;charset=UTF-8,%3csvg viewPort='0 0 12 12' version='1.1' xmlns='http://www.w3.org/2000/svg'%3e%3cline x1='1' y1='11' x2='11' y2='1' stroke='black' stroke-width='2'/%3e%3cline x1='1' y1='1' x2='11' y2='11' stroke='black' stroke-width='2'/%3e%3c/svg%3e");
}

.append-group-border {
  border: $border-color;
}

.search-icon-border {
  border-left: $border-color;
}

.cursor-default {
  cursor: default !important;
}
</style>
