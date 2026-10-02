#!/bin/sh
# Pass an absolute path to a built CLI. No database is contacted.
set -eu
cli=${1:?Usage: sh scripts/smoke.sh /absolute/path/to/grimleytk}
task_dir=$(mktemp -d)
trap 'rm -rf "$task_dir"' EXIT
cd "$task_dir"
"$cli" init
"$cli" create domain catalog --schema catalog --owner catalog-service
"$cli" create table catalog.products --description 'Main product table'
"$cli" create column catalog.products.id --type uuid --primary-key --nullable=false
"$cli" create column catalog.products.price --type numeric --nullable=false
"$cli" create domain wishlist --schema wishlist --owner wishlist-service
"$cli" create view wishlist.products_view --from catalog.products --columns id,price
"$cli" validate
"$cli" plan > first.plan
"$cli" plan > second.plan
cmp first.plan second.plan
cat first.plan
"$cli" show
"$cli" show domains
"$cli" show tables
"$cli" show reads
cp grimley.yaml original.yaml
if "$cli" init; then echo 'init unexpectedly overwrote config'; exit 1; fi
cmp grimley.yaml original.yaml
if "$cli" create view wishlist.products_view --from catalog.products --columns id; then
 echo 'create view unexpectedly overwrote view'; exit 1
fi
cmp grimley.yaml original.yaml
"$cli" apply --help > /dev/null
printf '\nCLI smoke test passed (no database connection).\n'
