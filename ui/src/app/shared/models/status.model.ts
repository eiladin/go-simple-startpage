import { Config } from './config.model';
import { StatusSite } from './statusSite.model';
import { Link } from './link.model';

/**
 * Status extends the Config class and replaces the `sites` property with StatusSite
 */
export class Status extends Config {
  /**
   * Replace the implementation of `sites` on the base class
   */
  public sites: StatusSite[];

    /**
   * Create an instance of {@link Config}
   * @param {string} network Network name/description
   * @param {Array<Link>} links Array of Link
   * @param {Array<Site>} sites Array of Site
   */
  constructor(
    network: string = '',
    links: Link[] = [],
    sites: StatusSite[] = []
  ) {
    super(network, links, []);
    sites.forEach((site) => {
      this.sites.push(Object.assign(new StatusSite(), site));
    });

  }

}
