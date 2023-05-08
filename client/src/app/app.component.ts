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
    console.log('AppComponent: ngAfterContentInit');
    this.appCoreService.appCore.onPublish((appCore) => {
      console.log('AppComponent: onPublish callback');
      appCore.globalState.subscribe(v => {
        console.log("global state update: ", v)
        if (v.indexOf("dummy") != -1) {
          appCore.globalState.get("dummy").subscribe((val: any) => {
            console.log("updated value - ", (val as StringValue).val);
          })
        }
      })
    })
  }
}