import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import {
  AddExprToCollectSpecGQL,
  GetAvailableVariablesGQL,
  GetCollectionGQL,
  GetGoroutinesGQL,
  GetTreeGQL,
  ProcessSnapshot,
  RemoveExprFromCollectSpecGQL
} from "../../graphql/graphql-codegen-generated";
import { Router } from "@angular/router";
import { WeightedTreeComponent } from 'traceviz/dist/ngx-traceviz-lib';
import { Update, ValueMap, } from "traceviz-client-core";
import { MatDrawer } from "@angular/material/sidenav";
import { ResizeEvent } from 'angular-resizable-element';
import { CheckedEventArg, TypeInfoComponent } from "./type-info.component";
import { MatSelect } from "@angular/material/select";
import { FlamegraphComponent, Frame as FlameFrame } from "../flamegraph/flamegraph.component";
import { BehaviorSubject, map, switchMap } from "rxjs";
import { StacksComponent } from "../stacks/stacks.component";

class Frame {
  constructor(public name: string, public file: string, public line: number) {
  };
}

@Component({
  selector: 'snapshot',
  templateUrl: './snapshot.component.html',
  styleUrls: ['snapshot.component.css'],
  // !!!
  // // Provide the AppCoreService at the level of this component, overriding the
  // // `providedIn: 'root'` specified on the declaration of AppCoreService (which
  // // asks for a single instance of AppCoreService to be injected everywhere in
  // // the app). By scoping an instance to this SnapshotComponent, we prevent the
  // // state from escaping this component. In turn, this means that every time
  // // that a user navigates to the snapshot page, she gets an empty state (i.e.
  // // any filtering they might have previously done on the page is gone; this is
  // // considered a good thing).
  // providers: [AppCoreService]
})
export class SnapshotComponent implements OnInit, AfterViewInit {
  // collectionID and snapshotID input properties are set by the router.
  @Input('colID') collectionID!: number;
  // The initial value of snapshotID$ will be overwritten by the time the
  // observable is read.
  protected snapshotID$ = new BehaviorSubject<number>(-1);

  @Input('snapID') set snapshotID(val: number) {
    console.log("!!! snapshotID being set to:", val);
    this.snapshotID$.next(val);
  }

  protected collectionName?: string;
  protected snapshots?: Partial<ProcessSnapshot>[];
  @ViewChild(WeightedTreeComponent) weightedTree: WeightedTreeComponent | undefined;
  @ViewChild('functionDrawer') frameDetailsSidebar!: MatDrawer;
  @ViewChild(TypeInfoComponent) typeInfo?: TypeInfoComponent;
  @ViewChild('snapshotsSelect') snapSelect!: MatSelect;
  @ViewChild(FlamegraphComponent) flamegraph!: FlamegraphComponent;
  @ViewChild(StacksComponent) stacks!: StacksComponent;

  protected selectedFrame?: Frame;
  // Data about the selected node. Each element is a string containing all the
  // captured variables from one frame (where all frames correspond to the
  // selected node).
  protected funcInfo?: string[];

  protected loadingAvailableVars: boolean = false;

  constructor(
    // !!!
    // private readonly appCoreService: AppCoreService,
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly getGoroutinesQuery: GetGoroutinesGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly getTreeQuery: GetTreeGQL,
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

    // !!!
    // this.appCoreService.appCore.globalState.set(
    //   "collection_id", new IntegerValue(this.collectionID));

    // this.getTreeQuery.fetch({colID: this.collectionID, snapID: this.snapshotID})
    //   .subscribe(value => {
    //     console.log("!!! query got result")
    //     this.flamegraph.data = JSON.parse(value.data.getTree);
    //   })

  }

  ngAfterViewInit() {
    console.log("!!! snapshot afterviewinit")
    // Bind the flamegraph data to updates of snapshot ID.
    this.flamegraph.data$ = this.snapshotID$.pipe(
      // !!! tap(val => console.log("snapshotID$ got", val)),
      switchMap(snapID =>
        this.getTreeQuery.fetch({colID: this.collectionID, snapID: snapID})
          .pipe(
            map(res => JSON.parse(res.data.getTree))
          )
      )
    )

    this.stacks.data$ = this.snapshotID$.pipe(
      switchMap(snapID => {
          type Args = {
            colID: number,
            snapID: number,
            gID: number | undefined,
          };
          const args: Args = {colID: this.collectionID, snapID: snapID, gID: undefined};
          const urlParts = document.URL.split('#');
          if (urlParts.length > 1) {
            if (urlParts[1].startsWith('g_')) {
              args.gID = Number(urlParts[1].slice(2));
              console.log("filtering for goroutine: ", args.gID);
            }
          }
          return this.getGoroutinesQuery.fetch(args)
            .pipe(
              map(res => res.data.goroutines)
            )
        }
      )
    )

    // !!! update the sidebar in response to snapshotID changes

  }

  // !!! remove
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
      console.log("!!! checked: ", this.selectedFrame!.name, ev.expr);
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

  showDetails(node: FlameFrame): void {
    console.log("showDetails", node.details);

    this.funcInfo = node.vars;
    this.selectedFrame = new Frame(node.details, node.file, node.line);
    console.log("querying for available vars for func: %s off: %d", node.details, node.pcoff);
    this.loadingAvailableVars = true;
    this.typeInfo!.dataSource.data = [];
    this.varsQuery.fetch({func: node.details, pcOff: node.pcoff})
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
    this.frameDetailsSidebar.open();
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
