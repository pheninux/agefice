GOOS=linux GOARCH=amd64 go build ./cmd/web
cp web agefice
rm -f web
scp agefice root@217.160.188.174:/var/www/go/deploy/agefice/
scp -r  '/mnt/c/Users/a706836/OneDrive - Worldline/Documents/GitHub/agefice/ui' root@217.160.188.174:/var/www/go/deploy/agefice/
scp -r '/mnt/c/Users/a706836/OneDrive - Worldline/Documents/GitHub/agefice/tls' root@217.160.188.174:/var/www/go/deploy/agefice/

#generate certif
#go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

