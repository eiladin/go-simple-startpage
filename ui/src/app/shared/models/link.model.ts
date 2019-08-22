/**
 * Link
 */
export class Link {
    /**
     * Create an instance of {@link Link}
     * @param {string} name Description of the link
     * @param {string} uri Uri
     * @param {number} sortOrder 0-based index for ordering of links
     */
    constructor(
        public name: string = '',
        public uri: string = '',
        public sortOrder: number = 0
    ) { }
}
