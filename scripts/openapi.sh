#!/bin/bash

set -e

# -------------------------------
# Generate Go server
# -------------------------------
openapi-generator-cli generate \
  -i openapi/crypto-balance-rest.yaml \
  -g go-server \
  -o /tmp/oapi \
  --additional-properties=packageName=cryptobalancerest

rm -rf openapi/servergen
mkdir -p openapi/servergen
mv /tmp/oapi/go/ openapi/servergen
rm -rf /tmp/oapi

# -------------------------------
# Generate Go client
# -------------------------------
openapi-generator-cli generate \
  -i openapi/crypto-balance-rest.yaml \
  -g go \
  -o /tmp/oapi \
  --additional-properties=packageName=cryptobalancerest

rm -rf openapi/clientgen
mkdir -p openapi/clientgen/go
mv /tmp/oapi/*.go openapi/clientgen/go
rm -rf /tmp/oapi

# -------------------------------
# Generate TypeScript client
# -------------------------------
openapi-generator-cli generate \
  -i openapi/crypto-balance-rest.yaml \
  -g typescript-fetch \
  -o /tmp/oapi \
  --additional-properties=npmName=@airgap-solution/crypto-balance-rest-client,supportsES6=true,withInterfaces=true

# Move generated sources into src/
mkdir -p openapi/clientgen/ts/src
rm -rf openapi/clientgen/ts/src/*
mv /tmp/oapi/* openapi/clientgen/ts/src
rm -rf /tmp/oapi

# -------------------------------
# package.json
# -------------------------------
if [ ! -f "openapi/clientgen/ts/package.json" ]; then
    echo "Creating package.json..."
    cat > openapi/clientgen/ts/package.json << 'EOF'
{
  "name": "@airgap-solution/crypto-balance-rest-client",
  "version": "1.0.0",
  "description": "TypeScript client for AirGap Crypto Balance REST API",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "files": [
    "dist/",
    "src/"
  ],
  "scripts": {
    "build": "tsc --declaration",
    "clean": "rm -rf dist/",
    "prepare": "yarn run build"
  },
  "devDependencies": {
    "@openapitools/openapi-generator-cli": "^2.7.0",
    "typescript": "^5.0.0"
  },
  "repository": {
    "type": "git",
    "url": "https://github.com/airgap-solution/crypto-balance-rest.git",
    "directory": "openapi/clientgen/ts"
  },
  "publishConfig": {
    "registry": "https://npm.pkg.github.com"
  }
}
EOF
fi

# -------------------------------
# tsconfig.json
# -------------------------------
if [ ! -f "openapi/clientgen/ts/tsconfig.json" ]; then
    echo "Creating tsconfig.json..."
    cat > openapi/clientgen/ts/tsconfig.json << 'EOF'
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "esnext",
    "lib": ["ES2020", "DOM"],
    "declaration": true,
    "declarationMap": true,
    "outDir": "./dist",
    "rootDir": "./src",            // ensures flattening: src/foo.ts -> dist/foo.js
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "moduleResolution": "node"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
EOF
fi

# -------------------------------
# .gitignore
# -------------------------------
if [ ! -f "openapi/clientgen/ts/.gitignore" ]; then
    echo "Creating .gitignore..."
    cat > openapi/clientgen/ts/.gitignore << 'EOF'
node_modules/
dist/
*.log
.DS_Store
EOF
fi

# -------------------------------
# Install dependencies and build
# -------------------------------
echo "Installing TypeScript client dependencies and building..."
cd openapi/clientgen/ts
yarn install
yarn build

# -------------------------------
# Safety check for typings
# -------------------------------
if [ ! -f "dist/index.d.ts" ]; then
    echo "⚠️ No top-level index.d.ts found, generating a shim..."
    echo "export * from './';" > dist/index.d.ts
fi

cd - > /dev/null

echo "✅ OpenAPI generation complete! TypeScript client built into dist/ with typings."
