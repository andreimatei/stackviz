import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from "@angular/router";
import {
  AllCollectionsGQL,
  Collection,
  CollectServiceSnapshotsGQL
} from "src/app/graphql/graphql-codegen-generated";
import { MatButtonModule } from "@angular/material/button";
import { QueryRef } from "apollo-angular";
import { map, Observable, tap } from "rxjs";

@Component({
  selector: 'collections-list',
  standalone: true,
  imports: [CommonModule, RouterModule, MatButtonModule],
  templateUrl: './collections-list.component.html',
  styles: []
})
export class CollectionsListComponent implements OnInit {
  protected loading: boolean;
  protected collections?: Collection[];
  private collectionsQueryInstance?: QueryRef<any, any>;
  protected querySubscription?: Observable<Collection[]>;

  constructor(
    private readonly collectionsQuery: AllCollectionsGQL,
    private readonly newCollectionQuery: CollectServiceSnapshotsGQL,
  ) {
    this.loading = false;
  }

  ngOnInit() {
    const queryInstance = this.collectionsQuery.watch()
    this.collectionsQueryInstance = queryInstance
    this.querySubscription = queryInstance.valueChanges.pipe(
      map(result => result.data.collections),
      tap(value => console.log("!!! got collections: ", value))
    )
  }

  collectSnapshots(): void {
    console.log("!!! collecting snapshots");
    this.newCollectionQuery.mutate().subscribe({
      next: (res) => {
        this.collectionsQueryInstance!.refetch().then(
          undefined,  // onfulfilled
          reason => console.log("failed to refresh collections: ", reason),
        );
      },
      error: (error) => console.log("failed to create collection: ", error)
    })
  }
}
