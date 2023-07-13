import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from "@angular/material/table";
import { CollectedVar, FrameInfo } from "../../graphql/graphql-codegen-generated";
import { RouterLink } from "@angular/router";
import { VarListComponent } from "src/app/components/var-list/var-list.component";
import { MatSlideToggleChange, MatSlideToggleModule } from "@angular/material/slide-toggle";
import { MatCardModule } from "@angular/material/card";

@Component({
  selector: 'app-backtrace',
  standalone: true,
  imports: [CommonModule, MatTableModule, RouterLink, VarListComponent, MatSlideToggleModule, MatCardModule],
  template: `
    <mat-card>
      <mat-card-header>
        <mat-slide-toggle #filterStack [checked]="false">
          Only stack frames with collected data
        </mat-slide-toggle>
      </mat-card-header>
      <mat-card-content>
        <table mat-table [dataSource]="filterStack.checked ? filteredFrames : frames">
          <ng-container matColumnDef="index">
            <th mat-header-cell *matHeaderCellDef></th>
            <td mat-cell class="cell-index" *matCellDef="let rIndex = index;"> {{rIndex}}</td>
          </ng-container>
          <ng-container matColumnDef="func">
            <th mat-header-cell *matHeaderCellDef>Backtrace</th>
            <td mat-cell *matCellDef="let frame">{{frame.Func}}</td>
          </ng-container>
          <ng-container matColumnDef="location">
            <th mat-header-cell *matHeaderCellDef>file:line</th>
            <td mat-cell *matCellDef="let frame">{{frame.File}}:{{frame.Line}}</td>
          </ng-container>
          <ng-container matColumnDef="vars">
            <th mat-header-cell *matHeaderCellDef>Data</th>
            <td mat-cell class="cell-index" *matCellDef="let frame">
              <app-var-list [collectionID]="collectionID" [vars]="frame.Vars"></app-var-list>
            </td>
          </ng-container>
          <tr mat-header-row *matHeaderRowDef="['index','func','location','vars']"></tr>
          <tr mat-row *matRowDef="let rowData; columns: ['index','func','location','vars']"></tr>
        </table>
        </mat-card-content>
    </mat-card>
  `,
  styleUrls: ['./backtrace.component.css']
})
export class BacktraceComponent {
  frames!: Frame[];
  filteredFrames!: Frame[];

  // collectionID identifies the collection that the stacktrace is part of. Used
  // to generate the correct links.
  @Input({required: true}) collectionID!: number;

  @Input({required: true}) set data(val: input) {
    // join the vars with the frames
    this.frames = val.frames.map((frame, idx) =>
      ({
        Func: frame.Func,
        File: frame.File,
        Line: frame.Line,
        Vars: val.vars.filter(v => v.FrameIdx == idx)
      })
    )
    this.filteredFrames = this.frames.filter(f => f.Vars.length > 0);
  }
}

interface input {
  vars: CollectedVar[];
  frames: FrameInfo[];
}

interface Frame {
  Func: string;
  File: string;
  Line: number;
  Vars: CollectedVar[];
}
