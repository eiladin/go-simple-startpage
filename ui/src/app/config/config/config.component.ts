import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';

import { Config } from '../../shared/models/config.model';
import { Link } from '../../shared/models/link.model';
import { Site } from '../../shared/models/site.model';
import { ConfigService } from '../services/config.service';
import { EasingLogic } from 'ngx-page-scroll-core';

/**
 * Configuration Component for adding sites and links
 */
@Component({
  selector: 'app-config',
  templateUrl: './config.component.html',
  styleUrls: ['./config.component.scss']
})
export class ConfigComponent implements OnInit {
  /**
   * Page title
   */
  public title = 'Configuration';

  /**
   * Configuration object
   */
  public config: Config = new Config();

  /**
   * list of expanded / collapsed panels
   */
  public panelOpenState = [];

  public myEasing: EasingLogic = (t: number, b: number, c: number, d: number): number => {
      // easeInOutExpo easing
      if (t === 0) { return b; }
      if (t === d) { return b + c; }
      if ((t /= d / 2) < 1) {
        return c / 2 * Math.pow(2, 10 * (t - 1)) + b;
      }
      return c / 2 * (-Math.pow(2, -10 * --t) + 2) + b;
  };


  /**
   * Create an instance of {@link ConfigComponent}
   * @param {ConfigService} configService
   * @param {Title} titleService
   */
  constructor(
    private configService: ConfigService,
    private titleService: Title
  ) { }

  /**
   * ngOnInit
   */
  ngOnInit() {
    this.titleService.setTitle(this.title);
    window.onscroll = () => this.hasScrolled();
    this.getConfig();
  }

  /**
   * Determine if page has been scrolled more than 100 px
   */
  public hasScrolled() {
    const el = document.getElementById('config-scroll-to-top');
    if (el) {
      if (document.body.scrollTop > 100 || document.documentElement.scrollTop > 100) {
        el.style.display = 'block';
      } else {
        el.style.display = 'none';
      }
    }
  }

  /**
   * Method to load the configuration object
   */
  private getConfig() {
    this.configService.get().subscribe((res: any) => {
      this.config = new Config(res.network, res.links, res.sites);
      this.config.sortChildren();
      this.panelOpenState = [];
      for (let idx = 0; idx < this.config.sites.length; idx++) {
        this.panelOpenState.push(false);
      }
    });
  }

  /**
   * Method to save the configuration object
   */
  public saveConfig() {
    this.configService.save(this.config).subscribe();
  }

  /**
   * Method to export the configuration object to a Json file
   */
  public exportConfig() {
    this.configService.exportJson(this.config);
  }

  /**
   * Method to import the configuration object from a Json file
   */
  public importConfig() {
    this.configService.importJson().then((conf) => {
      this.config = conf;
      this.configService.save(this.config).subscribe();
    });
  }

  /**
   * Called by the UI when a user clicks `add link`
   */
  public addLink() {
    const link = new Link();
    this.config.links.push(link);
  }

  /**
   * Called by the UI when a user clicks `add site`
   */
  public addSite() {
    const site = new Site();
    this.config.sites.push(site);
  }

  /**
   * Called by the UI when a user clicks `delete link`
   * @param {Link} link
   */
  public deleteLink(link) {
    const idx = this.config.links.indexOf(link);
    this.config.links.splice(idx, 1);
  }

  /**
   * Called by the UI when a user clicks `delete site`
   * @param {Site} site
   */
  public deleteSite(site: Site) {
    console.log(site);
    const idx = this.config.sites.indexOf(site);
    this.config.sites.splice(idx, 1);
  }
}
