# solc --abi --overwrite *.sol -o build 
# solc --bin --overwrite *.sol -o build
# for filename in ./build/*.abi; do
#      BASE=$(basename $filename)
#      echo $BASE
#      GOFILE=${BASE%.abi}
#      echo $GOFILE
#      abigen --abi=$filename --pkg=erc721 --out=./$GOFILE.go
# done
abigen --abi=./ERC721.abi --pkg=erc721 --out=ERC721.go
abigen --abi=./ERC721Enumerable.abi --pkg=erc721 --out=ERC721Enumerable.go
abigen --abi=./ERC721Metadata.abi --pkg=erc721 --out=ERC721Metadata.go