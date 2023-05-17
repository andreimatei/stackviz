import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatPaginatorModule } from '@angular/material/paginator';
import { CoreModule, DataTableModule, WeightedTreeModule, TextFieldModule } from 'traceviz/dist/ngx-traceviz-lib';
import { SnapshotComponent } from './snapshot.component';
import { ComponentsModule } from "../components.module";
import { StacksComponent } from "../stacks/stacks.component";
import { GraphQLModule } from '../../graphql/graphql.module';
import { RouterModule } from "@angular/router";


@NgModule({
  declarations: [
    SnapshotComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    CoreModule,
    ComponentsModule,
    DataTableModule,
    GraphQLModule,
    HttpClientModule,
    MatCardModule,
    MatPaginatorModule,
    TextFieldModule,
    StacksComponent,
    RouterModule,
    WeightedTreeModule,
  ],
  providers: [],
})
export class SnapshotModule { }
