import { AfterViewInit, Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import {
  AddExprToCollectSpecGQL,
  GetAvailableVariablesGQL,
  GetCollectionGQL,
  ProcessSnapshot,
  RemoveExprFromCollectSpecGQL
} from "../../graphql/graphql-codegen-generated";
import { Router } from "@angular/router";
import { AppCoreService, WeightedTreeComponent } from 'traceviz/dist/ngx-traceviz-lib';
import { Action, IntegerValue, Update, ValueMap, } from "traceviz-client-core";
import { MatDrawer } from "@angular/material/sidenav";
import { ResizeEvent } from 'angular-resizable-element';
import { CheckedEventArg, TypeInfoComponent } from "./type-info.component";
import { MatSelect } from "@angular/material/select";

class Frame {
  constructor(public name: string, public file: string, public line: number) {
  };
}

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
  // collectionID and snapshotID input properties are set by the router.
  @Input('colID') collectionID!: number;
  private _snapshotID!: number;
  @Input('snapID') set snapshotID(val: number) {
    this._snapshotID = val;
    // When setting the snapshot ID, we update it in the global state too.
    try {
      this.appCoreService.appCore.globalState.get("snapshot_id");
    } catch (error) {
      this.appCoreService.appCore.globalState.set("snapshot_id", new IntegerValue(0));
    }
    this.appCoreService.appCore.globalState.get(
      "snapshot_id").fold(new IntegerValue(val), false, true);
  }

  get snapshotID(): number {
    return this._snapshotID;
  }

  protected collectionName?: string;
  protected snapshots?: Partial<ProcessSnapshot>[];
  @ViewChild(WeightedTreeComponent) weightedTree: WeightedTreeComponent | undefined;
  @ViewChild('functionDrawer') frameDetailsSidebar!: MatDrawer;
  @ViewChild(TypeInfoComponent) typeInfo?: TypeInfoComponent;
  @ViewChild('snapshotsSelect') snapSelect!: MatSelect;

  protected selectedFrame?: Frame;
  // Data about the selected node. Each element is a string containing all the
  // captured variables from one frame (where all frames correspond to the
  // selected node).
  protected funcInfo?: string[];

  protected loadingAvailableVars: boolean = false;

  constructor(
    private readonly appCoreService: AppCoreService,
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly addExpr: AddExprToCollectSpecGQL,
    private readonly removeExpr: RemoveExprFromCollectSpecGQL,
    private readonly router: Router,
  ) {
    console.log('ctor');
  }

  ngOnInit(): void {
    console.log('ngOnInit');
    this.getCollectionQuery.fetch({colID: this.collectionID})
      .subscribe(results => {
        this.collectionName = results.data.collectionByID?.name;
        this.snapshots = results.data.collectionByID?.processSnapshots!;
      })

    this.appCoreService.appCore.globalState.set(
      "collection_id", new IntegerValue(this.collectionID));
  }

  ngAfterViewInit() {
    this.weightedTree!.interactionsDir?.get().withAction(
      new Action(WeightedTreeComponent.NODE, WeightedTreeComponent.CTRL_CLICK,
        new Call(this.onNodeCtrlClick.bind(this))));
  }

  onNodeCtrlClick(localState: ValueMap): void {
    if (localState.has('vars')) {
      this.funcInfo = localState.expectStringList('vars');
    }

    const funcName = localState.expectString('full_name');
    this.selectedFrame = new Frame(funcName, localState.expectString('file'), localState.expectNumber('line'))
    const pcOffset = localState.expectNumber('pc_off');
    console.log("querying for available vars for func: %s off: %d", funcName, pcOffset);
    this.loadingAvailableVars = true;
    this.typeInfo!.dataSource.data = [];
    this.varsQuery.fetch({func: funcName, pcOff: pcOffset})
      .subscribe(
        results => {
          console.log("got results");
          this.loadingAvailableVars = false;
          if (results.error) {
            console.log("err: ", results.error)
            return
          }
          this.typeInfo!.dataSource.initData(
            results.data.availableVars.Vars,
            results.data.availableVars.Types,
            results.data.frameInfo ? results.data.frameInfo.exprs : [],
          )
        })
    console.log("opening sidebar")
    this.frameDetailsSidebar.open();
  }

  closeSidebar() {
    this.frameDetailsSidebar.close()
  }

  onSelectedSnapshotChange(newValue: string) {
    let newSnapshotID = Number(newValue);
    // Navigate to the updated route. This will cause the snapshotID input
    // property to be updated by the router, which in turn causes everything to
    // be reloaded and re-rendered.
    this.router.navigateByUrl(`collections/${this.collectionID}/snap/${newSnapshotID}`);
  }


  checkedChange(ev: CheckedEventArg) {
    if (ev.checked) {
      this.addExpr.mutate({frame: this.selectedFrame!.name, expr: ev.expr}).subscribe({
        next: value => console.log(value.data?.addExprToCollectSpec.frames![0].exprs)
      })
    } else {
      this.removeExpr.mutate({frame: this.selectedFrame!.name, expr: ev.expr}).subscribe({
        next: value => console.log(value.data?.removeExprFromCollectSpec.frames![0].exprs)
      })
    }
  }

  public style: object = {
    "width": '500px',
    "text-color": 'red',
  };

  onResizeEnd(event: ResizeEvent): void {
    this.style = {
      width: `${event.rectangle.width}px`,
    };

  }

  validateResize(event: ResizeEvent): boolean {
    return true;
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
