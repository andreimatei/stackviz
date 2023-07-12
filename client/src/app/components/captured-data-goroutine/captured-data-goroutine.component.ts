import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatExpansionModule } from "@angular/material/expansion";
import { CollectedVar, } from "../../graphql/graphql-codegen-generated";
import { MatTableModule } from "@angular/material/table";
import { RouterLink } from "@angular/router";
import { GoroutineData } from "src/app/components/captured-data/captured-data.component";
import { MatSlideToggleModule } from "@angular/material/slide-toggle";
import { VarListComponent } from "src/app/components/var-list/var-list.component";

// CapturedDataGoroutineComponent displays the collected variables for one
// goroutine.
@Component({
  selector: 'app-captured-data-goroutine',
  standalone: true,
  imports: [CommonModule, MatExpansionModule, MatTableModule, RouterLink, MatSlideToggleModule, VarListComponent],
  template: `
    <mat-expansion-panel expanded>
      <mat-expansion-panel-header>
        <mat-panel-title class="header">
          Captured data for&nbsp;
          <a
            [routerLink]="['/collections', collectionID, 'snap', snapshotID]"
            [queryParams]="{filter: 'gid=' + goroutineID}"
            click="$event.stopPropagation() // stop bubbling that colapses the map-panel"
          >goroutine {{goroutineID}}</a>
        </mat-panel-title>
      </mat-expansion-panel-header>
      <mat-slide-toggle #fullStack [checked]="false">
        Full goroutine stack
      </mat-slide-toggle>
      <table mat-table [dataSource]="fullStack.checked ? allFrames : framesWithVars">
        <ng-container matColumnDef="frameIdx">
          <th mat-header-cell *matHeaderCellDef>Function</th>
          <td mat-cell *matCellDef="let f">
            {{f.Func}}
          </td>
        </ng-container>
        <ng-container matColumnDef="vars">
          <th mat-header-cell *matHeaderCellDef>Vars</th>
          <td mat-cell *matCellDef="let f">
            <app-var-list [collectionID]="collectionID" [snapshotID]="snapshotID"
                          [vars]="f.Vars"></app-var-list>
          </td>
        </ng-container>
        <tr mat-header-row *matHeaderRowDef="['frameIdx', 'vars']"></tr>
        <tr mat-row *matRowDef="let rowData; columns: ['frameIdx', 'vars']"></tr>
      </table>
    </mat-expansion-panel>
  `,
  styleUrls: ['./captured-data-goroutine.component.css']
})
export class CapturedDataGoroutineComponent {
  @Input({required: true}) collectionID!: number;
  @Input({required: true}) snapshotID!: number;
  @Input({required: true}) goroutineID!: number;

  allFrames!: frameInfo[];
  framesWithVars!: frameInfo[];

  @Input({required: true}) set data(val: GoroutineData) {
    this.allFrames = val.frames.map((f, idx) => {
      return {
        File: f.File,
        Line: f.Line,
        Func: f.Func,
        Vars: val.vars.filter(v => v.FrameIdx === idx),
      }
    })
    this.framesWithVars = this.allFrames.filter(f => f.Vars.length > 0);
  }
}

interface frameInfo {
  File: string;
  Func: string;
  Line: number;
  Vars: CollectedVar[],
}
