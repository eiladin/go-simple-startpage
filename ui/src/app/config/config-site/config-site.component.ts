import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Site } from '../../shared/model/site.model';
import { MatChipInputEvent } from '@angular/material/chips';
import { Tag } from '../../shared/model/tag.model';

/**
 * Component for adding/editing sites
 */
@Component({
  selector: 'app-config-site',
  templateUrl: './config-site.component.html',
  styleUrls: ['./config-site.component.scss']
})
export class ConfigSiteComponent implements OnInit {
  /**
   * Site list passed to the component
   * @example
   * &lt;app-config-site [sites]="sites"&gt;&lt;/app-config-site&gt;
   */
  @Input() site: Site;

  /**
   * Event that is emitted when the user clicks the delete button
   */
  @Output() delete: EventEmitter<any> = new EventEmitter();

  /**
   * Creates an instance of {@link ConfigSiteComponent}
   */
  constructor() { }

  /**
   * Angular ngOnInit function
   */
  ngOnInit() { }

  /**
   * Method to delete the site, called when the user clicks the delete button
   */
  public deleteSite() {
    this.delete.emit(this.site);
  }

  /**
   * Method to remove a tag from a site
   * @param {Tag} tag
   */
  public removeTag(tag: Tag) {
    const index = this.site.tags.indexOf(tag);
    if (index >= 0) {
      this.site.tags.splice(index, 1);
    }
  }

  /**
   * Method to add a tag to a site
   * @param {MatChipInputEvent} event
   */
  public addTag(event: MatChipInputEvent) {
    const input = event.input;
    const value = event.value;
    if ((value || '').trim()) {
      this.site.tags.push({ value: value.trim() });
    }
    if (input) {
      input.value = '';
    }
  }

  /**
   * Method to set isSupportedApp when the icon value is changed, if the icon has a dot (.) in the name
   * @param {string} value
   */
  public iconChanged(value: string) {
    this.site.isSupportedApp = value.indexOf('.') > -1;
  }


}
