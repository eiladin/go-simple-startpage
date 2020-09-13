import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { IAppConfig } from '../shared/models/appConfig.model';

@Injectable()
export class AppConfigService implements IAppConfig {
  private _appConfig: IAppConfig;

  public get version() { return this._appConfig.version || ''; }

  constructor(private http: HttpClient) { }

  public load(): Promise<any> {
      return this.http.get(environment.gateway + '/api/appconfig')
          .toPromise()
          .then((res: IAppConfig) => this._appConfig = res)
          .catch(err => console.log(err));
  }
}
