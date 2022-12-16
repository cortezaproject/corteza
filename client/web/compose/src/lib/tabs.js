import { NoID } from '@cortezaproject/corteza-js'
/**
  * If block has no ID, it is a new block and we need to use tempID
  * to find it in the list of blocks.
*/
export function fetchID (block) {
  return block.blockID === NoID ? block.meta.tempID : block.blockID
}
