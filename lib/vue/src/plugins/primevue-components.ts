import type { App, Plugin } from 'vue'
import CInputSwitch from '../components/input/CInputSwitch.vue'
import CInputRole from '../components/input/CInputRole.vue'
import CFieldPicker from '../components/input/CFieldPicker.vue'
import CPermissionsButton from '../components/permissions/CPermissionsButton.vue'
import CPermissionsDialog from '../components/permissions/CPermissionsDialog.vue'
import AutoComplete from 'primevue/autocomplete'
import Accordion from 'primevue/accordion'
import AccordionContent from 'primevue/accordioncontent'
import AccordionHeader from 'primevue/accordionheader'
import AccordionPanel from 'primevue/accordionpanel'
import Avatar from 'primevue/avatar'
import Badge from 'primevue/badge'
import BlockUI from 'primevue/blockui'
import Breadcrumb from 'primevue/breadcrumb'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Card from 'primevue/card'
import Checkbox from 'primevue/checkbox'
import Chip from 'primevue/chip'
import ColorPicker from 'primevue/colorpicker'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import DatePicker from 'primevue/datepicker'
import Divider from 'primevue/divider'
import Drawer from 'primevue/drawer'
import FloatLabel from 'primevue/floatlabel'
import Form from '@primevue/forms/form'
import FormField from '@primevue/forms/formfield'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import Menu from 'primevue/menu'
import Message from 'primevue/message'
import MultiSelect from 'primevue/multiselect'
import Panel from 'primevue/panel'
import Password from 'primevue/password'
import Popover from 'primevue/popover'
import Paginator from 'primevue/paginator'
import RadioButton from 'primevue/radiobutton'

import ProgressBar from 'primevue/progressbar'
import ProgressSpinner from 'primevue/progressspinner'
import SelectButton from 'primevue/selectbutton'
import Select from 'primevue/select'
import Slider from 'primevue/slider'
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
import FileUpload from 'primevue/fileupload'
import Tree from 'primevue/tree'
import Tooltip from 'primevue/tooltip'

/**
 * Globally registers commonly used PrimeVue components.
 * Install once per app — no per-file imports needed for these components.
 */
export const PrimeVueComponentsPlugin: Plugin = {
  install(app: App) {
    app.component('AutoComplete', AutoComplete)
    app.component('Accordion', Accordion)
    app.component('AccordionContent', AccordionContent)
    app.component('AccordionHeader', AccordionHeader)
    app.component('AccordionPanel', AccordionPanel)
    app.component('Avatar', Avatar)
    app.component('Badge', Badge)
    app.component('Breadcrumb', Breadcrumb)
    app.component('Button', Button)
    app.component('ButtonGroup', ButtonGroup)
    app.component('Card', Card)
    app.component('Checkbox', Checkbox)
    app.component('Chip', Chip)
    app.component('ColorPicker', ColorPicker)
    app.component('Column', Column)
    app.component('ConfirmDialog', ConfirmDialog)
    app.component('DataTable', DataTable)
    app.component('DatePicker', DatePicker)
    app.component('Dialog', Dialog)
    app.component('Divider', Divider)
    app.component('Drawer', Drawer)
    app.component('BlockUI', BlockUI)
    app.component('FloatLabel', FloatLabel)
    app.component('Form', Form)
    app.component('FormField', FormField)
    app.component('IconField', IconField)
    app.component('InputIcon', InputIcon)
    app.component('InputNumber', InputNumber)
    app.component('InputText', InputText)
    app.component('InputGroup', InputGroup)
    app.component('InputGroupAddon', InputGroupAddon)
    app.component('Menu', Menu)
    app.component('Message', Message)
    app.component('MultiSelect', MultiSelect)
    app.component('Paginator', Paginator)
    app.component('Panel', Panel)
    app.component('Password', Password)

    app.component('Popover', Popover)
    app.component('ProgressBar', ProgressBar)
    app.component('ProgressSpinner', ProgressSpinner)
    app.component('RadioButton', RadioButton)
    app.component('Select', Select)
    app.component('SelectButton', SelectButton)
    app.component('Slider', Slider)
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
    app.component('FileUpload', FileUpload)
    app.component('Tree', Tree)

    // Directives
    app.directive('tooltip', Tooltip)

    // Human shared components
    app.component('CInputSwitch', CInputSwitch)
    app.component('CInputRole', CInputRole)
    app.component('CFieldPicker', CFieldPicker)
    app.component('CPermissionsButton', CPermissionsButton)
    app.component('CPermissionsDialog', CPermissionsDialog)
  },
}
