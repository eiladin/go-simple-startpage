import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Link } from '../../shared/model/link.model';

/**
 * Component for adding/editing links
 */
@Component({
  selector: 'app-config-link',
  templateUrl: './config-link.component.html',
  styleUrls: ['./config-link.component.scss']
})
export class ConfigLinkComponent implements OnInit {

  /**
   * Parameter for {@link Link}
   */
  @Input() link: Link;

  /**
   * Creates an instance of {@link ConfigLinkComponent}
   */
  constructor() { }

  /**
   * Angular ngOnInit life-cycle method
   */
  ngOnInit() {
  }

}
