import {AfterContentInit, Component, ContentChild} from '@angular/core';

import { AppCoreService } from 'traceviz/dist/ngx-traceviz-lib';

/** The application component of the LogViz client. */
@Component({
  selector: 'logviz',
  templateUrl: './app.component.html',
  styleUrls: ['app.component.css'],
})
export class AppComponent implements AfterContentInit {
  constructor(public readonly appCoreService: AppCoreService) {
  }

  ngAfterContentInit(): void {}
}
