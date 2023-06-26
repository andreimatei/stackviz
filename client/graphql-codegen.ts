import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: "http://localhost:7410/graphql",
  documents: "src/app/graphql/collection.graphql",
  generates: {
    "src/app/graphql/graphql-codegen-generated.ts": {
      plugins: ["typescript", "typescript-operations", "typescript-apollo-angular"],
      config: {
        // https://the-guild.dev/graphql/codegen/plugins/typescript/typescript-apollo-angular#addexplicitoverride
        addExplicitOverride: true,
        scalars: {
          ID: 'number',
        }
      }
    },
  },
};

export default config;
