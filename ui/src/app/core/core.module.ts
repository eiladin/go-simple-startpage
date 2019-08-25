import { NgModule, ModuleWithProviders, APP_INITIALIZER, Optional, SkipSelf } from '@angular/core';
import { AppConfigService } from './app-config.service';
import { Title } from '@angular/platform-browser';

export function init(appConfigService: AppConfigService) {
  return () => {
      return appConfigService.load();
  };
}
@NgModule({ })
export class CoreModule { 
  public static forRoot(): ModuleWithProviders {
    return {
        ngModule: CoreModule,
        providers: [
            // Providers
            Title,
            AppConfigService,
            {
                'provide': APP_INITIALIZER,
                'useFactory': init,
                'deps': [AppConfigService],
                'multi': true
            }
        ]
    };
}
constructor( @Optional() @SkipSelf() parentModule: CoreModule) {
    if (parentModule) {
        throw new Error('CoreModule is already loaded. Import it in the AppModule only');
    }
}
}
