import { AfterContentInit, Component } from '@angular/core';
import { AllCollectionsGQL } from "./graphql/graphql-codegen-generated";

/** The application component of the LogViz client. */
@Component({
  selector: 'stacksviz',
  template: '<router-outlet></router-outlet>',
  styleUrls: ['app.component.css'],
})
export class AppComponent implements AfterContentInit {
  constructor(ac: AllCollectionsGQL) {
    ac.fetch().subscribe(results => console.log("!!! GraphQL results:", JSON.stringify(results.data)));
  }

  ngAfterContentInit(): void {
  }
}
