import { Component, Input } from '@angular/core';
import { CollectedVar, } from "../../graphql/graphql-codegen-generated";
import { RouterLink } from "@angular/router";
import { CommonModule } from "@angular/common";

// CapturedDataGoroutineComponent displays the collected variables for one
// goroutine.
@Component({
  selector: 'app-var-list',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <ul>
      <li *ngFor="let v of vars">
        {{v.Expr}} : {{v.Value}}
        <div *ngIf="v.Links.length > 0">
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
      </li>
    </ul>
  `,
})
export class VarListComponent {
  @Input({required: true}) collectionID!: number;
  @Input({required: true}) vars!: CollectedVar[];
}
