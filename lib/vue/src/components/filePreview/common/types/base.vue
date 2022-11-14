<script lang="js">
export default {
  props: {
    labels: {
      required: false,
      type: Object,
      default: () => ({}),
    },

    src: {
      required: true,
      type: [ String, Object ],
      default: undefined,
    },

    mime: {
      required: false,
      type: String,
      default: undefined,
    },

    meta: {
      type: Object,
      default: () => ({}),
    },

    name: {
      required: false,
      type: String,
      default: null,
    },

    previewStyle: {
      required: false,
      type: Object,
      default: () => ({}),
    },

    previewClass: {
      required: false,
      type: Array,
      default: () => [],
    },

    inline: {
      type: Boolean,
      default: false,
    },
  },
}

</script>
