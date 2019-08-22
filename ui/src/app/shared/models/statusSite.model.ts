import { Site } from './site.model';
/**
 * StatusSite extends the Site class and adds additional properties
 */
export class StatusSite extends Site {

    /**
     * Boolean that returns the current status (up = true/down = false)
     */
    public isUp: boolean;

    /**
     * IP address of the site determined when checking the site status on the server
     */
    public ip: string;

    /**
     * Boolean that returns if the site has been checked for availability
     */
    public isStatusLoaded = false;
}
