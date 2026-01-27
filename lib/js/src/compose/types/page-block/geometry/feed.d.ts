interface FeedOptions {
    color: string;
    prefilter: string;
    moduleID: string;
}
export type FeedInput = Partial<Feed> | Feed;
/**
 * Feed class represents an event feed for the given calendar
 */
export default class Feed {
  resource: string
  titleField: string
  geometryField: string
  displayMarker: boolean
  displayPolygon: boolean
  options: FeedOptions
  constructor(i?: FeedInput);
  apply(i?: FeedInput): void;
  isValid(): boolean;
}
export {}
