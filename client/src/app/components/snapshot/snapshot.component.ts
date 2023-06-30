import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import {
  AddExprToCollectSpecGQL, FrameSpec,
  GetAvailableVariablesGQL,
  GetCollectionGQL, GetCollectSpecGQL,
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
import {
  FlamegraphComponent,
  Frame as FlameFrame,
  VarInfo
} from "../flamegraph/flamegraph.component";
import {
  debounceTime,
  distinctUntilChanged,
  map,
  merge,
  Observable,
  pipe,
  Subject,
  tap
} from "rxjs";
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
    console.log("snapshot set to", val);
    this._snapshotID = val;
    this.snapshotID$.next(val);
  }

  get snapshotID(): number {
    return this._snapshotID;
  }

  @Input('filter') set filter(val: string) {
    this.onFilterChange(val, true);
  }

  onFilterChange(val: string, immediate: boolean) {
    this.filterVal = val;
    this.filter$.next(val);
    if (immediate) {
      this.filterChangeImmediate$.next(val);
    }
  }


  protected filter$ = new Subject<string>();
  protected filterDebounced$ = this.filter$.pipe(
    debounceTime(400),
    distinctUntilChanged(),
    tap(val => {
      const url = new URL(window.location.href);
      url.searchParams.set('filter', val);
      window.history.replaceState(null, "", url);
    })
  );
  protected filterVal?: string;
  protected filterChangeImmediate$ = new Subject<string>()
  protected filterChange$ = merge(
    this.filterDebounced$,
    this.filterChangeImmediate$,
  );

  protected collectionName?: string;
  protected snapshots?: Partial<ProcessSnapshot>[];


  protected capturedData$: any;

  @ViewChild('functionDrawer') frameDetailsSidebar!: MatDrawer;
  @ViewChild(TypeInfoComponent) typeInfo?: TypeInfoComponent;
  @ViewChild('snapshotsSelect') snapSelect!: MatSelect;
  @ViewChild(FlamegraphComponent) flamegraph!: FlamegraphComponent;
  @ViewChild(StacksComponent) stacks!: StacksComponent;

  protected selectedFrame?: Frame;
  // Data about the selected node. Each element is a string containing all the
  // captured variables from one frame (where all frames correspond to the
  // selected node).
  protected funcInfo?: VarInfo[][];

  protected loadingAvailableVars: boolean = false;

  private treeQuery!: ReturnType<GetTreeGQL["watch"]>;
  private goroutinesQuery!: ReturnType<GetGoroutinesGQL["watch"]>;
  protected readonly collectSpecQuery$!: Observable<Partial<FrameSpec>[]>;

  constructor(
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly getGoroutinesQuery: GetGoroutinesGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly getTreeQuery: GetTreeGQL,
    private readonly addExpr: AddExprToCollectSpecGQL,
    private readonly removeExpr: RemoveExprFromCollectSpecGQL,
    private readonly collectSpec: GetCollectSpecGQL,
    private readonly router: Router,
    private readonly route: ActivatedRoute,
  ) {
    this.collectSpecQuery$ = collectSpec.watch().valueChanges.pipe(
      map(val => val.data.collectSpec),
    );
  }

  ngOnInit(): void {
    this.getCollectionQuery.fetch({colID: this.collectionID})
      .subscribe(results => {
        this.collectionName = results.data.collectionByID?.name;
        this.snapshots = results.data.collectionByID?.processSnapshots!;
      });

    const parsedFilter = this.filterVal ? parseFilter(this.filterVal) : {};
    const args = {
      colID: this.collectionID,
      snapID: this.snapshotID,
      filter: parsedFilter.filter,
      gID: parsedFilter.gID,
    }

    this.treeQuery = this.getTreeQuery.watch(args);
    this.goroutinesQuery = this.getGoroutinesQuery.watch(args);

    // Refetch all the data when either the selected snapshot or the filter
    // changes.
    merge(
      this.filterChange$,
      this.snapshotID$,
    ).pipe(
      // Allow a small amount of time to coalesce the filterChange$ and
      // snapshotID$ signals that are sometimes set together (i.e. by the router
      // when navigating a link).
      debounceTime(30)
    ).subscribe((val) => {
      const parsedFilter = this.filterVal ? parseFilter(this.filterVal) : {};
      const args = {
        colID: this.collectionID,
        snapID: this.snapshotID,
        filter: parsedFilter.filter,
        gID: parsedFilter.gID,
      };
      console.log(`refetching with snapshot: ${this.snapshotID}, filter ${this.filterVal}`);
      this.treeQuery.refetch(args);
      this.goroutinesQuery.refetch(args);
    });

    this.capturedData$ = this.goroutinesQuery.valueChanges.pipe(
      map(res =>
        res.data.goroutines.Raw
          .filter(g => g.Vars && g.Vars.length > 0)
          .map(g => ({gid: g.ID, vars: g.Vars}))
      ),
      tap(v => console.log("!!! captured data update:", v)),
    );
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
  };

  onResizeEnd(event: ResizeEvent): void {
    this.style = {
      width: `${event.rectangle.width}px`,
    };
  }

  showDetails(node: FlameFrame): void {
    console.log("showDetails. funcInfo:", node.vars);
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
            results.data.collectSpec ? results.data.collectSpec[0].exprs : [],
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

function parseFilter(filter: string): parsedFilter {
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

