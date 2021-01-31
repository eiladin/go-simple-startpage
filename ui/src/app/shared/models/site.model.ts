import { Tag } from './tag.model';

export class Site {

    constructor(
        public name: string = '',
        public uri: string = '',
        public icon: string = '',
        public tags: Tag[] = [],
        public isSupportedApp: boolean = false
        ) { }
        
        mdiicon() {
            return "mdi-"+this.icon
        }
    }
