import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatPaginatorModule } from '@angular/material/paginator';
import { CoreModule, DataTableModule, WeightedTreeModule, TextFieldModule } from 'traceviz/dist/ngx-traceviz-lib';
import { AppComponent } from './app.component';
import { ComponentsModule } from "./components/components.module";
import { StacksComponent } from "./components/stacks/stacks.component";
import { GraphQLModule } from './graphql/graphql.module';
import { RouterModule, Routes } from '@angular/router';
import { CollectionsListComponent } from "./components/collections-list/collections-list.component";
import { SnapshotModule } from "./components/snapshot/snapshot.module";
import { SnapshotComponent } from "./components/snapshot/snapshot.component";

const routes: Routes = [
  { path: '', redirectTo: '/collections', pathMatch: 'full'},
  { path: 'collections', component: CollectionsListComponent },
  { path: 'snapshot', component: SnapshotComponent},
];

@NgModule({
  declarations: [
    AppComponent,
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
    RouterModule.forRoot(routes),
    StacksComponent,
    SnapshotModule,
    TextFieldModule,
    WeightedTreeModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
