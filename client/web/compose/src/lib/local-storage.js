export function getItem (key = '') {
  return JSON.parse(window.localStorage.getItem(key) || '[]')
}

export function setItem (key = '', value = []) {
  return window.localStorage.setItem(key, JSON.stringify(value))
}

export function removeItem (key = '') {
  return window.localStorage.removeItem(key)
}
