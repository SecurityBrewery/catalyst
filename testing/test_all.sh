reset_catalyst_data() {
    if [ -d "catalyst_data" ]; then
        rm -rf catalyst_data
    fi
    mkdir -p catalyst_data
}

# function to run playwright tests
run_playwright_tests() {
    cd ui
    bun test:e2e
    cd ..
}

run_upgradetest() {
    # iterate over all folders in ./testing/data
    for dir in ./testing/data/*/; do
        echo "Running tests with data from $dir"
        reset_catalyst_data
        cp "$dir"/data.db catalyst_data/data.db
        run_playwright_tests
    done
}

run_upgradetest
