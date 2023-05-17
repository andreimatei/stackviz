import { APOLLO_OPTIONS, ApolloModule } from 'apollo-angular';
import { HttpLink } from 'apollo-angular/http';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';
import { NgModule } from '@angular/core';
import { ApolloClientOptions, ApolloLink, InMemoryCache } from '@apollo/client/core';

const basicContext = setContext((_, { headers }) => {
  return {
    headers: {
      ...headers,
      Accept: 'charset=utf-8',
      authorization: `Bearer random token`,
      'Content-Type': 'application/json',
    },
  };
});

const errorLink = onError(({ graphQLErrors, networkError, response }) => {
  // React only on graphql errors
  if (graphQLErrors && graphQLErrors.length > 0) {
    if (
      (graphQLErrors[0] as any)?.statusCode >= 400 &&
      (graphQLErrors[0] as any)?.statusCode < 500
    ) {
      // handle client side error
      console.error(`[Client side error]: ${graphQLErrors[0].message}`);
    } else {
      // handle server side error
      console.error(`[Server side error]: ${graphQLErrors[0].message}`);
    }
  }
  if (networkError) {
    // handle network error
    console.error(`[Network error]: ${networkError.message}`);
  }
});

const uri = '/graphql'; // <-- add the URL of the GraphQL server here
export function createApollo(httpLink: HttpLink): ApolloClientOptions<any> {
  return {
    link: ApolloLink.from([basicContext, errorLink, httpLink.create({ uri })]),
    cache: new InMemoryCache(),
    connectToDevTools: true,  // !!! change it based on environment (dev vs prod)
    assumeImmutableResults: true,
    defaultOptions: {
      watchQuery: {
        errorPolicy: 'all',
      },
    },
  };
}

@NgModule({
  exports: [ApolloModule],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [HttpLink],
    },
  ],
})
export class GraphQLModule {}
