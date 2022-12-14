import { Apply, NoID } from '../../../../cast'
import { IsOf } from '../../../../guards'

type Text = {
    text: string;
};

type URL = {
    text: string;
    url: string;
    newWindow: boolean;
}

type Dropdown = {
    text: string;
    url: string;
    delimiter: boolean;
}

type Compose = {
    text: string;
    referenceId: string;
}

interface ListOptions {
    disabled: boolean;
    textColor: string;
    bgVariant: string;
    itemOption: Text | URL | Compose | {};
    dropdown: Dropdown[];
    dropdownType: string;
    dropdownText: string;
}

export type ListInput = Partial<List> | List

const defOptions = {
  disabled: false,
  textColor: '#0B344E',
  bgVariant: '',
  itemOption: {},
  dropdown: [],
  dropdownType: '',
  dropdownText: '',
}

export default class List {
  public navigationType = ''

  public options: ListOptions = { ...defOptions }

  constructor (i?: ListInput) {
    this.apply(i)
  }

  apply (i?: ListInput): void {
    if (!i) return

    Apply(this, i, String, 'navigationType')

    if (i.options) {
      this.options = { ...this.options, ...i.options }
    }
  }
}
