interface DropdownItem {
  label: string
  url: string
  delimiter: boolean
  target: string
}
interface Dropdown {
  items: DropdownItem[]
}
interface ItemOptions {
  label: string
  url: string
  target: string
  delimiter: boolean
  pageID: string
  displaySubPages: boolean
  dropdown: Dropdown
  align: string
}
interface NavigationItemOptions {
  enabled: boolean
  textColor: string
  backgroundColor: string
  item: ItemOptions
}
export type NavigationItemInput = Partial<NavigationItem> | NavigationItem
export default class NavigationItem {
  type: string
  options: NavigationItemOptions
  constructor(i?: NavigationItemInput)
  apply(i?: NavigationItemInput): void
}
export {}
