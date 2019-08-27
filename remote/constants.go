package remote

const PublicKeyExtension = ".pub"
const PrivateKeyExtension = ".pem"
const HostPublicKeyExtension = ".hostkey"

const DefaultSSHPort = 22

// Remote Temporary Path
const TempPath = "/tmp/osr"

// Remote Resources Path (Relative to home of user)
const ResourcesPath = "resources/"

// Maximum capacity in % of disks before alert it.
const CapacityWarning = 90
