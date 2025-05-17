PATH=$(pwd)

if [ -f /go/bin/boring-registry ]; then
    BORING_REGISTR=/go/bin/boring-registry
elif [ -f ~/go/bin/boring-registry ]; then
    BORING_REGISTR=~/go/bin/boring-registry
elif [ -x "$(which boring-registry)" ]; then
    BORING_REGISTR=$(which boring-registry)
else
    echo "there is no any known binary of boring-registry found!"
fi

$BORING_REGISTR upload provider sotoon \
    --filename-sha256sums $PATH/dist/terraform-provider-sotoon_*_SHA256SUMS \
    --namespace sotoon \
    --storage-s3-bucket=terraform-registry \
    --storage-s3-endpoint=https://s3.thr1.sotoon.ir \
    --storage-s3-region=US-East \
    --storage-s3-pathstyle=true
