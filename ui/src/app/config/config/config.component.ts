import { Component, ElementRef, Inject, OnInit, ViewChild } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { DOCUMENT } from '@angular/common';

import { Config } from '../../shared/models/config.model';
import { Link } from '../../shared/models/link.model';
import { Site } from '../../shared/models/site.model';
import { ConfigService } from '../services/config.service';
import { PageScrollInstance, PageScrollService, EasingLogic } from 'ngx-page-scroll-core';
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

  @ViewChild('cardContainer')
  public cardContainer: ElementRef;

  public myEasing: EasingLogic = (t: number, b: number, c: number, d: number): number => {
      // easeInOutExpo easing
      if (t === 0) { return b; }
      if (t === d) { return b + c; }
      if ((t /= d / 2) < 1) {
        return c / 2 * Math.pow(2, 10 * (t - 1)) + b;
      }
      return c / 2 * (-Math.pow(2, -10 * --t) + 2) + b;
  }

  constructor(
    @Inject(DOCUMENT) private document: any,
    private configService: ConfigService,
    private titleService: Title,
    private pageScrollService: PageScrollService
  ) { }

  ngOnInit() {
    this.titleService.setTitle(this.title);
    const el = document.getElementById('cardContainer');
    if (el) {
      el.onscroll = () => this.hasScrolled();
    }
    this.getConfig();
  }

  public hasScrolled() {
    const el = document.getElementById('config-scroll-to-top');
    const top = document.getElementById('cardContainer');
    if (el && top) {
      if (top.scrollTop > 100 || top.scrollTop > 100) {
        el.style.display = 'block';
      } else {
        el.style.display = 'none';
      }
    }
  }

  private getConfig() {
    this.configService.get().subscribe((res: any) => {
      this.config = new Config(res.network, res.links, res.sites);
      this.panelOpenState = [];
      for (let site of this.config.sites) {
        this.panelOpenState.push(false);
      }
    });
  }

  public saveConfig() {
    this.configService.save(this.config).subscribe();
  }

  public exportConfig() {
    this.configService.exportJson(this.config);
  }

  public importConfig() {
    this.configService.importJson().then((conf) => {
      this.config = conf;
      console.log(this.config)
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

  public addSite() {
    const site = new Site();
    this.config.sites.push(site);
  }

  public deleteLink(link: Link) {
    const idx = this.config.links.indexOf(link);
    this.config.links.splice(idx, 1);
  }

  public deleteSite(site: Site) {
    console.log(site);
    const idx = this.config.sites.indexOf(site);
    this.config.sites.splice(idx, 1);
  }

  public scrollToTop() {
    const instance: PageScrollInstance = this.pageScrollService.create({
      document: this.document,
      scrollTarget: "#top",
      scrollViews: [this.cardContainer.nativeElement],
      easingLogic: this.myEasing,
      advancedInlineOffsetCalculation: true,
    });
    this.pageScrollService.start(instance);
  }
}
