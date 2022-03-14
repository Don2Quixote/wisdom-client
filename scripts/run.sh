source scripts/colors.sh

check() {
    if test $? -ne 0
    then
    	printf "$*"
        exit 1
    fi

}

sudo docker build -f build/Dockerfile -t wisdom-client . > /dev/null
check "${RED}failed to build image${NC}\n"

sudo docker run --name wisdom-client --rm --network=host \
    -e WOW_HOST=localhost:4444 \
    -e MAX_COMPLEXITY=16 \
    wisdom-client
check "${RED}failed to launch container${NC}\n"

