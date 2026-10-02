// SMTP host and addresses with stray whitespace fail when sending, so the
// settings form trims them before saving or testing. The password is kept
// as typed since whitespace can be part of it.
const trim = v => (typeof v === 'string' ? v.trim() : v)

export function trimServer(server) {
  return {
    ...server,
    host: trim(server.host),
    user: trim(server.user),
    from: trim(server.from),
    tlsServerName: trim(server.tlsServerName),
  }
}
