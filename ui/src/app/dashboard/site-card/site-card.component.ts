import { Component, OnInit, Input } from '@angular/core';
import { StatusSite } from '../../shared/models/statusSite.model';

/**
 * View/Edit Site information
 *
 * @export
 * @class SiteCardComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-site-card',
  templateUrl: './site-card.component.html',
  styleUrls: ['./site-card.component.scss']
})
export class SiteCardComponent implements OnInit {
  /**
   * StatusSite to view/edit
   *
   * @type {StatusSite}
   * @memberof SiteCardComponent
   */
  @Input() site: StatusSite = new StatusSite();

  /**
   * Creates an instance of SiteCardComponent.
   * @param {StatusService} statusService
   * @memberof SiteCardComponent
   */
  constructor() { }

  /**
   * Angular life-cycle method
   *
   * @memberof SiteCardComponent
   */
  ngOnInit() { }

  /**
   * Navigates to the provided uri in a new window
   *
   * @param {string} uri
   * @memberof SiteCardComponent
   */
  navigate(uri: string) {
    window.open(uri);
  }
}
