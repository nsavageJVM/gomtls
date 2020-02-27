# Mutual TLS Starter
### Under construction

go get -u github.com/cloudflare/cfssl/cmd/cfssl
go get -u github.com/cloudflare/cfssl/cmd/cfssljson
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
cfssl gencert -ca=ca.pem  -ca-key=ca-key.pem  -config=ca-config.json \
 -hostname=localhost,127.0.0.1 -profile=massl server-csr.json | cfssljson -bare server

cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=massl \ client-csr.json | cfssljson -bare client
