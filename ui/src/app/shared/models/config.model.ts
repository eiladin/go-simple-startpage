import { Link } from './link.model';
import { Site } from './site.model';

export class Config {
  constructor(
    network: string = '',
    links: Link[] = [],
    sites: Site[] = []
  ) {
    if (links == null) {
      links = [];
    }
    if (sites == null) {
      sites = [];
    }
    sites.forEach((site) => {
      this.sites.push(Object.assign(new Site(), site));
    });

    links.forEach((link) => {
      this.links.push(Object.assign(new Link(), link));
    });

    this.network = network;
  }

  public network = '';
  public sites: Site[] = [];
  public links: Link[] = [];

  /**
   * Internal method to sort Links by sortOrder
   */
  private sortLinks() {
    this.links.sort((a, b) => {
      if (a.sortOrder < b.sortOrder) { return -1; }
      if (a.sortOrder > b.sortOrder) { return 1; }
      return 0;
    });
  }

  private sortSites() {
    this.sites.sort((a, b) => {
      if (a.sortOrder < b.sortOrder) { return -1; }
      if (a.sortOrder > b.sortOrder) { return 1; }
      return 0;
    });
  }

  public sortChildren() {
    this.sortLinks();
    this.sortSites();
  }

}
