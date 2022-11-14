<template>
  <b-col
    class="my-2"
    sm="6"
    md="4"
    lg="3"
  >
    <b-card
      no-body
      class="h-100 shadow-sm pt-3 mx-2"
      footer-bg-variant="white"
      footer-class="text-center pt-0"
      :class="{ 'shadow': hovered && isEnabled, 'namespace-item' : isEnabled, 'disabled' : !isEnabled}"
      @mouseover="hovered = true"
      @mouseleave="hovered = false"
    >
      <div
        class="circled-avatar d-flex align-items-center justify-content-center m-auto"
        :class="[namespace.meta.logoEnabled ? 'p-2' : 'bg-light p-3']"
      >
        <b-img
          v-if="namespace.meta.logoEnabled"
          :src="logo"
          :alt="namespace.name"
          class="mw-100 mh-100"
        />
        <h1
          v-else
          class="ns-initial m-auto text-uppercase text-secondary"
        >
          {{ namespace.initials }}
        </h1>
      </div>
      <b-card-body
        class="mw-100 text-center pb-3"
      >
        <div
          class="d-flex align-items-baseline"
          :class="{ 'h-100': !namespace.meta.description }"
        >
          <div
            class="d-flex flex-column justify-content-center w-100"
          >
            <h5
              :data-test-id="namespace.name"
            >
              {{ namespace.name }}
            </h5>
            <p
              v-if="namespace.meta.subtitle"
              class="d-inline-block my-1 text-secondary"
            >
              {{ namespace.meta.subtitle }}
            </p>
          </div>
        </div>

        <p
          v-if="namespace.meta.description"
          class="overflow-auto"
        >
          <small>{{ namespace.meta.description }}</small>
        </p>
        <b-link
          v-if="isEnabled"
          :to="{ name: 'pages', params: { slug: (namespace.slug || namespace.namespaceID) } }"
          data-test-id="link-visit-namespace"
          :aria-label="$t('visit') + ' ' + namespace.name"
          class="stretched-link"
        />
      </b-card-body>
    </b-card>
  </b-col>
</template>
<script>

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  props: {
    namespace: {
      type: Object,
      required: true,
    },
  },

  data () {
    return {
      processing: false,
      hovered: undefined,
      logoAttachment: undefined,
    }
  },

  computed: {
    isEnabled () {
      return !!this.namespace.enabled
    },

    canEdit () {
      return !!this.namespace.canUpdateNamespace
    },

    showFooter () {
      return this.isEnabled || this.canEdit
    },

    logo () {
      return this.namespace.meta.logo || this.$Settings.attachment('ui.mainLogo')
    },
  },
}
</script>

<style lang="scss" scoped>
$avatar-size: 110px;
$disabled-opacity: 0.6;

.namespace-item {
  &:hover {
    transition: all 0.2s ease;
    top: -1px;
  }
}

.circled-avatar {
  width: $avatar-size;
  height: $avatar-size;
  border-radius: 50%;
}

.disabled {
  opacity: $disabled-opacity;
}
</style>
