import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatPaginatorModule } from '@angular/material/paginator';
import { SnapshotComponent } from './snapshot.component';
import { ComponentsModule } from "../components.module";
import { StacksComponent } from "../stacks/stacks.component";
import { GraphQLModule } from '../../graphql/graphql.module';
import { RouterModule } from "@angular/router";
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from "@angular/material/input";
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatTreeModule } from "@angular/material/tree";
import { MatIconModule } from "@angular/material/icon";
import { MatButtonModule } from "@angular/material/button";
import { ResizableModule } from 'angular-resizable-element';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { TypeInfoComponent } from "./type-info.component";
import { FlamegraphComponent } from "../flamegraph/flamegraph.component";
import { FormsModule } from "@angular/forms";
import { MatTabsModule } from "@angular/material/tabs";
import { CapturedDataComponent } from "../captured-data/captured-data.component";


@NgModule({
  declarations: [
    SnapshotComponent,
    TypeInfoComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    ComponentsModule,
    GraphQLModule,
    HttpClientModule,
    MatButtonModule,
    MatCardModule,
    MatCheckboxModule,
    MatExpansionModule,
    MatIconModule,
    MatPaginatorModule,
    MatProgressBarModule,
    MatSelectModule,
    MatInputModule,
    MatSidenavModule,
    MatTreeModule,
    ResizableModule,
    RouterModule,
    StacksComponent,
    FlamegraphComponent,
    FormsModule,
    MatTabsModule,
    CapturedDataComponent,
  ],
  providers: [],
})
export class SnapshotModule {
}
