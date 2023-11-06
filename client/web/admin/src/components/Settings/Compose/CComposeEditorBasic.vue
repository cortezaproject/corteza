<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <b-form
      @submit.prevent="$emit('submit', basic)"
    >
      <div class="pb-3">
        <h5>{{ $t('attachments.page') }}</h5>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('attachments.max-size')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="basic['compose.page.attachments.max-size']"
                type="number"
                number
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('attachments.type.whitelist')"
              :description="$t('attachments.type.description')"
              label-class="text-primary"
              class="mb-0"
            >
              <b-form-input v-model="pageAttachmentWhitelist" />
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <div class="pb-3">
        <h5>{{ $t('attachments.record') }}</h5>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('attachments.max-size')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="basic['compose.record.attachments.max-size']"
                type="number"
                number
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('attachments.type.whitelist')"
              :description="$t('attachments.type.description')"
              label-class="text-primary"
              class="mb-0"
            >
              <b-form-input v-model="recordAttachmentWhitelist" />
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <div>
        <h5>{{ $t('attachments.icon') }}</h5>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('attachments.max-size')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="basic['compose.icon.attachments.max-size']"
                type="number"
                number
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('attachments.type.whitelist')"
              :description="$t('attachments.type.description')"
              label-class="text-primary"
              class="mb-0"
            >
              <b-form-input v-model="iconAttachmentWhitelist" />
            </b-form-group>
          </b-col>
        </b-row>
      </div>
    </b-form>

    <template #footer>
      <c-button-submit
        v-if="canManage"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', basic)"
      />
    </template>
  </b-card>
</template>

<script>

export default {
  name: 'CComposeEditorBasic',

  i18nOptions: {
    namespaces: 'compose.settings',
    keyPrefix: 'editor.basic',
  },

  props: {
    basic: {
      type: Object,
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

    canManage: {
      type: Boolean,
      required: true,
    },
  },

  computed: {
    pageAttachmentWhitelist: {
      get () {
        return (this.basic['compose.page.attachments.mimetypes'] || []).join(',')
      },

      set (value) {
        this.basic['compose.page.attachments.mimetypes'] = this.convertToExternal(value)
      },
    },

    recordAttachmentWhitelist: {
      get () {
        return (this.basic['compose.record.attachments.mimetypes'] || []).join(',')
      },

      set (value) {
        this.basic['compose.record.attachments.mimetypes'] = this.convertToExternal(value)
      },
    },

    iconAttachmentWhitelist: {
      get () {
        return (this.basic['compose.icon.attachments.mimetypes'] || []).join(',')
      },

      set (value) {
        this.basic['compose.icon.attachments.mimetypes'] = this.convertToExternal(value)
      },
    },
  },

  methods: {
    convertToExternal (value) {
      return (value || '').split(',').map(v => {
        return v.replace(/ /g, '')
      }).filter(v => {
        if (v.match(/^[-\w.]+\/[-\w/+.]+$/g)) {
          return v
        }
      })
    },
  },
}
</script>
