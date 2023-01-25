<script lang="js">
export default {
  props: {
    labels: {
      type: Object,
      default: () => ({}),
    },

    src: {
      required: true,
      type: [ String, Object ],
      default: undefined,
    },

    mime: {
      type: String,
      default: undefined,
    },

    meta: {
      type: Object,
      default: () => ({}),
    },

    name: {
      type: String,
      default: null,
    },

    previewStyle: {
      type: Object,
      default: () => ({}),
    },

    previewClass: {
      type: Array,
      default: () => [],
    },

    inline: {
      type: Boolean,
    },
  },
}

</script>
