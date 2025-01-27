import Vue from 'vue'
import { capitalize, uniq } from 'lodash'

import AutomationBase from './AutomationBase'
import AutomationConfigurator from './AutomationConfigurator'
import CalendarBase from './CalendarBase'
import CalendarConfigurator from './CalendarConfigurator'
import ChartBase from './ChartBase'
import ChartConfigurator from './ChartConfigurator'
import ContentBase from './ContentBase'
import ContentConfigurator from './ContentConfigurator'
import FileBase from './FileBase'
import FileConfigurator from './FileConfigurator'
import IFrameBase from './IFrameBase'
import IFrameConfigurator from './IFrameConfigurator'
import RecordBase from './RecordBase'
import RecordConfigurator from './RecordConfigurator'
import RecordEditor from './RecordEditor'
import RecordListBase from './RecordListBase'
import RecordListConfigurator from './RecordListConfigurator'
import RecordRevisionsBase from './RecordRevisionsBase'
import RecordRevisionsConfigurator from './RecordRevisionsConfigurator'
import RecordOrganizerBase from './RecordOrganizerBase'
import RecordOrganizerConfigurator from './RecordOrganizerConfigurator'
import SocialFeedBase from './SocialFeedBase'
import SocialFeedConfigurator from './SocialFeedConfigurator'
import MetricBase from './MetricBase'
import MetricConfigurator from './MetricConfigurator'
import CommentBase from './CommentBase'
import CommentConfigurator from './CommentConfigurator'
import ReportBase from './Report/Base'
import ReportConfigurator from './Report/Configurator'
import ProgressBase from './ProgressBase'
import ProgressConfigurator from './ProgressConfigurator'
import GeometryBase from './GeometryBase'
import GeometryConfigurator from './GeometryConfigurator/index'
import NavigationConfigurator from './Navigation/Configurator'
import NavigationBase from './Navigation/Base'
import TabsBase from './TabsBase'
import TabsConfigurator from './TabsConfigurator'

/**
 * List of all known page block components
 *
 */
const Registry = {
  AutomationBase,
  AutomationConfigurator,
  CalendarBase,
  CalendarConfigurator,
  ChartBase,
  ChartConfigurator,
  ContentBase,
  ContentConfigurator,
  FileBase,
  FileConfigurator,
  IFrameBase,
  IFrameConfigurator,
  RecordBase,
  RecordConfigurator,
  RecordEditor,
  RecordListBase,
  RecordListConfigurator,
  RecordRevisionsBase,
  RecordRevisionsConfigurator,
  RecordOrganizerBase,
  RecordOrganizerConfigurator,
  ReportBase,
  ReportConfigurator,
  SocialFeedBase,
  SocialFeedConfigurator,
  MetricBase,
  MetricConfigurator,
  CommentBase,
  CommentConfigurator,
  ProgressBase,
  ProgressConfigurator,
  GeometryBase,
  GeometryConfigurator,
  TabsBase,
  TabsConfigurator,
  NavigationConfigurator,
  NavigationBase,
}

const defaultMode = 'Base'

function GetComponent ({ block, mode = defaultMode }) {
  if (!block) {
    throw new Error('block prop missing')
  }

  const { kind } = block
  for (mode of uniq([capitalize(mode), defaultMode])) {
    if (mode === 'Editor' && block.options.referenceField && block.options.referenceModuleID) {
      mode = 'Base'
    }

    const cmpName = kind + mode
    if (Object.hasOwnProperty.call(Registry, cmpName)) {
      return Registry[cmpName]
    }
  }

  throw new Error('unknown block kind: ' + kind)
}

/**
 * Main entry point for PageBlock components
 *
 * It will look for combination of page block kind & mode (from props)
 * and render the component
 */
export default Vue.component('page-block', {
  functional: true,

  render (ce, ctx) {
    return ce(GetComponent(ctx.props), ctx.data, ctx.children)
  },
})
