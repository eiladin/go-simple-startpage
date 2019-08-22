import { Routes, RouterModule } from '@angular/router';
import { NgModule } from '@angular/core';

import { ConfigComponent } from './config/config.component';

const CONFIG_ROUTES: Routes = [
    { path: 'config', component: ConfigComponent }
];

@NgModule({
  imports: [RouterModule.forChild(CONFIG_ROUTES)],
  exports: [RouterModule]
})
export class ConfigRoutingModule { }
