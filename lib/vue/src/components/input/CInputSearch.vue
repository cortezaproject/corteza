<template>
  <div
    style="min-width: 150px;"
    :class="{ 'submittable': isSubmittable }"
    class="c-input-search d-flex position-relative"
  >
    <b-input
      ref="searchInput"
      data-test-id="input-search"
      :type="inputType"
      name="search"
      :value="localValue"
      :debounce="debounce"
      :disabled="disabled"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :size="size"
      class="text-truncate"
      @input="onInput"
      @keyup.enter="submitQuery"
    />

    <b-button
      v-if="clearable && localValue && !disabled"
      variant="link"
      class="close-button d-inline-flex align-items-center rounded-0 p-3"
      @click="onClear"
    >
      <font-awesome-icon
        :icon="['fas', 'times']"
        class="text-primary"
      />
    </b-button>

    <b-button
      v-if="showSubmittable"
      :variant="isSubmittable ? 'outline-light' : 'link'"
      :disabled="disabled"
      :class="{ 'border-0 cursor-default': !isSubmittable }"
      class="search-button d-inline-flex align-items-center rounded-0 border-light"
      @[isSubmittable]="submitQuery"
    >
      <font-awesome-icon
        :icon="['fas', 'search']"
        class="align-middle text-primary"
      />
    </b-button>
  </div>
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
      default: 'md',
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
      default: false,
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

  data () {
    return {
      localValue: this.value,
    }
  },

  computed: {
    inputType () {
      return this.clearable ? 'search' : 'text'
    },

    showSubmittable () {
      return !this.localValue || this.showSubmittableAndClearable
    },

    isSubmittable () {
      return this.submittable && !this.disabled ? 'click' : null
    },

    showSubmittableAndClearable () {
      return this.clearable && this.submittable
    },
  },

  watch: {
    value (value) {
      this.localValue = value
    },
  },

  methods: {
    onInput (value) {
      this.localValue = value

      if (!this.submittable) {
        this.$emit('input', value)
      }
    },

    submitQuery () {
      if (this.submittable) {
        this.$emit('search', this.$refs.searchInput.localValue)
      }
    },

    onClear () {
      this.localValue = ''
      if (!this.submittable) {
        this.$emit('input', '')
      }
      this.$nextTick(() => {
        this.$refs.searchInput.focus()
      })
    },
  },
}
</script>

<style lang="scss" scoped>
input:focus::placeholder {
  color: transparent;
}

.c-input-search {
  .search-button {
    position: absolute;
    right: 2px;
    top: 2px;
    bottom: 2px;
    z-index: 4;
    border-left-width: 2px;
  }

  .close-button {
    position: absolute;
    right: 1px;
    top: 1px;
    bottom: 1px;
    z-index: 5;
    border: none;
    background: none;

    &:hover {
      text-decoration: none;
    }
  }

  &.submittable .close-button {
    right: 48px;
  }

  .form-control {
    padding-right: 40px;
  }

  &.submittable .form-control {
    padding-right: 85px;
  }

  ::-webkit-search-cancel-button {
    -webkit-appearance: none;
    display: none;
  }
}

.cursor-default {
  cursor: default !important;
}
</style>
