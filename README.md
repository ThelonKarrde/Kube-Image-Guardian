# Kube-Image-Guardian
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FThelonKarrde%2FKube-Image-Guardian.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FThelonKarrde%2FKube-Image-Guardian?ref=badge_shield)


This is a tool which will help you to allow only certain list of repositories in order to enhance your security.

# Installation how to

First, you need to have your own CA as all communications within kubernetes goes over encrypted connection.  
You can use your own or here is an instruction on how to do it:  

Install [CFSSL](https://github.com/cloudflare/cfssl) Cloudflare's PKI toolkit, it will make life easier.

On mac: 
```
brew install cfssl
```

On Linux:
```
wget -q --show-progress --https-only --timestamping \
  https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 \
  https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
chmod +x cfssl_linux-amd64 cfssljson_linux-amd64
sudo mv cfssl_linux-amd64 /usr/local/bin/cfssl
sudo mv cfssljson_linux-amd64 /usr/local/bin/cfssljson
```

Next you need to build your own CA.  
Create a configuration `ca-config.json` file for it with content similar to this one:
```
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "server": {
        "usages": ["signing", "key encipherment", "server auth", "client auth"],
        "expiry": "8760h"
      }
    }
  }
}
```
Create certificate request for CA `ca-csr.json`:
```
{
  "CN": "Kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "UK",
      "L": "London",
      "O": "Kubernetes",
      "OU": "CA",
      "ST": "London"
    }
  ]
}
```
And generate you CA with a command:
```
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

Then let's create a server certificate.  

---

Create a server side certificate request `server-csr.json`:
```
{
  "CN": "admission",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "UK",
      "L": "London",
      "O": "Kubernetes",
      "OU": "Kubernetes",
      "ST": "London"
    }
  ]
}
```
And generate certificate itself:
```
cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -hostname=kube-image-guardian-webhook.kube-system.svc \
  -profile=server \
  server-csr.json | cfssljson -bare server
```
Note that you have to put to the `-hostname` parameter the address of your internal service if you change service related parameter in helm chart or if you are using none-helm installation.  
As output you will get two files:
```
server-key.pem
server.pem
```
---
Create an encoded version of CA bundle:
```
cat ca.pem | base64 - | tr -d '\n'
```
Put the result into `caBundle` variable of values.yaml file for helm chart.

---

Create a TLS secret for for kube-image-guardian:
```
kubectl -n kube-system create secret tls kube-image-guardian-certs --cert=server.pem --key=server-key.pem
```

Install chart (first with ignore policy):
```
helm install -n kube-system -f ./chart/kube-image-guardian/values.yaml --set validationWebhook.failurePolicy=Ignore kube-image-guardian ./chart/kube-image-guardian/
```

Once everything configured you can use default value for `failurePolicy`:
```
helm upgrade -n kube-system -f ./chart/kube-image-guardian/values.yaml kube-image-guardian ./chart/kube-image-guardian/
```

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FThelonKarrde%2FKube-Image-Guardian.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FThelonKarrde%2FKube-Image-Guardian?ref=badge_large)