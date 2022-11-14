<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group :label="$t('name.label')">
      <b-form-input
        v-model="name"
        data-test-id="input-name"
        class="mt-1"
        :placeholder="$t('name.placeholder')"
      />
    </b-form-group>

    <b-form-group :label="$t('slug.label')">
      <b-form-input
        v-model="slug"
        data-test-id="input-handle"
        class="mt-1"
        :state="slugState"
        :placeholder="$t('slug.placeholder')"
      />
      <b-form-invalid-feedback :state="slugState">
        {{ $t('slug.invalid-handle-characters') }}
      </b-form-invalid-feedback>
    </b-form-group>

    <div slot="footer">
      <b-button
        data-test-id="button-back"
        variant="outline-dark"
        class="float-left"
        @click="$emit('back')"
      >
        {{ $t('import.back') }}
      </b-button>

      <b-button
        data-test-id="button-import"
        variant="dark"
        :disabled="submitDisabled"
        class="float-right"
        @click="nextStep"
      >
        {{ $t('import.import') }}
      </b-button>
    </div>
  </b-card>
</template>

<script>
import { handleState } from 'corteza-webapp-compose/src/lib/handle'

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  props: {
    session: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

  data () {
    return {
      name: '',
      slug: '',
    }
  },

  computed: {
    submitDisabled () {
      return [this.nameState, this.slugState].includes(false)
    },

    nameState () {
      return this.name.length > 0 ? null : false
    },

    slugState () {
      return this.slug.length > 0 ? handleState(this.slug) : false
    },
  },

  methods: {
    nextStep () {
      // convert to api's structure
      const rtr = {
        name: this.name,
        slug: this.slug,
      }

      this.$emit('configured', rtr)
    },
  },
}
</script>
