import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatPaginatorModule } from '@angular/material/paginator';
import {
  CoreModule,
  DataTableModule,
  TextFieldModule,
  WeightedTreeModule
} from 'traceviz/dist/ngx-traceviz-lib';
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

@NgModule({
  declarations: [
    SnapshotComponent,
    TypeInfoComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    CoreModule,
    ComponentsModule,
    DataTableModule,
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
    TextFieldModule,
    WeightedTreeModule,
  ],
  providers: [],
})
export class SnapshotModule {
}
