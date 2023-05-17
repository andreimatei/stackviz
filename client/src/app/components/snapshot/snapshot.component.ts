import { AfterContentInit, Component } from '@angular/core';
import { AllCollectionsGQL } from "../../graphql/graphql-codegen-generated";

/** The application component of the LogViz client. */
@Component({
  selector: 'snapshot',
  templateUrl: './snapshot.component.html',
  styleUrls: ['snapshot.component.css'],
})
export class SnapshotComponent implements AfterContentInit {
  constructor(ac: AllCollectionsGQL) {
    ac.fetch().subscribe(results => console.log("!!! GraphQL results:", JSON.stringify(results.data)));
  }

  ngAfterContentInit(): void {
  }
}
