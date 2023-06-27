import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DataTableModule } from 'traceviz/dist/ngx-traceviz-lib';
import { MatTabsModule, } from '@angular/material/tabs'
import { Observable } from "rxjs";
import { GoroutineInfo } from "../../graphql/graphql-codegen-generated";
import { MatTableModule } from "@angular/material/table";

@Component({
  selector: 'app-stacks',
  standalone: true,
  imports: [CommonModule, DataTableModule, MatTabsModule, MatTableModule],
  template: `
    <div>
<!--      {{ numStacks }} stacks-->
<!--      ({{ numFilteredGoroutines }} filtered / {{ numTotalGoroutines }} total Goroutines),-->
<!--      {{ numBuckets }} buckets-->
      <hr>

      <mat-tab-group selectedIndex="2">
<!--        <mat-tab label="Aggregated">-->
<!--          <ul>-->
<!--            <li *ngFor="let stack of aggStacks">-->
<!--              {{ stack.properties.expectNumber("num_gs_in_bucket") }} goroutine(s):-->
<!--              &lt;!&ndash;<data-table [data]="stack" [style]="  "></data-table>&ndash;&gt;-->
<!--            </li>-->
<!--          </ul>-->
<!--        </mat-tab>-->
<!--        <mat-tab label="Raw">-->
<!--          <ng-template matTabContent>-->
<!--            <ul>-->
<!--              <li *ngFor="let stack of rawStacks">-->
<!--                Goroutine ID {{ stack.properties.expectNumber("g_id") }}:-->
<!--                &lt;!&ndash;<data-table [data]="stack" ></data-table>&ndash;&gt;-->
<!--              </li>-->
<!--            </ul>-->
<!--          </ng-template>-->
<!--         </mat-tab>-->
        <mat-tab label="Raw">
          #goroutines: {{goroutines ? goroutines.length : 0}}
          <ul>
            <li *ngFor="let g of goroutines" id="g_{{g.ID}}">
              <a id="g_{{g.ID}}">Goroutine {{ g.ID }}</a>

              <table mat-table [dataSource]="g.Vars" *ngIf="g.Vars && g.Vars.length > 0">
                <ng-container matColumnDef="vars">
                  <th mat-header-cell *matHeaderCellDef> Vars </th>
                  <td mat-cell *matCellDef="let v">
                    {{v.Value}}
                    <div *ngIf="v.Links?.length > 0">
                      Links:
                      <ul>
                        <li *ngFor="let l of v.Links">
<!--                          <a href="/collections/{{colID}}/snap/{{l.SnapshotID}}#g_{{l.GoroutineID}}">-->
                            Snapshot: {{l.SnapshotID}} Goroutine: {{l.GoroutineID}}
<!--                          </a>-->
                        </li>
                      </ul>
                    </div>
                  </td>
                </ng-container>
                <tr mat-header-row *matHeaderRowDef="['vars']"></tr>
                <tr mat-row *matRowDef="let rowData; columns: ['vars']"></tr>
               </table>

              <table mat-table [dataSource]="g.Frames">
                <ng-container matColumnDef="frames">
                  <th mat-header-cell *matHeaderCellDef> Backtrace </th>
                  <td mat-cell *matCellDef="let frame"> {{frame}} </td>
                </ng-container>
                <tr mat-header-row *matHeaderRowDef="['frames']"></tr>
                <tr mat-row *matRowDef="let rowData; columns: ['frames']"></tr>
              </table>
            </li>
          </ul>
        </mat-tab>
      </mat-tab-group>

    <a id="xxx"></a>
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
    data-table::part(table) {background: #f4f4f4; border: 1px solid #dcdcdc; width: unset; min-width: 350px;}

    th.mat-mdc-header-cell {
      background-color: gray;
    }
    tr.mat-mdc-row {
      height: 20px;
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

  protected goroutines!: GoroutineInfo[];

  // !!! backtracecols = ['frames'];

  @Input() set data$(val$: Observable<GoroutineInfo[]>) {
    val$.subscribe(
      value => {
        // Updating all the tables can take a few seconds, so in order to make
        // the user experience better, we first clear them, and only then create
        // the new ones.
        console.log("stacks got value: ", value);
        this.goroutines = [];
        setTimeout(() => {
          this.goroutines = value;
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
