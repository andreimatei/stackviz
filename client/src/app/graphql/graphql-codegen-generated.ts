import { gql } from 'apollo-angular';
import { Injectable } from '@angular/core';
import * as Apollo from 'apollo-angular';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  /**
   * Define a Relay Cursor type:
   * https://relay.dev/graphql/connections.htm#sec-Cursor
   */
  Cursor: any;
};

export type Collection = Node & {
  __typename?: 'Collection';
  id: Scalars['ID'];
  name: Scalars['String'];
  processSnapshots?: Maybe<Array<ProcessSnapshot>>;
};

/**
 * CreateCollectionInput is used for create Collection object.
 * Input was generated by ent.
 */
export type CreateCollectionInput = {
  name: Scalars['String'];
  processSnapshotIDs?: InputMaybe<Array<Scalars['ID']>>;
};

/**
 * CreateProcessSnapshotInput is used for create ProcessSnapshot object.
 * Input was generated by ent.
 */
export type CreateProcessSnapshotInput = {
  processID: Scalars['String'];
  snapshot: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  collectCollection?: Maybe<Collection>;
  createCollection?: Maybe<Collection>;
};


export type MutationCreateCollectionArgs = {
  input?: InputMaybe<CreateCollectionInput>;
};

/**
 * An object with an ID.
 * Follows the [Relay Global Object Identification Specification](https://relay.dev/graphql/objectidentification.htm)
 */
export type Node = {
  /** The id of the object. */
  id: Scalars['ID'];
};

/** Possible directions in which to order a list of items when provided an `orderBy` argument. */
export enum OrderDirection {
  /** Specifies an ascending order for a given `orderBy` argument. */
  Asc = 'ASC',
  /** Specifies a descending order for a given `orderBy` argument. */
  Desc = 'DESC'
}

/**
 * Information about pagination in a connection.
 * https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo
 */
export type PageInfo = {
  __typename?: 'PageInfo';
  /** When paginating forwards, the cursor to continue. */
  endCursor?: Maybe<Scalars['Cursor']>;
  /** When paginating forwards, are there more items? */
  hasNextPage: Scalars['Boolean'];
  /** When paginating backwards, are there more items? */
  hasPreviousPage: Scalars['Boolean'];
  /** When paginating backwards, the cursor to continue. */
  startCursor?: Maybe<Scalars['Cursor']>;
};

export type ProcessSnapshot = Node & {
  __typename?: 'ProcessSnapshot';
  id: Scalars['ID'];
  processID: Scalars['String'];
  snapshot: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  collectionByID?: Maybe<Collection>;
  collections: Array<Collection>;
  /** Fetches an object given its ID. */
  node?: Maybe<Node>;
  /** Lookup nodes by a list of IDs. */
  nodes: Array<Maybe<Node>>;
  processSnapshots: Array<ProcessSnapshot>;
};


export type QueryCollectionByIdArgs = {
  id: Scalars['ID'];
};


export type QueryNodeArgs = {
  id: Scalars['ID'];
};


export type QueryNodesArgs = {
  ids: Array<Scalars['ID']>;
};

export type AllCollectionsQueryVariables = Exact<{ [key: string]: never; }>;


export type AllCollectionsQuery = { __typename?: 'Query', collections: Array<{ __typename?: 'Collection', id: string, name: string }> };

export type GetCollectionQueryVariables = Exact<{
  colID: Scalars['ID'];
}>;


export type GetCollectionQuery = { __typename?: 'Query', collectionByID?: { __typename?: 'Collection', id: string, name: string, processSnapshots?: Array<{ __typename?: 'ProcessSnapshot', id: string, processID: string }> | null } | null };

export type CollectCollectionMutationVariables = Exact<{ [key: string]: never; }>;


export type CollectCollectionMutation = { __typename?: 'Mutation', collectCollection?: { __typename?: 'Collection', id: string } | null };

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
    override document = AllCollectionsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const GetCollectionDocument = gql`
    query GetCollection($colID: ID!) {
  collectionByID(id: $colID) {
    id
    name
    processSnapshots {
      id
      processID
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class GetCollectionGQL extends Apollo.Query<GetCollectionQuery, GetCollectionQueryVariables> {
    override document = GetCollectionDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const CollectCollectionDocument = gql`
    mutation CollectCollection {
  collectCollection {
    id
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class CollectCollectionGQL extends Apollo.Mutation<CollectCollectionMutation, CollectCollectionMutationVariables> {
    override document = CollectCollectionDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }