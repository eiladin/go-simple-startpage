import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';

import { ResponsiveModule } from 'ngx-responsive';

import { SharedModule } from '../shared/shared.module';
import { DashboardComponent } from './dashboard/dashboard.component';
import { DashboardRoutingModule } from './dashboard.routes';
import { SiteCardComponent } from './site-card/site-card.component';
import { SiteCardXsComponent } from './site-card-xs/site-card-xs.component';
import { SiteFilterPipe } from './site-filter/site-filter.pipe';

@NgModule({
  imports: [
    DashboardRoutingModule,
    CommonModule,
    SharedModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSidenavModule,
    MatToolbarModule,
    ResponsiveModule.forRoot()
  ],
  declarations: [DashboardComponent, SiteCardComponent, SiteCardXsComponent, SiteFilterPipe],
  providers: []
})
export class DashboardModule { }
