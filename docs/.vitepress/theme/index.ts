import DefaultTheme from 'vitepress/theme'
import type { Theme } from 'vitepress'
import CopyOrDownloadAsMarkdownButtons from 'vitepress-plugin-llms/vitepress-components/CopyOrDownloadAsMarkdownButtons.vue'
import Screenshot from './components/Screenshot.vue'
import './fonts.css'
import './style.css'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('CopyOrDownloadAsMarkdownButtons', CopyOrDownloadAsMarkdownButtons)
    app.component('Screenshot', Screenshot)
  },
} satisfies Theme
