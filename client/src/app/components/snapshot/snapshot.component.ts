import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import {
  GetAvailableVariablesGQL, GetAvailableVariablesQuery,
  GetCollectionGQL,
  ProcessSnapshot
} from "../../graphql/graphql-codegen-generated";
import { ActivatedRoute } from "@angular/router";
import { AppCoreService, WeightedTreeComponent } from 'traceviz/dist/ngx-traceviz-lib';
import { Action, IntegerValue, Update, ValueMap, } from "traceviz-client-core";
import { MatDrawer } from "@angular/material/sidenav";

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
export class SnapshotComponent implements OnInit, AfterViewInit {
  protected collectionID!: number;
  protected snapshotID!: number;
  protected collectionName?: string;
  protected snapshots?: ProcessSnapshot[];
  protected availableVars: GetAvailableVariablesQuery['availableVars']['Vars'];
  @ViewChild(WeightedTreeComponent) weightedTree: WeightedTreeComponent | undefined;
  @ViewChild('functionDrawer') input!: MatDrawer;

  // Data about the selected node. Each element is a string containing all the
  // captured variables from one frame (where all frames correspond to the
  // selected node).
  protected funcInfo?: string[];

  constructor(
    private readonly appCoreService: AppCoreService,
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly route: ActivatedRoute) {
  }

  ngOnInit(): void {
    // Get the collection ID and snapshot ID from the URL. The names of the URL
    // params are defined in the Routes collection.
    this.collectionID = Number(this.route.snapshot.paramMap.get('colID'));
    this.snapshotID = Number(this.route.snapshot.paramMap.get('snapID'));
    this.getCollectionQuery.fetch({colID: this.collectionID.toString()})
      .subscribe(results => {
        this.collectionName = results.data.collectionByID?.name;
        this.snapshots = results.data.collectionByID?.processSnapshots?.map(
          value => ({...value, snapshot: "",})
        )
      })
    this.appCoreService.appCore.globalState.set(
      "snapshot_id", new IntegerValue(this.snapshotID));
  }

  ngAfterViewInit() {
    this.weightedTree!.interactionsDir!.get().withAction(
      new Action(WeightedTreeComponent.NODE, WeightedTreeComponent.CTRL_CLICK,
        new Call(this.onNodeCtrlClick.bind(this))));
  }

  onNodeCtrlClick(localState: ValueMap): void {
    console.log(localState);
    console.log(this.input!);
    if (localState.has('vars')) {
      console.log("!!! vars: ", localState.get('vars').toString());
      this.funcInfo = localState.expectStringList('vars');
    }

    const funcName = localState.expectString('full_name');
    const pcOffset = localState.expectNumber('pc_off');
    console.log("!!! issuing available vars query: ", funcName, pcOffset);
    this.varsQuery.fetch({func: funcName, pcOff: pcOffset})
     .subscribe(
       results => {
       console.log("!!! got available vars: ", results.data.availableVars);
       this.availableVars = results.data.availableVars.Vars;
     })

    this.input.toggle();
  }

  onSelectedSnapshotChange(newValue: string) {
    let newSnapshotID = Number(newValue);
    this.appCoreService.appCore.globalState.get("snapshot_id").fold(new IntegerValue(newSnapshotID), false /* toggle */);
  }
}

// Call is an implementation of Update that calls the provided function.
class Call extends Update {
  constructor(private readonly handler: ((localState: ValueMap) => void)) {
    super();
  }

  override update(localState?: ValueMap | undefined) {
    this.handler(localState!)
  }

  get autoDocument(): string {
    throw new Error('Method not implemented.');
  }
}
