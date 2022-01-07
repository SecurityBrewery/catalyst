set -e

rm -rf generated

mkdir generated
spruce merge definition/*.yaml >generated/community.yml
spruce merge definition/*.yaml definition/enterprise/*.yaml >generated/catalyst.yml

echo generate caql parser and lexer
cd definition || exit
antlr -Dlanguage=Go -o ../generated/caql/parser CAQLParser.g4 CAQLLexer.g4
antlr -Dlanguage=JavaScript -o ../ui/src/suggestions/grammar CAQLParser.g4 CAQLLexer.g4
cd ..

echo generate json
openapi-generator generate -i generated/community.yml -o generated -g openapi
mv generated/openapi.json generated/community.json
openapi-generator generate -i generated/catalyst.yml -o generated -g openapi
mv generated/openapi.json generated/catalyst.json

# generate python client
# openapi-generator generate -i generated/community.yml -o generated/python -g python --package-name catalystpy --ignore-file-override .openapi-generator-ignore

echo generate server and tests
# go run ./generator/. ./generator
swachigo generated/community.yml generated

echo generate typescript client
openapi-generator generate -i generated/catalyst.yml -o ui/src/client -g typescript-axios --artifact-version 1.0.0-SNAPSHOT

rm -rf gen
rm -rf generated/models/old
rm -rf generated/.openapi-generator generated/.openapi-generator-ignore generated/README.md
# rm -rf generated/python/.openapi-generator generated/python/.gitlab-ci.yml generated/python/git_push.sh generated/python/.travis.yml generated/python/.gitignore generated/python/.openapi-generator-ignore
rm -rf ui/src/client/.openapi-generator ui/src/client/git_push.sh ui/src/client/.gitignore ui/src/client/.openapi-generator-ignore

go mod tidy
gci -w -local "github.com/SecurityBrewery/catalyst" .
