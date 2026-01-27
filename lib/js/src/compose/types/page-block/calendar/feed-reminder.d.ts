import { System as SystemAPI } from '../../../../api-clients'
import { User } from '../../../../system'
import { Event } from './shared'
interface FeedOptions {
    color: string;
}
interface Feed {
    options: FeedOptions;
}
interface Range {
    end: Date;
    start: Date;
}
/**
 * Loads & converts reminder resource into FC events
 * @param {SystemAPI} $SystemAPI SystemAPI provider
 * @param {User} user Current user
 * @param {Feed} feed Current feed
 * @param {Object} range Current date range
 * @returns {Promise<Array>} Resolves to a set of FC events to display
 */
export declare function ReminderFeed($SystemAPI: SystemAPI, user: User, feed: Feed, range: Range, options?: {}): Promise<Event[]>;
export {}
