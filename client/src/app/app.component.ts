import { Component } from '@angular/core';

@Component({
  selector: 'stacksviz',
  template: `
    <div class="app">
      <router-outlet></router-outlet>
    </div>
  `,
  styleUrls: ['app.component.css'],
})
export class AppComponent {
  constructor() {
  }
}
