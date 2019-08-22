import { Component } from '@angular/core';

/**
 * Footer component
 */
@Component({
    selector: 'app-footer',
    styleUrls: ['./footer.component.scss'],
    templateUrl: './footer.component.html'
})
export class FooterComponent {
    /**
     * Create an instance of {@link FooterComponent}
     * @param {AppConfigService} appConfig Application configuration service
     */
    version: number
    public constructor() { 
        this.version = 1.0
    }
}
