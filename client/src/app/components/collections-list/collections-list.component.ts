import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from "@angular/router";
import { AllCollectionsGQL, Collection } from "../../graphql/graphql-codegen-generated";

@Component({
  selector: 'collections-list',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './collections-list.component.html',
  styles: []
})
export class CollectionsListComponent {
  protected collections?: Collection[];

  constructor(collectionsQuery: AllCollectionsGQL) {
    collectionsQuery.fetch().subscribe(results => {
      this.collections = results.data.collections;
    });
  }

}
