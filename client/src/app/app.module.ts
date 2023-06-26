import { inject, NgModule } from '@angular/core';
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
import {
  ActivatedRouteSnapshot,
  Router,
  RouterModule,
  RouterStateSnapshot,
  Routes
} from '@angular/router';
import { CollectionsListComponent } from "./components/collections-list/collections-list.component";
import { SnapshotModule } from "./components/snapshot/snapshot.module";
import { SnapshotComponent } from "./components/snapshot/snapshot.component";
import { GetCollectionGQL } from "src/app/graphql/graphql-codegen-generated";
import { map } from "rxjs";

const routeToFirstSnapshot = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  const colID = Number(route.params['colID']);
  const query =  inject(GetCollectionGQL);
  const router = inject(Router);

  return query.fetch({colID: colID}).pipe(map(res => {
    const snap = res.data.collectionByID?.processSnapshots?.[0];
    if (!snap) {
      console.log("error: failed to get collection or first snapshot");
      return false
    }
    return router.createUrlTree(['collections/'+ colID + '/snap/' + snap?.id]);
  }))
};

const routes: Routes = [
  { path: 'collections', component: CollectionsListComponent },
  {
    path: 'collections/:colID',
    canActivate: [routeToFirstSnapshot],
    component: CollectionsListComponent,
  },
  { path: 'collections/:colID/snap/:snapID', component: SnapshotComponent},
  { path: '', redirectTo: '/collections', pathMatch: 'full'},
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
    RouterModule.forRoot(routes, { enableTracing: false }),
    StacksComponent,
    SnapshotModule,
    TextFieldModule,
    WeightedTreeModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
