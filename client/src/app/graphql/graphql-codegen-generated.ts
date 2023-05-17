import { gql } from 'apollo-angular';
import { Injectable } from '@angular/core';
import * as Apollo from 'apollo-angular';
export const AllCollectionsDocument = gql`
    query AllCollections {
  collections {
    id
    name
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class AllCollectionsGQL extends Apollo.Query<AllCollectionsQuery, AllCollectionsQueryVariables> {
    document = AllCollectionsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }