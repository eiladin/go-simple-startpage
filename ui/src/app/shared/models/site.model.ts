import { Tag } from './tag.model';

export class Site {

    constructor(
        public id: number = 0,
        public friendlyName: string = '',
        public uri: string = '',
        public icon: string = '',
        public sortOrder: number = 0,
        public tags: Tag[] = [],
        public isSupportedApp: boolean = false
    ) { }
}
