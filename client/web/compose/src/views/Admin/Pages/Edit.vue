<template>
  <div
    class="py-3"
  >
    <portal to="topbar-title">
      {{ $t('edit.edit') }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="page && page.canUpdatePage"
        size="sm"
      >
        <b-button
          data-test-id="button-page-builder"
          variant="primary"
          class="d-flex align-items-center"
          :to="{ name: 'admin.pages.builder' }"
        >
          {{ $t('label.pageBuilder') }}
          <font-awesome-icon
            :icon="['fas', 'tools']"
            class="ml-2"
          />
        </b-button>

        <page-translator
          :page.sync="page"
          :page-layouts.sync="layouts"
          button-variant="primary"
          style="margin-left:2px;"
        />

        <b-button
          v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.view'), container: '#body' }"
          variant="primary"
          :to="pageViewer"
          class="d-flex align-items-center"
          style="margin-left:2px;"
        >
          <font-awesome-icon
            :icon="['far', 'eye']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <b-container
      v-else
      fluid="xl"
    >
      <b-card
        no-body
        header-class="d-flex py-3 align-items-center border-bottom"
        class="shadow-sm"
      >
        <template #header>
          <b-dropdown
            v-if="page.canGrant || namespace.canGrant"
            data-test-id="dropdown-permissions"
            size="lg"
            variant="light"
            class="permissions-dropdown mr-1"
          >
            <template #button-content>
              <font-awesome-icon :icon="['fas', 'lock']" />
              <span>
                {{ $t('general:label.permissions') }}
              </span>
            </template>

            <b-dropdown-item-button>
              <c-permissions-button
                v-if="namespace.canGrant"
                :title="page.title || page.handle || page.pageID"
                :target="page.title || page.handle || page.pageID"
                :resource="`corteza::compose:page/${namespace.namespaceID}/${page.pageID}`"
                :button-label="$t('general:label.page')"
                :show-button-icon="false"
              />
            </b-dropdown-item-button>

            <b-dropdown-item-button>
              <c-permissions-button
                v-if="page.canGrant"
                :title="page.title || page.handle || page.pageID"
                :target="page.title || page.handle || page.pageID"
                :resource="`corteza::compose:page-layout/${namespace.namespaceID}/${page.pageID}/*`"
                :button-label="$t('general:label.pageLayout')"
                :show-button-icon="false"
                all-specific
              />
            </b-dropdown-item-button>
          </b-dropdown>
        </template>

        <b-row
          v-if="page"
          class="px-4 py-3"
        >
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('label.title')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="page.title"
                data-test-id="input-title"
                :placeholder="$t('placeholder.title')"
                required
                :state="titleState"
                class="mb-2"
              />
            </b-form-group>
          </b-col>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('label.handle')"
              label-class="text-primary"
            >
              <b-form-input
                v-model="page.handle"
                data-test-id="input-handle"
                :state="handleState"
                class="mb-2"
                :placeholder="$t('block.general.placeholder.handle')"
              />
              <b-form-invalid-feedback :state="handleState">
                {{ $t('block.general.invalid-handle-characters') }}
              </b-form-invalid-feedback>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
          >
            <b-form-group
              :label="$t('label.description')"
              label-class="text-primary"
            >
              <b-form-textarea
                v-model="page.description"
                data-test-id="input-description"
                :placeholder="$t('edit.pageDescription')"
                rows="4"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              label-class="d-flex align-items-center text-primary"
            >
              <template #label>
                {{ $t('icon.page') }}
                <b-button
                  v-b-tooltip.noninteractive.hover="{ title: $t('icon.configure'), container: '#body' }"
                  variant="outline-light"
                  class="d-flex align-items-center px-1 text-primary border-0 ml-1"
                  @click="openIconModal"
                >
                  <font-awesome-icon
                    :icon="['far', 'edit']"
                  />
                </b-button>
              </template>

              <img
                v-if="icon.src"
                :src="pageIcon"
                width="auto"
                height="50"
              >

              <span v-else>
                {{ $t('icon.noIcon') }}
              </span>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('edit.otherOptions')"
              label-class="text-primary"
            >
              <b-form-checkbox
                v-if="!isRecordPage"
                v-model="page.visible"
                data-test-id="checkbox-page-visibility"
              >
                {{ $t('edit.visible') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-model="page.config.navItem.expanded"
                data-test-id="checkbox-show-sub-pages-in-sidebar"
              >
                {{ $t('showSubPages') }}
              </b-form-checkbox>

              <b-form-checkbox
                v-if="isRecordPage"
                v-model="page.meta.notifications.enabled"
                data-test-id="checkbox-page-notifications-enabled"
              >
                {{ $t('edit.notifications.enabled') }}
              </b-form-checkbox>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
          >
            <hr>

            <b-form-group
              :label="$t('page-layout.layouts')"
              label-class="text-primary"
              class="mb-0"
            >
              <c-form-table-wrapper
                :labels="{ addButton: $t('general:label.add') }"
                @add-item="addLayout"
              >
                <b-table-simple
                  v-if="layouts.length > 0"
                  responsive="lg"
                  borderless
                  small
                >
                  <b-thead>
                    <tr>
                      <th
                        scope="col"
                        style="width: 40px;"
                      />
                      <th
                        class="text-primary"
                        scope="col"
                        style="min-width: 300px;"
                      >
                        {{ $t('page-layout.title') }}
                      </th>
                      <th
                        class="text-primary"
                        scope="col"
                        style="min-width: 300px;"
                      >
                        {{ $t('page-layout.handle') }}
                      </th>
                      <th
                        scope="col"
                        style="min-width: 100px;"
                      />
                    </tr>
                  </b-thead>

                  <draggable
                    v-model="layouts"
                    handle=".grab"
                    group="layouts"
                    tag="b-tbody"
                  >
                    <tr
                      v-for="(layout, index) in layouts"
                      :key="index"
                    >
                      <b-td
                        class="grab text-center align-middle"
                      >
                        <font-awesome-icon
                          :icon="['fas', 'bars']"
                          class="text-secondary"
                        />
                      </b-td>
                      <b-td
                        class="align-middle"
                      >
                        <b-input-group>
                          <b-form-input
                            v-model="layout.meta.title"
                            :state="layoutTitleState(layout.meta.title)"
                            @input="layout.meta.updated = true"
                          />
                          <b-input-group-append>
                            <page-layout-translator
                              :page-layout="layout"
                              :disabled="layout.pageLayoutID === '0'"
                              highlight-key="meta.title"
                            />
                          </b-input-group-append>
                        </b-input-group>
                      </b-td>
                      <b-td
                        class="align-middle"
                      >
                        <b-input-group>
                          <b-form-input
                            v-model="layout.handle"
                            :state="layoutHandleState(layout.handle)"
                            @input="layout.meta.updated = true"
                          />
                          <b-input-group-append>
                            <b-button
                              v-b-tooltip.noninteractive.hover="{ title: $t('page-layout.tooltip.configure'), container: '#body' }"
                              variant="extra-light"
                              class="d-flex align-items-center px-3"
                              @click="configureLayout(index)"
                            >
                              <font-awesome-icon
                                :icon="['fas', 'wrench']"
                              />
                            </b-button>
                            <b-button
                              v-b-tooltip.noninteractive.hover="{ title: $t('page-layout.tooltip.builder'), container: '#body' }"
                              variant="primary"
                              :disabled="layout.pageLayoutID === '0'"
                              class="d-flex align-items-center"
                              :to="{ name: 'admin.pages.builder', query: { layoutID: layout.pageLayoutID } }"
                            >
                              <font-awesome-icon
                                :icon="['fas', 'tools']"
                              />
                            </b-button>
                          </b-input-group-append>
                        </b-input-group>
                      </b-td>
                      <td
                        class="text-right align-middle"
                        style="min-width: 100px;"
                      >
                        <c-permissions-button
                          v-if="page.canGrant && layout.pageLayoutID !== '0'"
                          button-variant="outline-extra-light"
                          size="sm"
                          :title="layout.meta.title || layout.handle || layout.pageLayoutID"
                          :target="layout.meta.title || layout.handle || layout.pageLayoutID"
                          :tooltip="$t('permissions:resources.compose.page-layout.tooltip')"
                          :resource="`corteza::compose:page-layout/${layout.namespaceID}/${layout.pageID}/${layout.pageLayoutID}`"
                          class="text-dark border-0 mr-2"
                        />
                        <c-input-confirm
                          show-icon
                          @confirmed="removeLayout(index)"
                        />
                      </td>
                    </tr>
                  </draggable>
                </b-table-simple>
              </c-form-table-wrapper>
            </b-form-group>
          </b-col>
        </b-row>
      </b-card>
    </b-container>

    <b-modal
      v-if="layoutEditor.layout"
      :visible="!!layoutEditor.layout"
      :title="$t('page-layout.configure', { title: ((layoutEditor.layout || {}).meta || {}).title, interpolation: { escapeValue: false } })"
      :ok-title="$t('general:label.saveAndClose')"
      ok-variant="primary"
      :ok-disabled="!layoutEditor.layout.meta.title"
      cancel-variant="light"
      size="xl"
      scrollable
      no-fade
      @ok="updateLayout()"
      @cancel="layoutEditor.layout = undefined"
      @hide="layoutEditor.layout = undefined"
    >
      <h5 class="mb-3">
        {{ $t('page-layout.general') }}
      </h5>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('page-layout.title')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="layoutEditor.layout.meta.title"
                :state="layoutTitleState(layoutEditor.layout.meta.title)"
                @input="layoutEditor.layout.meta.updated = true"
              />

              <b-input-group-append>
                <page-layout-translator
                  :page-layout="layoutEditor.layout"
                  :disabled="layoutEditor.layout.pageLayoutID === '0'"
                  highlight-key="meta.title"
                />
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('page-layout.handle')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="layoutEditor.layout.handle"
              :state="layoutHandleState(layoutEditor.layout.handle)"
              @input="layoutEditor.layout.meta.updated = true"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <b-form-group
        v-if="isRecordPage"
        :label="$t('page-layout.useTitle')"
        label-class="text-primary ml-auto mt-2"
      >
        <c-input-checkbox
          v-model="layoutEditor.layout.config.useTitle"
          switch
          :labels="checkboxLabel"
        />
        <i18next
          path="page-layout.tooltip.title"
          tag="small"
          class="text-muted"
        >
          <code>${record.values.fieldName}</code>
          <code>${recordID}</code>
          <code>${ownerID}</code>
          <span><code>${userID}</code>, <code>${user.name}</code></span>
        </i18next>
      </b-form-group>

      <hr>

      <h5 class="mb-3">
        {{ $t('page-layout.visibility') }}
      </h5>

      <b-form-group
        label-class="d-flex align-items-center text-primary mb-0"
      >
        <template #label>
          {{ $t('page-layout.condition.label') }}
          <c-hint
            :tooltip="$t('page-layout.tooltip.performance.condition')"
            icon-class="text-warning"
          />
        </template>
        <b-input-group>
          <b-input-group-prepend>
            <b-button variant="extra-light">
              ƒ
            </b-button>
          </b-input-group-prepend>
          <b-form-input
            v-model="layoutEditor.layout.config.visibility.expression"
            :placeholder="$t('page-layout.condition.placeholder')"
          />
          <b-input-group-append>
            <b-button
              variant="outline-secondary"
              :href="documentationURL"
              class="d-flex justify-content-center align-items-center"
              target="_blank"
            >
              ?
            </b-button>
          </b-input-group-append>
        </b-input-group>

        <i18next
          v-if="isRecordPage"
          path="page-layout.condition.description.record-page"
          tag="small"
          class="text-muted"
        >
          <code>record.values.fieldName</code>
          <code>user.(userID/email...)</code>
          <code>screen.(width/height)</code>
          <code>isView/isCreate/isEdit</code>
          <code>user.userID == record.values.createdBy</code>
          <code>screen.width &lt; 1024</code>
        </i18next>

        <i18next
          v-else
          path="page-layout.condition.description.non-record-page"
          tag="small"
          class="text-muted"
        >
          <code>user.(userID/email...)</code>
          <code>screen.(width/height)</code>
          <code>user.email == "test@mail.com"</code>
          <code>screen.width &lt; 1024</code>
        </i18next>
      </b-form-group>

      <b-form-group
        :label="$t('page-layout.roles.label')"
        label-class="text-primary"
      >
        <b-spinner v-if="resolvingLayoutRoles" />

        <c-input-role
          v-else
          :value="getLayoutRoles()"
          :placeholder="$t('page-layout.roles.placeholder')"
          multiple
          @input="onLayoutRoleChange"
        />
      </b-form-group>

      <template v-if="isRecordPage">
        <hr>

        <h5 class="mb-3">
          {{ $t('page-layout.recordToolbar.label') }}
        </h5>

        <b-form-group
          :label="$t('page-layout.recordToolbar.buttons.label')"
          label-class="text-primary"
        >
          <b-form-checkbox
            v-model="layoutEditor.layout.config.buttons.back.enabled"
          >
            {{ $t('page-layout.recordToolbar.buttons.showBack') }}
          </b-form-checkbox>

          <b-form-checkbox
            v-model="layoutEditor.layout.config.buttons.delete.enabled"
          >
            {{ $t('page-layout.recordToolbar.buttons.showDelete') }}
          </b-form-checkbox>

          <b-form-checkbox
            v-model="layoutEditor.layout.config.buttons.clone.enabled"
          >
            {{ $t('page-layout.recordToolbar.buttons.showClone') }}
          </b-form-checkbox>

          <b-form-checkbox
            v-model="layoutEditor.layout.config.buttons.new.enabled"
          >
            {{ $t('page-layout.recordToolbar.buttons.showNew') }}
          </b-form-checkbox>

          <b-form-checkbox
            v-model="layoutEditor.layout.config.buttons.edit.enabled"
          >
            {{ $t('page-layout.recordToolbar.buttons.showEdit') }}
          </b-form-checkbox>

          <b-form-checkbox
            v-model="layoutEditor.layout.config.buttons.submit.enabled"
          >
            {{ $t('page-layout.recordToolbar.buttons.showSave') }}
          </b-form-checkbox>
        </b-form-group>

        <b-form-group
          :label="$t('page-layout.recordToolbar.actions.label')"
          label-class="text-primary"
          class="mb-0"
        >
          <c-form-table-wrapper
            :labels="{ addButton: $t('general:label.add') }"
            @add-item="addLayoutAction"
          >
            <b-table-simple
              v-if="layoutEditor.layout.config.actions.length > 0"
              responsive
              borderless
              small
              class="layout-actions"
            >
              <draggable
                v-model="layoutEditor.layout.config.actions"
                handle=".grab"
                group="actions"
                tag="b-tbody"
              >
                <tr
                  v-for="(action, index) in layoutEditor.layout.config.actions"
                  :key="index"
                  :class="{ 'border-top border-light': index > 0 }"
                >
                  <b-td style="width: 40px;">
                    <div
                      class="grab d-flex align-items-center justify-content-center"
                      style="height: calc(1.5em + 0.75rem + 45px);"
                    >
                      <font-awesome-icon
                        :icon="['fas', 'bars']"
                        class="text-secondary"
                      />
                    </div>
                  </b-td>
                  <b-td style="min-width: 250px;">
                    <b-form-group
                      :label="$t('page-layout.recordToolbar.actions.buttonLabel')"
                      label-class="text-primary"
                      class="mb-1"
                    >
                      <b-form-input
                        v-model="action.meta.label"
                        class="mb-1"
                      />
                    </b-form-group>
                    <b-form-group
                      v-if="action.kind === 'toLayout'"
                      :label="$t('page-layout.recordToolbar.actions.toLayout.label')"
                      label-class="text-primary"
                      class="mb-0"
                    >
                      <b-form-select
                        v-model="action.params.pageLayoutID"
                        :options="actionLayoutOptions"
                        value-field="pageLayoutID"
                        text-field="label"
                      />
                    </b-form-group>
                    <b-form-group
                      v-if="action.kind === 'toURL'"
                      :label="$t('page-layout.recordToolbar.actions.toURL.label')"
                      label-class="text-primary"
                      class="mb-0"
                    >
                      <b-form-input
                        v-model="action.params.url"
                        type="url"
                        :placeholder="$t('page-layout.recordToolbar.actions.toURL.placeholder')"
                      />
                    </b-form-group>
                  </b-td>
                  <b-td style="min-width: 250px;">
                    <b-form-group
                      :label="$t('page-layout.recordToolbar.actions.kind.label')"
                      label-class="text-primary"
                      class="mb-1"
                    >
                      <b-form-select
                        v-model="action.kind"
                        :options="actionKindOptions"
                        class="mb-1"
                        @change="onActionKindChange(action)"
                      />
                    </b-form-group>
                    <b-form-group
                      v-if="action.kind === 'toURL'"
                      :label="$t('page-layout.recordToolbar.actions.openIn.label')"
                      label-class="text-primary"
                      class="mb-0"
                    >
                      <b-form-select
                        v-model="action.params.openIn"
                        :options="actionOpenInOptions"
                      />
                    </b-form-group>
                  </b-td>
                  <b-td style="min-width: 150px;">
                    <b-form-group
                      :label="$t('page-layout.recordToolbar.actions.variant')"
                      label-class="text-primary"
                    >
                      <b-form-select
                        v-model="action.meta.style.variant"
                        :options="actionVariantOptions"
                      />
                    </b-form-group>
                  </b-td>
                  <b-td style="min-width: 100px;">
                    <b-form-group
                      :label="$t('page-layout.recordToolbar.actions.placement.label')"
                      label-class="text-primary"
                    >
                      <b-form-select
                        v-model="action.placement"
                        :options="actionPlacementOptions"
                      />
                    </b-form-group>
                  </b-td>
                  <b-td style="min-width: 80px;">
                    <b-form-group
                      :label="$t('page-layout.recordToolbar.actions.visible')"
                      label-class="text-primary text-center"
                    >
                      <div
                        class="d-flex align-items-center justify-content-center"
                        style="height: calc(1.5em + 0.75rem + 2px);"
                      >
                        <b-form-checkbox
                          v-model="action.enabled"
                          class="ml-2"
                        />
                      </div>
                    </b-form-group>
                  </b-td>
                  <b-td style="min-width: 80px;">
                    <div
                      class="d-flex align-items-center justify-content-end"
                      style="height: calc(1.5em + 0.75rem + 45px);"
                    >
                      <c-input-confirm
                        show-icon
                        class="ml-2"
                        @confirmed="removeLayoutAction(index)"
                      />
                    </div>
                  </b-td>
                </tr>
              </draggable>
            </b-table-simple>
          </c-form-table-wrapper>
        </b-form-group>
      </template>
    </b-modal>

    <b-modal
      v-model="showIconModal"
      :title="$t('icon.configure')"
      size="lg"
      label-class="text-primary"
      footer-class="d-flex align-items-center"
      no-fade
    >
      <template #modal-footer>
        <c-input-confirm
          v-if="attachments && selectedAttachmentID"
          :disabled="(attachments && !selectedAttachmentID) || processingIcon"
          size="md"
          variant="danger"
          @confirmed="deleteIcon"
        >
          {{ $t('icon.delete') }}
        </c-input-confirm>

        <div class="ml-auto">
          <b-button
            variant="light"
            @click="closeIconModal"
          >
            {{ $t('general:label.cancel') }}
          </b-button>

          <b-button
            variant="primary"
            class="ml-2"
            @click="saveIconModal"
          >
            {{ $t('general:label.saveAndClose') }}
          </b-button>
        </div>
      </template>

      <b-form-group
        :label="$t('icon.upload')"
        label-class="text-primary"
        class="mb-0"
      >
        <uploader
          :endpoint="endpoint"
          :accepted-files="['image/*']"
          :param-name="'icon'"
          @uploaded="uploadAttachment"
        />

        <b-form-group
          :label="$t('url.label')"
          label-class="text-primary"
          class="my-2"
        >
          <b-input-group>
            <b-input
              v-model="linkUrl"
              :disabled="isIconSet"
            />
            <b-input-group-append>
              <b-button
                v-b-modal.logo
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.preview-link'), container: '#body' }"
                :disabled="!linkUrl"
                variant="light"
                rounded
                class="d-flex align-items-center btn-light"
              >
                <font-awesome-icon :icon="['fas', 'external-link-alt']" />
              </b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
      </b-form-group>

      <template v-if="attachments.length > 0">
        <hr>

        <b-form-group
          :label="$t('icon.list')"
          label-class="text-primary"
        >
          <div
            v-if="processingIcon"
            class="d-flex align-items-center justify-content-center h-100"
          >
            <b-spinner />
          </div>

          <div
            v-else
            class="d-flex flex-wrap"
          >
            <img
              v-for="a in attachments"
              :key="a.attachmentID"
              :src="a.src"
              :alt="a.name"
              width="auto"
              height="50"
              :class="{ 'selected-icon': selectedAttachmentID === a.attachmentID }"
              class="rounded pointer m-2"
              @click="toggleSelectedIcon(a.attachmentID)"
            >
          </div>
        </b-form-group>
      </template>
    </b-modal>

    <b-modal
      id="logo"
      hide-header
      hide-footer
      centered
      no-fade
      body-class="p-1"
    >
      <b-img
        :src="linkUrl"
        fluid-grow
      />
    </b-modal>

    <portal to="admin-toolbar">
      <editor-toolbar
        :hide-delete="hideDelete"
        :hide-clone="hideClone"
        :hide-save="hideSave"
        :disable-save="disableSave"
        :processing="processing"
        :processing-save="processingSave"
        :processing-save-and-close="processingSaveAndClose"
        :processing-delete="processingDelete"
        :processing-clone="processingClone"
        @clone="handleClone()"
        @delete="handleDeletePage()"
        @save="handleSave()"
        @saveAndClose="handleSave({ closeOnSuccess: true })"
        @back="$router.push(previousPage || { name: 'admin.pages' })"
      >
        <template #delete>
          <b-dropdown
            v-if="showDeleteDropdown"
            data-test-id="dropdown-delete"
            size="lg"
            variant="danger"
            :text="$t('general:label.delete')"
          >
            <b-dropdown-item-button
              data-test-id="dropdown-item-delete-update-parent-of-sub-pages"
              @click="handleDeletePage('rebase')"
            >
              {{ $t('delete.rebase') }}
            </b-dropdown-item-button>
            <b-dropdown-item-button
              data-test-id="dropdown-item-delete-sub-pages"
              @click="handleDeletePage('cascade')"
            >
              {{ $t('delete.cascade') }}
            </b-dropdown-item-button>
          </b-dropdown>
        </template>
      </editor-toolbar>
    </portal>
  </div>
</template>

<script>
import { isEqual } from 'lodash'
import { mapGetters, mapActions } from 'vuex'
import EditorToolbar from 'corteza-webapp-compose/src/components/Admin/EditorToolbar'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import PageLayoutTranslator from 'corteza-webapp-compose/src/components/Admin/PageLayout/PageLayoutTranslator'
import pages from 'corteza-webapp-compose/src/mixins/pages'
import Uploader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Uploader'
import Draggable from 'vuedraggable'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { handle, components } from '@cortezaproject/corteza-vue'
const { CInputRole } = components

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  name: 'PageEdit',

  components: {
    EditorToolbar,
    PageTranslator,
    PageLayoutTranslator,
    Uploader,
    Draggable,
    CInputRole,
  },

  mixins: [
    pages,
  ],

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedComposePage(next)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedComposePage(next)
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    pageID: {
      type: String,
      required: true,
    },
  },

  data () {
    return {
      processing: false,
      processingIcon: false,
      processingSave: false,
      processingSaveAndClose: false,
      processingClone: false,
      processingDelete: false,

      page: new compose.Page(),
      initialPageState: new compose.Page(),

      showIconModal: false,
      attachments: [],
      selectedAttachmentID: '',
      linkUrl: '',

      layouts: [],

      layoutEditor: {
        index: undefined,
        layout: undefined,
      },

      resolvingLayoutRoles: false,
      resolvedRoles: {},

      removedLayouts: new Set(),

      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
      previousPage: 'ui/previousPage',
    }),

    titleState () {
      return this.page.title.length > 0 ? null : false
    },

    handleState () {
      return handle.handleState(this.page.handle)
    },

    pageViewer () {
      const { pageID } = this.page
      const name = this.isRecordPage ? 'page.record.create' : 'page'
      return { name, params: { pageID } }
    },

    isRecordPage () {
      return this.page && this.page.moduleID !== NoID
    },

    hasChildren () {
      return this.page ? this.pages.some(({ selfID }) => selfID === this.page.pageID) : false
    },

    disableSave () {
      return !this.page || [this.titleState, this.handleState].includes(false) || this.layouts.some(l => !l.meta.title || handle.handleState(l.handle) === false)
    },

    hideDelete () {
      return !this.page || this.hasChildren || !this.page.canDeletePage || !!this.page.deletedAt
    },

    hideSave () {
      return !this.page || !this.page.canUpdatePage
    },

    hideClone () {
      return !this.page || this.page.moduleID !== NoID
    },

    showDeleteDropdown () {
      return this.hasChildren && this.page.canDeletePage && !this.page.deletedAt
    },

    endpoint () {
      return this.$ComposeAPI.iconUploadEndpoint({
        namespaceID: this.namespace.namespaceID,
      })
    },

    icon: {
      get () {
        return this.page.config.navItem.icon || {}
      },

      set (icon) {
        this.$set(this.page.config.navItem, 'icon', icon)
      },
    },

    isIconSet () {
      return !!this.selectedAttachmentID
    },

    pageIcon () {
      if (!this.icon.src) {
        return
      }

      return this.icon.type === 'link' ? this.icon.src : this.makeAttachmentUrl(this.icon.src)
    },

    documentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/compose-configuration/page-layouts.html#visibility-condition`
    },

    currentLayoutRoles: {
      get () {
        if (!this.layoutEditor.layout) {
          return []
        }

        return this.layoutEditor.layout.config.visibility.roles
      },

      set (roles) {
        this.$set(this.layoutEditor.layout.config.visibility, 'roles', roles)
      },
    },

    actionKindOptions () {
      return [
        { value: 'toLayout', text: this.$t('page-layout.recordToolbar.actions.kind.toLayout') },
        { value: 'toURL', text: this.$t('page-layout.recordToolbar.actions.kind.toURL') },
      ]
    },

    actionLayoutOptions () {
      return [
        { pageLayoutID: '', label: this.$t('page-layout.recordToolbar.actions.toLayout.placeholder') },
        ...this.layouts.filter(({ pageLayoutID }) => pageLayoutID !== NoID)
          .map(({ pageLayoutID, handle, meta }) => ({ pageLayoutID, label: meta.title || handle || pageLayoutID })),
      ]
    },

    actionOpenInOptions () {
      return [
        { value: 'sameTab', text: this.$t('page-layout.recordToolbar.actions.openIn.sameTab') },
        { value: 'newTab', text: this.$t('page-layout.recordToolbar.actions.openIn.newTab') },
      ]
    },

    actionVariantOptions () {
      return [
        { value: 'primary', text: this.$t('general:variants.primary') },
        { value: 'secondary', text: this.$t('general:variants.secondary') },
        { value: 'success', text: this.$t('general:variants.success') },
        { value: 'warning', text: this.$t('general:variants.warning') },
        { value: 'danger', text: this.$t('general:variants.danger') },
        { value: 'info', text: this.$t('general:variants.info') },
        { value: 'light', text: this.$t('general:variants.light') },
        { value: 'dark', text: this.$t('general:variants.dark') },
      ]
    },

    actionPlacementOptions () {
      return [
        { value: 'start', text: this.$t('page-layout.recordToolbar.actions.placement.start') },
        { value: 'center', text: this.$t('page-layout.recordToolbar.actions.placement.center') },
        { value: 'end', text: this.$t('page-layout.recordToolbar.actions.placement.end') },
      ]
    },
  },

  watch: {
    pageID: {
      immediate: true,
      handler (pageID) {
        this.page = undefined
        this.initialPageState = undefined
        this.layouts = []

        this.removedLayouts = new Set()

        if (pageID) {
          this.processing = true

          const { namespaceID } = this.namespace
          this.findPageByID({ namespaceID, pageID, force: true }).then((page) => {
            this.page = page.clone()
            this.initialPageState = page.clone()
            return this.fetchAttachments()
          }).then(this.fetchLayouts)
            .finally(() => {
              this.processing = false
            }).catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
        }
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      findPageByID: 'page/findByID',
      updatePage: 'page/update',
      deletePage: 'page/delete',
      createPage: 'page/create',
      loadPages: 'page/load',
      findLayoutsByPageID: 'pageLayout/findByPageID',
      createPageLayout: 'pageLayout/create',
      updatePageLayout: 'pageLayout/update',
      deletePageLayout: 'pageLayout/delete',
    }),

    async fetchLayouts () {
      const { namespaceID } = this.namespace
      return this.findLayoutsByPageID({ namespaceID, pageID: this.pageID, force: true }).then(layouts => {
        this.layouts = layouts.map((layout) => new compose.PageLayout(layout))
      })
    },

    async resolveLayoutRoles () {
      if (this.currentLayoutRoles.length) {
        this.resolvingLayoutRoles = true

        Promise.all(this.currentLayoutRoles.map(roleID => {
          if (this.resolvedRoles[roleID]) {
            return Promise.resolve()
          }

          return this.$SystemAPI.roleRead({ roleID }).then(role => {
            this.resolvedRoles[roleID] = role
          })
        })).finally(() => {
          this.resolvingLayoutRoles = false
        })
      }
    },

    getLayoutRoles () {
      return this.currentLayoutRoles.map(roleID => this.resolvedRoles[roleID])
    },

    onLayoutRoleChange (roles) {
      roles.forEach(r => {
        if (!this.resolvedRoles[r.roleID]) {
          this.resolvedRoles[r.roleID] = r
        }
      })

      this.currentLayoutRoles = roles.map(r => r.roleID)
    },

    isRoleVisible ({ meta }) {
      return !(meta.context && meta.context.resourceTypes)
    },

    addLayout () {
      this.layouts.push(new compose.PageLayout({ namespaceID: this.namespace.namespaceID, pageID: this.pageID }))
    },

    updateLayout () {
      this.layoutEditor.layout.meta.updated = true
      this.layouts.splice(this.layoutEditor.index, 1, this.layoutEditor.layout)
      this.layoutEditor.index = undefined
      this.layoutEditor.layout = undefined
    },

    removeLayout (index) {
      const { pageLayoutID } = this.layouts[index] || {}
      if (pageLayoutID !== NoID) {
        this.removedLayouts.add(this.layouts[index])
      }

      this.layouts.splice(index, 1)
    },

    configureLayout (index) {
      this.layoutEditor.index = index
      this.layoutEditor.layout = new compose.PageLayout(this.layouts[index])

      this.resolveLayoutRoles()
    },

    async handleSaveLayouts () {
      // Delete first so old deleted handles don't interfere with new identical ones
      return Promise.all([...this.removedLayouts].map(this.deletePageLayout)).then(() => {
        return Promise.all(this.layouts.map(layout => {
          if (layout.pageLayoutID === NoID) {
            return this.createPageLayout(layout)
          } else if (layout.meta.updated) {
            return this.updatePageLayout(layout)
          }

          return Promise.resolve([])
        }))
      })
    },

    async handlePageLayoutReorder () {
      const { namespaceID } = this.namespace
      const pageIDs = this.layouts.map(({ pageLayoutID }) => pageLayoutID)

      return this.$ComposeAPI.pageLayoutReorder({ namespaceID, pageID: this.pageID, pageIDs }).then(() => {
        return this.$store.dispatch('pageLayout/load', { namespaceID, clear: true, force: true })
      })
    },

    handleSave ({ closeOnSuccess = false } = {}) {
      const toggleProcessing = () => {
        this.processing = !this.processing

        if (closeOnSuccess) {
          this.processingSaveAndClose = !this.processingSaveAndClose
        } else {
          this.processingSave = !this.processingSave
        }
      }

      /**
       * Pass a special tag alongside payload that
       * instructs store layer to add content-language header to the API request
       */
      const resourceTranslationLanguage = this.currentLanguage

      toggleProcessing()
      const { namespaceID } = this.namespace

      return this.saveIcon().then(icon => {
        this.page.config.navItem.icon = icon
        return this.updatePage({ namespaceID, ...this.page, resourceTranslationLanguage })
      }).then(page => {
        this.page = page.clone()
        this.initialPageState = page.clone()
        return this.handleSaveLayouts()
      }).then(this.handlePageLayoutReorder)
        .then(() => {
          this.fetchLayouts()
          this.removedLayouts = new Set()

          this.toastSuccess(this.$t('notification:page.saved'))
          if (closeOnSuccess) {
            this.$router.push(this.previousPage || { name: 'admin.pages' })
          }
        }).finally(() => {
          toggleProcessing()
        }).catch(this.toastErrorHandler(this.$t('notification:page.saveFailed')))
    },

    handleDeletePage (strategy = 'abort') {
      this.processingDelete = true

      this.deletePage({ ...this.page, strategy }).then(() => {
        setTimeout(() => {
          this.toastSuccess(this.$t('notification:page.deleted'))
          this.$router.push({ name: 'admin.pages' })
        }, 300)
      })
        .catch(this.toastErrorHandler(this.$t('notification:page.deleteFailed')))
        .finally(() => {
          setTimeout(() => {
            this.processingDelete = false
          }, 300)
        })
    },

    uploadAttachment ({ attachmentID }) {
      this.fetchAttachments()
      this.toggleSelectedIcon(attachmentID)
    },

    async fetchAttachments () {
      this.processingIcon = true

      return this.$ComposeAPI.iconList({ sort: 'id DESC' })
        .then(({ set: attachments = [] }) => {
          const baseURL = this.$ComposeAPI.baseURL
          this.attachments = []

          if (attachments.length === 0) {
            this.icon = {}
            this.initialPageState.config.navItem.icon = {}
          } else {
            attachments.forEach(a => {
              const src = !a.url.includes(baseURL) ? this.makeAttachmentUrl(a.url) : a.url
              this.attachments.push({ ...a, src })
            })
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.iconFetchFailed')))
        .finally(() => {
          this.processingIcon = false
        })
    },

    addLayoutAction () {
      this.layoutEditor.layout.addAction()
    },

    removeLayoutAction (index) {
      this.layoutEditor.layout.config.actions.splice(index, 1)
    },

    onActionKindChange (action) {
      if (action.kind === 'toURL' && !action.params.openIn) {
        this.$set(action.params, 'openIn', 'sameTab')
      }
    },

    async saveIcon () {
      return this.$ComposeAPI.pageUpdateIcon({
        namespaceID: this.namespace.namespaceID,
        pageID: this.pageID,
        type: this.icon.type || 'link',
        source: this.icon.src,
      })
    },

    toggleSelectedIcon (attachmentID = '') {
      this.selectedAttachmentID = this.selectedAttachmentID === attachmentID ? '' : attachmentID
    },

    openIconModal () {
      this.linkUrl = this.icon.type === 'link' ? this.icon.src : ''
      this.setCurrentIcon()
      this.showIconModal = true
    },

    saveIconModal () {
      const type = this.selectedAttachmentID ? 'attachment' : 'link'

      let src = this.linkUrl
      if (this.selectedAttachmentID) {
        src = (this.attachments.find(({ attachmentID }) => attachmentID === this.selectedAttachmentID) || {}).url
      }

      this.icon = { type, src }

      if (type === 'link' && !src) {
        this.icon = {}
      }

      this.showIconModal = false
    },

    deleteIcon () {
      this.processingIcon = true

      return this.$ComposeAPI.iconDelete({ iconID: this.selectedAttachmentID }).then(() => {
        return this.fetchAttachments().then(() => {
          this.setCurrentIcon()
          this.toastSuccess(this.$t('notification:page.iconDeleteSuccess'))
        })
      }).finally(() => {
        this.processingIcon = false
      }).catch(this.toastErrorHandler(this.$t('notification:page.iconDeleteFailed')))
    },

    closeIconModal () {
      this.linkUrl = this.icon.type === 'link' ? this.icon.src : ''
      this.setCurrentIcon()
      this.showIconModal = false
    },

    setCurrentIcon () {
      this.selectedAttachmentID = (this.attachments.find(a => a.url === this.icon.src) || {}).attachmentID

      if (!this.selectedAttachmentID) {
        this.icon = {}
      }
    },

    makeAttachmentUrl (src) {
      return `${this.$ComposeAPI.baseURL}${src}`
    },

    layoutTitleState (title) {
      return title ? null : false
    },

    layoutHandleState (layoutHandle) {
      return handle.handleState(layoutHandle)
    },

    checkUnsavedComposePage (next) {
      if (!this.page.deletedAt) {
        const layoutsStateChange = this.layouts.some((layout) => layout.meta.updated)
        const pageStateChange = !isEqual(this.page, this.initialPageState)

        return next((layoutsStateChange || pageStateChange) ? window.confirm(this.$t('unsavedChanges')) : true)
      }

      next()
    },

    setDefaultValues () {
      this.processing = false
      this.processingSaveAndClose = false
      this.processingSave = false
      this.processingClone = false
      this.processingDelete = false
      this.page = {}
      this.initialPageState = {}
      this.showIconModal = false
      this.attachments = []
      this.selectedAttachmentID = ''
      this.linkUrl = ''
      this.layouts = []
      this.layoutEditor = {}
      this.resolvedRoles = {}
      this.removedLayouts.clear()
      this.checkboxLabel = {}
    },
  },
}
</script>

<style lang="scss" scoped>
.selected-icon {
  outline: 2px solid var(--success);
}

.list-background {
  background-color: var(--body-bg);
}

.layout-actions {
  tr td {
    padding-bottom: 0.75rem;
  }

  tr:not(:first-child) td {
    padding-top: 0.75rem;
  }
}
</style>
