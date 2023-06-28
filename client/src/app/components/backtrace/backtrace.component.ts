import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from "@angular/material/table";
import { CollectedVar } from "../../graphql/graphql-codegen-generated";

@Component({
  selector: 'app-backtrace',
  standalone: true,
  imports: [CommonModule, MatTableModule],
  template: `
      <table mat-table [dataSource]="vars" *ngIf="vars && vars.length > 0">
          <ng-container matColumnDef="vars">
              <th mat-header-cell *matHeaderCellDef>Vars</th>
              <td mat-cell *matCellDef="let v">
                  {{v.Value}}
                  <div *ngIf="v.Links?.length > 0">
                      Links:
                      <ul>
                          <li *ngFor="let l of v.Links">
                              <a href="/collections/{{collectionID}}/snap/{{l.SnapshotID}}?filter={{encodeURIComponent('gid=')}}{{l.GoroutineID}}">
                              Snapshot: {{l.SnapshotID}} Goroutine: {{l.GoroutineID}}
                              </a>
                          </li>
                      </ul>
                  </div>
              </td>
          </ng-container>
          <tr mat-header-row *matHeaderRowDef="['vars']"></tr>
          <tr mat-row *matRowDef="let rowData; columns: ['vars']"></tr>
      </table>
      <table mat-table [dataSource]="frames">
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
          <tr mat-header-row *matHeaderRowDef="['index','func','location']"></tr>
          <tr mat-row *matRowDef="let rowData; columns: ['index','func','location']"></tr>
      </table>
  `,
  styleUrls: ['./backtrace.component.css']
})
export class BacktraceComponent {
  @Input({required: true}) vars!: CollectedVar[];
  @Input({required: true}) frames!: Frame[];
  // collectionID identifies the collection that the stacktrace is part of. Used
  // to generate the correct links.
  @Input({required: true}) collectionID!: number;
  protected readonly encodeURIComponent = encodeURIComponent;
}

interface Frame {
  Func: string;
  File: string;
  Line: number;
}
