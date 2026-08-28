module.exports = {
  create: () => ({
    request: () => Promise.reject(new Error('axios is disabled in C311 component tests')),
  }),
  CancelToken: { source: () => ({ token: {}, cancel: () => {} }) },
}
