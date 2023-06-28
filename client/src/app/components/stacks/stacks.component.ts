import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DataTableModule } from 'traceviz/dist/ngx-traceviz-lib';
import { MatTabsModule, } from '@angular/material/tabs'
import { Observable } from "rxjs";
import { GoroutineInfo, GoroutinesGroup, SnapshotInfo } from "../../graphql/graphql-codegen-generated";
import { MatTableModule } from "@angular/material/table";
import { BacktraceComponent } from "../backtrace/backtrace.component";

@Component({
  selector: 'app-stacks',
  standalone: true,
  imports: [CommonModule, DataTableModule, MatTabsModule, MatTableModule, BacktraceComponent],
  template: `
      <div>
          {{ goroutines?.length }} goroutines
          <!-- !!! reimplement filtering -->
          <!--({{ numFilteredGoroutines }} filtered / {{ numTotalGoroutines }} total Goroutines),-->
          {{ goroutineGroups?.length }} buckets
          <hr>

          <mat-tab-group selectedIndex="0">
              <mat-tab label="Aggregated">
                  <ng-template matTabContent>
                      <ul>
                          <li *ngFor="let g of goroutineGroups">
                              {{g.IDs.length}} goroutines in goroutine group
                              Goroutine IDs {{g.IDs}}
                              <app-backtrace [vars]="g.Vars" [frames]="g.Frames"
                                             [collectionID]="colID"></app-backtrace>
                          </li>
                      </ul>
                  </ng-template>
              </mat-tab>

              <mat-tab label="Raw">
                  #goroutines: {{goroutines ? goroutines.length : 0}}
                  <ul>
                      <li *ngFor="let g of goroutines" id="g_{{g.ID}}">
                          <a id="g_{{g.ID}}">Goroutine {{ g.ID }}</a>
                          <app-backtrace [vars]="g.Vars" [frames]="g.Frames"
                                         [collectionID]="colID"></app-backtrace>
                      </li>
                  </ul>
              </mat-tab>
          </mat-tab-group>
      </div>
  `,
  styles: [`ul {
      list-style-type: none; /* Remove bullets */
      padding: 0; /* Remove padding */
      margin: 0; /* Remove margins */
    }
    li {
      border-bottom: 1px solid black;
    }
    `
  ]
})
export class StacksComponent {
  // @ContentChild(DataSeriesQueryDirective) dataSeriesQueryDir?: DataSeriesQueryDirective;
  // @ContentChild(InteractionsDirective) interactionsDir?: InteractionsDirective;

  // TODO(andrei): mark these as required.
  @Input() colID!: number;
  @Input() snapID?: number | null;
  //
  // private unsubscribe = new Subject<void>();
  // protected rawStacks?: ResponseNode[];
  // protected aggStacks?: ResponseNode[];
  // protected numStacks?: number;
  // protected numBuckets?: number;
  // protected numFilteredGoroutines?: number;
  // protected numTotalGoroutines?: number;

  protected goroutines?: GoroutineInfo[];
  protected goroutineGroups?: GoroutinesGroup[];

  // !!! backtracecols = ['frames'];

  @Input() set data$(val$: Observable<SnapshotInfo>) {
    val$.subscribe(
      value => {
        // Updating all the tables can take a few seconds, so in order to make
        // the user experience better, we first clear them, and only then create
        // the new ones.
        console.log("stacks got value: ", value);
        this.goroutineGroups = value.Aggregated;
        this.goroutines = [];
        setTimeout(() => {
          this.goroutines = value.Raw;
        }, 0);
      }
    )
  }

  constructor() {
  }

  // ngAfterContentInit(): void {
  //   // !!!
  //   if (!this.snapID) {
  //     console.log("!!! stacks: short circuit");
  //     return;
  //   }
  //
  //   // this.appCoreService.appCore.onPublish((appCore: AppCore) => {
  //   //   if (this.dataSeriesQueryDir === undefined) {
  //   //     appCore.err(new ConfigurationError(`stacks is missing required 'data-series' child.`)
  //   //       .from(SOURCE)
  //   //       .at(Severity.ERROR));
  //   //     return;
  //   //   }
  //   //   let dataSeriesQuery = this.dataSeriesQueryDir?.dataSeriesQuery;
  //   //
  //   //   // Handle new data series.
  //   //   dataSeriesQuery?.response
  //   //     .pipe(takeUntil(this.unsubscribe))
  //   //     .subscribe((response: ResponseNode) => {
  //   //       this.numTotalGoroutines = response.properties.expectNumber('num_total_goroutines');
  //   //       this.numFilteredGoroutines = response.properties.expectNumber('num_filtered_goroutines');
  //   //       this.numBuckets = response.children[0].properties.expectNumber('num_buckets');
  //   //       this.aggStacks = response.children[0].children;
  //   //       this.rawStacks = response.children[1].children;
  //   //       this.numStacks = this.rawStacks.length;
  //   //     })
  //   // });
  //
  //   type Args = {
  //     colID: number,
  //     snapID: number,
  //     gID: number | undefined,
  //   };
  //   const args: Args = {colID: this.colID, snapID: this.snapID, gID: undefined};
  //   const urlParts = document.URL.split('#');
  //   if (urlParts.length > 1) {
  //     if (urlParts[1].startsWith('g_')) {
  //       args.gID = Number(urlParts[1].slice(2));
  //       console.log("filtering for goroutine: ", args.gID);
  //     }
  //   }
  //
  //   console.log(`calling GetGoroutines with args:`, args);
  //   this.getGoroutinesQuery
  //     .fetch(args)
  //     .subscribe(
  //       res => {
  //         this.goroutines = res.data.goroutines;
  //         this.mustScroll = true;
  //         console.log("got goroutines:", this.goroutines.length)
  //       },
  //       error => {
  //         console.log("!!! error getting Goroutines: ", error);
  //       })
  // }

  // ngOnDestroy(): void {
  //   this.unsubscribe.next();
  //   this.unsubscribe.complete();
  // }

  // ngAfterViewChecked(): void {
  //   if (!this.mustScroll) {
  //     return
  //   }
  //   const urlParts = document.URL.split('#');
  //   if (urlParts.length > 1) {
  //     document.getElementById(urlParts[1])?.scrollIntoView({behavior: "smooth"})
  //     console.log('!!! scrolled to ', urlParts[1])
  //     this.mustScroll = false;
  //   }
  // }
}
