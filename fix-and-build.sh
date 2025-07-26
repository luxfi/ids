#\!/bin/bash
set -e

echo "Fixing circular dependencies and building..."

# First, ensure all modules have their dependencies as replace directives
cd /Users/z/work/lux

# Add replace directives back temporarily for building
echo "Adding temporary replace directives..."

cd ids
echo "replace github.com/luxfi/crypto => ../crypto" >> go.mod

cd ../crypto
echo "replace github.com/luxfi/ids => ../ids" >> go.mod

cd ../database
echo "replace github.com/luxfi/crypto => ../crypto" >> go.mod
echo "replace github.com/luxfi/ids => ../ids" >> go.mod

cd ../evm
echo "replace github.com/luxfi/node => ../node" >> go.mod
echo "replace github.com/luxfi/database => ../database" >> go.mod
echo "replace github.com/luxfi/ids => ../ids" >> go.mod
echo "replace github.com/luxfi/crypto => ../crypto" >> go.mod
echo "replace github.com/luxfi/geth => ../geth" >> go.mod

cd ../node
echo "replace github.com/luxfi/database => ../database" >> go.mod
echo "replace github.com/luxfi/ids => ../ids" >> go.mod
echo "replace github.com/luxfi/crypto => ../crypto" >> go.mod
echo "replace github.com/luxfi/geth => ../geth" >> go.mod
echo "replace github.com/luxfi/evm => ../evm" >> go.mod

cd ../cli
echo "replace github.com/luxfi/node => ../node" >> go.mod
echo "replace github.com/luxfi/geth => ../geth" >> go.mod
echo "replace github.com/luxfi/evm => ../evm" >> go.mod

cd ..

# Now build in dependency order
MODULES=(geth ids crypto database evm node cli)

for module in "${MODULES[@]}"; do
    echo "Building $module..."
    cd "$module"
    go mod tidy
    go build ./...
    echo "✓ $module built successfully"
    cd ..
done

echo "All modules built successfully\!"
