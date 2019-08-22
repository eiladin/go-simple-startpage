import { Link } from './link.model';
import { Site } from './site.model';

/**
 * Config
 */
export class Config {
  /**
   * Create an instance of {@link Config}
   * @param {string} network Network name/description
   * @param {Array<Link>} links Array of Link
   * @param {Array<Site>} sites Array of Site
   */
  constructor(
    network: string = '',
    links: Link[] = [],
    sites: Site[] = []
  ) {
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

  /**
  * Internal method to sort Sites by sortOrder
  */
  private sortSites() {
    this.sites.sort((a, b) => {
      if (a.sortOrder < b.sortOrder) { return -1; }
      if (a.sortOrder > b.sortOrder) { return 1; }
      return 0;
    });
  }

  /**
   * Internal method to sort Links and Sites
   */
  public sortChildren() {
    this.sortLinks();
    this.sortSites();
  }

}
