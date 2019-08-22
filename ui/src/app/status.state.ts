import { State, Action, StateContext, Selector } from '@ngxs/store';
import { Status } from './shared/models/status.model';
import { HttpClient } from '@angular/common/http';
import { tap, catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { StatusSite } from './shared/models/statusSite.model';

export class LoadData {
  static type = 'LoadData';
}

export class UpdateStatus {
  static type = 'UpdateStatus';
}

/**
 * Interface for tracking the Count of sites
 */
export interface ICounts {
  /**
   * Total number of sites
   */
  total: number;
  /**
   * Number of sites where isUp is truthy
   */
  up: number;
  /**
   * Number of sites where isUp is not truthy
   */
  down: number;
}

@State<Status>({
  name: 'status',
  defaults: new Status()
})
export class StatusState {

  @Selector()
  static getCounts(state: Status): ICounts {
      const total = state.sites.length;
      const down = state.sites.reduce((count: number, site: StatusSite) => site.isStatusLoaded && !site.isUp ? count + 1 : count, 0);
      const up = total - down;
      return { total, up, down };
  }

  @Selector()
  static isLoading(state: Status): boolean {
    return state.sites.filter(item => item.isStatusLoaded === true).length !== state.sites.length;
  }

  constructor(private http: HttpClient) { }

  @Action(LoadData)
  loadData(ctx: StateContext<Status>, action: LoadData) {
    return this.http.get<Status>('api/status/').pipe(
      tap(newStatus => {
        const state = ctx.getState();
        ctx.setState(Object.assign(new Status(newStatus.network, newStatus.links, newStatus.sites)));
        ctx.dispatch(new UpdateStatus());
      }),
      catchError(error => {
        console.error(error);
        return of(error);
      })
    );
  }

  @Action(UpdateStatus)
  updateStatus({ getState, patchState }: StateContext<Status>, action: UpdateStatus) {
    getState().sites
      .map(site =>
        this.http.post<StatusSite>('api/status', site)
          .pipe(
            tap(resultSite => {
              patchState({
                sites: getState().sites.map(origSite => {
                  if (origSite.friendlyName === resultSite.friendlyName) {
                    origSite = Object.assign(new StatusSite(), resultSite);
                    origSite.isStatusLoaded = true;
                  }
                  return origSite;
                })
              });
            })
          ).subscribe()
      );
  }

}
