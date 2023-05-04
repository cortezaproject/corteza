export function checkValidURL (url: string): string {
  if (url.indexOf('://') > 0) { // do not modify the link if it has sth before "://"", e.g. ftp, http, etc.
    return url
  } else {
    return 'http://' + url
  }
}
