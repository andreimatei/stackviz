import { Component } from '@angular/core';

import { AppCoreService } from 'traceviz/dist/ngx-traceviz-lib';

/** The application component of the LogViz client. */
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['app.component.css'],
})
export class AppComponent {
  constructor(public appCoreService: AppCoreService) {
  }
}
