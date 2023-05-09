import { Component } from '@angular/core';

import { AppCoreService } from 'traceviz/dist/ngx-traceviz-lib';
import { StringValue } from 'traceviz-core';

/** The application component of the LogViz client. */
@Component({
  selector: 'logviz',
  templateUrl: './app.component.html',
  styleUrls: ['app.component.css'],
})
export class AppComponent {
  constructor(public readonly appCoreService: AppCoreService) {
  }

  ngAfterContentInit(): void {
  }
}