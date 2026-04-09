<template>
  <span>
    <b-button
      :id="popoverId"
      variant="link"
      class="text-dark font-weight-bold text-decoration-none"
    >
      <font-awesome-icon :icon="['far', 'face-smile']" />
    </b-button>

    <b-popover
      :target="popoverId"
      triggers="click blur"
      placement="bottom"
      container="body"
      custom-class="emoji-picker-popover border-light"
      @shown="onShown"
    >
      <c-emoji-picker
        ref="picker"
        :emojis="allEmojis"
        :labels="labels.emojiPicker || {}"
        :show-quick-reactions="false"
        @select="onSelect"
      />
    </b-popover>
  </span>
</template>

<script>
// Use ~corteza-vue/ alias to avoid circular self-reference to dist
import CEmojiPicker from '~corteza-vue/components/CEmojiPicker.vue'

let popoverCounter = 0

export default {
  name: 'TMarkEmoji',

  components: {
    CEmojiPicker,
  },

  props: {
    editor: {
      type: Object,
      required: true,
    },

    format: {
      type: Object,
      default: () => ({}),
    },

    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      popoverId: `emoji-picker-${++popoverCounter}`,
    }
  },

  computed: {
    allEmojis () {
      return this.editor.storage?.emoji?.emojis || []
    },
  },

  methods: {
    onShown () {
      this.$nextTick(() => {
        if (this.$refs.picker) {
          this.$refs.picker.reset()
        }
      })
    },

    onSelect (emoji) {
      if (emoji && emoji.name) {
        this.editor.chain().focus().insertContent({
          type: 'emoji',
          attrs: { name: emoji.name },
        }).run()
      }
      this.$root.$emit('bv::hide::popover', this.popoverId)
    },
  },
}
</script>

<style lang="scss">
.emoji-picker-popover {
  background: var(--white, #fff);
  max-width: 300px;

  .popover-body {
    padding: 0;
  }

  .arrow::after {
    border-bottom-color: var(--white, #fff);
  }
}
</style>
