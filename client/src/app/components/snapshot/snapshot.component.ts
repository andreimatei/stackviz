import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import {
  AddExprToCollectSpecGQL,
  AddFlightRecorderEventToCollectSpecGQL,
  FrameSpec,
  GetAvailableVariablesGQL,
  GetCollectionGQL,
  GetFrameSpecsGQL,
  GetSnapshotGQL,
  GetTreeGQL,
  ProcessSnapshot,
  RemoveExprFromCollectSpecGQL,
  SyncFlightRecorderGQL
} from "../../graphql/graphql-codegen-generated";
import { Router, RouterLink } from "@angular/router";
import { MatDrawer, MatSidenavModule } from "@angular/material/sidenav";
import { ResizableModule, ResizeEvent } from 'angular-resizable-element';
import {
  CheckedEvent,
  FlightRecorderEvent,
  FlightRecorderEventSpec,
  goroutineIDKey,
  TypeInfoComponent
} from "./type-info.component";
import { MatSelect, MatSelectModule } from "@angular/material/select";
import {
  FlamegraphComponent,
  Frame as FlameFrame,
  VarInfo
} from "../flamegraph/flamegraph.component";
import {
  debounceTime,
  distinctUntilChanged,
  from,
  map,
  merge,
  Observable,
  of,
  Subject,
  tap
} from "rxjs";
import { StacksComponent } from "../stacks/stacks.component";
import { CapturedDataComponent, GoroutineData } from "../captured-data/captured-data.component";
import { MatTabsModule } from '@angular/material/tabs';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatOptionModule } from '@angular/material/core';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { AsyncPipe, KeyValuePipe, NgFor, NgIf, NgStyle } from '@angular/common';

class Frame {
  constructor(public name: string, public file: string, public line: number) {
  };
}

@Component({
  selector: 'snapshot',
  templateUrl: './snapshot.component.html',
  styleUrls: ['snapshot.component.css'],
  standalone: true,
  imports: [
    MatSidenavModule,
    ResizableModule,
    NgStyle,
    MatButtonModule,
    MatIconModule,
    MatExpansionModule,
    NgIf,
    MatProgressBarModule,
    TypeInfoComponent,
    NgFor,
    CapturedDataComponent,
    RouterLink,
    MatFormFieldModule,
    MatSelectModule,
    MatOptionModule,
    MatInputModule,
    ReactiveFormsModule,
    FormsModule,
    FlamegraphComponent,
    MatTabsModule,
    StacksComponent,
    AsyncPipe,
    KeyValuePipe,
  ],
})
export class SnapshotComponent implements OnInit, AfterViewInit {
  // collectionID and snapshotID input properties are set by the router.
  @Input('colID') collectionID!: number;
  // The id of the specification that produced the collection.
  protected collectSpecID!: number;

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

  protected collectionName!: string;
  protected snapshots?: Partial<ProcessSnapshot>[];


  protected capturedData$!: Observable<GoroutineData[]>;

  @ViewChild('functionDrawer') frameDetailsSidebar!: MatDrawer;
  @ViewChild(TypeInfoComponent) typeInfo!: TypeInfoComponent;
  @ViewChild('snapshotsSelect') snapSelect!: MatSelect;
  @ViewChild(FlamegraphComponent) flamegraph!: FlamegraphComponent;
  @ViewChild(StacksComponent) stacks!: StacksComponent;

  protected selectedFrame?: Frame;

  // Data about the selected frame. Map from goroutine ID to variables collected
  // for the selected frame in that goroutine.
  protected funcInfo?: Map<number, VarInfo[]>;
  protected flightRecorderEventSpecs?: FlightRecorderEventSpec[];

  protected loadingAvailableVars: boolean = false;

  private treeQuery!: ReturnType<GetTreeGQL["watch"]>;
  private goroutinesQuery!: ReturnType<GetSnapshotGQL["watch"]>;
  protected frameSpecsQuery$!: Observable<Partial<FrameSpec>[]>;

  constructor(
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly getGoroutinesQuery: GetSnapshotGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly getTreeQuery: GetTreeGQL,
    private readonly addExpr: AddExprToCollectSpecGQL,
    private readonly removeExpr: RemoveExprFromCollectSpecGQL,
    private readonly frameSpecsQuery: GetFrameSpecsGQL,
    private readonly addFlightRecorderEventSpecQuery: AddFlightRecorderEventToCollectSpecGQL,
    private readonly syncFlightRecorderQuery: SyncFlightRecorderGQL,
    private readonly router: Router,
  ) {
  }

