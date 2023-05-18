import { AfterContentInit, Component, OnInit } from '@angular/core';
import {
  AllCollectionsGQL,
  GetCollectionGQL,
  ProcessSnapshot
} from "../../graphql/graphql-codegen-generated";
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
  protected collectionName?: string;
  protected snapshots?: ProcessSnapshot[];

  constructor(private getCollectionQuery: GetCollectionGQL,
              private route: ActivatedRoute) {}

  ngOnInit(): void {
    this.snapshotID = Number(this.route.snapshot.paramMap.get('id'));

    this.getCollectionQuery.fetch({colID: this.snapshotID.toString()}).subscribe(results => {
      this.collectionName = results.data.collectionByID?.name;
      this.snapshots = results.data.collectionByID?.processSnapshots?.map(
        value => ({...value, snapshot: "",})
      )
    })
  }

  ngAfterContentInit(): void {
  }
}
