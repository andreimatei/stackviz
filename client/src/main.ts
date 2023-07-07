import { AppComponent } from './app/app.component';
import { SnapshotComponent } from './app/components/snapshot/snapshot.component';
import { map } from 'rxjs';
import { GetCollectionGQL } from 'src/app/graphql/graphql-codegen-generated';
import { importProvidersFrom, inject } from '@angular/core';
import {
  CollectionsListComponent
} from './app/components/collections-list/collections-list.component';
import {
  ActivatedRouteSnapshot,
  provideRouter,
  Router,
  RouterStateSnapshot,
  Routes,
  withComponentInputBinding
} from '@angular/router';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatCardModule } from '@angular/material/card';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { GraphQLModule } from './app/graphql/graphql.module';
import { provideAnimations } from '@angular/platform-browser/animations';
import { bootstrapApplication, BrowserModule } from '@angular/platform-browser';

const routeToFirstSnapshot = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  const colID = Number(route.params['colID']);
  const query = inject(GetCollectionGQL);
  const router = inject(Router);

  return query.fetch({colID: colID}).pipe(map(res => {
    const snap = res.data.collectionByID?.processSnapshots?.[0];
    if (!snap) {
      console.log("error: failed to get collection or first snapshot");
      return false
    }
    return router.createUrlTree(['collections/' + colID + '/snap/' + snap?.id]);
  }))
};

const routes: Routes = [
  {path: 'collections', component: CollectionsListComponent},
  {
    path: 'collections/:colID',
    canActivate: [routeToFirstSnapshot],
    component: CollectionsListComponent,
  },
  {path: 'collections/:colID/snap/:snapID', component: SnapshotComponent},
  {path: '', redirectTo: '/collections', pathMatch: 'full'},
];

bootstrapApplication(AppComponent, {
  providers: [
    importProvidersFrom(BrowserModule, GraphQLModule, MatCardModule, MatPaginatorModule),
    provideAnimations(),
    provideHttpClient(withInterceptorsFromDi()),
    provideRouter(routes, withComponentInputBinding())
  ]
})
  .catch(err => console.error(err));
