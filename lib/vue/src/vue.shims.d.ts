
declare module '*.vue' {
  import Vue from 'vue'
  export default Vue
}

declare module 'vue/types/vue' {
  import { apiClients, EventBus } from '@cortezaproject/corteza-js'
  import { plugins } from './plugins'
  interface Vue {
    $SystemAPI: typeof apiClients.System;
    $ComposeAPI: typeof apiClients.Compose;
    $UIHooks: typeof plugins.UiHooks;
    $EventBus: typeof EventBus;
  }
}
