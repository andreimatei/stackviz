import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTabsModule, } from '@angular/material/tabs'
import { Observable } from "rxjs";
import {
  CollectedVar,
  FrameInfo,
  GoroutineInfo,
  SnapshotInfo
} from "../../graphql/graphql-codegen-generated";
import { MatTableModule } from "@angular/material/table";
import { BacktraceComponent } from "../backtrace/backtrace.component";

@Component({
  selector: 'app-stacks',
  standalone: true,
  imports: [CommonModule, MatTabsModule, MatTableModule, BacktraceComponent],
  template: `
    <div>
      {{ goroutines?.length }} goroutines
      <!-- TODO(andrei): retrieve the total/filtered number from the server -->
      <!--({{ numFilteredGoroutines }} filtered / {{ numTotalGoroutines }} total Goroutines),-->
      {{ goroutineGroups?.length }} buckets
      <hr>

      <mat-tab-group selectedIndex="0">
        <mat-tab label="Aggregated">
          <ng-template matTabContent>
            <ul>
              <li *ngFor="let group of goroutineGroups">
                {{group.goroutineIDs.length}} goroutines in goroutine group
                Goroutine IDs {{group.goroutineIDs}}
                <app-backtrace [data]="{vars: group.vars, frames: group.frames}"
                               [collectionID]="colID"></app-backtrace>
              </li>
            </ul>
          </ng-template>
        </mat-tab>

        <mat-tab label="Raw">
          <ul>
            <li *ngFor="let g of goroutines" id="g_{{g.ID}}">
              <a id="g_{{g.ID}}">Goroutine {{ g.ID }}</a>

              <div *ngIf="g.Data.FlightRecorderData">
                Recorded data:
                <ul>
                  <li *ngFor="let recData of g.Data.FlightRecorderData">
                    {{ recData }}
                  </li>
                </ul>
              </div>

              <app-backtrace [data]="{vars: g.Data.Vars, frames: g.Frames}"
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
  @Input({required: true}) colID!: number;

  protected goroutines?: GoroutineInfo[];
  protected goroutineGroups?: groupInfo[];

  @Input() set data$(val$: Observable<SnapshotInfo>) {
    val$.subscribe(
      si => {
        console.log("stacks got value: ", si);
        this.goroutineGroups = si.Aggregated.map(
          group => ({
            goroutineIDs: group.IDs,
            frames: group.Frames,
            // Flatten the variables in each group.
            vars: group.Data.reduce((prev, v) => prev.concat(v.Vars), [] as CollectedVar[])
          })
        )
        // Updating all the tables can take a few seconds, so in order to make
        // the user experience better, we first clear them, and only then create
        // the new ones.
        this.goroutines = [];
        setTimeout(() => {
          this.goroutines = si.Raw;
        }, 0);
      }
    )
  }

  constructor() {
  }
}

interface groupInfo {
  goroutineIDs: number[];
  vars: CollectedVar[];
  frames: FrameInfo[];
}
