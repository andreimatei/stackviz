import { Component, ContentChild, Input, ViewEncapsulation } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  AppCoreService,
  DataSeriesQueryDirective,
  DataTableModule,
  InteractionsDirective
} from 'traceviz/dist/ngx-traceviz-lib';
import { MatTabsModule, } from '@angular/material/tabs'
import { AppCore, ConfigurationError, ResponseNode, Severity } from 'traceviz-client-core';
import { Subject } from "rxjs";
import { takeUntil } from 'rxjs/operators';
import {
  GetGoroutinesGQL,
  GetGoroutinesQuery,
  GetGoroutinesQueryVariables, GoroutineInfo
} from "../../graphql/graphql-codegen-generated";

const SOURCE = 'data-table';

@Component({
  selector: 'app-stacks',
  standalone: true,
  imports: [CommonModule, DataTableModule, MatTabsModule],
  template: `
    <div>
      {{ numStacks }} stacks
      ({{ numFilteredGoroutines }} filtered / {{ numTotalGoroutines }} total Goroutines),
      {{ numBuckets }} buckets
      <hr>

      <mat-tab-group>
        <mat-tab label="Aggregated">
          <ul>
            <li *ngFor="let stack of aggStacks">
              {{ stack.properties.expectNumber("num_gs_in_bucket") }} goroutine(s):
              <data-table [data]="stack" [style]="  "></data-table>
            </li>
          </ul>
        </mat-tab>
        <mat-tab label="Raw">
          <ng-template matTabContent>
            <ul>
              <li *ngFor="let stack of rawStacks">
                Goroutine ID {{ stack.properties.expectNumber("g_id") }}:
                <data-table [data]="stack" ></data-table>
              </li>
            </ul>
          </ng-template>
         </mat-tab>
        <mat-tab label="Raw 2">
          <ul>
            <li *ngFor="let g of goroutines">
              Goroutine {{ g.ID }}
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
    data-table::part(table) {background: #f4f4f4; border: 1px solid #dcdcdc; width: unset; min-width: 350px;}
    `
  ]
})
export class StacksComponent {
  @ContentChild(DataSeriesQueryDirective) dataSeriesQueryDir?: DataSeriesQueryDirective;
  @ContentChild(InteractionsDirective) interactionsDir?: InteractionsDirective;

  // TODO(andrei): mark these as required.
  @Input() colID!: number
  @Input() snapID!: number

  private unsubscribe = new Subject<void>();
  protected rawStacks?: ResponseNode[];
  protected aggStacks?: ResponseNode[];
  protected numStacks?: number;
  protected numBuckets?: number;
  protected numFilteredGoroutines?: number;
  protected numTotalGoroutines?: number;

  protected goroutines!: GoroutineInfo[];

  constructor(
    private readonly appCoreService: AppCoreService,
    private readonly getGoroutinesQuery: GetGoroutinesGQL,
  ) {
  }

  ngAfterContentInit(): void {
    this.appCoreService.appCore.onPublish((appCore: AppCore) => {
      if (this.dataSeriesQueryDir === undefined) {
        appCore.err(new ConfigurationError(`stacks is missing required 'data-series' child.`)
          .from(SOURCE)
          .at(Severity.ERROR));
        return;
      }
      let dataSeriesQuery = this.dataSeriesQueryDir?.dataSeriesQuery;

      // Handle new data series.
      dataSeriesQuery?.response
        .pipe(takeUntil(this.unsubscribe))
        .subscribe((response: ResponseNode) => {
          this.numTotalGoroutines = response.properties.expectNumber('num_total_goroutines');
          this.numFilteredGoroutines = response.properties.expectNumber('num_filtered_goroutines');
          this.numBuckets = response.children[0].properties.expectNumber('num_buckets');
          this.aggStacks = response.children[0].children;
          this.rawStacks = response.children[1].children;
          this.numStacks = this.rawStacks.length;
        })
    });

    this.getGoroutinesQuery
      .fetch({colID: this.colID, snapID: this.snapID})
      .subscribe(res => {
        this.goroutines = res.data.goroutines;
        console.log("got goroutines:", this.goroutines.length)
      })
  }

  ngOnDestroy(): void {
    this.unsubscribe.next();
    this.unsubscribe.complete();
  }
}
