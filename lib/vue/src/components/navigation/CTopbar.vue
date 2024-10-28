<template>
  <div class="header-navigation d-flex flex-wrap align-items-center py-2 px-3 gap-2">
    <div
      class="sidebar-spacer"
      :class="{
        'expanded': sidebarPinned,
      }"
    />

    <h2 class="title mb-0">
      <slot name="title" />
    </h2>

    <div class="tools-wrapper ml-auto">
      <slot name="tools" />
    </div>

    <div class="d-flex align-items-center ml-auto">
      <b-button
        v-if="!hideAppSelector && !settings.hideAppSelector"
        data-test-id="app-selector"
        variant="outline-extra-light"
        :href="appSelectorURL"
        class="text-dark border-0 px-1"
      >
        {{ labels.appMenu }}
      </b-button>

      <b-dropdown
        v-if="!settings.hideHelp"
        data-test-id="dropdown-helper"
        size="lg"
        variant="outline-extra-light"
        toggle-class="text-decoration-none text-dark rounded-circle border-0 w-100"
        menu-class="topbar-dropdown-menu border-0 shadow-sm text-dark mt-2"
        right
        no-caret
        class="nav-icon mx-1 text-sm-nowrap"
      >
        <template #button-content>
          <div
            class="d-flex align-items-center justify-content-center"
          >
            <font-awesome-icon
              class="m-0 h5"
              :icon="['far', 'question-circle']"
            />
            <span class="sr-only">
              {{ labels.helpForum }}
            </span>
          </div>
        </template>
        <div>
          <slot name="help-dropdown" />
        </div>
        <b-dropdown-item
          v-for="(helpLink, index) in helpLinks"
          :key="index"
          :href="helpLink.url | checkValidURL"
          :target="helpLink.newTab ? '_blank' : ''"
        >
          {{ helpLink.handle }}
        </b-dropdown-item>
        <b-dropdown-item
          v-if="!settings.hideForumLink"
          data-test-id="dropdown-helper-forum"
          href="https://forum.cortezaproject.org/"
          target="_blank"
        >
          {{ labels.helpForum }}
        </b-dropdown-item>
        <b-dropdown-item
          v-if="!settings.hideDocumentationLink"
          data-test-id="dropdown-helper-docs"
          :href="documentationURL"
          target="_blank"
        >
          {{ labels.helpDocumentation }}
        </b-dropdown-item>
        <b-dropdown-item
          v-if="!settings.hideFeedbackLink"
          data-test-id="dropdown-helper-feedback"
          href="mailto:info@cortezaproject.org"
          target="_blank"
        >
          {{ labels.helpFeedback }}
        </b-dropdown-item>
        <b-dropdown-divider
          v-if="!onlyVersion"
        />
        <b-dropdown-item
          disabled
          class="small"
        >
          {{ labels.helpVersion }}
          <br>
          {{ frontendVersion }}
        </b-dropdown-item>
      </b-dropdown>

      <b-dropdown
        v-if="!settings.hideProfile"
        data-test-id="dropdown-profile"
        :variant="avatarExists ? 'link' : 'outline-extra-light'"
        :toggle-class="`nav-icon text-decoration-none text-dark rounded-circle border ${avatarExists ? 'p-0' : ''}`"
        size="lg"
        right
        menu-class="topbar-dropdown-menu border-0 shadow-sm text-dark mt-2"
        no-caret
        class="nav-user-icon"
        @hide="preventDropdownClose"
      >
        <template #button-content>
          <div
            v-if="avatarExists"
            class="avatar d-flex h-100"
            :style="{
              'background-image': avatarExists  ? `url(${profileAvatarUrl})` : 'none',
            }"
          />

          <div
            v-else
            class="d-flex align-items-center justify-content-center"
          >
            <font-awesome-icon
              class="m-0 h5"
              :icon="['far', 'user']"
            />
            <span class="sr-only">
              {{ labels.helpForum }}
            </span>
          </div>
        </template>

        <b-dropdown-text
          data-test-id="dropdown-item-username"
          class="text-muted mb-2"
        >
          {{ labels.userSettingsLoggedInAs }}
        </b-dropdown-text>

        <div>
          <slot name="avatar-dropdown" />
        </div>

        <b-dropdown-item
          v-for="(profileLink, index) in profileLinks"
          :key="index"
          :href="profileLink.url | checkValidURL"
          :target="profileLink.newTab ? '_blank' : ''"
        >
          {{ profileLink.handle }}
        </b-dropdown-item>

        <b-dropdown-item
          v-if="!settings.hideProfileLink"
          data-test-id="dropdown-profile-user"
          :href="userProfileURL"
          target="_blank"
        >
          {{ labels.userSettingsProfile }}
        </b-dropdown-item>

        <b-dropdown-item
          v-if="!settings.hideChangePasswordLink"
          data-test-id="dropdown-profile-change-password"
          :href="changePasswordURL"
          target="_blank"
        >
          {{ labels.userSettingsChangePassword }}
        </b-dropdown-item>

        <b-dropdown
        v-if="!settings.hideThemeSelector"
          id="theme-dropleft"
          variant="link"
          text="Theme"
          dropleft
          no-caret
          toggle-class="text-decoration-none text-left dropdown-item rounded-0"
          class="d-flex"
          @show="isThemeDropdownVisible = true"
          @hide="isThemeDropdownVisible = false"
          @click.prevent.stop
        >
          <b-dropdown-item
            v-for="theme in themes"
            :key="theme.id"
            :disabled="currentTheme === theme.id"
            @click="saveThemeMode(theme.id)"
          >
            {{ theme.label }}
          </b-dropdown-item>
        </b-dropdown>

        <b-dropdown-divider />

        <b-dropdown-item
          data-test-id="dropdown-profile-logout"
          href=""
          @click="$auth.logout()"
          class="mt-2"
        >
          {{ labels.userSettingsLogout }}
        </b-dropdown-item>
      </b-dropdown>
    </div>
  </div>
