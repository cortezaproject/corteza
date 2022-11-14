import cmpInput from './C3Input.vue'
import cmpTextarea from './C3Textarea.vue'
import cmpCheckbox from './C3Checkbox.vue'
import cmpSelect from './C3Select.vue'
import { Component } from 'vue'
import _ from 'lodash'

interface Handler {
  value (props: object): unknown;
  update (props: object, value: unknown): void;
}

interface Props {
  [_: string]: unknown;
}

interface Control extends Handler {
  component: Component;
  props?: Props;
}

interface Specs {
  handler: string | Handler;
  props?: Props;
}

function makeHandler (prop: string): Handler {
  const path = prop.split('.')
  return {
    value: (props: Props): unknown => _.get(props, path),
    update: (props: Props, value: unknown): void => { _.set(props, path, value) },
  }
}

export function generic (component: Component, { handler, props }: Specs): Control {
  if (typeof handler === 'string') {
    handler = makeHandler(handler)
  }

  return {
    component,
    props,
    ...handler,
  }
}

export function input (label: string, handler: string | Handler): Control {
  return generic(cmpInput, { handler, props: { label } })
}

export function textarea (label: string, handler: string | Handler): Control {
  return generic(cmpTextarea, { handler, props: { label } })
}

export function checkbox (label: string, handler: string | Handler): Control {
  return generic(cmpCheckbox, { handler, props: { label } })
}

export function select (label: string, handler: string | Handler, options: object): Control {
  return generic(cmpSelect, {
    handler,
    props: {
      label,
      options,
    },
  })
}
