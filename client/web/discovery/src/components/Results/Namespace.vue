<template>
  <b-overlay>
    <b-card-header
      header-bg-variant="white"
      class="border-bottom"
    >
      <div class="d-flex align-items-center mb-3 justify-content-between">
        <h5
          class="text-primary text-capitalize text-truncate mr-2 mb-0"
        >
          {{ hit.value.name || hit.value.slug }}
        </h5>
        <b-avatar
          v-b-tooltip.hover
          :title="$t('types.namespace')"
          size="sm"
          icon="code-square"
          class="align-center bg-light text-dark"
        />
      </div>
      <div class="d-flex justify-content-between small">
        <slot name="header" />
      </div>
    </b-card-header>

    <b-card-body class="pb-0">
      <div
        v-for="(item, name, i) in limitData()"
        :key="i"
        class="d-flex flex-column mb-3"
      >
        <label
          class="text-capitalize text-primary mb-0"
        >
          {{ name }}
        </label>
        <div class="mt-1">
          <text-highlight
            :queries="query"
            highlight-style="padding: 0 0.05rem;"
          >
            {{ item }}
          </text-highlight>
        </div>
      </div>
    </b-card-body>
  </b-overlay>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'filters',
  },

  extends: base,

  computed: {
    blacklistedFields () {
      return [
        ...this.defaultBlacklistedFields,
        'meta',
        'namespace',
      ]
    },
  },
}
</script>
