import { Component } from '@angular/core';
import { AppConfigService } from 'src/app/core/app-config.service';

@Component({
    selector: 'app-footer',
    styleUrls: ['./footer.component.scss'],
    templateUrl: './footer.component.html'
})
export class FooterComponent {
    public constructor(public appConfig: AppConfigService) { }
}
