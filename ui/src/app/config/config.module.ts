import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatToolbarModule } from '@angular/material/toolbar';
import { ResponsiveModule } from 'ngx-responsive';
import { NgxPageScrollCoreModule } from 'ngx-page-scroll-core';
import { NgxPageScrollModule } from 'ngx-page-scroll';

import { SharedModule } from '../shared/shared.module';
import { ConfigComponent } from './config/config.component';
import { ConfigService } from './services/config.service';
import { ConfigRoutingModule } from './config.routes';
import { ConfigSiteComponent } from './config-site/config-site.component';
import { ConfigLinkComponent } from './config-link/config-link.component';

@NgModule({
  imports: [
    ConfigRoutingModule,
    SharedModule,
    MatButtonModule,
    MatCardModule,
    MatChipsModule,
    MatExpansionModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatSlideToggleModule,
    MatSnackBarModule,
    MatToolbarModule,
    ResponsiveModule.forRoot(),
    NgxPageScrollCoreModule,
    NgxPageScrollModule,
  ],
  declarations: [ConfigComponent, ConfigSiteComponent, ConfigLinkComponent],
  providers: [ConfigService]
})
export class ConfigModule { }
