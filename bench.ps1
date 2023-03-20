Write-Host "Running benchmarks"

Write-Host "Running alpine benchmarks"
docker run --rm -it $(docker build -q -f .\Dockerfile_alpine .) > _alpine_results.txt

Write-Host "Running ubuntu benchmarks"
docker run --rm -it $(docker build -q -f .\Dockerfile_ubuntu .) > _ubuntu_results.txt

Write-Host "Running debian benchmarks"
docker run --rm -it $(docker build -q -f .\Dockerfile_debian .) > _debian_results.txt
