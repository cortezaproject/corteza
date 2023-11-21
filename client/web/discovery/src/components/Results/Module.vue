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
          <span
            v-if="hit.value.namespace.name || hit.value.namespace.handle"
          >
            {{ hit.value.namespace.name || hit.value.namespace.handle }}
          </span>
          <span
            v-if="hit.value.namespace.name || hit.value.namespace.handle"
            class="mx-2"
          >
            <b-icon
              icon="arrow-right"
              font-scale="1"
            />
          </span>
          <span>
            {{ hit.value.name || hit.value.handle }}
          </span>
        </h5>

        <span class="text-nowrap">
          <b-badge
            v-if="Object.keys(hit.value.labels || { }).includes('federation')"
            variant="light"
            class="mr-1 h5 p-2 mb-0"
          >
            {{ $t('general:federated') }}
          </b-badge>
          <b-avatar
            v-b-tooltip.hover="{ title: $t('types.module'), container: '#body' }"
            size="sm"
            icon="list-ul"
            class="align-center bg-light text-dark"
          />
        </span>
      </div>

      <div class="d-flex justify-content-between small">
        <slot name="header" />
      </div>
    </b-card-header>

    <b-card-body class="pb-0">
      <div
        v-for="(value, name, i) in limitData()"
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
            {{ value }}
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
        'fields',
        'namespace',
        'labels',
        'module',
      ]
    },
  },
}
</script>
