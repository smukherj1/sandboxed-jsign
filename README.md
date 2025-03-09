
# Dependencies

```
# Jsign 7.1 signing utility JAR.
wget https://github.com/ebourg/jsign/releases/download/7.1/jsign-7.1.jar

# gVisor (https://gvisor.dev/docs/user_guide/install/#install-latest)
ARCH=$(uname -m)
RELEASE="20250224.0"
URL=https://storage.googleapis.com/gvisor/releases/release/${RELEASE}/${ARCH}

wget ${URL}/runsc ${URL}/runsc.sha512 \
 ${URL}/containerd-shim-runsc-v1 ${URL}/containerd-shim-runsc-v1.sha512

sha512sum -c runsc.sha512 \
 -c containerd-shim-runsc-v1.sha512

rm -f *.sha512
chmod a+rx runsc containerd-shim-runsc-v1
```