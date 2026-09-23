import { defineConfig } from 'vitepress'
import llmstxt, { copyOrDownloadAsMarkdownButtons } from 'vitepress-plugin-llms'

export default defineConfig({
  title: 'Human Docs',
  description: 'Build apps, automations and AI agents on Human.',
  cleanUrls: true,
  // Quickstart points at the reader's own local instance.
  ignoreDeadLinks: 'localhostLinks',
  lastUpdated: true,

  // Generated pages keep the .gen marker on disk and drop it from the URL.
  rewrites: {
    'reference/environment.gen.md': 'reference/environment.md',
    'reference/permissions/:name.gen.md': 'reference/permissions/:name.md',
  },

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/icon.svg' }],
    ['meta', { name: 'theme-color', content: '#09344E' }],
  ],

  markdown: {
    config(md) {
      md.use(copyOrDownloadAsMarkdownButtons)
    },
  },

  vite: {
    plugins: [
      llmstxt({ title: 'Human', description: 'Build apps, automations and AI agents on Human.' }),
    ],
  },

  themeConfig: {
    logo: { light: '/logo.svg', dark: '/logo-dark.svg', alt: 'Human' },
    siteTitle: false,

    nav: [
      { text: 'Get started', link: '/get-started/', activeMatch: '^/get-started/' },
      { text: 'Guides', link: '/guides/build-an-app', activeMatch: '^/guides/' },
      { text: 'Reference', link: '/reference/environment', activeMatch: '^/reference/' },
    ],

    sidebar: [
      {
        text: 'Get started',
        items: [
          { text: 'What is Human', link: '/get-started/' },
          { text: 'Quickstart', link: '/get-started/quickstart' },
          { text: 'Core concepts', link: '/get-started/concepts' },
        ],
      },
      {
        text: 'Guides',
        items: [
          { text: 'Build an app', link: '/guides/build-an-app' },
          { text: 'Agents', link: '/guides/agents' },
        ],
      },
      {
        text: 'Reference',
        items: [
          { text: 'Environment variables', link: '/reference/environment' },
          {
            text: 'Permissions',
            link: '/reference/permissions/',
            collapsed: true,
            items: [
              { text: 'System', link: '/reference/permissions/system' },
              { text: 'Compose', link: '/reference/permissions/compose' },
              { text: 'Automation', link: '/reference/permissions/automation' },
              { text: 'Federation', link: '/reference/permissions/federation' },
            ],
          },
        ],
      },
    ],

    search: { provider: 'local' },
    outline: { level: [2, 3], label: 'On this page' },
    docFooter: { prev: 'Previous', next: 'Next' },
    footer: { message: 'Released under the Apache 2.0 License.' },
  },
})
