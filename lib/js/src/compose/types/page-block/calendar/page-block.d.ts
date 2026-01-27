import { PageBlock } from '../base'
import Feed, { FeedInput } from './feed'
import { ReminderFeed } from './feed-reminder'
import { RecordFeed } from './feed-record'
interface Header {
    left: string;
    center: string;
    right: string;
}
interface CalendarOptionsHeader {
    hide: boolean;
    views: string[];
    hidePrevNext: boolean;
    hideToday: boolean;
    hideTitle: boolean;
}
interface Options {
    defaultView: string;
    feeds: Array<Feed>;
    header: Partial<CalendarOptionsHeader>;
    locale: string;
    refreshRate: number;
    showRefresh: boolean;
    magnifyOption: string;
    eventDisplayOption: string;
}
/**
 * Helper class to help define calendar's functionality
 */
export declare class PageBlockCalendar extends PageBlock {
  readonly kind = 'Calendar'
  options: Options
  static feedResources: Readonly<{
        record: 'compose:record';
        reminder: 'system:reminder';
    }>
  constructor(i?: PageBlock | Partial<PageBlock>);
  applyOptions(o?: Partial<Options>): void;
  /**
     * Generates a header object of fullcalendar
     * @returns {Object}
     */
  getHeader(): Header | undefined;
  /**
     * Provides a list of available views.
     * @note When adding new ones, make sure included plugins support it.
     * @returns {Array}
     */
  static availableViews(): Array<string>;
  /**
     * Reorder views according to available views array order.
     * @param {Array} views Array of views to filter & sort
     */
  reorderViews(views?: string[]): Array<string>;
  /**
     * Converts old < V4 view names to >= V4 view names.
     * @note It wil preserve fields that don't need to/can't be converted
     * @param {string} views converted view name
     */
  static handleLegacyView(views?: string): string;
  /**
     * Converts old < V4 view names to >= V4 view names.
     * @note It wil preserve fields that don't need to/can't be converted
     * @param {string[]} views converted view names
     */
  static handleLegacyViews(views: string[]): string[];
  static makeFeed(f?: FeedInput): Feed;
  static ReminderFeed: typeof ReminderFeed
  static RecordFeed: typeof RecordFeed
}
export {}