</template>

<script>
import { library } from '@fortawesome/fontawesome-svg-core'
import { faMoon, faSun} from '@fortawesome/free-solid-svg-icons'

library.add(faMoon, faSun)

export default {
  data() {
    return {
      currentTheme: 'light',
      isThemeDropdownVisible: false,
    }
  },

  props: {
    sidebarPinned: {
      type: Boolean,
      required: true,
      default: false,
    },

    hideAppSelector: {
      type: Boolean,
      required: false,
      default: false,
    },

    appSelectorURL: {
      type: String,
      default: '../'
    },

    settings: {
      type: Object,
      required: true,
    },

    labels: {
      type: Object,
      required: true,
    },
  },

  computed: {
    userProfileURL () {
      return this.$auth.cortezaAuthURL
    },

    changePasswordURL () {
      return `${this.$auth.cortezaAuthURL}/change-password`
    },

    documentationURL () {
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/index.html`
    },

    helpLinks () {
      const { helpLinks = [] } = this.settings || {}
      return (helpLinks || []).filter(({ handle, url }) => handle && url)
    },

    profileLinks () {
      const { profileLinks = [] } = this.settings || {}
      return (profileLinks || []).filter(({ handle, url }) => handle && url)
    },

    onlyVersion () {
      const {
        hideForumLink,
        hideDocumentationLink,
        hideFeedbackLink,
      } = this.settings || {}

      return !this.helpLinks.length && hideForumLink && hideDocumentationLink && hideFeedbackLink
    },

    frontendVersion () {
      /* eslint-disable no-undef */
      return VERSION
    },

    profileAvatarUrl () {
      return `${this.$SystemAPI.baseURL}/attachment/avatar/${this.$auth.user.meta.avatarID}/original/profile-photo-avatar`
    },

    avatarExists () {
      return this.$auth.user.meta.avatarID !== "0" && this.$auth.user.meta.avatarID
    },

    themes () {
      return [
        {
          id: 'light',
          label: this.labels.lightTheme,
        },
        {
          id: 'dark',
          label: this.labels.darkTheme,
        },
      ]
    },
  },

  watch: {
    '$auth.user.meta.theme': {
      immediate: true,
      handler (theme) {
        this.currentTheme = theme
      },
    },
  },

  methods: {
    async saveThemeMode (theme) {
      this.currentTheme = theme
      this.$set(this.$auth.user.meta, 'theme', theme)

      this.$SystemAPI.userUpdate(this.$auth.user).then(() => {
        document.getElementsByTagName('html')[0].setAttribute('data-color-mode', theme)
      }).catch(console.error)
    },

    preventDropdownClose (e) {
      if (this.isThemeDropdownVisible) {
        e.preventDefault()
      }
    },
  },
}
</script>

<style lang="scss" scoped>
$nav-icon-size: calc(var(--topbar-height) - 24px);
$nav-user-icon-size: calc(var(--topbar-height) - 16px);

.nav-icon {
  width: $nav-icon-size;
  height: $nav-icon-size;
}

.nav-user-icon {
  min-width: $nav-user-icon-size;
  min-height: $nav-user-icon-size;
}

.header-navigation {
  width: 100vw;
  min-height: var(--topbar-height);
  background-color: var(--topbar-bg);

  .sidebar-spacer.expanded {
    min-width: calc(var(--sidebar-width) - 50px);
  }
}

.avatar {
  border-radius: 50%;
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center;

  &:hover {
    opacity: 0.8;
    transition: opacity .25s ease-in-out;
    -moz-transition: opacity .25s ease-in-out;
    -webkit-transition: opacity .25s ease-in-out;
  }
}

.title {
  display: flex;
  align-items: center;
  min-height: $nav-user-icon-size;
  padding-left: 42px;

  .vue-portal-target {
    display: -webkit-box; /* For Safari and old versions of Chrome */
    display: -ms-flexbox; /* For old versions of IE */
    -webkit-box-orient: vertical; /* For Safari and old versions of Chrome */
    -webkit-line-clamp: 3; /* Maximum number of lines to display */
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.tools-wrapper {
  flex-grow: 1;

  .vue-portal-target {
    display: flex;
    justify-content: end;
    align-items: center;
    flex-wrap: wrap;
  }
}
</style>

<style lang="scss">
.topbar-dropdown-menu {
  z-index: 1100;
}

#theme-dropleft {
  .btn {
    font-family: var(--font-regular);
  }
}
</style>
