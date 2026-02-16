import type { App, Plugin } from 'vue'

import Accordion from 'primevue/accordion'
import AccordionContent from 'primevue/accordioncontent'
import AccordionHeader from 'primevue/accordionheader'
import AccordionPanel from 'primevue/accordionpanel'
import Avatar from 'primevue/avatar'
import Badge from 'primevue/badge'
import Breadcrumb from 'primevue/breadcrumb'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Checkbox from 'primevue/checkbox'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import DatePicker from 'primevue/datepicker'
import Divider from 'primevue/divider'
import FloatLabel from 'primevue/floatlabel'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import Paginator from 'primevue/paginator'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import Select from 'primevue/select'
import Tab from 'primevue/tab'
import TabList from 'primevue/tablist'
import TabPanel from 'primevue/tabpanel'
import TabPanels from 'primevue/tabpanels'
import Tabs from 'primevue/tabs'
import Tag from 'primevue/tag'
import Textarea from 'primevue/textarea'
import TieredMenu from 'primevue/tieredmenu'
import Toast from 'primevue/toast'
import ToggleSwitch from 'primevue/toggleswitch'
import Tree from 'primevue/tree'

/**
 * Globally registers commonly used PrimeVue components.
 * Install once per app — no per-file imports needed for these components.
 */
export const PrimeVueComponentsPlugin: Plugin = {
  install(app: App) {
    app.component('Accordion', Accordion)
    app.component('AccordionContent', AccordionContent)
    app.component('AccordionHeader', AccordionHeader)
    app.component('AccordionPanel', AccordionPanel)
    app.component('Avatar', Avatar)
    app.component('Badge', Badge)
    app.component('Breadcrumb', Breadcrumb)
    app.component('Button', Button)
    app.component('Card', Card)
    app.component('Checkbox', Checkbox)
    app.component('Column', Column)
    app.component('ConfirmDialog', ConfirmDialog)
    app.component('DataTable', DataTable)
    app.component('DatePicker', DatePicker)
    app.component('Dialog', Dialog)
    app.component('Divider', Divider)
    app.component('FloatLabel', FloatLabel)
    app.component('IconField', IconField)
    app.component('InputIcon', InputIcon)
    app.component('InputNumber', InputNumber)
    app.component('InputText', InputText)
    app.component('Menu', Menu)
    app.component('Message', Message)
    app.component('Paginator', Paginator)
    app.component('ProgressSpinner', ProgressSpinner)
    app.component('Select', Select)
    app.component('Tab', Tab)
    app.component('TabList', TabList)
    app.component('TabPanel', TabPanel)
    app.component('TabPanels', TabPanels)
    app.component('Tabs', Tabs)
    app.component('Tag', Tag)
    app.component('Textarea', Textarea)
    app.component('TieredMenu', TieredMenu)
    app.component('Toast', Toast)
    app.component('ToggleSwitch', ToggleSwitch)
    app.component('Tree', Tree)
  },
}
