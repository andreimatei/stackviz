import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: "http://localhost:7410/query",
  documents: "src/app/graphql/collection.graphql",
  generates: {
    "src/app/graphql/graphql-codegen-generated.ts": {
      plugins: ["typescript", "typescript-operations", "typescript-apollo-angular"],
      config: {
        addExplicitOverride: true
      }
    },
  },
};

export default config;