  ngOnInit(): void {
    this.getCollectionQuery.fetch({colID: this.collectionID})
      .subscribe(results => {
        this.collectionName = results.data.collectionByID!.name;
        this.collectSpecID = results.data.collectionByID!.collectSpec;
        this.snapshots = results.data.collectionByID?.processSnapshots!;

        this.frameSpecsQuery$ = this.frameSpecsQuery.fetch({collectSpecID: this.collectSpecID}).pipe(
          map(val => val.data.frameSpecsWhere),
        )
        // !!!
        // this.frameSpecsQuery$ = frameSpecsQuery.watch({collect_spec_id: this.collectSpecID}).valueChanges.pipe(
        //   map(val => val.data.frameSpecsWhere),
        // );
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
        res.data.getSnapshot.Raw
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
      tap(res => {
        console.log("loading flight recorder data:",
          typeof res.data.getSnapshot.FlightRecorderData,
          res.data.getSnapshot.FlightRecorderData,
        );

        const u = res.data.getSnapshot.FlightRecorderData;
        const m = new Map(Object.entries(u));
        console.log("m:", m);
        console.log("u:", u);
        console.log("m:", m.get('2312'));

      }),
      map(res => res.data.getSnapshot)
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

  checkedChange(ev: CheckedEvent) {
    if (ev.checked) {
      this.addExpr.mutate({frame: this.selectedFrame!.name, expr: ev.expr}).subscribe({
        next: value => console.log(value.data?.addExprToCollectSpec.frames![0].collectExpressions)
      })
    } else {
      this.removeExpr.mutate({frame: this.selectedFrame!.name, expr: ev.expr}).subscribe({
        next: value => console.log(value.data?.removeExprFromCollectSpec.frames![0].collectExpressions)
      })
    }
  }

  flightRecorderChange(ev: FlightRecorderEvent) {
    if (ev.deleted) {
      // !!! TODO
    } else {
      console.log("adding flight recorder event:, frame:", ev, this.selectedFrame!.name);
      let key = ev.key == goroutineIDKey ? "goroutine_id" : ev.key;
      this.addFlightRecorderEventSpecQuery
        .mutate({
          collectSpecID: this.collectSpecID,
          frame: this.selectedFrame!.name,
          expr: ev.expr,
          keyExpr: key,
        })
        .subscribe({
          next: value =>
            console.log("flight recorder events are now: ",
              value.data?.addFlightRecorderEventToFrameSpec.flightRecorderEvents)
        })
    }
  }

  syncFlightRecorder() {
    console.log("syncing flight recorder");
    this.syncFlightRecorderQuery.mutate({collectSpecID: this.collectSpecID}).subscribe({
      next: value => {
        console.log("sync done");
      }
    });
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
    this.typeInfo.dataSource.data = [];
    this.varsQuery.fetch({func: node.details, pcOff: node.pcoff})
      .subscribe(
        results => {
          console.log("got results", results);
          this.loadingAvailableVars = false;
          if (results.error) {
            console.log("err: ", results.error)
            return
          }
          if (results.data.frameSpecsWhere.length > 1) {
            console.log("expected at most one frameSpec, got: ", results.data.frameSpecsWhere)
            return
          }
          const frameSpec = results.data.frameSpecsWhere[0];
          const collectExprs = frameSpec ? frameSpec.collectExpressions : [];
          if (frameSpec) {
            this.flightRecorderEventSpecs = frameSpec.flightRecorderEvents.map(
              e => JSON.parse(e) as FlightRecorderEventSpec);
            console.log("loaded flightRecorderEvents:", this.flightRecorderEventSpecs);
          } else {
            this.flightRecorderEventSpecs = [];
          }
          this.typeInfo.dataSource.initData(
            results.data.availableVars.Vars,
            results.data.availableVars.Types,
            collectExprs,
          )
        })
    this.frameDetailsSidebar.open();
  }

  validateResize(event: ResizeEvent): boolean {
    return true;
  }

  protected readonly from = from;
  protected readonly of = of;
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
