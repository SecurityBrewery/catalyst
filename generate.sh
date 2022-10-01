set -e

rm -rf generated

mkdir generated
spruce merge definition/*.yaml >generated/community.yml
spruce merge definition/*.yaml definition/enterprise/*.yaml >generated/catalyst.yml

echo generate caql parser and lexer
cd definition || exit
# antlr 4.10.1
antlr -Dlanguage=Go -o ../generated/caql/parser CAQLParser.g4 CAQLLexer.g4
antlr -Dlanguage=JavaScript -o ../ui/src/suggestions/grammar CAQLParser.g4 CAQLLexer.g4
cd ..

echo generate json
# openapi-generator 6.0.0
openapi-generator generate -i generated/community.yml -o generated -g openapi
mv generated/openapi.json generated/community.json
openapi-generator generate -i generated/catalyst.yml -o generated -g openapi
mv generated/openapi.json generated/catalyst.json

echo generate server and tests
swagger-go-chi generated/community.yml generated
rm -rf generated/auth generated/cli
find generated -type f -name "*.go" -print0 | xargs -0 sed -i '' -e 's#"github.com/go-chi/chi"#"github.com/go-chi/chi/v5"#g'

echo generate typescript client
openapi-generator generate -i generated/catalyst.yml -o ui/src/client -g typescript-axios --artifact-version 1.0.0-SNAPSHOT

rm -rf gen
rm -rf generated/models/old
rm -rf generated/.openapi-generator generated/.openapi-generator-ignore generated/README.md
rm -rf ui/src/client/.openapi-generator ui/src/client/git_push.sh ui/src/client/.gitignore ui/src/client/.openapi-generator-ignore

go mod tidy
gci write --Section Standard --Section Default --Section "Prefix(github.com/SecurityBrewery/catalyst)" .
cd internal/maut
gci write --Section Standard --Section Default --Section "Prefix(github.com/jonas-plum/maut)" .
cd ../..
