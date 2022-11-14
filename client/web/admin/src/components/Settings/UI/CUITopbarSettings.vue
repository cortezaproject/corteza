<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <b-form
      @submit.prevent="$emit('submit', settings)"
    >
      <h5>{{ $t('general') }}</h5>
      <b-form-group>
        <b-form-checkbox
          v-model="topbarSettings.hideAppSelector"
        >
          {{ $t('app-selector.hide') }}
        </b-form-checkbox>

        <b-form-checkbox
          v-model="topbarSettings.hideHelp"
        >
          {{ $t('help.hide') }}
        </b-form-checkbox>

        <b-form-checkbox
          v-model="topbarSettings.hideProfile"
        >
          {{ $t('profile.hide') }}
        </b-form-checkbox>
      </b-form-group>

      <div>
        <hr>
        <b-row>
          <b-col
            cols="12"
            lg="3"
          >
            <h5>{{ $t('help.title') }}</h5>
            <b-form-group>
              <b-form-checkbox
                v-model="topbarSettings.hideForumLink"
              >
                {{ $t('help.hide-forum-link') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-model="topbarSettings.hideDocumentationLink"
              >
                {{ $t('help.hide-documentation-link') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-model="topbarSettings.hideFeedbackLink"
              >
                {{ $t('help.hide-feedback-link') }}
              </b-form-checkbox>
            </b-form-group>
          </b-col>
          <b-col>
            <b-form-group
              label-class="pb-1"
            >
              <template #label>
                {{ $t('links.title') }}
                <b-button
                  variant="link"
                  class="text-decoration-none"
                  @click="topbarSettings.helpLinks.push({ handle: '', url: '', newTab: true })"
                >
                  + Add
                </b-button>
              </template>
              <b-table
                :fields="links.fields"
                :items="topbarSettings.helpLinks"
                table-variant="light"
                responsive="sm"
                striped
                small
                class="mb-0"
              >
                <template #cell(handle)="data">
                  <b-form-input
                    v-model="data.item.handle"
                    size="sm"
                  />
                </template>
                <template #cell(url)="data">
                  <b-form-input
                    v-model="data.item.url"
                    type="url"
                    size="sm"
                  />
                </template>
                <template #cell(newTab)="data">
                  <b-form-checkbox
                    v-model="data.item.newTab"
                  />
                </template>
                <template #cell(actions)="data">
                  <b-button
                    variant="outline-danger"
                    class="border-0 px-1"
                    @click="topbarSettings.helpLinks.splice(data.index, 1)"
                  >
                    <font-awesome-icon
                      :icon="['far', 'trash-alt']"
                    />
                  </b-button>
                </template>
              </b-table>
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <div>
        <hr>
        <b-row>
          <b-col
            cols="12"
            lg="3"
          >
            <h5>{{ $t('profile.title') }}</h5>
            <b-form-group>
              <b-form-checkbox
                v-model="topbarSettings.hideProfileLink"
              >
                {{ $t('profile.hide-profile-link') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-model="topbarSettings.hideChangePasswordLink"
              >
                {{ $t('profile.hide-change-password-link') }}
              </b-form-checkbox>
            </b-form-group>
          </b-col>
          <b-col>
            <b-form-group
              label-class="pb-1"
            >
              <template #label>
                {{ $t('links.title') }}
                <b-button
                  variant="link"
                  class="text-decoration-none"
                  @click="topbarSettings.profileLinks.push({ handle: '', url: '', newTab: true })"
                >
                  + Add
                </b-button>
              </template>
              <b-table
                :fields="links.fields"
                :items="topbarSettings.profileLinks"
                table-variant="light"
                responsive="sm"
                striped
                small
                class="mb-0"
              >
                <template #cell(handle)="data">
                  <b-form-input
                    v-model="data.item.handle"
                    size="sm"
                  />
                </template>
                <template #cell(url)="data">
                  <b-form-input
                    v-model="data.item.url"
                    type="url"
                    size="sm"
                  />
                </template>
                <template #cell(newTab)="data">
                  <b-form-checkbox
                    v-model="data.item.newTab"
                  />
                </template>
                <template #cell(actions)="data">
                  <b-button
                    variant="outline-danger"
                    class="border-0 px-1"
                    @click="topbarSettings.profileLinks.splice(data.index, 1)"
                  >
                    <font-awesome-icon
                      :icon="['far', 'trash-alt']"
                    />
                  </b-button>
                </template>
              </b-table>
            </b-form-group>
          </b-col>
        </b-row>
      </div>
    </b-form>

    <template #footer>
      <c-submit-button
        class="float-right"
        :disabled="!canManage"
        :processing="processing"
        :success="success"
        @submit="onSubmit"
      />
    </template>
  </b-card>
</template>

<script>
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CUITopbarSettings',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.topbar',
  },

  components: {
    CSubmitButton,
  },

  props: {
    settings: {
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

  data () {
    return {
      topbarSettings: {},

      links: {
        fields: [
          {
            key: 'handle',
            label: this.$t('links.handle'),
            thStyle: { width: '30%' },
          },
          {
            key: 'url',
            label: this.$t('links.url'),
            thStyle: { width: '55%' },
          },
          {
            key: 'newTab',
            label: this.$t('links.new-tab'),
            thClass: 'text-center',
            tdClass: 'text-center align-middle',
            thStyle: { width: '10%' },
          },
          {
            key: 'actions',
            label: '',
            tdClass: 'text-right align-middle',
            thStyle: { width: '1%' },
          },
        ],
      },
    }
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        this.topbarSettings = settings['ui.topbar'] || {}

        if (!this.topbarSettings.helpLinks) {
          this.topbarSettings.helpLinks = []
        }

        if (!this.topbarSettings.profileLinks) {
          this.topbarSettings.profileLinks = []
        }
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', { 'ui.topbar': this.topbarSettings })
    },
  },
}
</script>
