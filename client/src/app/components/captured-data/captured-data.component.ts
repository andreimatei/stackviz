import { Component, Input, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatExpansionModule } from "@angular/material/expansion";
import { combineLatest, map, Observable, startWith, tap } from "rxjs";
import { CollectedVar, FrameInfo, FrameSpec } from "../../graphql/graphql-codegen-generated";
import { MatTableModule } from "@angular/material/table";
import { RouterLink } from "@angular/router";
import { MatChipListbox, MatChipsModule } from "@angular/material/chips";
import {
  CapturedDataGoroutineComponent
} from "src/app/components/captured-data-goroutine/captured-data-goroutine.component";

// CapturedDataComponent is a component that displays the captured data for a
// list of goroutines.
@Component({
  selector: 'app-captured-data',
  standalone: true,
  imports: [CommonModule, MatExpansionModule, MatTableModule, RouterLink, MatChipsModule, CapturedDataGoroutineComponent],
  template: `
    <mat-chip-listbox #expressionChips multiple>
      <ng-container *ngFor="let frameSpec of (frameSpecs$ | async)">
        <mat-chip-option
            selected
            id={{expr}}
            *ngFor="let expr of frameSpec.collectExpressions">
          {{expr}}
        </mat-chip-option>
      </ng-container>
    </mat-chip-listbox>

    <app-captured-data-goroutine *ngFor="let goroutine of (filteredData$ | async)"
                                 [collectionID]="collectionID"
                                 [snapshotID]="snapshotID"
                                 [goroutineID]="goroutine.gid"
                                 [data]="goroutine"
    ></app-captured-data-goroutine>
  `,
  styleUrls: ['./captured-data.component.css']
})
export class CapturedDataComponent {
  @Input({required: true}) data$!: Observable<GoroutineData[]>;
  @Input({required: true}) frameSpecs$!: Observable<Partial<FrameSpec>[]>;
  @Input({required: true}) collectionID!: number;
  @Input({required: true}) snapshotID!: number;

  @ViewChild('expressionChips') expressionChips!: MatChipListbox;
  // filteredData$ filters the data$ for the goroutines that have some captured
  // data.
  protected filteredData$!: Observable<GoroutineData[]>;

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
  frames: FrameInfo[],
  vars: CollectedVar[],
}
