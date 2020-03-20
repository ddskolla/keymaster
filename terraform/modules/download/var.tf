
variable "local_path" {
  type = string
  description = "Local path for downloaded file"
}

variable "file_mode" {
  type = string
  description = "UNIX permissions to apply to downloaded file"
  default = "0755"
}

variable "urls" {
  type = map(string)
  description = "A map from platform type to download URL"
}

variable "checksums" {
  type = map(string)
  description = "A map from platform type to expected SHA256 checksum"
}
