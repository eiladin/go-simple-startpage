import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Observable, of } from 'rxjs';
import { mergeMap } from 'rxjs/operators';

import { Config } from '../../shared/models/config.model';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';


/**
 * Service used to load, save, import and export service and link configuration data.
 */
@Injectable()
export class ConfigService {
    /**
     * Config API url fragment
     */
    private configApi = environment.gateway + '/api/config';

    /**
     * Create an instance of {@link ConfigService}
     * @param {HttpClient} httpClient
     * @param {MatSnackBar} snackBar
     */
    constructor(
        private http: HttpClient,
        private snackBar: MatSnackBar) {
    }

    /**
     * Get the configuration object that describes what will be monitored
     */
    public get() {
        return this.http.get<Config>(this.configApi);
    }

    /**
     * Save the configuration to the server
     * @param {Config} config
     */
    public save(config: Config): Observable<Config> {
        const alphaSort = config.sites.sort((a, b) => {
            if (a.friendlyName > b.friendlyName) { return 1; }
            if (b.friendlyName < a.friendlyName) { return -1; }
            return 0;
        });

        config.sites.forEach((curr, idx, arr) => {
            curr.sortOrder = idx;
        });
        config.links.forEach((curr, idx, arr) => {
            curr.sortOrder = idx;
        });
        return this.http.post<Config>(this.configApi, config)
            .pipe(
                mergeMap((val, idx) => {
                    this.snackBar.open('Configuration Saved', undefined, { duration: 2000 });
                    return of(val);
                })
            );
    }

    /**
     * Export configuration to Json file
     * @param {Config} config
     */
    public exportJson(config: Config) {
        const filename = 'config.json';
        const data = JSON.stringify(config);
        const blob = new Blob([data], { type: 'text/json' });
        if (window.navigator && window.navigator.msSaveOrOpenBlob) {
            window.navigator.msSaveOrOpenBlob(blob, filename);
        } else {
            const anchor = document.createElement('a');
            anchor.setAttribute('href', window.URL.createObjectURL(blob));
            anchor.setAttribute('download', filename);
            document.body.appendChild(anchor);
            anchor.click();
            document.body.removeChild(anchor);
        }
        this.snackBar.open('Configuration Exported', undefined, { duration: 2000 });
    }

    /**
     * Import configuration from Json file
     */
    public importJson(): Promise<Config> {
        return new Promise<Config>((resolve, reject) => {
            const fileInput = document.createElement('input');
            fileInput.setAttribute('type', 'file');

            function doImport(ev: MouseEvent) {
                const file = this.files[0];
                const reader = new FileReader();
                reader.onloadend = (evt) => resolve(<Config>JSON.parse((<any>evt.target).result));
                reader.readAsText(file);
            }

            fileInput.addEventListener('change', doImport, false);
            fileInput.click();
        });
    }
}
