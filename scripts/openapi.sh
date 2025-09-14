#!/bin/bash

set -e

# Generate Go server
openapi-generator-cli generate \
  -i openapi/crypto-balance-rest.yaml \
  -g go-server \
  -o /tmp/oapi \
  --additional-properties=packageName=cryptobalancerest

rm -rf openapi/servergen
mkdir openapi/servergen
mv /tmp/oapi/go/ openapi/servergen
rm -rf /tmp/oapi

# Generate Go client
openapi-generator-cli generate \
  -i openapi/crypto-balance-rest.yaml \
  -g go \
  -o /tmp/oapi \
  --additional-properties=packageName=cryptobalancerest

rm -rf openapi/clientgen
mkdir -p openapi/clientgen/go
mv /tmp/oapi/*.go openapi/clientgen/go
rm -rf /tmp/oapi

# Generate TypeScript client (fetch-based, works in Expo/React Native)
openapi-generator-cli generate \
  -i openapi/crypto-balance-rest.yaml \
  -g typescript-fetch \
  -o /tmp/oapi \
  --additional-properties=npmName=@airgap-solution/crypto-balance-rest-client,supportsES6=true,withInterfaces=true

mkdir -p openapi/clientgen/ts
rm -rf openapi/clientgen/ts/src
mv /tmp/oapi openapi/clientgen/ts/src
rm -rf /tmp/oapi

# Create package.json if it doesn't exist
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
    "build": "tsc",
    "clean": "rm -rf dist/",
    "prepublishOnly": "yarn run build"
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

# Create tsconfig.json if it doesn't exist
if [ ! -f "openapi/clientgen/ts/tsconfig.json" ]; then
    echo "Creating tsconfig.json..."
    cat > openapi/clientgen/ts/tsconfig.json << 'EOF'
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "esnext",
    "lib": ["ES2020", "DOM"],
    "declaration": true,
    "outDir": "./dist",
    "rootDir": "./src",
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

# Create .gitignore if it doesn't exist
if [ ! -f "openapi/clientgen/ts/.gitignore" ]; then
    echo "Creating .gitignore..."
    cat > openapi/clientgen/ts/.gitignore << 'EOF'
node_modules/
dist/
*.log
.DS_Store
EOF
fi

# Install dependencies and build TypeScript client with Yarn
echo "Installing TypeScript client dependencies and building..."
cd openapi/clientgen/ts
yarn install
yarn build
cd - > /dev/null

echo "OpenAPI generation complete!"
