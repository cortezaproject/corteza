<template>
  <div id="link-popover-container">
    <b-button
      id="link-popover"
      variant="link"
      class="text-dark font-weight-bold text-decoration-none"
      @click="showPopover"
    >
      <span :class="activeClasses(format.attrs)">
        <font-awesome-icon icon="link" />
      </span>
    </b-button>

    <b-popover
      v-if="currentValue"
      ref="popover"
      :show.sync="visible"
      triggers="focus"
      target="link-popover"
      placement="auto"
      container="link-popover-container"
    >
      <b-input-group style="min-width: 250px;">
        <b-form-input
          v-model="attrs.href"
          type="url"
          autofocus
          :placeholder="labels.urlPlaceholder"
          @keydown.enter.prevent.stop="link"
          @keydown.esc.prevent.stop="close"
        />
        <b-input-group-append>
          <b-button
            variant="outline-success"
            @click="link"
          >
            {{ labels.ok }}
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </b-popover>
  </div>
</template>

<script>
import base from './base.vue'

/**
 * Component is used to display link formatters. It provides an interface to
 * input the URL that should be applied.
 */
export default {
  name: 'TMarkLink',
  extends: base,

  props: {
    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      visible: false,
      attrs: { href: null, target: '_self' },
    }
  },

  computed: {
    /**
     * Does a simple check if entered URL is valid.
     * @todo Improve this
     * @returns {Boolean}
     */
    urlValid () {
      if (!this.attrs.href) {
        return false
      }
      return !!this.attrs.href
    },
  },

  methods: {
    /**
     * Helper to show the popup & determine if a link already exists
     */
    showPopover () {
      if (this.currentValue) {
        this.visible = true
        this.attrs = { ...this.getMarkAttrs(this.format.type) }
      }
    },

    /**
     * Helper to submit the given link
     */
    link () {
      this.onClick(this.format.type, this.attrs)
      this.close()
    },

    /**
     * Helper to close the popup & reset the state
     */
    close () {
      this.attrs.href = null
      this.visible = false
    },
  },
}

</script>
