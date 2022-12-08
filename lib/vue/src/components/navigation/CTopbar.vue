<template>
  <div class="header-navigation d-flex align-items-center sticky-top pr-3">
    <div
      class="spacer"
      :class="{
        'expanded': sidebarPinned,
      }"
    />

    <h2
      class="d-none d-sm-inline-block text-truncate mb-0"
    >
      <slot name="title" />
    </h2>


    <div class="ml-auto text-sm-nowrap">
      <slot name="tools" />
    </div>

    <b-button
      v-if="!hideAppSelector && !settings.hideAppSelector"
      variant="outline-light"
      :href="appSelectorURL"
      size="lg"
      class="d-flex align-items-center justify-content-center text-dark border-0 nav-icon rounded-circle"
    >
      <font-awesome-icon
        class="m-0 h5"
        :icon="['fas', 'grip-horizontal']"
      />
    </b-button>

    <b-dropdown
      v-if="!settings.hideHelp"
      data-test-id="dropdown-helper"
      size="lg"
      variant="outline-light"
      class="nav-icon mx-1"
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
        :href="helpLink.url"
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
        href="mailto:info@crust.tech"
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
      data-v-onboarding="profile"
      size="lg"
      variant="outline-light"
      class="nav-user-icon"
      toggle-class="nav-icon text-decoration-none text-dark rounded-circle border"
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
        :href="profileLink.url"
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
  width: $nav-user-icon-size;
  height: $nav-user-icon-size;
}

.header-navigation {
  width: 100vw;
  height: $header-height;
  background-color: #F3F3F5 !important;
  padding-left: calc(3.5rem + 6px);
}

.topbar-dropdown-menu {
  max-height: 80vh;
  overflow-y: auto;
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
</style>

<style lang="scss">
.topbar-tools {
  .vue-portal-target {
    display: flex;
    align-items: center;
  }
}
</style>
