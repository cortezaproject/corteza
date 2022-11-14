[
  'CortezaAPI',
].forEach((cfg) => {
  if (window[cfg] === undefined) {
    throw new Error(`Missing or invalid configuration. 
          Make sure there is a public/config.js configuration file with window.${cfg} entry.`)
  }
})
