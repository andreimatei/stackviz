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
import { ActivatedRoute, Router } from "@angular/router";
import { MatDrawer } from "@angular/material/sidenav";
import { ResizeEvent } from 'angular-resizable-element';
import { CheckedEventArg, TypeInfoComponent } from "./type-info.component";
import { MatSelect } from "@angular/material/select";
import { FlamegraphComponent, Frame as FlameFrame } from "../flamegraph/flamegraph.component";
import { debounceTime, distinctUntilChanged, map, merge, Subject } from "rxjs";
import { StacksComponent } from "../stacks/stacks.component";

class Frame {
  constructor(public name: string, public file: string, public line: number) {
  };
}

@Component({
  selector: 'snapshot',
  templateUrl: './snapshot.component.html',
  styleUrls: ['snapshot.component.css'],
})
export class SnapshotComponent implements OnInit, AfterViewInit {
  // collectionID and snapshotID input properties are set by the router.
  @Input('colID') collectionID!: number;

  protected snapshotID$ = new Subject<number>();
  private _snapshotID!: number;
  @Input('snapID') set snapshotID(val: number) {
    console.log("!!! snapshotID being set to:", val);
    this._snapshotID = val;
    this.snapshotID$.next(val);
  }

  get snapshotID(): number {
    return this._snapshotID;
  }

  protected filter$ = new Subject<string>();
  protected filter?: string;

  protected collectionName?: string;
  protected snapshots?: Partial<ProcessSnapshot>[];
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

  private treeQuery!: ReturnType<GetTreeGQL["watch"]>;
  private goroutinesQuery!: ReturnType<GetGoroutinesGQL["watch"]>;

  constructor(
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly getGoroutinesQuery: GetGoroutinesGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly getTreeQuery: GetTreeGQL,
    private readonly addExpr: AddExprToCollectSpecGQL,
    private readonly removeExpr: RemoveExprFromCollectSpecGQL,
    private readonly router: Router,
    private readonly route: ActivatedRoute,
  ) {
  }

  ngOnInit(): void {
    const filterParam = this.route.snapshot.queryParamMap.get('filter');
    console.log("filter: ", filterParam);
    if (filterParam) {
      this.onFilterChange(filterParam);
    }

    this.getCollectionQuery.fetch({colID: this.collectionID})
      .subscribe(results => {
        this.collectionName = results.data.collectionByID?.name;
        this.snapshots = results.data.collectionByID?.processSnapshots!;
      });

    const parsedFilter = this.filter ? parseFilter(this.filter) : {};
    const args = {
      colID: this.collectionID,
      snapID: this.snapshotID,
      filter: parsedFilter.filter,
      gID: parsedFilter.gID,
    }

    console.log("!!! watching with args:", args);
    this.treeQuery = this.getTreeQuery.watch(args);
    this.goroutinesQuery = this.getGoroutinesQuery.watch(args);

    // Refetch all the data when either the selected snapshot or the filter
    // changes.
    merge(
      this.filter$.pipe(
        debounceTime(400),
        distinctUntilChanged(),
      ),
      this.snapshotID$,
    ).subscribe(val => {
      const parsedFilter = this.filter ? parseFilter(this.filter) : {};
      const args = {
        colID: this.collectionID,
        snapID: this.snapshotID,
        filter: parsedFilter.filter,
        gID: parsedFilter.gID,
      };
      console.log("refetching with filter:", this.filter);
      this.treeQuery.refetch(args);
      this.goroutinesQuery.refetch(args);
    });
  }

  ngAfterViewInit() {
    // Update the child components when we get new data results.
    this.flamegraph.data$ = this.treeQuery.valueChanges.pipe(
      map(res => JSON.parse(res.data.getTree))
    );
    this.stacks.data$ = this.goroutinesQuery.valueChanges.pipe(
      map(res => res.data.goroutines)
    );
    // TODO(andrei): update the sidebar in response to snapshotID changes
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

  onFilterChange(val: string) {
    this.filter = val;
    this.filter$.next(val);
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

interface parsedFilter {
  gID?: number;
  filter?: string;
}

function parseFilter(filter : string ): parsedFilter {
  const gidPrefix = "gid="
  if (filter.startsWith(gidPrefix)) {
    return {
      gID: Number(filter.slice(gidPrefix.length))
    };
  }
  return {
    filter: filter
  };
}

