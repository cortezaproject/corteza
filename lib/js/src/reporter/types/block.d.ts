import { DisplayElement } from './display-elements'
export declare class Block {
  blockID: string
  title: string
  description: string
  layout: string
  elements: Array<DisplayElement>
  xywh: number[]
  constructor(p: Partial<Block>)
}
