import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from "@angular/router";
import {
  AllCollectionsGQL,
  CollectCollectionGQL,
  Collection
} from "src/app/graphql/graphql-codegen-generated";
import { MatButtonModule } from "@angular/material/button";
import { QueryRef } from "apollo-angular";
import { map, tap, Observable } from "rxjs";

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
  protected querySubscription?: Observable<Array<{ id: string, name: string }>>;

  constructor(
    private readonly collectionsQuery: AllCollectionsGQL,
    private readonly newCollectionQuery: CollectCollectionGQL,
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
    console.log("!!! CollectCollection");
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
