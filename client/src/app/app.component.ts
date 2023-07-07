import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
    selector: 'stacksviz',
    template: `
    <div class="app">
      <router-outlet></router-outlet>
    </div>
  `,
    styleUrls: ['app.component.css'],
    standalone: true,
    imports: [RouterOutlet],
})
export class AppComponent {
  constructor() {
  }
}
