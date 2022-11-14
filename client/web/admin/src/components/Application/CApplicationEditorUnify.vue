<template>
  <b-card
    data-test-id="card-application-selector"
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit', unify)"
    >
      <b-form-group
        :label="$t('name.label')"
        :description="$t('name.description')"
        label-cols="2"
      >
        <b-form-input
          v-model="unify.name"
          data-test-id="input-name"
        />
      </b-form-group>

      <b-form-group
        :label="$t('logo.label')"
        :description="$t('logo.description')"
        label-cols="2"
      >
        <template #label>
          <div
            class="d-flex align-items-center"
          >
            {{ $t('logo.label') }}
            <b-button
              v-if="showLogoPreview"
              v-b-modal.logo
              data-test-id="button-logo-show"
              variant="link"
              class="d-flex align-items-center border-0 p-0 ml-2"
            >
              <font-awesome-icon
                :icon="['fas', 'eye']"
              />
            </b-button>

            <b-button
              v-if="showLogoPreview"
              data-test-id="button-logo-reset"
              variant="light"
              size="sm"
              class="py-0 ml-2"
              @click="resetLogo()"
            >
              {{ $t('logo.reset') }}
            </b-button>
          </div>
        </template>
        <b-form-file
          v-model="unifyAssets.logo"
          data-test-id="file-logo-upload"
          accept="image/*"
          :placeholder="$t('logo.placeholder')"
        />
      </b-form-group>

      <b-modal
        id="logo"
        hide-header
        hide-footer
        centered
        body-class="p-1"
      >
        <b-img
          data-test-id="img-logo-preview"
          :src="unify.logo"
          fluid-grow
        />
      </b-modal>

      <b-form-group
        :label="$t('url.label')"
        :description="$t('url.description')"
        label-cols="2"
      >
        <b-form-input
          v-model="unify.url"
          data-test-id="input-url"
        />
      </b-form-group>

      <b-form-group
        label-cols="2"
      >
        <b-form-checkbox
          v-model="unify.listed"
          data-test-id="checkbox-listed"
        >
          {{ $t('listed') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        label-cols="2"
      >
        <b-form-checkbox
          v-model="unify.pinned"
          data-test-id="checkbox-pinned"
          :disabled="!canPin"
        >
          {{ $t('pinned') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        :label="$t('config.label')"
        :description="$t('config.description')"
        class="mb-0"
      >
        <b-form-textarea
          v-model="unify.config"
          data-test-id="textarea-config"
          :state="configState"
          rows="10"
        />
      </b-form-group>
    </b-form>

    <template #header>
      <h3
        data-test-id="card-title"
        class="m-0"
      >
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        data-test-id="button-submit"
        class="float-right"
        :processing="processing"
        :success="success"
        :disabled="disabled"
        @submit="$emit('submit', { unify, unifyAssets })"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CApplicationEditorUnify',

  i18nOptions: {
    namespaces: 'system.applications',
    keyPrefix: 'editor.unify',
  },

  components: {
    CSubmitButton,
  },

  props: {
    unify: {
      type: Object,
      required: true,
    },

    application: {
      type: Object,
      required: true,
    },

    canPin: {
      type: Boolean,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },
  },

  data () {
    return {
      unifyAssets: {
        icon: undefined,
        logo: undefined,
      },
    }
  },

  computed: {
    disabled () {
      return this.validConfig === false
    },

    validConfig () {
      if (!this.unify) {
        return null
      }

      try {
        if ((this.unify.config || '').trim() !== '') {
          JSON.parse(this.unify.config)
        }
        return null
      } catch (e) {
        return false
      }
    },

    configState () {
      if (((this.unify || {}).config || '').trim() === '') {
        return null
      } else {
        return this.validConfig
      }
    },

    showLogoPreview () {
      return this.unify.logoID !== NoID
    },
  },

  created () {
    this.unify.name = this.unify.name ? this.unify.name : this.application.name
  },

  methods: {
    resetLogo () {
      this.unify.logo = undefined
      this.unify.logoID = NoID
    },
  },
}
</script>
