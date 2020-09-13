import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ResponsiveState } from 'ngx-responsive';

import { Status } from '../../shared/model/status.model';
import { SiteFilter } from '../site-filter/site-filter.model';
import { Subscription, Observable } from 'rxjs';
import { Select, Store } from '@ngxs/store';
import { StatusState, LoadData, ICounts } from '../../status.state';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {

  public filterArgs = new SiteFilter();
  public loading = false;
  @Select(StatusState) public status$: Observable<any>;
  @Select(StatusState.getCounts) public counts$: Observable<ICounts>;
  @Select(StatusState.isLoading) public isLoading$: Observable<boolean>;
  private subscriptions: Subscription[] = [];
  public sideNavMode: string;
  public sideNavOpened: boolean;

  constructor(
    private store: Store,
    private titleService: Title,
    private responsiveState: ResponsiveState
  ) { }

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

  ngOnDestroy() {
    for (const sub of this.subscriptions) {
      sub.unsubscribe();
    }
  }

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

  clearSearch() {
    this.filterArgs.data = '';
  }

  filterToDown() {
    this.filterArgs.data = 'status:down';
  }
}
