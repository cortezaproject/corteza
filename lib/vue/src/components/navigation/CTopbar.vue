<template>
  <div class="header-navigation d-flex flex-wrap flex-sm-nowrap align-items-center sticky-top pr-3">
    <div
      class="spacer"
      :class="{
        'expanded': sidebarPinned,
      }"
    />

    <h2
      class="title py-2 py-sm-0 mr-0 mr-sm-2 d-sm-inline-block text-truncate order-2 order-sm-0 mb-0"
    >
      <slot name="title" />
    </h2>

    <div class="tools-wrapper ml-auto text-sm-nowrap order-3 order-sm-0 pt-2 pb-3 py-sm-0">
      <slot name="tools" />
    </div>

    <div class="d-flex align-items-center order-1 order-sm-0 ml-auto ml-sm-0">
      <b-button
      v-if="!hideAppSelector && !settings.hideAppSelector"
      variant="outline-light"
      :href="appSelectorURL"
      size="lg"
      class="d-flex align-items-center justify-content-center text-dark border-0 nav-icon rounded-circle text-sm-nowrap"
    >
      <font-awesome-icon
        class="m-0 h5"
        :icon="['fas', 'grip-horizontal']"
      />
    </b-button>
    </div>

    <div class="d-flex align-items-center order-1 order-sm-0">
      <b-dropdown
        v-if="!settings.hideHelp"
        data-test-id="dropdown-helper"
        size="lg"
        variant="outline-light"
        class="nav-icon mx-1 text-sm-nowrap"
        toggle-class="text-decoration-none text-dark rounded-circle border-0 w-100"
        menu-class="topbar-dropdown-menu border-0 shadow-sm text-dark font-weight-bold mt-2"
        right
        no-caret
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
    </div>

    <div class="d-flex align-items-center order-1 order-sm-0 flex-grow-sm-0">
      <b-dropdown
        v-if="!settings.hideProfile"
        data-test-id="dropdown-profile"
        data-v-onboarding="profile"
        :variant="avatarExists ? 'link' : 'outline-light'"
        :toggle-class="`nav-icon text-decoration-none text-dark rounded-circle border ${avatarExists ? 'p-0' : ''}`"
        size="lg"
        right
        menu-class="topbar-dropdown-menu border-0 shadow-sm text-dark font-weight-bold mt-2"
        no-caret
        class="nav-user-icon"
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
        <b-dropdown-text class="text-muted mb-2">
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
export default {
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
  }
}
</script>

<style lang="scss" scoped>
$header-height: 64px;
$nav-width: 320px;
$nav-icon-size: 40px;
$nav-user-icon-size: 50px;

.icon-logo {
  height: calc(#{$header-height} / 2);
  background-repeat: no-repeat;
  background-position: center;
}

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
  height: $header-height;
  background-color: #F9FAFB !important;
  padding-left: calc(3.5rem + 6px);
}

.header-navigation .title > * {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.topbar-dropdown-menu {
  max-height: 80vh;
  overflow-y: auto;
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

.spacer {
  min-width: 0px;
  -webkit-transition: min-width 0.15s ease-in-out;
  -moz-transition: min-width 0.15s ease-in-out;
  -o-transition: min-width 0.15s ease-in-out;
  transition: min-width 0.15s ease-in-out;

  &.expanded {
    min-width: calc(#{$nav-width} - 42px);
    -webkit-transition: min-width 0.2s ease-in-out;
    -moz-transition: min-width 0.2s ease-in-out;
    -o-transition: min-width 0.2s ease-in-out;
    transition: min-width 0.2s ease-in-out;
  }
}

@media (max-width:576px) {
  .header-navigation {
    height: auto;
    padding-left: 1.5rem;
    .title {
      flex-basis: 100%;
    }
  }
  .tools-wrapper {
    flex-grow: 1;
    .vue-portal-target {
      display: flex;
      justify-content: end;
    }

    ~ div {
      height: 64px;
    }
  }
}
</style>

<style lang="scss">
.topbar-tools {
  .vue-portal-target {
    display: flex;
    align-items: center;
  }
}
</style>
