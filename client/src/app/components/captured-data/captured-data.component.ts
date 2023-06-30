import { Component, Input, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatExpansionModule } from "@angular/material/expansion";
import { combineLatest, map, Observable, startWith, tap } from "rxjs";
import { CollectedVar, FrameSpec } from "../../graphql/graphql-codegen-generated";
import { MatTableModule } from "@angular/material/table";
import { RouterLink } from "@angular/router";
import { MatChipListbox, MatChipsModule } from "@angular/material/chips";

@Component({
  selector: 'app-captured-data',
  standalone: true,
  imports: [CommonModule, MatExpansionModule, MatTableModule, RouterLink, MatChipsModule],
  template: `
    <mat-chip-listbox #expressionChips multiple>
      <ng-container *ngFor="let frameSpec of (collectionSpec$ | async)">
        <mat-chip-option
          selected
          id={{expr}}
          *ngFor="let expr of frameSpec.exprs">
          {{expr}}
        </mat-chip-option>
      </ng-container>
    </mat-chip-listbox>

    <mat-expansion-panel expanded *ngFor="let goroutine of (filteredData$ | async)">
      <mat-expansion-panel-header>
        <mat-panel-title class="header">
          Data for goroutine {{goroutine.gid}}
        </mat-panel-title>
      </mat-expansion-panel-header>
      <table mat-table [dataSource]="goroutine.vars">
        <ng-container matColumnDef="expr">
          <th mat-header-cell *matHeaderCellDef>Expression</th>
          <td mat-cell *matCellDef="let v">
            {{v.Expr}}
          </td>
        </ng-container>
        <ng-container matColumnDef="val">
          <th mat-header-cell *matHeaderCellDef>Value</th>
          <td mat-cell *matCellDef="let v">
            {{v.Value}}
          </td>
        </ng-container>
        <ng-container matColumnDef="links">
          <th mat-header-cell *matHeaderCellDef>Links</th>
          <td mat-cell *matCellDef="let v">
            <div *ngIf="v.Links?.length > 0">
              <ul>
                <li *ngFor="let l of v.Links">
                  <a
                    [routerLink]="['/collections', collectionID, 'snap',l.SnapshotID]"
                    [queryParams]="{filter: 'gid=' + l.GoroutineID}"
                  >
                    Snapshot: {{l.SnapshotID}} Goroutine: {{l.GoroutineID}}
                  </a>
                </li>
              </ul>
            </div>
          </td>
        </ng-container>
        <tr mat-header-row *matHeaderRowDef="['expr', 'val', 'links']"></tr>
        <tr mat-row *matRowDef="let rowData; columns: ['expr','val', 'links']"></tr>
      </table>
    </mat-expansion-panel>
  `,
  styleUrls: ['./captured-data.component.css']
})
export class CapturedDataComponent {
  @Input({required: true}) data$!: Observable<GoroutineData[]>;
  @Input({required: true}) collectionSpec$!: Observable<Partial<FrameSpec>[]>;
  @ViewChild('expressionChips') expressionChips!: MatChipListbox;
  protected filteredData$!: Observable<GoroutineData[]>;
  @Input({required: true}) collectionID!: number;

  constructor() {
  }

  ngAfterViewInit() {
    // Note: I'm doing Promise.resolve().then() to defer the execution after the
    // change detection cycle finished. Otherwise, I get the error about the
    // template changing after change detection. We need to set up this pipe in
    // ngAfterViewInit, though, because we need this.expressionChips to be
    // populated.
    Promise.resolve(null).then(() => {
    // Filter the captured vars whenever we load new data or when the chip
    // selection changes.
    this.filteredData$ =
      combineLatest([
        this.data$,
        // Provide an initial value because the chips listbox does not provide
        // an event for the initial load.
        this.expressionChips.chipSelectionChanges.pipe(startWith(undefined)),
      ]).pipe(
        map(([data, selection]) => {
          if (selection == undefined) {
            return data;
          }
          const selected = this.expressionChips.value as string[];
          return data.map(gd => {
            let cp = {...gd};
            cp.vars = cp.vars.map(v => ({...v}));
            cp.vars = cp.vars.filter(
              (v: CollectedVar) => selected.includes(v.Expr))
            return cp;
          }).filter(gd => gd.vars.length > 0);
        }),
      );
    })
  }
}

export interface GoroutineData {
  gid: number,
  vars: CollectedVar[],
}
