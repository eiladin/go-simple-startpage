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

  constructor(
    network: string = '',
    links: Link[] = [],
    sites: StatusSite[] = []
  ) {
    super(network, links, []);
    if (sites !== null) {
      sites.forEach((site) => {
        this.sites.push(Object.assign(new StatusSite(), site));
      });
    }
  }

}
