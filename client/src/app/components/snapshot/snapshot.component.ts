import { AfterContentInit, Component, OnInit } from '@angular/core';
import {
  AllCollectionsGQL,
  GetCollectionGQL,
  ProcessSnapshot
} from "../../graphql/graphql-codegen-generated";
import { ActivatedRoute } from "@angular/router";
import { AppCoreService } from 'traceviz/dist/ngx-traceviz-lib';
import { IntegerValue } from "traceviz-client-core";

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
  protected collectionID?: number;
  protected snapshotID?: number;
  protected collectionName?: string;
  protected snapshots?: ProcessSnapshot[];

  constructor(
    private readonly appCoreService: AppCoreService,
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly route: ActivatedRoute) {
  }

  ngOnInit(): void {
    // Get the collection ID and snapshot ID from the URL. The names of the URL
    // params are defined in the Routes collection.
    this.collectionID = Number(this.route.snapshot.paramMap.get('colID'));
    this.snapshotID = Number(this.route.snapshot.paramMap.get('snapID'));
    console.log("!!! reading collection ", this.collectionID)
    this.getCollectionQuery.fetch({colID: this.collectionID.toString()}).subscribe(results => {
      this.collectionName = results.data.collectionByID?.name;
      this.snapshots = results.data.collectionByID?.processSnapshots?.map(
        value => ({...value, snapshot: "",})
      )
      console.log("!!! snapshots: ", this.snapshots)
    })
    // TODO(andrei): This is a hack; I don't know how to set the initial value
    // for a field in globalState. Even if I attempt to do it in an
    // appCore.onPublish() callback, I think it's still too early; the
    // respective state does not seem to exist by then (maybe because this
    // callback fires before the one that creates them based on the children in
    // the template.
    setTimeout(() => {
      this.appCoreService.appCore.globalState.get("snapshot_id").fold(new IntegerValue(this.snapshotID!), false /* toggle */);
    }, 100)
    // this.appCoreService.appCore.onPublish(
    //   (appCore) => {
    //     appCore.globalState.get("snapshot_id").fold(new IntegerValue(this.snapshotID), false /* toggle */);
    //   })
  }

  ngAfterContentInit(): void {
  }

  onSelectedSnapshotChange(newValue: string) {
    let newSnapshotID = Number(newValue);
    this.appCoreService.appCore.globalState.get("snapshot_id").fold(new IntegerValue(newSnapshotID), false /* toggle */);
  }
}
