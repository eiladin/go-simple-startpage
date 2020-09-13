import { Component, OnInit, Input } from '@angular/core';
import { StatusSite } from '../../shared/models/statusSite.model';

/**
 * Site Card Component used on the dashboard for phones
 */
@Component({
  selector: 'app-site-card-xs',
  templateUrl: './site-card-xs.component.html',
  styleUrls: ['./site-card-xs.component.scss']
})
export class SiteCardXsComponent implements OnInit {
  /**
   * Controls display of the loading indicator
   */
  public isLoading = true;

  /**
   * parameter for {@link Site}
   */
  @Input() site: StatusSite = new StatusSite();

  /**
   * Create an instance of {@link SiteCardXsComponent}
   */
  constructor() { }

  /**
   * Angular ngOnInit function
   */
  ngOnInit() { }

  /**
  * Open the provided uri in a new tab
  * @param {string} uri
  */
  navigate(uri: string) {
    window.open(uri);
  }

}
