import { Component, ContentChild } from '@angular/core';
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

  constructor(private readonly appCoreService: AppCoreService) {
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
          this.stacks = response.children;
        })
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe.next();
    this.unsubscribe.complete();
  }
}
