import { Site } from './site.model';

export class StatusSite extends Site {
    public isUp: boolean;
    public ip: string;
    public isStatusLoaded = false;
}
