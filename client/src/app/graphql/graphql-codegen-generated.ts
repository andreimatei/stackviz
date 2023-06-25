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

export type CollectSpec = Node & {
  __typename?: 'CollectSpec';
  frames?: Maybe<Array<FrameInfo>>;
  id: Scalars['ID'];
};

export type CollectedVar = {
  __typename?: 'CollectedVar';
  Links: Array<Link>;
  Value: Scalars['String'];
};

export type Collection = Node & {
  __typename?: 'Collection';
  id: Scalars['ID'];
  name: Scalars['String'];
  processSnapshots?: Maybe<Array<ProcessSnapshot>>;
};

/**
 * CreateCollectSpecInput is used for create CollectSpec object.
 * Input was generated by ent.
 */
export type CreateCollectSpecInput = {
  frameIDs?: InputMaybe<Array<Scalars['ID']>>;
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
 * CreateFrameInfoInput is used for create FrameInfo object.
 * Input was generated by ent.
 */
export type CreateFrameInfoInput = {
  exprs: Array<Scalars['String']>;
  frame: Scalars['String'];
};

/**
 * CreateProcessSnapshotInput is used for create ProcessSnapshot object.
 * Input was generated by ent.
 */
export type CreateProcessSnapshotInput = {
  framesOfInterest?: InputMaybe<Scalars['String']>;
  processID: Scalars['String'];
  snapshot: Scalars['String'];
};

export type FieldInfo = {
  __typename?: 'FieldInfo';
  Embedded: Scalars['Boolean'];
  Name: Scalars['String'];
  Type: Scalars['String'];
};

export type FrameInfo = Node & {
  __typename?: 'FrameInfo';
  exprs: Array<Scalars['String']>;
  frame: Scalars['String'];
  id: Scalars['ID'];
};

export type GoroutineInfo = {
  __typename?: 'GoroutineInfo';
  Frames: Array<Scalars['String']>;
  ID: Scalars['Int'];
  Vars: Array<CollectedVar>;
};

export type Link = {
  __typename?: 'Link';
  FrameIdx: Scalars['Int'];
  GoroutineID: Scalars['Int'];
  SnapshotID: Scalars['Int'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addExprToCollectSpec: CollectSpec;
  collectCollection?: Maybe<Collection>;
  createCollection?: Maybe<Collection>;
  removeExprFromCollectSpec: CollectSpec;
};


export type MutationAddExprToCollectSpecArgs = {
  expr: Scalars['String'];
  frame: Scalars['String'];
};


export type MutationCreateCollectionArgs = {
  input?: InputMaybe<CreateCollectionInput>;
};


export type MutationRemoveExprFromCollectSpecArgs = {
  expr: Scalars['String'];
  frame: Scalars['String'];
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
  framesOfInterest?: Maybe<Scalars['String']>;
  id: Scalars['ID'];
  processID: Scalars['String'];
  snapshot: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  availableVars: VarsAndTypes;
  collectSpecs: Array<CollectSpec>;
  collectionByID?: Maybe<Collection>;
  collections: Array<Collection>;
  frameInfo?: Maybe<FrameInfo>;
  frameInfos: Array<FrameInfo>;
  goroutines: Array<GoroutineInfo>;
  /** Fetches an object given its ID. */
  node?: Maybe<Node>;
  /** Lookup nodes by a list of IDs. */
  nodes: Array<Maybe<Node>>;
  processSnapshots: Array<ProcessSnapshot>;
  typeInfo: TypeInfo;
};


export type QueryAvailableVarsArgs = {
  func: Scalars['String'];
  pcOff: Scalars['Int'];
};


export type QueryCollectionByIdArgs = {
  id: Scalars['ID'];
};


export type QueryFrameInfoArgs = {
  func: Scalars['String'];
};


export type QueryGoroutinesArgs = {
  colID: Scalars['Int'];
  snapID: Scalars['Int'];
};


export type QueryNodeArgs = {
  id: Scalars['ID'];
};


export type QueryNodesArgs = {
  ids: Array<Scalars['ID']>;
};


export type QueryTypeInfoArgs = {
  name: Scalars['String'];
};

export type TypeInfo = {
  __typename?: 'TypeInfo';
  Fields?: Maybe<Array<Maybe<FieldInfo>>>;
  FieldsNotLoaded: Scalars['Boolean'];
  Name: Scalars['String'];
};

export type VarInfo = {
  __typename?: 'VarInfo';
  FormalParameter: Scalars['Boolean'];
  LoclistAvailable: Scalars['Boolean'];
  Name: Scalars['String'];
  Type: Scalars['String'];
};

export type VarsAndTypes = {
  __typename?: 'VarsAndTypes';
  Types: Array<TypeInfo>;
  Vars: Array<VarInfo>;
};

export type AllCollectionsQueryVariables = Exact<{ [key: string]: never; }>;


export type AllCollectionsQuery = { __typename?: 'Query', collections: Array<{ __typename?: 'Collection', id: string, name: string }> };

export type GetCollectionQueryVariables = Exact<{
  colID: Scalars['ID'];
}>;


export type GetCollectionQuery = { __typename?: 'Query', collectionByID?: { __typename?: 'Collection', id: string, name: string, processSnapshots?: Array<{ __typename?: 'ProcessSnapshot', id: string, processID: string }> | null } | null };

export type CollectCollectionMutationVariables = Exact<{ [key: string]: never; }>;


export type CollectCollectionMutation = { __typename?: 'Mutation', collectCollection?: { __typename?: 'Collection', id: string } | null };

export type AddExprToCollectSpecMutationVariables = Exact<{
  frame: Scalars['String'];
  expr: Scalars['String'];
}>;


export type AddExprToCollectSpecMutation = { __typename?: 'Mutation', addExprToCollectSpec: { __typename?: 'CollectSpec', frames?: Array<{ __typename?: 'FrameInfo', frame: string, exprs: Array<string> }> | null } };

export type RemoveExprFromCollectSpecMutationVariables = Exact<{
  frame: Scalars['String'];
  expr: Scalars['String'];
}>;


export type RemoveExprFromCollectSpecMutation = { __typename?: 'Mutation', removeExprFromCollectSpec: { __typename?: 'CollectSpec', frames?: Array<{ __typename?: 'FrameInfo', frame: string, exprs: Array<string> }> | null } };

export type GetAvailableVariablesQueryVariables = Exact<{
  func: Scalars['String'];
  pcOff: Scalars['Int'];
}>;


export type GetAvailableVariablesQuery = { __typename?: 'Query', availableVars: { __typename?: 'VarsAndTypes', Vars: Array<{ __typename?: 'VarInfo', Name: string, Type: string, FormalParameter: boolean, LoclistAvailable: boolean }>, Types: Array<{ __typename?: 'TypeInfo', Name: string, FieldsNotLoaded: boolean, Fields?: Array<{ __typename?: 'FieldInfo', Name: string, Type: string, Embedded: boolean } | null> | null }> }, frameInfo?: { __typename?: 'FrameInfo', exprs: Array<string> } | null };

export type GetTypeInfoQueryVariables = Exact<{
  name: Scalars['String'];
}>;


export type GetTypeInfoQuery = { __typename?: 'Query', typeInfo: { __typename?: 'TypeInfo', Name: string, FieldsNotLoaded: boolean, Fields?: Array<{ __typename?: 'FieldInfo', Name: string, Type: string, Embedded: boolean } | null> | null } };

export type GetGoroutinesQueryVariables = Exact<{
  colID: Scalars['Int'];
  snapID: Scalars['Int'];
}>;


export type GetGoroutinesQuery = { __typename?: 'Query', goroutines: Array<{ __typename?: 'GoroutineInfo', ID: number, Frames: Array<string>, Vars: Array<{ __typename?: 'CollectedVar', Value: string, Links: Array<{ __typename?: 'Link', SnapshotID: number, GoroutineID: number, FrameIdx: number }> }> }> };

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
export const AddExprToCollectSpecDocument = gql`
    mutation AddExprToCollectSpec($frame: String!, $expr: String!) {
  addExprToCollectSpec(frame: $frame, expr: $expr) {
    frames {
      frame
      exprs
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class AddExprToCollectSpecGQL extends Apollo.Mutation<AddExprToCollectSpecMutation, AddExprToCollectSpecMutationVariables> {
    override document = AddExprToCollectSpecDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const RemoveExprFromCollectSpecDocument = gql`
    mutation RemoveExprFromCollectSpec($frame: String!, $expr: String!) {
  removeExprFromCollectSpec(frame: $frame, expr: $expr) {
    frames {
      frame
      exprs
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class RemoveExprFromCollectSpecGQL extends Apollo.Mutation<RemoveExprFromCollectSpecMutation, RemoveExprFromCollectSpecMutationVariables> {
    override document = RemoveExprFromCollectSpecDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const GetAvailableVariablesDocument = gql`
    query GetAvailableVariables($func: String!, $pcOff: Int!) {
  availableVars(func: $func, pcOff: $pcOff) {
    Vars {
      Name
      Type
      FormalParameter
      LoclistAvailable
    }
    Types {
      Name
      Fields {
        Name
        Type
        Embedded
      }
      FieldsNotLoaded
    }
  }
  frameInfo(func: $func) {
    exprs
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class GetAvailableVariablesGQL extends Apollo.Query<GetAvailableVariablesQuery, GetAvailableVariablesQueryVariables> {
    override document = GetAvailableVariablesDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const GetTypeInfoDocument = gql`
    query GetTypeInfo($name: String!) {
  typeInfo(name: $name) {
    Name
    FieldsNotLoaded
    Fields {
      Name
      Type
      Embedded
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class GetTypeInfoGQL extends Apollo.Query<GetTypeInfoQuery, GetTypeInfoQueryVariables> {
    override document = GetTypeInfoDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const GetGoroutinesDocument = gql`
    query GetGoroutines($colID: Int!, $snapID: Int!) {
  goroutines(colID: $colID, snapID: $snapID) {
    ID
    Frames
    Vars {
      Value
      Links {
        SnapshotID
        GoroutineID
        FrameIdx
      }
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class GetGoroutinesGQL extends Apollo.Query<GetGoroutinesQuery, GetGoroutinesQueryVariables> {
    override document = GetGoroutinesDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }