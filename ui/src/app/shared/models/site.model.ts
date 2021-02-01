export class Site {

    constructor(
        public name: string = '',
        public uri: string = '',
        public icon: string = '',
        public tags: string[] = [],
        public isSupportedApp: boolean = false
        ) { }
        
        mdiicon() {
            return "mdi-"+this.icon
        }
    }
