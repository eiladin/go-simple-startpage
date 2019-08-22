import { Tag } from './tag.model';
/**
 * Site Class
 *
 * @export
 * @class Site
 */
export class Site {
    /**
     * Creates an instance of Site.
     * @param {string} [friendlyName=''] display name
     * @param {string} [uri=''] uri
     * @param {string} [icon=''] material icon
     * @param {number} [sortOrder=0] overall position
     * @param {Tag[]} [tags=[]] array of tags
     * @param {isSupportedApp} [isSupportedApp=false] is supported app
     * @memberof Site
     */
    constructor(
        public friendlyName: string = '',
        public uri: string = '',
        public icon: string = '',
        public sortOrder: number = 0,
        public tags: Tag[] = [],
        public isSupportedApp: boolean = false
    ) { }

    get imageName() {
        return `/supportedapps/${this.icon}`;
    }
}
