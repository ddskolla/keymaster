
locals {
  download_command = <<EOM
case "$OSTYPE" in
  darwin*)
    URL="${var.urls["darwin"]}"
    CHECKSUM="${var.checksums["darwin"]}"
    SHA_CMD="shasum -a 256" ;;
  linux*)
    URL="${var.urls["linux"]}"
    CHECKSUM="${var.checksums["linux"]}"
    SHA_CMD="sha256sum" ;;
  *)
    echo "Unsupported platform for file download"
    exit 1 ;;
esac

LOCAL_FILE="${var.local_path}"
FILE_MODE="${var.file_mode}"

# Download file if missing or wrong checksum
if $SHA_CMD -c <(echo "$CHECKSUM  $LOCAL_FILE") >/dev/null 2>&1; then
  echo "File was up to date: $LOCAL_FILE"
else
  echo "Local file needs refresh: $LOCAL_FILE"
  echo "Downloading: $URL"
  curl --fail -L "$URL" -o "$LOCAL_FILE"
  echo "Expecting checksum: $CHECKSUM"
  $SHA_CMD -c <(echo "$CHECKSUM  $LOCAL_FILE")
  chmod "$FILE_MODE" $LOCAL_FILE
  echo "Download completed."
fi
EOM
}

resource "null_resource" "download" {
  triggers = {
    command = sha256(local.download_command)
    # This isn't ideal, it will only converge on the second apply. The uuid() call
    # here is to force a state update with the actual hash after the download.
    exists = fileexists(var.local_path) ? filesha256(var.local_path) : uuid()
  }
  provisioner "local-exec" {
    interpreter = ["/bin/bash", "-e", "-c"]
    command = local.download_command
  }
}
