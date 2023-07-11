import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTabsModule, } from '@angular/material/tabs'
import { Observable } from "rxjs";
import {
  GoroutineInfo,
  GoroutinesGroup,
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
          <ul>
            <li *ngFor="let g of goroutines" id="g_{{g.goroutineInfo.ID}}">
              <a id="g_{{g.goroutineInfo.ID}}">Goroutine {{ g.goroutineInfo.ID }}</a>

              <div *ngIf="g.flightRecorderData">
                Recorded data:
                <ul>
                  <li *ngFor="let recData of g.flightRecorderData">
                    {{ recData }}
                  </li>
                </ul>
              </div>

              <app-backtrace [vars]="g.goroutineInfo.Vars" [frames]="g.goroutineInfo.Frames"
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
  // !!! @Input({required: true}) snapID: number | null = null;

  protected goroutines?: goroutineInfoWithFlightRecorderData[];
  protected goroutineGroups?: GoroutinesGroup[];

  @Input() set data$(val$: Observable<SnapshotInfo>) {
    val$.subscribe(
      si => {
        console.log("stacks got value: ", si);
        this.goroutineGroups = si.Aggregated;
        // Updating all the tables can take a few seconds, so in order to make
        // the user experience better, we first clear them, and only then create
        // the new ones.
        this.goroutines = [];
        setTimeout(() => {
          this.goroutines = si.Raw.map(
            g => {
              let gData: string[] | undefined;
              if (si.FlightRecorderData) {
                gData = si.FlightRecorderData[g.ID.toString()];
              }
              const flightRecorderData = new Map(Object.entries(si.FlightRecorderData));
              return {
                goroutineInfo: g,
                flightRecorderData: gData,
              };
            }
          );
        }, 0);
      }
    )
  }

  constructor() {
  }
}

interface goroutineInfoWithFlightRecorderData {
  goroutineInfo: GoroutineInfo;
  flightRecorderData?: string[];
}
