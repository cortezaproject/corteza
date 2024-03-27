<template>
  <b-card
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
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
              :label="$t('links.title')"
              label-class="text-primary"
            >
              <c-form-table-wrapper
                :labels="{ addButton: $t('general:label.add') }"
                @add-item="topbarSettings.helpLinks.push({ handle: '', url: '', newTab: true })"
              >
                <b-table
                  :fields="links.fields"
                  :items="topbarSettings.helpLinks"
                  thead-tr-class="text-primary"
                  responsive="sm"
                  small
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
                    <c-input-confirm
                      show-icon
                      @confirmed="topbarSettings.helpLinks.splice(data.index, 1)"
                    />
                  </template>
                </b-table>
              </c-form-table-wrapper>
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

              <b-form-checkbox
                v-model="topbarSettings.hideThemeSelector"
              >
                {{ $t('profile.hide-theme-selector') }}
              </b-form-checkbox>
            </b-form-group>
          </b-col>

          <b-col>
            <b-form-group
              :label="$t('links.title')"
              label-class="text-primary"
            >
              <c-form-table-wrapper
                :labels="{ addButton: $t('general:label.add') }"
                @add-item="topbarSettings.profileLinks.push({ handle: '', url: '', newTab: true })"
              >
                <b-table
                  :fields="links.fields"
                  :items="topbarSettings.profileLinks"
                  thead-tr-class="text-primary"
                  responsive="sm"
                  small
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
                    <c-input-confirm
                      show-icon
                      @confirmed="topbarSettings.profileLinks.splice(data.index, 1)"
                    />
                  </template>
                </b-table>
              </c-form-table-wrapper>
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
        @submit="onSubmit"
      />
    </template>
  </b-card>
</template>

<script>

export default {
  name: 'CUITopbarSettings',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.topbar',
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
          this.$set(this.topbarSettings, 'helpLinks', [])
        }

        if (!this.topbarSettings.profileLinks) {
          this.$set(this.topbarSettings, 'profileLinks', [])
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
