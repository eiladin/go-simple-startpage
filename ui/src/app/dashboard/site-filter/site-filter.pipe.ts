import { Pipe, PipeTransform } from '@angular/core';

import { SiteFilter } from './site-filter.model';
import { StatusSite } from '../../shared/models/statusSite.model';

/**
 * Pipe to filter sites by name, used by search
 */
@Pipe({
  name: 'siteFilter',
  pure: false
})
export class SiteFilterPipe implements PipeTransform {

  transform(items: StatusSite[], filter: SiteFilter): any {
    if (!items || !filter || !filter.data) {
      return items;
    }

    return items.filter(item => {
      const name = item.friendlyName || '';
      const ip = item.ip || '';
      const isUp = !!item.isUp;

      if (filter.data.toLowerCase() === 'status:down') {
        return isUp === false;
      }
      if (filter.data.toLowerCase() === 'status:up') {
        return isUp === true;
      }

      if ((name.toLowerCase().indexOf(filter.data.toLowerCase()) !== -1) ||
          (ip.indexOf(filter.data) !== -1)) {
        return true;
      }

      if (item.tags) {
        return (item.tags.filter(tag => tag.value.toLowerCase().indexOf(filter.data.toLowerCase()) !== -1)).length;
      }

      return 0;
    });
  }

}
