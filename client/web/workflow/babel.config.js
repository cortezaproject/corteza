module.exports = {
  presets: [
    '@vue/cli-plugin-babel/preset',
  ],
  env: {
    test: {
      plugins: [
        ['istanbul', { useInlineSourceMaps: false }],
      ],
    },
  },
}
