#!/bin/sh
echo "clean.........................up"
rm *.pem
rm *.csr
echo "gen ca"
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
echo "gen server cert"
cfssl gencert -ca=ca.pem  -ca-key=ca-key.pem  -config=ca-config.json -hostname=localhost,127.0.0.1 -profile=massl server-csr.json | cfssljson -bare server
echo "gen client cert"
 cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=massl client-csr.json | cfssljson -bare client