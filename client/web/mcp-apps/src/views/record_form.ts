import Checkbox from 'primevue/checkbox'
import DatePicker from 'primevue/datepicker'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import MultiSelect from 'primevue/multiselect'
import RadioButton from 'primevue/radiobutton'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import ToggleSwitch from 'primevue/toggleswitch'
import { mount } from '../shared/mount'
import RecordForm from './RecordForm.vue'

// What the webapp's field editors expect to find registered.
mount(RecordForm, 'human-record-form', {
  Checkbox,
  DatePicker,
  InputNumber,
  InputText,
  MultiSelect,
  RadioButton,
  Select,
  Tag,
  ToggleSwitch,
})
