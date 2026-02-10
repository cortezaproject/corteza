<template>
  <div
    class="hit-container p-2"
    @click="$emit('click')"
  >
    <b-list-group-item
      :id="`search-hit-${hit.value.recordID}`"
      class="search-card border rounded bg-white p-2 position-relative"
    >
      <b-button
        class="open-in-new-tab border-0 text-secondary px-2"
        :title="labels.openInNewTab || 'Open in new tab'"
        size="sm"
        variant="outline-extra-light"
        @click.stop="$emit('open-new-tab')"
      >
        <font-awesome-icon
          :icon="['fas', 'external-link-alt']"
        />
      </b-button>

      <label class="text-primary font-weight-bold mb-0">
        {{ hit.highlight.label }}
      </label>
      <p class="text-dark mb-0 text-truncate">
        {{ hit.highlight.value }}
      </p>
    </b-list-group-item>
  </div>
</template>

<script>
import { library } from '@fortawesome/fontawesome-svg-core'
import { faExternalLinkAlt } from '@fortawesome/free-solid-svg-icons'

library.add(faExternalLinkAlt)

export default {
  name: 'RecordItem',

  props: {
    hit: {
      type: Object,
      required: true,
    },
    labels: {
      type: Object,
      default: () => ({}),
    },
  },
}
</script>

<style lang="scss" scoped>
.hit-container {
  transition: background-color 0.2s ease;
  cursor: pointer;

  .search-card {
    transition: all 0.2s ease;
  }

  &:hover {
    background-color: var(--light) !important;

    .open-in-new-tab {
      opacity: 1 !important;
      pointer-events: auto;
    }
  }
}

.open-in-new-tab {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  z-index: 2;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s ease;
}
</style>
