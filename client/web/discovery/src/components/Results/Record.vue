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
            {{ hit.value.module.name || hit.value.module.handle }}
          </span>
        </h5>

        <span class="text-nowrap">
          <b-badge
            v-if="Object.keys(hit.value.labels || { }).includes('federation')"
            variant="light"
            class="mr-1 h5 p-2 mb-0"
          >
            Federated
          </b-badge>
          <b-avatar
            v-b-tooltip.hover
            title="Record"
            size="sm"
            icon="file-earmark-text"
            class="align-center bg-light text-dark"
            style="z-index: 1;"
          />
        </span>
      </div>

      <div class="d-flex justify-content-between small">
        <slot name="header" />
      </div>
    </b-card-header>

    <b-card-body class="pb-0">
      <div
        v-if="limitData().length"
      >
        <div
          v-for="(item, i) in limitData()"
          :key="i"
          class="d-flex flex-column mb-3"
        >
          <label
            class="text-capitalize text-primary mb-0"
          >
            {{ item.label || item.name }}
          </label>
          <p class="multiline mt-1 mb-0">
            <text-highlight
              :queries="query"
              highlight-style="padding: 0 0.05rem;"
            >
              {{ item.value }}
            </text-highlight>
          </p>
        </div>
      </div>

      <p
        v-else
      >
        No values
      </p>
    </b-card-body>
  </b-overlay>
</template>

<script>
import base from './base'

export default {
  extends: base,

  methods: {
    limitData () {
      const { values = [] } = this.hit.value

      return (values || []).map(({ name, label, value = [] }) => {
        if (value) {
          value = value.map(v => {
            return v.toString().includes('{"coordinates":[') ? ((JSON.parse(v || '{}') || {}).coordinates || []).join(', ') : v
          }).join('\n')
        }

        return { name, label, value }
      })
    },
  },
}
</script>
