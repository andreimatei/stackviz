import { Component, Input } from "@angular/core";
import { CommonModule } from "@angular/common";
import { MatExpansionModule } from "@angular/material/expansion";
import { VarInfo } from "src/app/components/flamegraph/flamegraph.component";
import { VarListComponent } from "src/app/components/var-list/var-list.component";
import { MatCardModule } from "@angular/material/card";
import { RouterLink } from "@angular/router";

@Component({
  selector: 'app-frame-data',
  standalone: true,
  imports: [CommonModule, MatExpansionModule, VarListComponent, MatCardModule, RouterLink],
  templateUrl: './frame-data.html',
})
export class FrameDataComponent {
  @Input({required: true}) data!: Map<number, VarInfo[]>;
  @Input({required: true}) snapshotID!: number;
  @Input({required: true}) collectionID!: number;
}
