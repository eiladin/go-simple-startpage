import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ResponsiveState } from 'ngx-responsive';

import { Status } from '../../shared/models/status.model';
import { SiteFilter } from '../site-filter/site-filter.model';
import { Subscription, Observable } from 'rxjs';
import { Select, Store } from '@ngxs/store';
import { StatusState, LoadData, ICounts } from '../../status.state';

/**
 * Dashboard component showing the current status of sites
 */
@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  /**
   * Filter arguments for use with {@link SiteFilterPipe}
   */
  public filterArgs = new SiteFilter();

  /**
   * Internal variable used to track if the site status is loading
   */
  public loading = false;

  /**
   * Observable status from {@link StatusState}
   */
  @Select(StatusState) public status$: Observable<any>;

  /**
   * Observable site counts from {@link StatusService}
   */
  @Select(StatusState.getCounts) public counts$: Observable<ICounts>;

  @Select(StatusState.isLoading) public isLoading$: Observable<boolean>;

  /**
   * Internal variable used to track the subscriptions
   */
  private subscriptions: Subscription[] = [];

  /**
   * Set the sidenav mode property (side/over)
   */
  public sideNavMode: string;

  /**
   * Sets the sidenav initial state (open/close)
   */
  public sideNavOpened: boolean;

  /**
   * Create an instance of {@link DashboardComponent}
   * @param {StatusService} statusService
   * @param {Title} titleService
   * @param {ResponsiveState} responsiveState
   */
  constructor(
    private store: Store,
    private titleService: Title,
    private responsiveState: ResponsiveState
  ) { }

  /**
   * ngOnInit
   */
  ngOnInit() {
    this.store.dispatch(new LoadData());

    this.subscriptions.push(this.status$.subscribe( (state: Status) =>
      this.titleService.setTitle(state.network)
    ));

    this.subscriptions.push(this.responsiveState.device$.subscribe(state => {
      if (state === 'desktop') {
        this.sideNavMode = 'side';
        this.sideNavOpened = true;
      } else {
        this.sideNavMode = 'over';
        this.sideNavOpened = false;
      }
    }));
  }

  /**
   * Angular ngOnDestroy life cycle method
   */
  ngOnDestroy() {
    for (const sub of this.subscriptions) {
      sub.unsubscribe();
    }
  }

  /**
   * Swap themes between light <-> dark
   */
  changeTheme() {
    const body = document.body;
    const meta = document.querySelector('meta[name=theme-color]');
    if (!body || !meta) { return; }
    if (body.classList.contains('dark-theme')) {
      body.classList.remove('dark-theme');
      meta.setAttribute('content', '#3f51b5');
    } else {
      body.classList.add('dark-theme');
      meta.setAttribute('content', '#607d8b');
    }
  }

  /**
   * Clear search field
   */
  clearSearch() {
    this.filterArgs.data = '';
  }

  filterToDown() {
    this.filterArgs.data = 'status:down';
  }
}
