import { AfterContentInit, Component, OnInit } from '@angular/core';
import { AllCollectionsGQL } from "../../graphql/graphql-codegen-generated";
import { ActivatedRoute } from "@angular/router";
import { AppCoreService } from 'traceviz/dist/ngx-traceviz-lib';

@Component({
  selector: 'snapshot',
  templateUrl: './snapshot.component.html',
  styleUrls: ['snapshot.component.css'],
  // Provide the AppCoreService at the level of this component, overriding the
  // `providedIn: 'root'` specified on the declaration of AppCoreService (which
  // asks for a single instance of AppCoreService to be injected everywhere in
  // the app). By scoping an instance to this SnapshotComponent, we prevent the
  // state from escaping this component. In turn, this means that every time
  // that a user navigates to the snapshot page, she gets an empty state (i.e.
  // any filtering they might have previously done on the page is gone; this is
  // considered a good thing.
  providers: [AppCoreService]
})
export class SnapshotComponent implements AfterContentInit, OnInit {
  protected snapshotID?: number;
  protected snapshotName?: string;

  constructor(private ac: AllCollectionsGQL, private route: ActivatedRoute) {}

  ngOnInit(): void {
    this.snapshotID = Number(this.route.snapshot.paramMap.get('id'));
    this.ac.fetch().subscribe(
      results => {
        results.data.collections.forEach(c => {
            if (Number(c.id) == this.snapshotID) {
              this.snapshotName = c.name;
            }
          })
      });
  }

  ngAfterContentInit(): void {
  }
}
