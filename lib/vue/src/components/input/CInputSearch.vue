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
      :value="value"
      :debounce="debounce"
      :disabled="disabled"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :size="size"
      class="pr-0 text-truncate"
      @input="onInput"
      @update="search"
      @keyup.enter="submitQuery"
    />

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

  watch: {
    value (value) {
      this.localValue = value
    },
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

  methods: {
    onInput (value) {
      this.localValue = value
    },

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

  ::-webkit-search-cancel-button {
    -webkit-appearance: none;
    height: 1em;
    width: 1em;
    background: var(--primary);
    -webkit-mask-image: url("data:image/svg+xml;charset=UTF-8,%3csvg viewPort='0 0 12 12' version='1.1' xmlns='http://www.w3.org/2000/svg'%3e%3cline x1='1' y1='11' x2='11' y2='1' stroke='black' stroke-width='2'/%3e%3cline x1='1' y1='1' x2='11' y2='11' stroke='black' stroke-width='2'/%3e%3c/svg%3e");
    mask-image: url("data:image/svg+xml;charset=UTF-8,%3csvg viewPort='0 0 12 12' version='1.1' xmlns='http://www.w3.org/2000/svg'%3e%3cline x1='1' y1='11' x2='11' y2='1' stroke='black' stroke-width='2'/%3e%3cline x1='1' y1='1' x2='11' y2='11' stroke='black' stroke-width='2'/%3e%3c/svg%3e");
    cursor: pointer;
    margin-right: 13px;
    margin-left: 5px;
  }

  &.submittable {
    ::-webkit-search-cancel-button {
      margin-right: 56px;
    }
  }
}

.cursor-default {
  cursor: default !important;
}
</style>
