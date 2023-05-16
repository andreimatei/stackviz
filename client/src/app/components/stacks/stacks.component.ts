import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ComponentRef,
  ContentChild
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  AppCoreService,
  DataSeriesQueryDirective,
  DataTableModule,
  InteractionsDirective
} from 'traceviz/dist/ngx-traceviz-lib';
import { AppCore, ConfigurationError, ResponseNode, Severity } from 'traceviz-client-core';
import { Subject } from "rxjs";
import { takeUntil } from 'rxjs/operators';

const SOURCE = 'data-table';

@Component({
  selector: 'stacks',
  standalone: true,
  imports: [CommonModule, DataTableModule],
  template: `
    <div>
      {{ numStacks }} stacks ({{ numFilteredGoroutines }} filtered / {{ numTotalGoroutines }} total Goroutines)
      <ul>
        <li *ngFor="let stack of stacks">
            <data-table [data]="stack" ></data-table>
        </li>
      </ul>
    </div>
  `,
  styles: [`ul {
      list-style-type: none; /* Remove bullets */
      padding: 0; /* Remove padding */
      margin: 0; /* Remove margins */
    }
    li {
        border-bottom: 1px solid black;
    }`
  ]
})
export class StacksComponent {
  @ContentChild(DataSeriesQueryDirective) dataSeriesQueryDir?: DataSeriesQueryDirective;
  @ContentChild(InteractionsDirective) interactionsDir?: InteractionsDirective;

  private unsubscribe = new Subject<void>();
  protected stacks?: ResponseNode[];
  protected numStacks?: number;
  protected numFilteredGoroutines?: number;
  protected numTotalGoroutines?: number;

  constructor(
    private readonly appCoreService: AppCoreService,
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
          this.numStacks = response.children.length;
          this.numTotalGoroutines = response.properties.expectNumber('num_total_goroutines');
          this.numFilteredGoroutines = response.properties.expectNumber('num_filtered_goroutines');
          this.stacks = response.children;
        })
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe.next();
    this.unsubscribe.complete();
  }
}
